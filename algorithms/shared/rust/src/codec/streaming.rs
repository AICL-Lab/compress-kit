use crate::codec::buffered::{BufferedDecoder, BufferedEncoder, DecodeFunc, EncodeFunc};

pub fn streaming_encoder(encode: EncodeFunc) -> BufferedEncoder {
    BufferedEncoder::new(encode)
}

pub fn streaming_decoder(decode: DecodeFunc) -> BufferedDecoder {
    BufferedDecoder::new(decode)
}

#[cfg(test)]
mod tests {
    use crate::codec::error::CodecError;
    use crate::codec::{streaming_decoder, streaming_encoder};

    #[test]
    fn streaming_encoder_wraps_encode_func() {
        fn encode(input: &[u8]) -> Result<Vec<u8>, CodecError> {
            Ok([b"rust:".as_slice(), input].concat())
        }

        let mut encoder = streaming_encoder(encode);
        let out = crate::codec::buffer::encode_buffer(&mut encoder, b"abc").unwrap();
        assert_eq!(out, b"rust:abc");
    }

    #[test]
    fn streaming_decoder_wraps_decode_func() {
        fn decode(input: &[u8]) -> Result<Vec<u8>, CodecError> {
            Ok([b"rust-dec:".as_slice(), input].concat())
        }

        let mut decoder = streaming_decoder(decode);
        let out = crate::codec::buffer::decode_buffer(&mut decoder, b"abc").unwrap();
        assert_eq!(out, b"rust-dec:abc");
    }
}
