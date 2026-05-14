use crate::codec::buffer_policy::ResizingBuffer;
use crate::codec::encoder::{Decoder, Encoder};
use crate::codec::error::{CodecError, MAX_INPUT_SIZE, MAX_OUTPUT_SIZE};

fn encode_buffer_limit(input_len: usize) -> Result<usize, CodecError> {
    input_len
        .checked_mul(8)
        .and_then(|len| len.checked_add(2048))
        .ok_or(CodecError::SizeLimit)
}

/// EncodeBuffer is a convenience function that encodes input using the streaming API.
/// Equivalent to: new encoder → Process(input) → Finish() → collect output.
///
/// Returns the complete encoded output or an error.
pub fn encode_buffer(encoder: &mut dyn Encoder, input: &[u8]) -> Result<Vec<u8>, CodecError> {
    if input.len() > MAX_INPUT_SIZE {
        return Err(CodecError::SizeLimit);
    }

    let encode_limit = encode_buffer_limit(input.len())?;
    let mut runner = ResizingBuffer::new(
        input.len().saturating_mul(2).saturating_add(2048),
        encode_limit,
    );

    runner.run(&mut |output| encoder.process(input, output))?;
    runner.run(&mut |output| encoder.finish(output))?;
    Ok(runner.into_vec())
}

pub(crate) fn decode_buffer_with_limit(
    decoder: &mut dyn Decoder,
    input: &[u8],
    limit: usize,
) -> Result<Vec<u8>, CodecError> {
    if input.len() > MAX_INPUT_SIZE {
        return Err(CodecError::SizeLimit);
    }

    let mut runner = ResizingBuffer::new(input.len().saturating_add(1024), limit);
    runner.run(&mut |output| decoder.process(input, output))?;
    runner.run(&mut |output| decoder.finish(output))?;
    Ok(runner.into_vec())
}

/// DecodeBuffer is a convenience function that decodes input using the streaming API.
/// Equivalent to: new decoder → Process(input) → Finish() → collect output.
///
/// Returns the complete decoded output or an error.
pub fn decode_buffer(decoder: &mut dyn Decoder, input: &[u8]) -> Result<Vec<u8>, CodecError> {
    decode_buffer_with_limit(decoder, input, MAX_OUTPUT_SIZE)
}

#[cfg(test)]
mod tests {
    use super::encode_buffer;
    use super::decode_buffer_with_limit;
    use crate::codec::encoder::{Decoder, Encoder};
    use crate::codec::error::{CodecError, State};

    struct LimitHitProcessDecoder {
        process_calls: usize,
        finish_calls: usize,
    }

    impl Decoder for LimitHitProcessDecoder {
        fn process(&mut self, _: &[u8], _: &mut [u8]) -> Result<usize, CodecError> {
            self.process_calls += 1;
            Err(CodecError::BufTooSmall)
        }

        fn flush(&mut self, _: &mut [u8]) -> Result<usize, CodecError> {
            Ok(0)
        }

        fn finish(&mut self, _: &mut [u8]) -> Result<usize, CodecError> {
            self.finish_calls += 1;
            Err(CodecError::Other("finish should not be called".into()))
        }

        fn reset(&mut self) {}

        fn state(&self) -> State {
            State::Streaming
        }
    }

    #[test]
    fn decode_buffer_stops_growing_at_decode_limit_boundary() {
        let mut decoder = LimitHitProcessDecoder {
            process_calls: 0,
            finish_calls: 0,
        };

        let err = decode_buffer_with_limit(&mut decoder, b"", 1024).unwrap_err();

        assert_eq!(err, CodecError::SizeLimit);
        assert_eq!(decoder.process_calls, 1);
        assert_eq!(decoder.finish_calls, 0);
    }

    #[test]
    fn encode_buffer_retries_finish_after_buffer_growth() {
        struct RetryEncoder {
            calls: usize,
        }

        impl Encoder for RetryEncoder {
            fn process(&mut self, _: &[u8], _: &mut [u8]) -> Result<usize, CodecError> {
                Ok(0)
            }

            fn flush(&mut self, _: &mut [u8]) -> Result<usize, CodecError> {
                Ok(0)
            }

            fn finish(&mut self, out: &mut [u8]) -> Result<usize, CodecError> {
                self.calls += 1;
                if self.calls == 1 {
                    return Err(CodecError::BufTooSmall);
                }
                out[..6].copy_from_slice(b"abcdef");
                Ok(6)
            }

            fn reset(&mut self) {}

            fn state(&self) -> State {
                State::Streaming
            }
        }

        let mut encoder = RetryEncoder { calls: 0 };
        let out = encode_buffer(&mut encoder, b"ignored").unwrap();
        assert_eq!(out, b"abcdef");
    }
}
