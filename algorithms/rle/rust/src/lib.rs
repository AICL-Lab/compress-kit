use compresskit_codec::codec::{CodecError, Decoder, Encoder, State};

fn write_u32_le(buf: &mut Vec<u8>, v: u32) {
    buf.extend_from_slice(&v.to_le_bytes());
}

fn read_u32_le(data: &[u8], pos: &mut usize) -> Option<u32> {
    if *pos + 4 > data.len() {
        return None;
    }
    let v = u32::from_le_bytes([data[*pos], data[*pos + 1], data[*pos + 2], data[*pos + 3]]);
    *pos += 4;
    Some(v)
}

pub fn encode_memory(input: &[u8]) -> Result<Vec<u8>, CodecError> {
    let mut output = Vec::new();
    if input.is_empty() {
        return Ok(output);
    }

    let mut current = input[0];
    let mut count: u32 = 1;

    for &b in &input[1..] {
        if b == current && count < u32::MAX {
            count += 1;
        } else {
            write_u32_le(&mut output, count);
            output.push(current);
            current = b;
            count = 1;
        }
    }

    write_u32_le(&mut output, count);
    output.push(current);

    Ok(output)
}

pub fn decode_memory(input: &[u8]) -> Result<Vec<u8>, CodecError> {
    let mut output = Vec::new();
    let mut pos = 0;

    while pos < input.len() {
        let count = read_u32_le(input, &mut pos).ok_or(CodecError::Truncated)?;
        if count == 0 {
            return Err(CodecError::Corrupt);
        }
        if pos >= input.len() {
            return Err(CodecError::Truncated);
        }
        let value = input[pos];
        pos += 1;

        if output.len() as u64 + count as u64 > compresskit_codec::codec::MAX_OUTPUT_SIZE as u64 {
            return Err(CodecError::SizeLimit);
        }

        for _ in 0..count {
            output.push(value);
        }
    }

    Ok(output)
}

pub struct StreamingEncoder {
    state: State,
    buffer: Vec<u8>,
}

impl StreamingEncoder {
    pub fn new() -> Self {
        StreamingEncoder {
            state: State::Ready,
            buffer: Vec::new(),
        }
    }
}

impl Default for StreamingEncoder {
    fn default() -> Self {
        Self::new()
    }
}

impl Encoder for StreamingEncoder {
    fn process(&mut self, input: &[u8], _output: &mut [u8]) -> Result<usize, CodecError> {
        match self.state {
            State::Ready | State::Flushing => {
                self.state = State::Streaming;
                self.buffer.extend_from_slice(input);
                Ok(0)
            }
            State::Streaming => {
                self.buffer.extend_from_slice(input);
                Ok(0)
            }
            State::Finished => {
                self.state = State::Error;
                Err(CodecError::InvalidState)
            }
            State::Error => Err(CodecError::InvalidState),
        }
    }

    fn flush(&mut self, _output: &mut [u8]) -> Result<usize, CodecError> {
        match self.state {
            State::Ready => Ok(0),
            State::Streaming => {
                self.state = State::Flushing;
                Ok(0)
            }
            State::Flushing => Ok(0),
            State::Finished => {
                self.state = State::Error;
                Err(CodecError::InvalidState)
            }
            State::Error => Err(CodecError::InvalidState),
        }
    }

    fn finish(&mut self, output: &mut [u8]) -> Result<usize, CodecError> {
        match self.state {
            State::Ready | State::Streaming | State::Flushing => {
                let encoded = encode_memory(&self.buffer)?;
                if output.len() < encoded.len() {
                    return Err(CodecError::BufTooSmall);
                }
                output[..encoded.len()].copy_from_slice(&encoded);
                self.state = State::Finished;
                Ok(encoded.len())
            }
            State::Finished => {
                self.state = State::Error;
                Err(CodecError::InvalidState)
            }
            State::Error => Err(CodecError::InvalidState),
        }
    }

    fn reset(&mut self) {
        self.state = State::Ready;
        self.buffer.clear();
    }

    fn state(&self) -> State {
        self.state
    }
}

pub struct StreamingDecoder {
    state: State,
    buffer: Vec<u8>,
}

impl StreamingDecoder {
    pub fn new() -> Self {
        StreamingDecoder {
            state: State::Ready,
            buffer: Vec::new(),
        }
    }
}

impl Default for StreamingDecoder {
    fn default() -> Self {
        Self::new()
    }
}

impl Decoder for StreamingDecoder {
    fn process(&mut self, input: &[u8], _output: &mut [u8]) -> Result<usize, CodecError> {
        match self.state {
            State::Ready | State::Flushing => {
                self.state = State::Streaming;
                self.buffer.extend_from_slice(input);
                Ok(0)
            }
            State::Streaming => {
                self.buffer.extend_from_slice(input);
                Ok(0)
            }
            State::Finished => {
                self.state = State::Error;
                Err(CodecError::InvalidState)
            }
            State::Error => Err(CodecError::InvalidState),
        }
    }

    fn flush(&mut self, _output: &mut [u8]) -> Result<usize, CodecError> {
        match self.state {
            State::Ready => Ok(0),
            State::Streaming => {
                self.state = State::Flushing;
                Ok(0)
            }
            State::Flushing => Ok(0),
            State::Finished => {
                self.state = State::Error;
                Err(CodecError::InvalidState)
            }
            State::Error => Err(CodecError::InvalidState),
        }
    }

    fn finish(&mut self, output: &mut [u8]) -> Result<usize, CodecError> {
        match self.state {
            State::Ready | State::Streaming | State::Flushing => {
                let decoded = decode_memory(&self.buffer)?;
                if output.len() < decoded.len() {
                    return Err(CodecError::BufTooSmall);
                }
                output[..decoded.len()].copy_from_slice(&decoded);
                self.state = State::Finished;
                Ok(decoded.len())
            }
            State::Finished => {
                self.state = State::Error;
                Err(CodecError::InvalidState)
            }
            State::Error => Err(CodecError::InvalidState),
        }
    }

    fn reset(&mut self) {
        self.state = State::Ready;
        self.buffer.clear();
    }

    fn state(&self) -> State {
        self.state
    }
}
