// Minimal in-memory Arithmetic Coding encode/decode API for testing.
// Main CLI remains in main.rs. This library provides simple helpers.

use std::io;

use compresskit_codec::codec::BitWriter;

const SYMBOL_LIMIT: usize = 257;
const EOF_SYMBOL: u32 = (SYMBOL_LIMIT - 1) as u32;
const MAX_TOTAL: u32 = 1 << 24;
const MAX_OUTPUT_SIZE: usize = 1024 * 1024 * 1024; // 1 GiB
const STATE_BITS: u64 = 32;
const FULL_RANGE: u64 = 1u64 << STATE_BITS;
const HALF_RANGE: u64 = FULL_RANGE >> 1;
const FIRST_QUARTER: u64 = HALF_RANGE >> 1;
const THIRD_QUARTER: u64 = FIRST_QUARTER * 3;

fn scale_frequencies(freq: &mut [u32]) {
    let total: u64 = freq.iter().map(|&f| f as u64).sum();
    if total == 0 {
        for f in freq.iter_mut() {
            *f = 1;
        }
        return;
    }
    if total <= MAX_TOTAL as u64 {
        return;
    }
    let mut new_total = 0u64;
    for f in freq.iter_mut() {
        if *f == 0 {
            continue;
        }
        let mut scaled = (*f as u64 * MAX_TOTAL as u64) / total;
        if scaled == 0 {
            scaled = 1;
        }
        *f = scaled as u32;
        new_total += scaled;
    }
    if new_total == 0 {
        let base = MAX_TOTAL / freq.len() as u32;
        for f in freq.iter_mut() {
            *f = if base == 0 { 1 } else { base };
        }
    }
}

fn build_cumulative(freq: &[u32]) -> Vec<u32> {
    let mut cumulative = vec![0u32; freq.len() + 1];
    for (i, &f) in freq.iter().enumerate() {
        cumulative[i + 1] = cumulative[i] + f;
    }
    cumulative
}

pub fn encode(input: &[u8]) -> Result<Vec<u8>, io::Error> {
    let mut freq = vec![0u32; SYMBOL_LIMIT];
    for &b in input {
        freq[b as usize] += 1;
    }
    freq[EOF_SYMBOL as usize] = 1;
    scale_frequencies(&mut freq);
    let cumulative = build_cumulative(&freq);

    let mut output = Vec::new();
    output.extend_from_slice(b"AENC");
    output.extend_from_slice(&(SYMBOL_LIMIT as u32).to_le_bytes());
    for &f in &freq {
        output.extend_from_slice(&f.to_le_bytes());
    }

    let mut low = 0u64;
    let mut high = FULL_RANGE - 1;
    let mut pending_bits = 0u64;
    let mut bit_writer = BitWriter::new();

    for &b in input {
        encode_symbol(
            b as u32,
            &cumulative,
            &mut low,
            &mut high,
            &mut pending_bits,
            &mut bit_writer,
        );
    }
    encode_symbol(
        EOF_SYMBOL,
        &cumulative,
        &mut low,
        &mut high,
        &mut pending_bits,
        &mut bit_writer,
    );

    pending_bits += 1;
    if low < FIRST_QUARTER {
        emit_bit(0, &mut bit_writer, &mut pending_bits);
    } else {
        emit_bit(1, &mut bit_writer, &mut pending_bits);
    }

    output.extend(bit_writer.finish());

    Ok(output)
}

fn emit_bit(bit: u8, bit_writer: &mut BitWriter, pending_bits: &mut u64) {
    bit_writer.write_bit_u8(bit);
    for _ in 0..*pending_bits {
        bit_writer.write_bit_u8(if bit == 0 { 1 } else { 0 });
    }
    *pending_bits = 0;
}

fn encode_symbol(
    sym: u32,
    cumulative: &[u32],
    low: &mut u64,
    high: &mut u64,
    pending_bits: &mut u64,
    bit_writer: &mut BitWriter,
) {
    let range = *high - *low + 1;
    let total = cumulative[cumulative.len() - 1] as u64;
    let sym_low = cumulative[sym as usize] as u64;
    let sym_high = cumulative[sym as usize + 1] as u64;

    *high = *low + (range * sym_high) / total - 1;
    *low += (range * sym_low) / total;

    loop {
        if *high < HALF_RANGE {
            emit_bit(0, bit_writer, pending_bits);
            *low <<= 1;
            *high = (*high << 1) | 1;
        } else if *low >= HALF_RANGE {
            emit_bit(1, bit_writer, pending_bits);
            *low = (*low - HALF_RANGE) << 1;
            *high = ((*high - HALF_RANGE) << 1) | 1;
        } else if *low >= FIRST_QUARTER && *high < THIRD_QUARTER {
            *pending_bits += 1;
            *low = (*low - FIRST_QUARTER) << 1;
            *high = ((*high - FIRST_QUARTER) << 1) | 1;
        } else {
            break;
        }
    }
}

pub fn decode(input: &[u8]) -> Result<Vec<u8>, io::Error> {
    if input.len() < 8 {
        return Err(io::Error::new(
            io::ErrorKind::InvalidData,
            "input too short",
        ));
    }
    if &input[0..4] != b"AENC" {
        return Err(io::Error::new(io::ErrorKind::InvalidData, "invalid magic"));
    }

    let count = u32::from_le_bytes([input[4], input[5], input[6], input[7]]) as usize;
    if count != SYMBOL_LIMIT {
        return Err(io::Error::new(
            io::ErrorKind::InvalidData,
            "invalid symbol count",
        ));
    }

    let mut pos = 8;
    let mut freq = vec![0u32; count];
    for f in freq.iter_mut() {
        if pos + 4 > input.len() {
            return Err(io::Error::new(
                io::ErrorKind::InvalidData,
                "truncated freq table",
            ));
        }
        *f = u32::from_le_bytes([input[pos], input[pos + 1], input[pos + 2], input[pos + 3]]);
        pos += 4;
    }

    let cumulative = build_cumulative(&freq);
    let total = cumulative[cumulative.len() - 1] as u64;

    let mut low = 0u64;
    let mut high = FULL_RANGE - 1;
    let mut code = 0u64;

    let mut bit_buffer = Vec::new();
    for &byte in &input[pos..] {
        for i in (0..8).rev() {
            bit_buffer.push((byte >> i) & 1);
        }
    }

    for i in 0..STATE_BITS {
        code <<= 1;
        if (i as usize) < bit_buffer.len() {
            code |= bit_buffer[i as usize] as u64;
        }
    }
    let mut bit_pos = STATE_BITS as usize;

    let mut output = Vec::new();
    loop {
        let range = high - low + 1;
        let offset = code - low;
        let value = ((offset + 1) * total - 1) / range;

        let mut sym = 0usize;
        for i in 0..cumulative.len() - 1 {
            if value >= cumulative[i] as u64 && value < cumulative[i + 1] as u64 {
                sym = i;
                break;
            }
        }

        if sym == EOF_SYMBOL as usize {
            break;
        }
        if output.len() >= MAX_OUTPUT_SIZE {
            return Err(io::Error::new(
                io::ErrorKind::InvalidData,
                "output size exceeds maximum allowed (1 GiB)",
            ));
        }
        output.push(sym as u8);

        let sym_low = cumulative[sym] as u64;
        let sym_high = cumulative[sym + 1] as u64;
        high = low + (range * sym_high) / total - 1;
        low += (range * sym_low) / total;

        loop {
            if high < HALF_RANGE {
                low <<= 1;
                high = (high << 1) | 1;
                code <<= 1;
                if bit_pos < bit_buffer.len() {
                    code |= bit_buffer[bit_pos] as u64;
                    bit_pos += 1;
                }
            } else if low >= HALF_RANGE {
                low = (low - HALF_RANGE) << 1;
                high = ((high - HALF_RANGE) << 1) | 1;
                code = (code - HALF_RANGE) << 1;
                if bit_pos < bit_buffer.len() {
                    code |= bit_buffer[bit_pos] as u64;
                    bit_pos += 1;
                }
            } else if low >= FIRST_QUARTER && high < THIRD_QUARTER {
                low = (low - FIRST_QUARTER) << 1;
                high = ((high - FIRST_QUARTER) << 1) | 1;
                code = (code - FIRST_QUARTER) << 1;
                if bit_pos < bit_buffer.len() {
                    code |= bit_buffer[bit_pos] as u64;
                    bit_pos += 1;
                }
            } else {
                break;
            }
        }
    }

    Ok(output)
}

// Streaming adapters
use compresskit_codec::codec::{
    io_error_to_codec_error, BufferedDecoder, BufferedEncoder, CodecError, Decoder, Encoder,
};

/// Creates a new streaming Arithmetic encoder.
pub fn new_encoder() -> impl Encoder {
    BufferedEncoder::new(arithmetic_encode)
}

/// Creates a new streaming Arithmetic decoder.
pub fn new_decoder() -> impl Decoder {
    BufferedDecoder::new(arithmetic_decode)
}

fn arithmetic_encode(input: &[u8]) -> Result<Vec<u8>, CodecError> {
    encode(input).map_err(io_error_to_codec_error)
}

fn arithmetic_decode(input: &[u8]) -> Result<Vec<u8>, CodecError> {
    decode(input).map_err(io_error_to_codec_error)
}

#[cfg(test)]
mod tests {
    use super::{decode, encode};

    #[test]
    fn roundtrip_zero_byte_without_confusing_it_with_eof() {
        let encoded = encode(&[0x00]).unwrap();
        let decoded = decode(&encoded).unwrap();

        assert_eq!(decoded, vec![0x00]);
    }
}
