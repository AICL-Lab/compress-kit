use compresskit_codec::codec::{decode_buffer, encode_buffer, CodecError, Decoder, Encoder, State};

const ENCODE_SUFFIX: &[u8] = b"tail";
const ENCODE_BODY: &[u8] = b"body";
const DECODE_SUFFIX: &[u8] = b"done";
const DECODE_BODY: &[u8] = b"data";

struct RetryProcessEncoder {
    process_calls: usize,
    process_output_lens: Vec<usize>,
    finish_calls: usize,
}

impl Encoder for RetryProcessEncoder {
    fn process(&mut self, _: &[u8], output: &mut [u8]) -> Result<usize, CodecError> {
        self.process_calls += 1;
        match self.process_calls {
            1 => {
                self.process_output_lens.push(output.len());
                output.fill(b'p');
                Err(CodecError::BufTooSmall)
            }
            2 => {
                self.process_output_lens.push(output.len());
                output[..ENCODE_BODY.len()].copy_from_slice(ENCODE_BODY);
                Ok(ENCODE_BODY.len())
            }
            _ => Err(CodecError::Other("unexpected extra process retry".into())),
        }
    }

    fn flush(&mut self, _: &mut [u8]) -> Result<usize, CodecError> {
        Ok(0)
    }

    fn finish(&mut self, _: &mut [u8]) -> Result<usize, CodecError> {
        self.finish_calls += 1;
        Ok(0)
    }

    fn reset(&mut self) {}

    fn state(&self) -> State {
        State::Streaming
    }
}

struct RetryProcessDecoder {
    process_calls: usize,
    process_output_lens: Vec<usize>,
    finish_calls: usize,
}

impl Decoder for RetryProcessDecoder {
    fn process(&mut self, _: &[u8], output: &mut [u8]) -> Result<usize, CodecError> {
        self.process_calls += 1;
        match self.process_calls {
            1 => {
                self.process_output_lens.push(output.len());
                output.fill(b'd');
                Err(CodecError::BufTooSmall)
            }
            2 => {
                self.process_output_lens.push(output.len());
                output[..DECODE_BODY.len()].copy_from_slice(DECODE_BODY);
                Ok(DECODE_BODY.len())
            }
            _ => Err(CodecError::Other("unexpected extra process retry".into())),
        }
    }

    fn flush(&mut self, _: &mut [u8]) -> Result<usize, CodecError> {
        Ok(0)
    }

    fn finish(&mut self, _: &mut [u8]) -> Result<usize, CodecError> {
        self.finish_calls += 1;
        Ok(0)
    }

    fn reset(&mut self) {}

    fn state(&self) -> State {
        State::Streaming
    }
}

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

struct LimitHitProcessEncoder {
    process_calls: usize,
    finish_calls: usize,
}

impl Encoder for LimitHitProcessEncoder {
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
fn encode_buffer_retries_process_transactionally() {
    let mut encoder = RetryProcessEncoder {
        process_calls: 0,
        process_output_lens: Vec::new(),
        finish_calls: 0,
    };

    let output = encode_buffer(&mut encoder, b"x").expect("encode_buffer should succeed");

    assert_eq!(encoder.process_calls, 2);
    assert_eq!(encoder.process_output_lens.len(), 2);
    assert!(encoder.process_output_lens[1] > encoder.process_output_lens[0]);
    assert_eq!(encoder.finish_calls, 1);
    assert_eq!(output, ENCODE_BODY);
}

#[test]
fn decode_buffer_retries_process_transactionally() {
    let mut decoder = RetryProcessDecoder {
        process_calls: 0,
        process_output_lens: Vec::new(),
        finish_calls: 0,
    };

    let output = decode_buffer(&mut decoder, b"x").expect("decode_buffer should succeed");

    assert_eq!(decoder.process_calls, 2);
    assert_eq!(decoder.process_output_lens.len(), 2);
    assert!(decoder.process_output_lens[1] > decoder.process_output_lens[0]);
    assert_eq!(decoder.finish_calls, 1);
    assert_eq!(output, DECODE_BODY);
}

#[test]
fn encode_buffer_retries_finish_transactionally() {
    let mut encoder = RetryFinishEncoder {
        finish_calls: 0,
        finish_output_lens: Vec::new(),
    };

    let output = encode_buffer(&mut encoder, b"x").expect("encode_buffer should succeed");

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

    let output = decode_buffer(&mut decoder, b"x").expect("decode_buffer should succeed");

    assert_eq!(decoder.finish_calls, 2);
    assert_eq!(decoder.finish_output_lens.len(), 2);
    assert!(decoder.finish_output_lens[1] > decoder.finish_output_lens[0]);
    assert_eq!(output, DECODE_SUFFIX);
}

#[test]
fn encode_buffer_stops_growing_at_encode_limit_boundary() {
    let mut encoder = LimitHitProcessEncoder {
        process_calls: 0,
        finish_calls: 0,
    };

    let err = encode_buffer(&mut encoder, b"").unwrap_err();

    assert_eq!(err, CodecError::SizeLimit);
    assert_eq!(encoder.process_calls, 1);
    assert_eq!(encoder.finish_calls, 0);
}
