use crate::codec::error::CodecError;

pub(crate) fn grow_buffer(current_len: usize, limit: usize) -> usize {
    if current_len < 1024 {
        if limit < 1024 {
            return limit;
        }
        return 1024;
    }

    let next = current_len.saturating_mul(2);
    if next < current_len {
        return limit;
    }
    if next > limit {
        return limit;
    }
    next
}

pub(crate) struct ResizingBuffer {
    buf: Vec<u8>,
    written: usize,
    limit: usize,
}

impl ResizingBuffer {
    pub(crate) fn new(initial_size: usize, limit: usize) -> Self {
        let size = initial_size.min(limit);
        Self {
            buf: vec![0; size],
            written: 0,
            limit,
        }
    }

    #[cfg(test)]
    pub(crate) fn capacity(&self) -> usize {
        self.buf.len()
    }

    pub(crate) fn run<F>(&mut self, step: &mut F) -> Result<(), CodecError>
    where
        F: FnMut(&mut [u8]) -> Result<usize, CodecError>,
    {
        loop {
            match step(&mut self.buf[self.written..]) {
                Ok(n) => {
                    self.written = self.written.checked_add(n).ok_or(CodecError::SizeLimit)?;
                    if self.written > self.limit {
                        return Err(CodecError::SizeLimit);
                    }
                    return Ok(());
                }
                Err(CodecError::BufTooSmall) => {
                    if self.written > self.limit || self.buf.len() >= self.limit {
                        return Err(CodecError::SizeLimit);
                    }

                    let new_size = grow_buffer(self.buf.len(), self.limit);
                    if new_size <= self.buf.len() {
                        return Err(CodecError::SizeLimit);
                    }

                    self.buf.resize(new_size, 0);
                }
                Err(err) => return Err(err),
            }
        }
    }

    pub(crate) fn into_vec(mut self) -> Vec<u8> {
        self.buf.truncate(self.written);
        self.buf
    }
}

#[cfg(test)]
mod tests {
    use super::ResizingBuffer;
    use crate::codec::error::CodecError;

    #[test]
    fn resizing_buffer_retries_after_buf_too_small() {
        let mut runner = ResizingBuffer::new(1, 8);
        let mut calls = 0;

        runner
            .run(&mut |out| {
                calls += 1;
                if calls == 1 {
                    return Err(CodecError::BufTooSmall);
                }
                out[..3].copy_from_slice(b"def");
                Ok(3)
            })
            .unwrap();

        assert_eq!(runner.into_vec(), b"def");
    }

    #[test]
    fn resizing_buffer_clamps_initial_size_to_limit() {
        let runner = ResizingBuffer::new(8, 3);
        assert_eq!(runner.capacity(), 3);
    }
}
