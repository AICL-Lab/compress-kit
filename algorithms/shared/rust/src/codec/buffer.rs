use crate::codec::encoder::{Decoder, Encoder};
use crate::codec::error::{CodecError, MAX_INPUT_SIZE, MAX_OUTPUT_SIZE};

fn grow_buffer(current_len: usize, limit: usize) -> usize {
    if current_len == 0 {
        if limit < 1024 {
            return limit;
        }
        return 1024;
    }
    let next = current_len.saturating_mul(2);
    if next < current_len {
        return limit;
    }
    if next > limit {
        return limit;
    }
    next
}

/// EncodeBuffer is a convenience function that encodes input using the streaming API.
/// Equivalent to: new encoder → Process(input) → Finish() → collect output.
///
/// Returns the complete encoded output or an error.
pub fn encode_buffer(encoder: &mut dyn Encoder, input: &[u8]) -> Result<Vec<u8>, CodecError> {
    if input.len() > MAX_INPUT_SIZE {
        return Err(CodecError::SizeLimit);
    }

    // Allocate output buffer using a conservative estimate.
    // Use 2x input size + 2KB overhead as a reasonable initial allocation.
    let initial_size = input.len().saturating_mul(2).saturating_add(2048);
    let mut out_buf = vec![0u8; initial_size];
    let mut total_written = 0;

    // Process input
    loop {
        match encoder.process(input, &mut out_buf[total_written..]) {
            Ok(n) => {
                total_written += n;
                break;
            }
            Err(CodecError::BufTooSmall) => {
                if out_buf.len() >= MAX_OUTPUT_SIZE {
                    return Err(CodecError::SizeLimit);
                }
                let new_size = grow_buffer(out_buf.len(), MAX_OUTPUT_SIZE);
                if new_size <= out_buf.len() {
                    return Err(CodecError::SizeLimit);
                }
                out_buf.resize(new_size, 0);
            }
            Err(e) => return Err(e),
        }
    }

    // Finish encoding
    loop {
        match encoder.finish(&mut out_buf[total_written..]) {
            Ok(n) => {
                total_written += n;
                break;
            }
            Err(CodecError::BufTooSmall) => {
                if out_buf.len() >= MAX_OUTPUT_SIZE {
                    return Err(CodecError::SizeLimit);
                }
                let new_size = grow_buffer(out_buf.len(), MAX_OUTPUT_SIZE);
                if new_size <= out_buf.len() {
                    return Err(CodecError::SizeLimit);
                }
                out_buf.resize(new_size, 0);
            }
            Err(e) => return Err(e),
        }
    }

    out_buf.truncate(total_written);
    Ok(out_buf)
}

/// DecodeBuffer is a convenience function that decodes input using the streaming API.
/// Equivalent to: new decoder → Process(input) → Finish() → collect output.
///
/// Returns the complete decoded output or an error.
pub fn decode_buffer(decoder: &mut dyn Decoder, input: &[u8]) -> Result<Vec<u8>, CodecError> {
    // Allocate output buffer.
    // Decode typically expands, so start with input size and grow as needed.
    let initial_size = input.len().saturating_add(1024);
    let mut out_buf = vec![0u8; initial_size];
    let mut total_written = 0;

    // Process input
    loop {
        match decoder.process(input, &mut out_buf[total_written..]) {
            Ok(n) => {
                total_written += n;
                break;
            }
            Err(CodecError::BufTooSmall) => {
                if total_written > MAX_OUTPUT_SIZE {
                    return Err(CodecError::SizeLimit);
                }
                if out_buf.len() >= MAX_OUTPUT_SIZE {
                    return Err(CodecError::SizeLimit);
                }
                let new_size = grow_buffer(out_buf.len(), MAX_OUTPUT_SIZE);
                if new_size <= out_buf.len() {
                    return Err(CodecError::SizeLimit);
                }
                out_buf.resize(new_size, 0);
            }
            Err(e) => return Err(e),
        }
    }

    // Finish decoding
    loop {
        match decoder.finish(&mut out_buf[total_written..]) {
            Ok(n) => {
                total_written += n;
                break;
            }
            Err(CodecError::BufTooSmall) => {
                if out_buf.len() >= MAX_OUTPUT_SIZE {
                    return Err(CodecError::SizeLimit);
                }
                let new_size = grow_buffer(out_buf.len(), MAX_OUTPUT_SIZE);
                if new_size <= out_buf.len() {
                    return Err(CodecError::SizeLimit);
                }
                out_buf.resize(new_size, 0);
            }
            Err(e) => return Err(e),
        }
    }

    if total_written > MAX_OUTPUT_SIZE {
        return Err(CodecError::SizeLimit);
    }

    out_buf.truncate(total_written);
    Ok(out_buf)
}
