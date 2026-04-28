use std::env;
use std::fs::File;
use std::io::{self, BufReader, BufWriter, Read, Write};
use std::process;

// Arithmetic coding Rust implementation.
// File format is fully compatible with C++/Go implementations, supports cross-language encode/decode verification.
// Magic: AENC (4 bytes)
// Frequency table: count(4 bytes LE) + count × freq(4 bytes LE)
// Arithmetic coding bitstream

const SYMBOL_LIMIT: usize = 257;
const EOF_SYMBOL: u32 = (SYMBOL_LIMIT - 1) as u32;
const MAX_TOTAL: u32 = 1 << 24;
const MAX_INPUT_SIZE: u64 = 4 * 1024 * 1024 * 1024; // 4 GiB max

const STATE_BITS: u64 = 32;
const FULL_RANGE: u64 = 1u64 << STATE_BITS;
const HALF_RANGE: u64 = FULL_RANGE >> 1;
const FIRST_QUARTER: u64 = HALF_RANGE >> 1;
const THIRD_QUARTER: u64 = FIRST_QUARTER * 3;

// ---------------------------------------------------------------------------
// BitWriter / BitReader
// ---------------------------------------------------------------------------

struct BitWriter<W: Write> {
    writer: W,
    buffer: u8,
    bits_in_buffer: u8,
}

impl<W: Write> BitWriter<W> {
    fn new(writer: W) -> Self {
        BitWriter {
            writer,
            buffer: 0,
            bits_in_buffer: 0,
        }
    }

    fn write_bit(&mut self, bit: u8) -> io::Result<()> {
        self.buffer = (self.buffer << 1) | (bit & 1);
        self.bits_in_buffer += 1;
        if self.bits_in_buffer == 8 {
            self.writer.write_all(&[self.buffer])?;
            self.bits_in_buffer = 0;
            self.buffer = 0;
        }
        Ok(())
    }

    fn flush(&mut self) -> io::Result<()> {
        if self.bits_in_buffer > 0 {
            self.buffer <<= 8 - self.bits_in_buffer;
            self.writer.write_all(&[self.buffer])?;
            self.bits_in_buffer = 0;
            self.buffer = 0;
        }
        self.writer.flush()
    }
}

struct BitReader<R: Read> {
    reader: R,
    current_byte: u8,
    bits_remaining: u8,
    reached_eof: bool,
}

impl<R: Read> BitReader<R> {
    fn new(reader: R) -> Self {
        BitReader {
            reader,
            current_byte: 0,
            bits_remaining: 0,
            reached_eof: false,
        }
    }

    fn read_bit(&mut self) -> u8 {
        if self.bits_remaining == 0 {
            let mut buf = [0u8; 1];
            match self.reader.read(&mut buf) {
                Ok(0) | Err(_) => {
                    self.reached_eof = true;
                    return 0;
                }
                Ok(_) => {
                    self.current_byte = buf[0];
                    self.bits_remaining = 8;
                }
            }
        }
        self.bits_remaining -= 1;
        (self.current_byte >> self.bits_remaining) & 1
    }
}

// ---------------------------------------------------------------------------
// ArithmeticEncoder
// ---------------------------------------------------------------------------

struct ArithmeticEncoder<W: Write> {
    writer: BitWriter<W>,
    low: u64,
    high: u64,
    pending_bits: u64,
}

impl<W: Write> ArithmeticEncoder<W> {
    fn new(writer: BitWriter<W>) -> Self {
        ArithmeticEncoder {
            writer,
            low: 0,
            high: FULL_RANGE - 1,
            pending_bits: 0,
        }
    }

    fn encode_symbol(&mut self, symbol: u32, cumulative: &[u32]) -> io::Result<()> {
        let range = self.high - self.low + 1;
        let total = *cumulative.last().ok_or_else(|| {
            io::Error::new(
                io::ErrorKind::InvalidData,
                "empty cumulative frequency table",
            )
        })? as u64;
        let sym_low = cumulative[symbol as usize] as u64;
        let sym_high = cumulative[symbol as usize + 1] as u64;

        self.high = self.low + (range * sym_high) / total - 1;
        self.low = self.low + (range * sym_low) / total;

        loop {
            if self.high < HALF_RANGE {
                self.output_bit(0)?;
            } else if self.low >= HALF_RANGE {
                self.output_bit(1)?;
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

    fn finish(&mut self) -> io::Result<()> {
        self.pending_bits += 1;
        if self.low < FIRST_QUARTER {
            self.output_bit(0)?;
        } else {
            self.output_bit(1)?;
        }
        self.writer.flush()
    }

    fn output_bit(&mut self, bit: u8) -> io::Result<()> {
        self.writer.write_bit(bit)?;
        let complement = bit ^ 1;
        while self.pending_bits > 0 {
            self.writer.write_bit(complement)?;
            self.pending_bits -= 1;
        }
        Ok(())
    }
}

// ---------------------------------------------------------------------------
// ArithmeticDecoder
// ---------------------------------------------------------------------------

struct ArithmeticDecoder<R: Read> {
    reader: BitReader<R>,
    low: u64,
    high: u64,
    code: u64,
}

impl<R: Read> ArithmeticDecoder<R> {
    fn new(mut reader: BitReader<R>) -> Self {
        let mut code: u64 = 0;
        for _ in 0..STATE_BITS {
            code = (code << 1) | reader.read_bit() as u64;
        }
        ArithmeticDecoder {
            reader,
            low: 0,
            high: FULL_RANGE - 1,
            code,
        }
    }

    fn decode_symbol(&mut self, cumulative: &[u32]) -> io::Result<u32> {
        let range = self.high - self.low + 1;
        let total = *cumulative.last().ok_or_else(|| {
            io::Error::new(
                io::ErrorKind::InvalidData,
                "empty cumulative frequency table",
            )
        })? as u64;
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

// ---------------------------------------------------------------------------
