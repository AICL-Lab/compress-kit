use compresskit_codec::codec::{CodecError, Decoder, Encoder, State};
use std::io::{self, Read, Write};

const SYMBOL_LIMIT: usize = 257;
const EOF_SYMBOL: u32 = (SYMBOL_LIMIT - 1) as u32;
const MAX_TOTAL: u32 = 1 << 24;

const STATE_BITS: u64 = 32;
const FULL_RANGE: u64 = 1u64 << STATE_BITS;
const HALF_RANGE: u64 = FULL_RANGE >> 1;
const FIRST_QUARTER: u64 = HALF_RANGE >> 1;
const THIRD_QUARTER: u64 = FIRST_QUARTER * 3;

// BitWriter for in-memory encoding
struct BitWriter {
    buffer: Vec<u8>,
    current_byte: u8,
    bits_in_buffer: u8,
}

impl BitWriter {
    fn new() -> Self {
        BitWriter {
            buffer: Vec::new(),
            current_byte: 0,
            bits_in_buffer: 0,
        }
    }

    fn write_bit(&mut self, bit: u8) {
        self.current_byte = (self.current_byte << 1) | (bit & 1);
        self.bits_in_buffer += 1;
        if self.bits_in_buffer == 8 {
            self.buffer.push(self.current_byte);
            self.bits_in_buffer = 0;
            self.current_byte = 0;
        }
    }

    fn flush(&mut self) {
        if self.bits_in_buffer > 0 {
            self.current_byte <<= 8 - self.bits_in_buffer;
            self.buffer.push(self.current_byte);
            self.bits_in_buffer = 0;
            self.current_byte = 0;
        }
    }

    fn into_bytes(mut self) -> Vec<u8> {
        self.flush();
        self.buffer
    }
}

// BitReader for in-memory decoding
struct BitReader<'a> {
    data: &'a [u8],
    pos: usize,
    current_byte: u8,
    bits_remaining: u8,
}

impl<'a> BitReader<'a> {
    fn new(data: &'a [u8]) -> Self {
        BitReader {
            data,
            pos: 0,
            current_byte: 0,
            bits_remaining: 0,
        }
    }

    fn read_bit(&mut self) -> u8 {
        if self.bits_remaining == 0 {
            if self.pos >= self.data.len() {
                return 0;
            }
            self.current_byte = self.data[self.pos];
            self.pos += 1;
            self.bits_remaining = 8;
        }
        self.bits_remaining -= 1;
        (self.current_byte >> self.bits_remaining) & 1
    }
}

struct ArithmeticEncoder {
    writer: BitWriter,
    low: u64,
    high: u64,
    pending_bits: u64,
}

impl ArithmeticEncoder {
    fn new() -> Self {
        ArithmeticEncoder {
            writer: BitWriter::new(),
            low: 0,
            high: FULL_RANGE - 1,
            pending_bits: 0,
        }
    }

    fn encode_symbol(&mut self, symbol: u32, cumulative: &[u32]) -> Result<(), CodecError> {
        let range = self.high - self.low + 1;
        let total = *cumulative
            .last()
            .ok_or(CodecError::Corrupt)? as u64;
        let sym_low = cumulative[symbol as usize] as u64;
        let sym_high = cumulative[symbol as usize + 1] as u64;

        self.high = self.low + (range * sym_high) / total - 1;
        self.low = self.low + (range * sym_low) / total;

        loop {
            if self.high < HALF_RANGE {
                self.output_bit(0);
            } else if self.low >= HALF_RANGE {
                self.output_bit(1);
                self.low -= HALF_RANGE;
                self.high -= HALF_RANGE;
            } else if self.low >= FIRST_QUARTER && self.high < THIRD_QUARTER {
                self.pending_bits += 1;
                self.low -= FIRST_QUARTER;
                self.high -= FIRST_QUARTER;
            } else {
                break;
            }
            self.low <<= 1;
            self.high = (self.high << 1) | 1;
        }
        Ok(())
    }

    fn finish(&mut self) {
        self.pending_bits += 1;
        if self.low < FIRST_QUARTER {
            self.output_bit(0);
        } else {
            self.output_bit(1);
        }
        self.writer.flush();
    }

    fn output_bit(&mut self, bit: u8) {
        self.writer.write_bit(bit);
        let complement = bit ^ 1;
        while self.pending_bits > 0 {
            self.writer.write_bit(complement);
            self.pending_bits -= 1;
        }
    }

    fn into_bytes(self) -> Vec<u8> {
        self.writer.into_bytes()
    }
}

struct ArithmeticDecoder<'a> {
    reader: BitReader<'a>,
    low: u64,
    high: u64,
    code: u64,
}

impl<'a> ArithmeticDecoder<'a> {
    fn new(reader: BitReader<'a>) -> Self {
        let mut decoder = ArithmeticDecoder {
            reader,
            low: 0,
            high: FULL_RANGE - 1,
            code: 0,
        };
        for _ in 0..STATE_BITS {
            decoder.code = (decoder.code << 1) | decoder.reader.read_bit() as u64;
        }
        decoder
    }

    fn decode_symbol(&mut self, cumulative: &[u32]) -> Result<u32, CodecError> {
        let range = self.high - self.low + 1;
        let total = *cumulative
            .last()
            .ok_or(CodecError::Corrupt)? as u64;
        let offset = self.code - self.low;
        let value = ((offset + 1) * total - 1) / range;

        let mut lo: u32 = 0;
        let mut hi: u32 = cumulative.len() as u32 - 1;
        while lo + 1 < hi {
            let mid = lo + (hi - lo) / 2;
            if cumulative[mid as usize] as u64 > value {
                hi = mid;
            } else {
                lo = mid;
            }
        }
        let symbol = lo;

        let sym_low = cumulative[symbol as usize] as u64;
        let sym_high = cumulative[symbol as usize + 1] as u64;

        self.high = self.low + (range * sym_high) / total - 1;
        self.low = self.low + (range * sym_low) / total;

        loop {
            if self.high < HALF_RANGE {
                // nothing
            } else if self.low >= HALF_RANGE {
                self.low -= HALF_RANGE;
                self.high -= HALF_RANGE;
                self.code -= HALF_RANGE;
            } else if self.low >= FIRST_QUARTER && self.high < THIRD_QUARTER {
                self.low -= FIRST_QUARTER;
                self.high -= FIRST_QUARTER;
                self.code -= FIRST_QUARTER;
            } else {
                break;
            }
            self.low <<= 1;
            self.high = (self.high << 1) | 1;
            self.code = (self.code << 1) | self.reader.read_bit() as u64;
        }

        Ok(symbol)
    }
}

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
    let mut new_total: u64 = 0;
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
        let mut base = MAX_TOTAL / freq.len() as u32;
        if base == 0 {
            base = 1;
        }
        for f in freq.iter_mut() {
            *f = base;
        }
    }
}

fn build_frequencies(data: &[u8]) -> Vec<u32> {
    let mut freq = vec![0u32; SYMBOL_LIMIT];
    for &b in data {
        freq[b as usize] += 1;
    }
    freq[EOF_SYMBOL as usize] = 1;
    scale_frequencies(&mut freq);
    freq
}

fn build_cumulative(freq: &[u32]) -> Vec<u32> {
    let mut cumulative = vec![0u32; freq.len() + 1];
    for (i, &f) in freq.iter().enumerate() {
        cumulative[i + 1] = cumulative[i] + f;
    }
    if let Some(&last) = cumulative.last() {
        if last == 0 {
            for i in 0..freq.len() {
                cumulative[i + 1] = (i + 1) as u32;
            }
        }
    }
    cumulative
}

fn write_u32_le(buf: &mut Vec<u8>, v: u32) {
    buf.extend_from_slice(&v.to_le_bytes());
}

fn read_u32_le(data: &[u8], pos: &mut usize) -> Option<u32> {
    if *pos + 4 > data.len() {
        return None;
    }
    let v = u32::from_le_bytes([data[*pos], data[*pos + 1], data[*pos + 2], data[*pos + 3]]);
    *pos += 4;
    Some(v)
}

pub fn encode_memory(input: &[u8]) -> Result<Vec<u8>, CodecError> {
    let freq = build_frequencies(input);
    let cumulative = build_cumulative(&freq);

    let mut out = Vec::new();
    out.extend_from_slice(b"AENC");
    write_u32_le(&mut out, freq.len() as u32);
    for &f in &freq {
        write_u32_le(&mut out, f);
    }

    let mut encoder = ArithmeticEncoder::new();
    for &b in input {
        encoder.encode_symbol(b as u32, &cumulative)?;
    }
    encoder.encode_symbol(EOF_SYMBOL, &cumulative)?;
    encoder.finish();

    out.extend_from_slice(&encoder.into_bytes());
    Ok(out)
}

pub fn decode_memory(input: &[u8]) -> Result<Vec<u8>, CodecError> {
    if input.len() < 4 {
        return Err(CodecError::Truncated);
    }
    if &input[0..4] != b"AENC" {
        return Err(CodecError::Corrupt);
    }

    let mut pos = 4;
    let count = read_u32_le(input, &mut pos).ok_or(CodecError::Truncated)? as usize;
    if count != SYMBOL_LIMIT {
        return Err(CodecError::Corrupt);
    }

    let mut freq = vec![0u32; count];
    for f in freq.iter_mut() {
        *f = read_u32_le(input, &mut pos).ok_or(CodecError::Truncated)?;
    }

    let cumulative = build_cumulative(&freq);
    let bit_reader = BitReader::new(&input[pos..]);
    let mut decoder = ArithmeticDecoder::new(bit_reader);

    let mut output = Vec::new();
    loop {
        let sym = decoder.decode_symbol(&cumulative)?;
        if sym == EOF_SYMBOL {
            break;
        }
        output.push(sym as u8);
    }

    Ok(output)
}

pub struct StreamingEncoder {
    state: State,
    buffer: Vec<u8>,
}

impl StreamingEncoder {
    pub fn new() -> Self {
        StreamingEncoder {
            state: State::Ready,
            buffer: Vec::new(),
        }
    }
}

impl Default for StreamingEncoder {
    fn default() -> Self {
        Self::new()
    }
}

impl Encoder for StreamingEncoder {
    fn process(&mut self, input: &[u8], _output: &mut [u8]) -> Result<usize, CodecError> {
        match self.state {
            State::Ready | State::Flushing => {
                self.state = State::Streaming;
                self.buffer.extend_from_slice(input);
                Ok(0)
            }
            State::Streaming => {
                self.buffer.extend_from_slice(input);
                Ok(0)
            }
            State::Finished => {
                self.state = State::Error;
                Err(CodecError::InvalidState)
            }
            State::Error => Err(CodecError::InvalidState),
        }
    }

    fn flush(&mut self, _output: &mut [u8]) -> Result<usize, CodecError> {
        match self.state {
            State::Ready => Ok(0),
            State::Streaming => {
                self.state = State::Flushing;
                Ok(0)
            }
            State::Flushing => Ok(0),
            State::Finished => {
                self.state = State::Error;
                Err(CodecError::InvalidState)
            }
            State::Error => Err(CodecError::InvalidState),
        }
    }

    fn finish(&mut self, output: &mut [u8]) -> Result<usize, CodecError> {
        match self.state {
            State::Ready | State::Streaming | State::Flushing => {
                let encoded = encode_memory(&self.buffer)?;
                if output.len() < encoded.len() {
                    return Err(CodecError::BufTooSmall);
                }
                output[..encoded.len()].copy_from_slice(&encoded);
                self.state = State::Finished;
                Ok(encoded.len())
            }
            State::Finished => {
                self.state = State::Error;
                Err(CodecError::InvalidState)
            }
            State::Error => Err(CodecError::InvalidState),
        }
    }

    fn reset(&mut self) {
        self.state = State::Ready;
        self.buffer.clear();
    }

    fn state(&self) -> State {
        self.state
    }
}

pub struct StreamingDecoder {
    state: State,
    buffer: Vec<u8>,
}

impl StreamingDecoder {
    pub fn new() -> Self {
        StreamingDecoder {
            state: State::Ready,
            buffer: Vec::new(),
        }
    }
}

impl Default for StreamingDecoder {
    fn default() -> Self {
        Self::new()
    }
}

impl Decoder for StreamingDecoder {
    fn process(&mut self, input: &[u8], _output: &mut [u8]) -> Result<usize, CodecError> {
        match self.state {
            State::Ready | State::Flushing => {
                self.state = State::Streaming;
                self.buffer.extend_from_slice(input);
                Ok(0)
            }
            State::Streaming => {
                self.buffer.extend_from_slice(input);
                Ok(0)
            }
            State::Finished => {
                self.state = State::Error;
                Err(CodecError::InvalidState)
            }
            State::Error => Err(CodecError::InvalidState),
        }
    }

    fn flush(&mut self, _output: &mut [u8]) -> Result<usize, CodecError> {
        match self.state {
            State::Ready => Ok(0),
            State::Streaming => {
                self.state = State::Flushing;
                Ok(0)
            }
            State::Flushing => Ok(0),
            State::Finished => {
                self.state = State::Error;
                Err(CodecError::InvalidState)
            }
            State::Error => Err(CodecError::InvalidState),
        }
    }

    fn finish(&mut self, output: &mut [u8]) -> Result<usize, CodecError> {
        match self.state {
            State::Ready | State::Streaming | State::Flushing => {
                let decoded = decode_memory(&self.buffer)?;
                if output.len() < decoded.len() {
                    return Err(CodecError::BufTooSmall);
                }
                output[..decoded.len()].copy_from_slice(&decoded);
                self.state = State::Finished;
                Ok(decoded.len())
            }
            State::Finished => {
                self.state = State::Error;
                Err(CodecError::InvalidState)
            }
            State::Error => Err(CodecError::InvalidState),
        }
    }

    fn reset(&mut self) {
        self.state = State::Ready;
        self.buffer.clear();
    }

    fn state(&self) -> State {
        self.state
    }
}
