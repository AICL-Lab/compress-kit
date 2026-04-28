pub mod codec {
    pub mod buffer;
    pub mod encoder;
    pub mod error;
    pub mod write;

    pub use buffer::{decode_buffer, encode_buffer};
    pub use encoder::{Decoder, Encoder};
    pub use error::{CodecError, State, MAX_INPUT_SIZE, MAX_OUTPUT_SIZE};
    pub use write::WriteEncoder;
}
