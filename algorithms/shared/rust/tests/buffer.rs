use compresskit_codec::codec::{decode_buffer, encode_buffer, CodecError, Decoder, Encoder, State};

const ENCODE_SUFFIX: &[u8] = b"tail";
const DECODE_SUFFIX: &[u8] = b"done";

struct RetryFinishEncoder {
    finish_calls: usize,
    prefix_len: usize,
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
                self.prefix_len = output.len();
                output.fill(b'a');
                Err(CodecError::BufTooSmall)
            }
            2 => {
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
    prefix_len: usize,
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
                self.prefix_len = output.len();
                output.fill(b'z');
                Err(CodecError::BufTooSmall)
            }
            2 => {
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
fn encode_buffer_preserves_finish_retry_prefix() {
    let mut encoder = RetryFinishEncoder {
        finish_calls: 0,
        prefix_len: 0,
    };

    let output = encode_buffer(&mut encoder, b"x").unwrap();

    assert_eq!(encoder.finish_calls, 2);
    assert_eq!(output.len(), encoder.prefix_len + ENCODE_SUFFIX.len());
    assert!(output[..encoder.prefix_len].iter().all(|&byte| byte == b'a'));
    assert_eq!(&output[encoder.prefix_len..], ENCODE_SUFFIX);
}

#[test]
fn decode_buffer_preserves_finish_retry_prefix() {
    let mut decoder = RetryFinishDecoder {
        finish_calls: 0,
        prefix_len: 0,
    };

    let output = decode_buffer(&mut decoder, b"x").unwrap();

    assert_eq!(decoder.finish_calls, 2);
    assert_eq!(output.len(), decoder.prefix_len + DECODE_SUFFIX.len());
    assert!(output[..decoder.prefix_len].iter().all(|&byte| byte == b'z'));
    assert_eq!(&output[decoder.prefix_len..], DECODE_SUFFIX);
}
