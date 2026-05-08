// Lifecycle tests for the streaming API foundation.
// These tests verify the contract defined in the design document.

use compresskit_codec::codec::{CodecError, Decoder, Encoder, State};

#[test]
fn test_lifecycle_basic() {
    // Basic mock encoder to test lifecycle state machine
    struct MockEncoder(State);
    impl Encoder for MockEncoder {
        fn process(&mut self, _: &[u8], _: &mut [u8]) -> Result<usize, CodecError> {
            match self.0 {
                State::Ready | State::Flushing => {
                    self.0 = State::Streaming;
                    Ok(0)
                }
                State::Streaming => Ok(0),
                State::Finished => {
                    self.0 = State::Error;
                    Err(CodecError::InvalidState)
                }
                State::Error => Err(CodecError::InvalidState),
            }
        }
        fn flush(&mut self, _: &mut [u8]) -> Result<usize, CodecError> {
            match self.0 {
                State::Ready => Ok(0),
                State::Streaming => {
                    self.0 = State::Flushing;
                    Ok(0)
                }
                State::Flushing => Ok(0),
                State::Finished => {
                    self.0 = State::Error;
                    Err(CodecError::InvalidState)
                }
                State::Error => Err(CodecError::InvalidState),
            }
        }
        fn finish(&mut self, _: &mut [u8]) -> Result<usize, CodecError> {
            match self.0 {
                State::Ready | State::Streaming | State::Flushing => {
                    self.0 = State::Finished;
                    Ok(0)
                }
                State::Finished => {
                    self.0 = State::Error;
                    Err(CodecError::InvalidState)
                }
                State::Error => Err(CodecError::InvalidState),
            }
        }
        fn reset(&mut self) {
            self.0 = State::Ready;
        }
        fn state(&self) -> State {
            self.0
        }
    }

    let mut enc = MockEncoder(State::Ready);
    assert_eq!(enc.state(), State::Ready);

    // Process should move to Streaming
    enc.process(b"test", &mut []).unwrap();
    assert_eq!(enc.state(), State::Streaming);

    // Flush should move to Flushing
    enc.flush(&mut []).unwrap();
    assert_eq!(enc.state(), State::Flushing);

    // Finish should move to Finished
    enc.finish(&mut []).unwrap();
    assert_eq!(enc.state(), State::Finished);

    // Post-finish calls should error and move to Error
    assert_eq!(
        enc.process(b"fail", &mut []).unwrap_err(),
        CodecError::InvalidState
    );
    assert_eq!(enc.state(), State::Error);

    // Reset should return to Ready
    enc.reset();
    assert_eq!(enc.state(), State::Ready);
}

#[test]
fn test_constants() {
    use compresskit_codec::codec::{MAX_INPUT_SIZE, MAX_OUTPUT_SIZE};
    assert_eq!(MAX_INPUT_SIZE, 4 * 1024 * 1024 * 1024);
    assert_eq!(MAX_OUTPUT_SIZE, 1024 * 1024 * 1024);
}

// Integration tests for actual algorithm adapters

#[test]
fn test_huffman_implements_traits() {
    let _enc: Box<dyn Encoder> = Box::new(huffman::new_encoder());
    let _dec: Box<dyn Decoder> = Box::new(huffman::new_decoder());
}

#[test]
fn test_arithmetic_implements_traits() {
    let _enc: Box<dyn Encoder> = Box::new(arithmetic::new_encoder());
    let _dec: Box<dyn Decoder> = Box::new(arithmetic::new_decoder());
}

#[test]
fn test_rle_implements_traits() {
    let _enc: Box<dyn Encoder> = Box::new(rle::new_encoder());
    let _dec: Box<dyn Decoder> = Box::new(rle::new_decoder());
}

#[test]
fn test_range_implements_traits() {
    let _enc: Box<dyn Encoder> = Box::new(rangecoder::new_encoder());
    let _dec: Box<dyn Decoder> = Box::new(rangecoder::new_decoder());
}

#[test]
fn test_huffman_roundtrip_via_traits() {
    let input = b"hello world";
    let mut enc = huffman::new_encoder();
    let mut output = vec![0u8; 4096];
    enc.process(input, &mut output).unwrap();
    let written = enc.finish(&mut output).unwrap();

    let mut dec = huffman::new_decoder();
    let mut decoded = vec![0u8; 4096];
    dec.process(&output[..written], &mut decoded).unwrap();
    let dec_written = dec.finish(&mut decoded).unwrap();

    assert_eq!(&decoded[..dec_written], input);
}

#[test]
fn test_arithmetic_roundtrip_via_traits() {
    let input = b"test data";
    let mut enc = arithmetic::new_encoder();
    let mut output = vec![0u8; 4096];
    enc.process(input, &mut output).unwrap();
    let written = enc.finish(&mut output).unwrap();

    let mut dec = arithmetic::new_decoder();
    let mut decoded = vec![0u8; 4096];
    dec.process(&output[..written], &mut decoded).unwrap();
    let dec_written = dec.finish(&mut decoded).unwrap();

    assert_eq!(&decoded[..dec_written], input);
}

#[test]
fn test_rle_roundtrip_via_traits() {
    let input = b"aaabbbccc";
    let mut enc = rle::new_encoder();
    let mut output = vec![0u8; 4096];
    enc.process(input, &mut output).unwrap();
    let written = enc.finish(&mut output).unwrap();

    let mut dec = rle::new_decoder();
    let mut decoded = vec![0u8; 4096];
    dec.process(&output[..written], &mut decoded).unwrap();
    let dec_written = dec.finish(&mut decoded).unwrap();

    assert_eq!(&decoded[..dec_written], input);
}
