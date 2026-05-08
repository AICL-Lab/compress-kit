//! Buffered encoder/decoder implementations.
//!
//! These are deep modules that implement the Encoder/Decoder traits using
//! single-shot encode/decode functions. They handle all state machine logic,
//! input buffering, and output management, allowing algorithm implementations
//! to focus solely on the encoding/decoding algorithm.
//!
//! This is the preferred way to implement streaming adapters for algorithms
//! that require complete input (Huffman, Arithmetic, Range, RLE, etc.).

use crate::codec::encoder::{Decoder, Encoder};
use crate::codec::error::{CodecError, State, MAX_INPUT_SIZE, MAX_OUTPUT_SIZE};

/// Function signature for single-shot encoding.
pub type EncodeFunc = fn(input: &[u8]) -> Result<Vec<u8>, CodecError>;

/// Function signature for single-shot decoding.
pub type DecodeFunc = fn(input: &[u8]) -> Result<Vec<u8>, CodecError>;

/// BufferedEncoder implements the Encoder trait using a single-shot encode function.
///
/// It handles all state machine logic, input buffering, and output management.
/// Algorithm implementations only need to provide the core encode function.
///
/// # Example
///
/// ```ignore
/// fn huffman_encode(input: &[u8]) -> Result<Vec<u8>, CodecError> {
///     // ... encoding logic ...
/// }
///
/// let encoder = BufferedEncoder::new(huffman_encode);
/// ```
pub struct BufferedEncoder {
    state: State,
    buffer: Vec<u8>,
    total_input: usize,
    encode: EncodeFunc,
}

impl BufferedEncoder {
    /// Creates a new buffered encoder backed by a single-shot encode function.
    ///
    /// The encode function will be called once during `finish()` with all buffered input.
    pub fn new(encode: EncodeFunc) -> Self {
        BufferedEncoder {
            state: State::Ready,
            buffer: Vec::new(),
            total_input: 0,
            encode,
        }
    }
}

impl Encoder for BufferedEncoder {
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
                let encoded = (self.encode)(&self.buffer)?;
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

/// BufferedDecoder implements the Decoder trait using a single-shot decode function.
///
/// It handles all state machine logic, input buffering, and output management.
/// Algorithm implementations only need to provide the core decode function.
///
/// # Example
///
/// ```ignore
/// fn huffman_decode(input: &[u8]) -> Result<Vec<u8>, CodecError> {
///     // ... decoding logic ...
/// }
///
/// let decoder = BufferedDecoder::new(huffman_decode);
/// ```
pub struct BufferedDecoder {
    state: State,
    buffer: Vec<u8>,
    total_input: usize,
    decode: DecodeFunc,
}

impl BufferedDecoder {
    /// Creates a new buffered decoder backed by a single-shot decode function.
    ///
    /// The decode function will be called once during `finish()` with all buffered input.
    pub fn new(decode: DecodeFunc) -> Self {
        BufferedDecoder {
            state: State::Ready,
            buffer: Vec::new(),
            total_input: 0,
            decode,
        }
    }
}

impl Decoder for BufferedDecoder {
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
                let decoded = (self.decode)(&self.buffer)?;
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

#[cfg(test)]
mod tests {
    use super::*;

    fn test_encode(input: &[u8]) -> Result<Vec<u8>, CodecError> {
        // Simple identity encode for testing
        let mut output = Vec::new();
        output.extend_from_slice(b"ENC");
        output.extend_from_slice(input);
        Ok(output)
    }

    fn test_decode(input: &[u8]) -> Result<Vec<u8>, CodecError> {
        if input.len() < 3 || &input[0..3] != b"ENC" {
            return Err(CodecError::Corrupt);
        }
        Ok(input[3..].to_vec())
    }

    #[test]
    fn encoder_lifecycle() {
        let mut encoder = BufferedEncoder::new(test_encode);
        let mut output = vec![0u8; 1024];

        assert_eq!(encoder.state(), State::Ready);

        // Process some input
        let n = encoder.process(b"hello", &mut output).unwrap();
        assert_eq!(n, 0);
        assert_eq!(encoder.state(), State::Streaming);

        // Finish encoding
        let n = encoder.finish(&mut output).unwrap();
        assert_eq!(n, 8); // "ENC" + "hello"
        assert_eq!(&output[..n], b"ENChello");
        assert_eq!(encoder.state(), State::Finished);

        // Cannot process after finish
        assert!(matches!(
            encoder.process(b"more", &mut output),
            Err(CodecError::InvalidState)
        ));

        // Reset and reuse
        encoder.reset();
        assert_eq!(encoder.state(), State::Ready);
    }

    #[test]
    fn decoder_lifecycle() {
        let mut decoder = BufferedDecoder::new(test_decode);
        let mut output = vec![0u8; 1024];

        assert_eq!(decoder.state(), State::Ready);

        // Process encoded input
        let n = decoder.process(b"ENChello", &mut output).unwrap();
        assert_eq!(n, 0);
        assert_eq!(decoder.state(), State::Streaming);

        // Finish decoding
        let n = decoder.finish(&mut output).unwrap();
        assert_eq!(n, 5);
        assert_eq!(&output[..n], b"hello");
        assert_eq!(decoder.state(), State::Finished);
    }

    #[test]
    fn buffer_too_small_is_transactional() {
        let mut encoder = BufferedEncoder::new(test_encode);
        let mut small_output = [0u8; 2]; // Too small for "ENC" + "hello"

        encoder.process(b"hello", &mut []).unwrap();
        assert_eq!(encoder.state(), State::Streaming);

        // Should fail with BufTooSmall without changing state
        let result = encoder.finish(&mut small_output);
        assert!(matches!(result, Err(CodecError::BufTooSmall)));
        assert_eq!(encoder.state(), State::Streaming); // State unchanged

        // Should succeed with larger buffer
        let mut large_output = [0u8; 1024];
        let n = encoder.finish(&mut large_output).unwrap();
        assert_eq!(n, 8);
    }
}
