// Minimal in-memory RLE encode/decode API for testing.
// Main CLI remains in main.rs. This library provides simple helpers.

use std::io;

const MAX_OUTPUT_SIZE: usize = 1024 * 1024 * 1024;
const MAX_INPUT_SIZE: usize = 4 * 1024 * 1024 * 1024;

/// Magic number for RLE format identification
const RLE_MAGIC: &[u8; 4] = b"RLE\x00";

pub fn encode(input: &[u8]) -> Result<Vec<u8>, io::Error> {
    if input.len() > MAX_INPUT_SIZE {
        return Err(io::Error::new(
            io::ErrorKind::InvalidData,
            "input exceeds maximum size (4 GiB)",
        ));
    }
    let mut output = Vec::with_capacity(input.len() / 8 + 4);

    // Write magic number
    output.extend_from_slice(RLE_MAGIC);

    if input.is_empty() {
        return Ok(output);
    }

    let mut current = input[0];
    let mut count: u32 = 1;

    for &b in &input[1..] {
        if b == current && count < u32::MAX {
            count += 1;
        } else {
            output.extend_from_slice(&count.to_le_bytes());
            output.push(current);
            current = b;
            count = 1;
        }
    }

    output.extend_from_slice(&count.to_le_bytes());
    output.push(current);
    Ok(output)
}

pub fn decode(input: &[u8]) -> Result<Vec<u8>, io::Error> {
    if input.len() > MAX_INPUT_SIZE {
        return Err(io::Error::new(
            io::ErrorKind::InvalidData,
            "input exceeds maximum size (4 GiB)",
        ));
    }
    // Verify magic number
    if input.len() < 4 {
        return Err(io::Error::new(
            io::ErrorKind::InvalidData,
            "input too short: missing magic number",
        ));
    }
    if &input[0..4] != RLE_MAGIC {
        return Err(io::Error::new(
            io::ErrorKind::InvalidData,
            "invalid RLE file: bad magic number",
        ));
    }

    let mut output = Vec::new();
    let mut pos = 4; // Start after magic number

    while pos < input.len() {
        if pos + 5 > input.len() {
            return Err(io::Error::new(
                io::ErrorKind::InvalidData,
                "RLE data truncated: incomplete count+value pair",
            ));
        }

        let count =
            u32::from_le_bytes([input[pos], input[pos + 1], input[pos + 2], input[pos + 3]]);
        let value = input[pos + 4];
        pos += 5;

        // Validate count is not zero
        if count == 0 {
            return Err(io::Error::new(
                io::ErrorKind::InvalidData,
                "invalid RLE data: count should not be 0",
            ));
        }

        let new_len = output.len() + count as usize;
        if new_len > MAX_OUTPUT_SIZE {
            return Err(io::Error::new(
                io::ErrorKind::InvalidData,
                "output size exceeds maximum allowed (1 GiB)",
            ));
        }

        output.resize(new_len, value);
    }

    Ok(output)
}

// Streaming adapters
use compresskit_codec::codec::{
    streaming_decoder, streaming_encoder, CodecError, Decoder, Encoder,
};

/// Creates a new streaming RLE encoder.
pub fn new_encoder() -> impl Encoder {
    streaming_encoder(rle_encode)
}

/// Creates a new streaming RLE decoder.
pub fn new_decoder() -> impl Decoder {
    streaming_decoder(rle_decode)
}

fn rle_encode(input: &[u8]) -> Result<Vec<u8>, CodecError> {
    encode(input).map_err(rle_io_error_to_codec_error)
}

fn rle_decode(input: &[u8]) -> Result<Vec<u8>, CodecError> {
    decode(input).map_err(rle_io_error_to_codec_error)
}

/// RLE-specific error conversion that handles SizeLimit cases.
fn rle_io_error_to_codec_error(e: io::Error) -> CodecError {
    match e.kind() {
        io::ErrorKind::UnexpectedEof => CodecError::Truncated,
        io::ErrorKind::InvalidData => {
            let msg = e.to_string();
            if msg.contains("truncated") || msg.contains("too short") {
                CodecError::Truncated
            } else if msg.contains("limit") || msg.contains("exceeds") {
                CodecError::SizeLimit
            } else if msg.contains("invalid") || msg.contains("bad") {
                CodecError::Corrupt
            } else {
                CodecError::Other(msg)
            }
        }
        _ => CodecError::Other(e.to_string()),
    }
}
