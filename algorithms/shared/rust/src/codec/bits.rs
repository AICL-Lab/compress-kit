//! Bit-level I/O utilities for entropy coding algorithms.
//!
//! This module provides `BitWriter` and `BitReader` for efficient bit-level
//! operations used in Huffman and Arithmetic coding.

/// A bit writer that packs bits into bytes (MSB first).
///
/// Bits are written from the most significant bit (MSB) to the least significant bit (LSB)
/// within each byte. This is the standard convention for most entropy coding algorithms.
///
/// # Example
///
/// ```
/// use compresskit_codec::codec::bits::BitWriter;
///
/// let mut writer = BitWriter::new();
/// writer.write_bit(true);
/// writer.write_bit(false);
/// writer.write_bit(true);
/// let bytes = writer.finish();
/// assert_eq!(bytes, vec![0b10100000]);
/// ```
pub struct BitWriter {
    buffer: Vec<u8>,
    current_byte: u8,
    bit_pos: u8, // 0-7, next bit position to write (0 = MSB)
}

impl BitWriter {
    /// Creates a new empty bit writer.
    pub fn new() -> Self {
        BitWriter {
            buffer: Vec::new(),
            current_byte: 0,
            bit_pos: 0,
        }
    }

    /// Creates a new bit writer with pre-allocated capacity.
    pub fn with_capacity(capacity: usize) -> Self {
        BitWriter {
            buffer: Vec::with_capacity(capacity),
            current_byte: 0,
            bit_pos: 0,
        }
    }

    /// Writes a single bit (MSB first).
    ///
    /// The bit is written at the current position, starting from the MSB.
    pub fn write_bit(&mut self, bit: bool) {
        if bit {
            self.current_byte |= 1 << (7 - self.bit_pos);
        }
        self.bit_pos += 1;
        if self.bit_pos == 8 {
            self.buffer.push(self.current_byte);
            self.current_byte = 0;
            self.bit_pos = 0;
        }
    }

    /// Writes a bit as a u8 (0 or 1).
    ///
    /// Convenience method for algorithms that use `u8` for bits.
    pub fn write_bit_u8(&mut self, bit: u8) {
        self.write_bit(bit != 0);
    }

    /// Writes multiple bits from a slice of bools.
    pub fn write_bits(&mut self, bits: &[bool]) {
        for &bit in bits {
            self.write_bit(bit);
        }
    }

    /// Flushes any remaining bits in the current byte.
    ///
    /// If there are partial bits (less than 8), they are packed into a final byte
    /// with the remaining bits set to 0.
    pub fn flush(&mut self) {
        if self.bit_pos > 0 {
            self.buffer.push(self.current_byte);
            self.current_byte = 0;
            self.bit_pos = 0;
        }
    }

    /// Finishes writing and returns the complete byte buffer.
    ///
    /// This is equivalent to calling `flush()` followed by returning the buffer.
    pub fn finish(mut self) -> Vec<u8> {
        self.flush();
        self.buffer
    }

    /// Returns the number of complete bytes written so far.
    pub fn byte_len(&self) -> usize {
        self.buffer.len()
    }

    /// Returns the total number of bits written (including partial byte).
    pub fn bit_len(&self) -> usize {
        self.buffer.len() * 8 + self.bit_pos as usize
    }
}

impl Default for BitWriter {
    fn default() -> Self {
        Self::new()
    }
}

/// A bit reader that unpacks bits from bytes (MSB first).
///
/// Bits are read from the most significant bit (MSB) to the least significant bit (LSB)
/// within each byte.
pub struct BitReader<'a> {
    data: &'a [u8],
    byte_pos: usize,
    bit_pos: u8, // 0-7, current bit position within byte (0 = MSB)
}

impl<'a> BitReader<'a> {
    /// Creates a new bit reader from a byte slice.
    pub fn new(data: &'a [u8]) -> Self {
        BitReader {
            data,
            byte_pos: 0,
            bit_pos: 0,
        }
    }

    /// Reads a single bit.
    ///
    /// Returns `true` for 1, `false` for 0.
    /// Returns `false` after reaching end of data.
    pub fn read_bit(&mut self) -> bool {
        if self.byte_pos >= self.data.len() {
            return false;
        }

        let bit = (self.data[self.byte_pos] >> (7 - self.bit_pos)) & 1 == 1;
        self.bit_pos += 1;

        if self.bit_pos == 8 {
            self.byte_pos += 1;
            self.bit_pos = 0;
        }

        bit
    }

    /// Reads a single bit as u8 (0 or 1).
    pub fn read_bit_u8(&mut self) -> u8 {
        self.read_bit() as u8
    }

    /// Reads multiple bits into a slice of bools.
    pub fn read_bits(&mut self, bits: &mut [bool]) {
        for bit in bits.iter_mut() {
            *bit = self.read_bit();
        }
    }

    /// Returns true if all data has been consumed.
    pub fn is_empty(&self) -> bool {
        self.byte_pos >= self.data.len()
    }

    /// Returns the number of bits remaining (including partial byte).
    pub fn bits_remaining(&self) -> usize {
        if self.byte_pos >= self.data.len() {
            0
        } else {
            (self.data.len() - self.byte_pos) * 8 - self.bit_pos as usize
        }
    }

    /// Returns the current byte position.
    pub fn byte_position(&self) -> usize {
        self.byte_pos
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn bit_writer_basic() {
        let mut writer = BitWriter::new();
        writer.write_bit(true);  // 1
        writer.write_bit(false); // 0
        writer.write_bit(true);  // 1
        writer.write_bit(false); // 0
        writer.write_bit(true);  // 1
        writer.write_bit(false); // 0
        writer.write_bit(true);  // 1
        writer.write_bit(false); // 0
        let bytes = writer.finish();
        assert_eq!(bytes, vec![0b10101010]);
    }

    #[test]
    fn bit_writer_partial_byte() {
        let mut writer = BitWriter::new();
        writer.write_bit(true);
        writer.write_bit(true);
        writer.write_bit(false);
        let bytes = writer.finish();
        assert_eq!(bytes, vec![0b11000000]);
    }

    #[test]
    fn bit_writer_multiple_bytes() {
        let mut writer = BitWriter::new();
        // First byte: 11111111
        for _ in 0..8 {
            writer.write_bit(true);
        }
        // Second byte: 00000000
        for _ in 0..8 {
            writer.write_bit(false);
        }
        // Third byte (partial): 101
        writer.write_bit(true);
        writer.write_bit(false);
        writer.write_bit(true);

        let bytes = writer.finish();
        assert_eq!(bytes, vec![0b11111111, 0b00000000, 0b10100000]);
    }

    #[test]
    fn bit_reader_basic() {
        let data = vec![0b10101010];
        let mut reader = BitReader::new(&data);

        assert!(reader.read_bit());  // 1
        assert!(!reader.read_bit()); // 0
        assert!(reader.read_bit());  // 1
        assert!(!reader.read_bit()); // 0
        assert!(reader.read_bit());  // 1
        assert!(!reader.read_bit()); // 0
        assert!(reader.read_bit());  // 1
        assert!(!reader.read_bit()); // 0
        assert!(reader.is_empty());
    }

    #[test]
    fn bit_reader_partial_byte() {
        let data = vec![0b11000000];
        let mut reader = BitReader::new(&data);

        assert!(reader.read_bit());
        assert!(reader.read_bit());
        assert!(!reader.read_bit());
        // After reading 3 bits, there are still 5 bits remaining
        assert_eq!(reader.bits_remaining(), 5);
        assert!(!reader.is_empty());
    }

    #[test]
    fn roundtrip() {
        let original = vec![true, false, true, true, false, true, false, false, true];

        // Write
        let mut writer = BitWriter::new();
        for &bit in &original {
            writer.write_bit(bit);
        }
        let bytes = writer.finish();

        // Read
        let mut reader = BitReader::new(&bytes);
        let mut decoded = vec![false; original.len()];
        reader.read_bits(&mut decoded);

        assert_eq!(decoded, original);
    }

    #[test]
    fn write_bit_u8() {
        let mut writer = BitWriter::new();
        writer.write_bit_u8(1);
        writer.write_bit_u8(0);
        writer.write_bit_u8(1);
        let bytes = writer.finish();
        assert_eq!(bytes, vec![0b10100000]);
    }

    #[test]
    fn bit_len_tracking() {
        let mut writer = BitWriter::new();
        assert_eq!(writer.bit_len(), 0);
        assert_eq!(writer.byte_len(), 0);

        writer.write_bit(true);
        assert_eq!(writer.bit_len(), 1);
        assert_eq!(writer.byte_len(), 0);

        writer.write_bit(false);
        writer.write_bit(true);
        writer.write_bit(false);
        writer.write_bit(true);
        writer.write_bit(false);
        writer.write_bit(true);
        assert_eq!(writer.bit_len(), 7);
        assert_eq!(writer.byte_len(), 0);

        writer.write_bit(false);
        assert_eq!(writer.bit_len(), 8);
        assert_eq!(writer.byte_len(), 1);

        writer.write_bit(true);
        assert_eq!(writer.bit_len(), 9);
        assert_eq!(writer.byte_len(), 1);
    }
}
