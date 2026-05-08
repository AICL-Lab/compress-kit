pub mod codec {
    pub mod buffer;
    pub mod buffered;
    pub mod bits;
    pub mod encoder;
    pub mod error;
    pub mod write;

    pub use buffer::{decode_buffer, encode_buffer};
    pub use buffered::{BufferedDecoder, BufferedEncoder, DecodeFunc, EncodeFunc};
    pub use bits::{BitReader, BitWriter};
    pub use encoder::{Decoder, Encoder};
    pub use error::{
        io_error_to_codec_error, map_io_error, CodecError, State, MAX_INPUT_SIZE, MAX_OUTPUT_SIZE,
    };
    pub use write::WriteEncoder;
}
