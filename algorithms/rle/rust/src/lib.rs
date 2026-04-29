// Minimal in-memory RLE encode/decode API for testing.
// Main CLI remains in main.rs. This library provides simple helpers.

use std::io;

const MAX_OUTPUT_SIZE: usize = 1024 * 1024 * 1024;

pub fn encode(input: &[u8]) -> Result<Vec<u8>, io::Error> {
    if input.is_empty() {
        return Ok(Vec::new());
    }

    let mut output = Vec::new();
    let mut current = input[0];
    let mut count: u32 = 1;

    for &b in &input[1..] {
        if b == current && count < u32::MAX {
            count += 1;
        } else {
            output.extend_from_slice(&count.to_le_bytes());
            output.push(current);
            current = b;
            count = 1;
        }
    }

    output.extend_from_slice(&count.to_le_bytes());
    output.push(current);
    Ok(output)
}

pub fn decode(input: &[u8]) -> Result<Vec<u8>, io::Error> {
    let mut output = Vec::new();
    let mut pos = 0;

    while pos < input.len() {
        if pos + 5 > input.len() {
            return Err(io::Error::new(
                io::ErrorKind::InvalidData,
                "RLE data truncated: incomplete count+value pair",
            ));
        }

        let count =
            u32::from_le_bytes([input[pos], input[pos + 1], input[pos + 2], input[pos + 3]]);
        let value = input[pos + 4];
        pos += 5;

        let new_len = output.len() + count as usize;
        if new_len > MAX_OUTPUT_SIZE {
            return Err(io::Error::new(
                io::ErrorKind::InvalidData,
                "output size exceeds maximum allowed (1 GiB)",
            ));
        }

        for _ in 0..count {
            output.push(value);
        }
    }

    Ok(output)
}

// Streaming adapters
use compresskit_codec::codec::{
    CodecError, Decoder as CodecDecoder, Encoder as CodecEncoder, State, MAX_INPUT_SIZE,
};

fn io_error_to_codec_error(e: io::Error) -> CodecError {
    match e.kind() {
        io::ErrorKind::UnexpectedEof => CodecError::Truncated,
        io::ErrorKind::InvalidData => {
            let msg = e.to_string();
            if msg.contains("truncated") || msg.contains("too short") {
                CodecError::Truncated
            } else if msg.contains("limit") || msg.contains("exceeds") {
                CodecError::SizeLimit
            } else if msg.contains("invalid") || msg.contains("bad") {
                CodecError::Corrupt
            } else {
                CodecError::Other(msg)
            }
        }
        _ => CodecError::Other(e.to_string()),
    }
}

pub struct StreamingEncoder {
    state: State,
    buffer: Vec<u8>,
    total_input: usize,
}

impl StreamingEncoder {
    pub fn new() -> Self {
        StreamingEncoder {
            state: State::Ready,
            buffer: Vec::new(),
            total_input: 0,
        }
    }
}

impl Default for StreamingEncoder {
    fn default() -> Self {
        Self::new()
    }
}

impl CodecEncoder for StreamingEncoder {
    fn process(&mut self, input: &[u8], _output: &mut [u8]) -> Result<usize, CodecError> {
        match self.state {
            State::Ready | State::Flushing => {
                if self.total_input > MAX_INPUT_SIZE.saturating_sub(input.len()) {
                    self.state = State::Error;
                    return Err(CodecError::SizeLimit);
                }
                self.state = State::Streaming;
                self.buffer.extend_from_slice(input);
                self.total_input += input.len();
                Ok(0)
            }
            State::Streaming => {
                if self.total_input > MAX_INPUT_SIZE.saturating_sub(input.len()) {
                    self.state = State::Error;
                    return Err(CodecError::SizeLimit);
                }
                self.buffer.extend_from_slice(input);
                self.total_input += input.len();
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
                let encoded = encode(&self.buffer).map_err(io_error_to_codec_error)?;
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
        self.total_input = 0;
    }

    fn state(&self) -> State {
        self.state
    }
}

pub struct StreamingDecoder {
    state: State,
    buffer: Vec<u8>,
    total_input: usize,
}

impl StreamingDecoder {
    pub fn new() -> Self {
        StreamingDecoder {
            state: State::Ready,
            buffer: Vec::new(),
            total_input: 0,
        }
    }
}

impl Default for StreamingDecoder {
    fn default() -> Self {
        Self::new()
    }
}

impl CodecDecoder for StreamingDecoder {
    fn process(&mut self, input: &[u8], _output: &mut [u8]) -> Result<usize, CodecError> {
        match self.state {
            State::Ready | State::Flushing => {
                if self.total_input > MAX_INPUT_SIZE.saturating_sub(input.len()) {
                    self.state = State::Error;
                    return Err(CodecError::SizeLimit);
                }
                self.state = State::Streaming;
                self.buffer.extend_from_slice(input);
                self.total_input += input.len();
                Ok(0)
            }
            State::Streaming => {
                if self.total_input > MAX_INPUT_SIZE.saturating_sub(input.len()) {
                    self.state = State::Error;
                    return Err(CodecError::SizeLimit);
                }
                self.buffer.extend_from_slice(input);
                self.total_input += input.len();
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
                let decoded = decode(&self.buffer).map_err(io_error_to_codec_error)?;
                if decoded.len() > MAX_OUTPUT_SIZE {
                    self.state = State::Error;
                    return Err(CodecError::SizeLimit);
                }
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
        self.total_input = 0;
    }

    fn state(&self) -> State {
        self.state
    }
}
