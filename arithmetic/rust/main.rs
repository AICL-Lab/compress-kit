use std::env;
use std::fs::File;
use std::io::{self, BufReader, BufWriter, Read, Write};
use std::process;

// 算术编码 Rust 实现。
// 文件格式与 C++/Go 实现完全一致，支持交叉编解码验证。
// Magic: AENC (4 bytes)
// 频率表: count(4 bytes LE) + count × freq(4 bytes LE)
// 算术编码比特流

const SYMBOL_LIMIT: usize = 257;
const EOF_SYMBOL: u32 = (SYMBOL_LIMIT - 1) as u32;
const MAX_TOTAL: u32 = 1 << 24;

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
        let total = *cumulative.last().unwrap() as u64;
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

    fn decode_symbol(&mut self, cumulative: &[u32]) -> u32 {
        let range = self.high - self.low + 1;
        let total = *cumulative.last().unwrap() as u64;
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

        symbol
    }
}

// ---------------------------------------------------------------------------
// 频率表处理
// ---------------------------------------------------------------------------

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

fn build_frequencies_from_file(path: &str) -> io::Result<Vec<u32>> {
    let mut freq = vec![0u32; SYMBOL_LIMIT];
    let file = File::open(path)
        .map_err(|e| io::Error::new(e.kind(), format!("无法打开输入文件用于读取: {path}: {e}")))?;
    let mut reader = BufReader::new(file);
    let mut buf = [0u8; 4096];
    loop {
        match reader.read(&mut buf) {
            Ok(0) => break,
            Ok(n) => {
                for &b in &buf[..n] {
                    freq[b as usize] += 1;
                }
            }
            Err(_) => break,
        }
    }
    freq[EOF_SYMBOL as usize] = 1;
    scale_frequencies(&mut freq);
    Ok(freq)
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

fn write_frequencies<W: Write>(writer: &mut W, freq: &[u32]) -> io::Result<()> {
    let count = freq.len() as u32;
    writer.write_all(&count.to_le_bytes())?;
    for &v in freq {
        writer.write_all(&v.to_le_bytes())?;
    }
    Ok(())
}

fn read_frequencies<R: Read>(reader: &mut R) -> io::Result<Vec<u32>> {
    let mut count_bytes = [0u8; 4];
    reader
        .read_exact(&mut count_bytes)
        .map_err(|e| io::Error::new(e.kind(), format!("读取频率表失败: {e}")))?;
    let count = u32::from_le_bytes(count_bytes) as usize;
    if count != SYMBOL_LIMIT {
        return Err(io::Error::new(
            io::ErrorKind::InvalidData,
            format!("频率表大小异常: {count}"),
        ));
    }
    let mut freq = vec![0u32; count];
    for f in freq.iter_mut() {
        let mut arr = [0u8; 4];
        reader
            .read_exact(&mut arr)
            .map_err(|e| io::Error::new(e.kind(), format!("读取频率表失败: {e}")))?;
        *f = u32::from_le_bytes(arr);
    }
    Ok(freq)
}

// ---------------------------------------------------------------------------
// 压缩 / 解压
// ---------------------------------------------------------------------------

fn compress_file(input_path: &str, output_path: &str) -> io::Result<()> {
    let freq = build_frequencies_from_file(input_path)?;
    let cumulative = build_cumulative(&freq);

    let input_file = File::open(input_path)?;
    let mut reader = BufReader::new(input_file);
    let output_file = File::create(output_path)?;
    let mut writer = BufWriter::new(output_file);

    writer.write_all(b"AENC")?;
    write_frequencies(&mut writer, &freq)?;

    let bit_writer = BitWriter::new(writer);
    let mut encoder = ArithmeticEncoder::new(bit_writer);

    let mut buf = [0u8; 4096];
    loop {
        let n = reader.read(&mut buf)?;
        if n == 0 {
            break;
        }
        for &b in &buf[..n] {
            encoder.encode_symbol(b as u32, &cumulative)?;
        }
    }
    encoder.encode_symbol(EOF_SYMBOL, &cumulative)?;
    encoder.finish()?;
    Ok(())
}

fn decompress_file(input_path: &str, output_path: &str) -> io::Result<()> {
    let file = File::open(input_path)?;
    let mut reader = BufReader::new(file);
    let mut magic = [0u8; 4];
    reader.read_exact(&mut magic)?;
    if &magic != b"AENC" {
        return Err(io::Error::new(
            io::ErrorKind::InvalidData,
            "输入文件格式非法",
        ));
    }
    let freq = read_frequencies(&mut reader)?;
    let cumulative = build_cumulative(&freq);

    let output_file = File::create(output_path)?;
    let mut writer = BufWriter::new(output_file);

    let bit_reader = BitReader::new(reader);
    let mut decoder = ArithmeticDecoder::new(bit_reader);

    loop {
        let sym = decoder.decode_symbol(&cumulative);
        if sym == EOF_SYMBOL {
            break;
        }
        writer.write_all(&[sym as u8])?;
    }

    writer.flush()?;
    Ok(())
}

pub fn arithmetic_encode_file(input_path: &str, output_path: &str) -> io::Result<()> {
    compress_file(input_path, output_path)
}

pub fn arithmetic_decode_file(input_path: &str, output_path: &str) -> io::Result<()> {
    decompress_file(input_path, output_path)
}

#[cfg(test)]
mod tests {
    use super::*;
    use std::fs;
    use std::path::PathBuf;
    use std::time::{SystemTime, UNIX_EPOCH};

    fn make_paths(prefix: &str) -> (PathBuf, PathBuf, PathBuf, PathBuf) {
        let mut dir = std::env::temp_dir();
        let stamp = SystemTime::now()
            .duration_since(UNIX_EPOCH)
            .unwrap()
            .as_nanos();
        dir.push(format!(
            "encoding_arith_{}_{}_{}",
            prefix,
            std::process::id(),
            stamp
        ));
        fs::create_dir_all(&dir).unwrap();
        let input = dir.join("input.bin");
        let encoded = dir.join("encoded.aenc");
        let output = dir.join("output.bin");
        (dir, input, encoded, output)
    }

    #[test]
    fn roundtrip_bytes() {
        let (dir, input, encoded, output) = make_paths("roundtrip");
        let mut data = b"arithmetic-rust-test-data-".repeat(256);
        data.extend_from_slice(&[0, 1, 2, 3, 254, 255]);
        fs::write(&input, &data).unwrap();

        arithmetic_encode_file(input.to_str().unwrap(), encoded.to_str().unwrap()).unwrap();
        arithmetic_decode_file(encoded.to_str().unwrap(), output.to_str().unwrap()).unwrap();

        let decoded = fs::read(&output).unwrap();
        assert_eq!(decoded, data);
        fs::remove_dir_all(dir).unwrap();
    }

    #[test]
    fn roundtrip_empty() {
        let (dir, input, encoded, output) = make_paths("empty");
        fs::write(&input, Vec::<u8>::new()).unwrap();

        arithmetic_encode_file(input.to_str().unwrap(), encoded.to_str().unwrap()).unwrap();
        arithmetic_decode_file(encoded.to_str().unwrap(), output.to_str().unwrap()).unwrap();

        let decoded = fs::read(&output).unwrap();
        assert!(decoded.is_empty());
        fs::remove_dir_all(dir).unwrap();
    }

    #[test]
    fn roundtrip_single_byte() {
        let (dir, input, encoded, output) = make_paths("single");
        fs::write(&input, [0x42]).unwrap();

        arithmetic_encode_file(input.to_str().unwrap(), encoded.to_str().unwrap()).unwrap();
        arithmetic_decode_file(encoded.to_str().unwrap(), output.to_str().unwrap()).unwrap();

        let decoded = fs::read(&output).unwrap();
        assert_eq!(decoded, vec![0x42]);
        fs::remove_dir_all(dir).unwrap();
    }

    #[test]
    fn roundtrip_all_bytes() {
        let (dir, input, encoded, output) = make_paths("allbytes");
        let data: Vec<u8> = (0..=255).collect();
        fs::write(&input, &data).unwrap();

        arithmetic_encode_file(input.to_str().unwrap(), encoded.to_str().unwrap()).unwrap();
        arithmetic_decode_file(encoded.to_str().unwrap(), output.to_str().unwrap()).unwrap();

        let decoded = fs::read(&output).unwrap();
        assert_eq!(decoded, data);
        fs::remove_dir_all(dir).unwrap();
    }
}

fn main() {
    let args: Vec<String> = env::args().collect();
    if args.len() != 4 {
        eprintln!("用法: {} encode|decode input output", args[0]);
        process::exit(1);
    }
    let mode = &args[1];
    let input_path = &args[2];
    let output_path = &args[3];

    let result = match mode.as_str() {
        "encode" => arithmetic_encode_file(input_path, output_path),
        "decode" => arithmetic_decode_file(input_path, output_path),
        _ => {
            eprintln!("未知模式，应为 encode 或 decode");
            process::exit(1);
        }
    };

    if let Err(e) = result {
        eprintln!("运行失败: {e}");
        process::exit(1);
    }
}
