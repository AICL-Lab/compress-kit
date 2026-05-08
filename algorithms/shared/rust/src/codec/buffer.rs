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

fn encode_buffer_limit(input_len: usize) -> Result<usize, CodecError> {
    input_len
        .checked_mul(8)
        .and_then(|len| len.checked_add(2048))
        .ok_or(CodecError::SizeLimit)
}

type BufferStep<'a> = dyn FnMut(&mut [u8]) -> Result<usize, CodecError> + 'a;

fn run_buffer_step(
    out_buf: &mut Vec<u8>,
    total_written: usize,
    limit: usize,
    step: &mut BufferStep<'_>,
) -> Result<usize, CodecError> {
    let mut total_written = total_written;

    loop {
        match step(&mut out_buf[total_written..]) {
            Ok(n) => {
                total_written = total_written.checked_add(n).ok_or(CodecError::SizeLimit)?;
                if total_written > limit {
                    return Err(CodecError::SizeLimit);
                }
                return Ok(total_written);
            }
            Err(CodecError::BufTooSmall) => {
                if total_written > limit || out_buf.len() >= limit {
                    return Err(CodecError::SizeLimit);
                }

                let new_size = grow_buffer(out_buf.len(), limit);
                if new_size <= out_buf.len() {
                    return Err(CodecError::SizeLimit);
                }

                out_buf.resize(new_size, 0);
            }
            Err(err) => return Err(err),
        }
    }
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

    // Allocate output buffer using a conservative estimate.
    // Use 2x input size + 2KB overhead as a reasonable initial allocation.
    let initial_size = input
        .len()
        .saturating_mul(2)
        .saturating_add(2048)
        .min(encode_limit);
    let mut out_buf = vec![0u8; initial_size];
    let mut process = |output: &mut [u8]| encoder.process(input, output);
    let mut total_written = run_buffer_step(&mut out_buf, 0, encode_limit, &mut process)?;

    let mut finish = |output: &mut [u8]| encoder.finish(output);
    total_written = run_buffer_step(&mut out_buf, total_written, encode_limit, &mut finish)?;

    out_buf.truncate(total_written);
    Ok(out_buf)
}

/// DecodeBuffer is a convenience function that decodes input using the streaming API.
/// Equivalent to: new decoder → Process(input) → Finish() → collect output.
///
/// Returns the complete decoded output or an error.
pub fn decode_buffer(decoder: &mut dyn Decoder, input: &[u8]) -> Result<Vec<u8>, CodecError> {
    if input.len() > MAX_INPUT_SIZE {
        return Err(CodecError::SizeLimit);
    }

    // Allocate output buffer.
    // Decode typically expands, so start with input size and grow as needed.
    let initial_size = input.len().saturating_add(1024);
    let mut out_buf = vec![0u8; initial_size];
    let mut process = |output: &mut [u8]| decoder.process(input, output);
    let mut total_written = run_buffer_step(&mut out_buf, 0, MAX_OUTPUT_SIZE, &mut process)?;

    let mut finish = |output: &mut [u8]| decoder.finish(output);
    total_written = run_buffer_step(&mut out_buf, total_written, MAX_OUTPUT_SIZE, &mut finish)?;

    out_buf.truncate(total_written);
    Ok(out_buf)
}
