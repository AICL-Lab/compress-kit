use std::env;
use std::fs::File;
use std::io::{self, BufReader, BufWriter, Read, Write};
use std::process;

use compresskit_codec::codec::{decode_buffer, encode_buffer};

// Simple Run-Length encoding implementation.
// Format: repeatedly write 4-byte little-endian count + 1-byte value until input ends.
// All three language implementations use the same format for cross-validation and benchmarking.

// Maximum output size limit (1 GiB) to prevent decompression bomb attacks
const MAX_OUTPUT_SIZE: u64 = 1024 * 1024 * 1024;

fn write_u32_le<W: Write>(w: &mut W, v: u32) -> io::Result<()> {
    let bytes = v.to_le_bytes();
    w.write_all(&bytes)
}

// Read a 32-bit little-endian unsigned integer from stream.
// Returns Ok(Some(v)) on successful read;
// Returns Ok(None) on normal EOF (no bytes read);
// Returns Err(...) on I/O or truncation error.
fn read_u32_le<R: Read>(r: &mut R) -> io::Result<Option<u32>> {
    let mut buf = [0u8; 4];
    let mut read = 0usize;
    while read < 4 {
        match r.read(&mut buf[read..]) {
            Ok(0) => {
                if read == 0 {
                    // Normal EOF
                    return Ok(None);
                } else {
                    return Err(io::Error::new(
                        io::ErrorKind::UnexpectedEof,
                        "RLE data truncated: cannot read complete count field",
                    ));
                }
            }
            Ok(n) => {
                read += n;
            }
            Err(e) if e.kind() == io::ErrorKind::Interrupted => {
                continue;
            }
            Err(e) => return Err(e),
        }
    }
    Ok(Some(u32::from_le_bytes(buf)))
}

// Perform Run-Length encoding on entire file.
pub fn rle_encode_file(input_path: &str, output_path: &str) -> io::Result<()> {
    let input = File::open(input_path).map_err(|e| {
        io::Error::new(
            e.kind(),
            format!("cannot open input file for reading: {input_path}: {e}"),
        )
    })?;
    let mut reader = BufReader::new(input);

    let output = File::create(output_path).map_err(|e| {
        io::Error::new(
            e.kind(),
            format!("cannot open output file for writing: {output_path}: {e}"),
        )
    })?;
    let mut writer = BufWriter::new(output);

    let mut first = [0u8; 1];
    let n = reader.read(&mut first)?;
    if n == 0 {
        // Empty file
        writer.flush()?;
        return Ok(());
    }
    let mut current = first[0];
    let mut count: u32 = 1;

    let mut buf = [0u8; 4096];

    loop {
        let n = reader.read(&mut buf)?;
        if n == 0 {
            break;
        }
        for &b in &buf[..n] {
            if b == current && count < u32::MAX {
                count += 1;
            } else {
                write_u32_le(&mut writer, count)?;
                writer.write_all(&[current])?;
                current = b;
                count = 1;
            }
        }
    }

    // Write last run
    write_u32_le(&mut writer, count)?;
    writer.write_all(&[current])?;
    writer.flush()?;
    Ok(())
}

// Decode RLE encoded file back to original byte stream.
pub fn rle_decode_file(input_path: &str, output_path: &str) -> io::Result<()> {
    let input = File::open(input_path).map_err(|e| {
        io::Error::new(
            e.kind(),
            format!("cannot open input file for reading: {input_path}: {e}"),
        )
    })?;
    let mut reader = BufReader::new(input);

    let output = File::create(output_path).map_err(|e| {
        io::Error::new(
            e.kind(),
            format!("cannot open output file for writing: {output_path}: {e}"),
        )
    })?;
    let mut writer = BufWriter::new(output);

    const BUF_SIZE: usize = 4096;
    let mut buf = [0u8; BUF_SIZE];
    let mut total_written: u64 = 0;

    loop {
        let count_opt = read_u32_le(&mut reader)?;
        let count = match count_opt {
            Some(c) => c,
            None => break, // Normal EOF
        };
        if count == 0 {
            return Err(io::Error::new(
                io::ErrorKind::InvalidData,
                "invalid RLE data: count should not be 0",
            ));
        }

        // Check output size limit
        if total_written + count as u64 > MAX_OUTPUT_SIZE {
            return Err(io::Error::new(
                io::ErrorKind::Other,
                format!("output size limit exceeded (max {} bytes)", MAX_OUTPUT_SIZE),
            ));
        }

        let mut value_buf = [0u8; 1];
        reader
            .read_exact(&mut value_buf)
            .map_err(|e| io::Error::new(e.kind(), "RLE data truncated: missing value byte"))?;
        let value = value_buf[0];

        let mut remaining = count;
        while remaining > 0 {
            let chunk = remaining.min(BUF_SIZE as u32) as usize;
            buf[..chunk].fill(value);
            writer.write_all(&buf[..chunk])?;
            total_written += chunk as u64;
            remaining -= chunk as u32;
        }
    }

    writer.flush()?;
    Ok(())
}

fn main() {
    let args: Vec<String> = env::args().collect();
    if args.len() != 4 {
        eprintln!("usage: {} encode|decode input output", args[0]);
        process::exit(1);
    }

    let mode = &args[1];
    let input_path = &args[2];
    let output_path = &args[3];

    let result = match mode.as_str() {
        "encode" => run_encode(input_path, output_path),
        "decode" => run_decode(input_path, output_path),
        _ => {
            eprintln!("unknown mode, expected encode or decode");
            process::exit(1);
        }
    };

    if let Err(e) = result {
        eprintln!("execution failed: {e}");
        process::exit(1);
    }
}

fn run_encode(input_path: &str, output_path: &str) -> io::Result<()> {
    let input = std::fs::read(input_path).map_err(|e| {
        io::Error::new(
            e.kind(),
            format!("cannot open input file for reading: {input_path}: {e}"),
        )
    })?;
    let mut encoder = rle::StreamingEncoder::new();
    let encoded = encode_buffer(&mut encoder, &input)
        .map_err(|e| io::Error::new(io::ErrorKind::Other, e.to_string()))?;
    std::fs::write(output_path, encoded).map_err(|e| {
        io::Error::new(
            e.kind(),
            format!("cannot open output file for writing: {output_path}: {e}"),
        )
    })
}

fn run_decode(input_path: &str, output_path: &str) -> io::Result<()> {
    let input = std::fs::read(input_path).map_err(|e| {
        io::Error::new(
            e.kind(),
            format!("cannot open input file for reading: {input_path}: {e}"),
        )
    })?;
    let mut decoder = rle::StreamingDecoder::new();
    let decoded = decode_buffer(&mut decoder, &input)
        .map_err(|e| io::Error::new(io::ErrorKind::Other, e.to_string()))?;
    std::fs::write(output_path, decoded).map_err(|e| {
        io::Error::new(
            e.kind(),
            format!("cannot open output file for writing: {output_path}: {e}"),
        )
    })
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
            "encoding_rle_{}_{}_{}",
            prefix,
            std::process::id(),
            stamp
        ));
        fs::create_dir_all(&dir).unwrap();
        let input = dir.join("input.bin");
        let encoded = dir.join("encoded.rle");
        let output = dir.join("output.bin");
        (dir, input, encoded, output)
    }

    #[test]
    fn roundtrip_bytes() {
        let (dir, input, encoded, output) = make_paths("roundtrip");
        let mut data = vec![0xAA; 2048];
        data.extend_from_slice(&b"run-length-rust-test-data-".repeat(128));
        data.extend(std::iter::repeat(0x00).take(512));
        fs::write(&input, &data).unwrap();

        rle_encode_file(input.to_str().unwrap(), encoded.to_str().unwrap()).unwrap();
        rle_decode_file(encoded.to_str().unwrap(), output.to_str().unwrap()).unwrap();

        let decoded = fs::read(&output).unwrap();
        assert_eq!(decoded, data);
        fs::remove_dir_all(dir).unwrap();
    }

    #[test]
    fn decode_rejects_zero_count() {
        let (dir, input, _, output) = make_paths("invalid");
        fs::write(&input, [0u8, 0, 0, 0, 0x42]).unwrap();

        let result = rle_decode_file(input.to_str().unwrap(), output.to_str().unwrap());
        assert!(result.is_err());
        fs::remove_dir_all(dir).unwrap();
    }
}
