use compresskit_codec::codec::{decode_buffer, encode_buffer, CodecError, Decoder, Encoder, State};

const ENCODE_SUFFIX: &[u8] = b"tail";
const DECODE_SUFFIX: &[u8] = b"done";

struct RetryFinishEncoder {
    finish_calls: usize,
    finish_output_lens: Vec<usize>,
}

impl Encoder for RetryFinishEncoder {
    fn process(&mut self, _: &[u8], _: &mut [u8]) -> Result<usize, CodecError> {
        Ok(0)
    }

    fn flush(&mut self, _: &mut [u8]) -> Result<usize, CodecError> {
        Ok(0)
    }

    fn finish(&mut self, output: &mut [u8]) -> Result<usize, CodecError> {
        self.finish_calls += 1;
        match self.finish_calls {
            1 => {
                self.finish_output_lens.push(output.len());
                output.fill(b'a');
                Err(CodecError::BufTooSmall)
            }
            2 => {
                self.finish_output_lens.push(output.len());
                output[..ENCODE_SUFFIX.len()].copy_from_slice(ENCODE_SUFFIX);
                Ok(ENCODE_SUFFIX.len())
            }
            _ => Err(CodecError::Other("unexpected extra finish retry".into())),
        }
    }

    fn reset(&mut self) {}

    fn state(&self) -> State {
        State::Streaming
    }
}

struct RetryFinishDecoder {
    finish_calls: usize,
    finish_output_lens: Vec<usize>,
}

impl Decoder for RetryFinishDecoder {
    fn process(&mut self, _: &[u8], _: &mut [u8]) -> Result<usize, CodecError> {
        Ok(0)
    }

    fn flush(&mut self, _: &mut [u8]) -> Result<usize, CodecError> {
        Ok(0)
    }

    fn finish(&mut self, output: &mut [u8]) -> Result<usize, CodecError> {
        self.finish_calls += 1;
        match self.finish_calls {
            1 => {
                self.finish_output_lens.push(output.len());
                output.fill(b'z');
                Err(CodecError::BufTooSmall)
            }
            2 => {
                self.finish_output_lens.push(output.len());
                output[..DECODE_SUFFIX.len()].copy_from_slice(DECODE_SUFFIX);
                Ok(DECODE_SUFFIX.len())
            }
            _ => Err(CodecError::Other("unexpected extra finish retry".into())),
        }
    }

    fn reset(&mut self) {}

    fn state(&self) -> State {
        State::Streaming
    }
}

#[test]
fn encode_buffer_retries_finish_transactionally() {
    let mut encoder = RetryFinishEncoder {
        finish_calls: 0,
        finish_output_lens: Vec::new(),
    };

    let output = encode_buffer(&mut encoder, b"x").unwrap();

    assert_eq!(encoder.finish_calls, 2);
    assert_eq!(encoder.finish_output_lens.len(), 2);
    assert!(encoder.finish_output_lens[1] > encoder.finish_output_lens[0]);
    assert_eq!(output, ENCODE_SUFFIX);
}

#[test]
fn decode_buffer_retries_finish_transactionally() {
    let mut decoder = RetryFinishDecoder {
        finish_calls: 0,
        finish_output_lens: Vec::new(),
    };

    let output = decode_buffer(&mut decoder, b"x").unwrap();

    assert_eq!(decoder.finish_calls, 2);
    assert_eq!(decoder.finish_output_lens.len(), 2);
    assert!(decoder.finish_output_lens[1] > decoder.finish_output_lens[0]);
    assert_eq!(output, DECODE_SUFFIX);
}
