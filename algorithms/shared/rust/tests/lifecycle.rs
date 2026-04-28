// Lifecycle tests for the streaming API foundation.
// These tests verify the contract defined in the design document.

use compresskit_codec::codec::{CodecError, Encoder, State};

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
    assert_eq!(enc.process(b"fail", &mut []).unwrap_err(), CodecError::InvalidState);
    assert_eq!(enc.state(), State::Error);

    // Reset should return to Ready
    enc.reset();
    assert_eq!(enc.state(), State::Ready);
}

#[test]
fn test_constants() {
    use compresskit_codec::codec::{MAX_INPUT_SIZE, MAX_OUTPUT_SIZE};
    assert_eq!(MAX_INPUT_SIZE, 4 * 1024 * 1024 * 1024);
    assert_eq!(MAX_OUTPUT_SIZE, 1 * 1024 * 1024 * 1024);
}

// Note: Full integration tests with actual Huffman/Arithmetic/Range/RLE encoders
// are in each algorithm's own test module.
