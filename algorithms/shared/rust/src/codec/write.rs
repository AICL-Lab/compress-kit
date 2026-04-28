use crate::codec::encoder::Encoder;
use crate::codec::error::{CodecError, MAX_OUTPUT_SIZE};
use std::io::{self, Write};

/// WriteEncoder wraps an Encoder to implement io::Write.
/// Calls to write() are forwarded to encoder.process().
/// Call flush() or drop to invoke encoder.finish() and flush final output.
pub struct WriteEncoder<'a, W: Write> {
    encoder: &'a mut dyn Encoder,
    writer: W,
    out_buf: Vec<u8>,
    out_buf_offset: usize,
}

impl<'a, W: Write> WriteEncoder<'a, W> {
    /// Creates a WriteEncoder that wraps the given encoder
    /// and writes encoded output to w.
    pub fn new(encoder: &'a mut dyn Encoder, writer: W) -> Self {
        WriteEncoder {
            encoder,
            writer,
            out_buf: vec![0u8; 64 * 1024], // 64KB output buffer
            out_buf_offset: 0,
        }
    }

    fn grow_buffer(current_len: usize, limit: usize) -> usize {
        if current_len == 0 {
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

    fn flush_out_buf(&mut self) -> io::Result<()> {
        if self.out_buf_offset > 0 {
            self.writer.write_all(&self.out_buf[..self.out_buf_offset])?;
            self.out_buf_offset = 0;
        }
        Ok(())
    }

    /// Finishes encoding and flushes all output.
    pub fn close(&mut self) -> io::Result<()> {
        loop {
            match self
                .encoder
                .finish(&mut self.out_buf[self.out_buf_offset..])
            {
                Ok(written) => {
                    self.out_buf_offset += written;
                    break;
                }
                Err(CodecError::BufTooSmall) => {
                    self.flush_out_buf()?;
                    let new_size = Self::grow_buffer(self.out_buf.len(), MAX_OUTPUT_SIZE);
                    if new_size <= self.out_buf.len() {
                        return Err(io::Error::new(
                            io::ErrorKind::Other,
                            CodecError::SizeLimit,
                        ));
                    }
                    self.out_buf.resize(new_size, 0);
                    self.out_buf_offset = 0;
                }
                Err(e) => return Err(io::Error::new(io::ErrorKind::Other, e)),
            }
        }

        self.flush_out_buf()
    }
}

impl<'a, W: Write> Write for WriteEncoder<'a, W> {
    fn write(&mut self, buf: &[u8]) -> io::Result<usize> {
        // Process input through encoder
        loop {
            match self
                .encoder
                .process(buf, &mut self.out_buf[self.out_buf_offset..])
            {
                Ok(written) => {
                    self.out_buf_offset += written;
                    break;
                }
                Err(CodecError::BufTooSmall) => {
                    self.flush_out_buf()?;
                    let new_size = Self::grow_buffer(self.out_buf.len(), MAX_OUTPUT_SIZE);
                    if new_size <= self.out_buf.len() {
                        return Err(io::Error::new(
                            io::ErrorKind::Other,
                            CodecError::SizeLimit,
                        ));
                    }
                    self.out_buf.resize(new_size, 0);
                    self.out_buf_offset = 0;
                }
                Err(e) => return Err(io::Error::new(io::ErrorKind::Other, e)),
            }
        }

        // Flush output buffer if getting full
        if self.out_buf_offset > self.out_buf.len() / 2 {
            self.flush_out_buf()?;
        }

        Ok(buf.len())
    }

    fn flush(&mut self) -> io::Result<()> {
        self.flush_out_buf()?;
        self.writer.flush()
    }
}
