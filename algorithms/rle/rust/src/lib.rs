// Minimal in-memory RLE encode/decode API for testing.
// Main CLI remains in main.rs. This library provides simple helpers.

use std::io;

const MAX_OUTPUT_SIZE: usize = 1024 * 1024 * 1024;

pub fn encode(input: &[u8]) -> Result<Vec<u8>, io::Error> {
    if input.is_empty() {
        return Ok(Vec::new());
    }
    
    let mut output = Vec::new();
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
    let mut output = Vec::new();
    let mut pos = 0;
    
    while pos < input.len() {
        if pos + 5 > input.len() {
            return Err(io::Error::new(
                io::ErrorKind::InvalidData,
                "RLE data truncated: incomplete count+value pair",
            ));
        }
        
        let count = u32::from_le_bytes([
            input[pos],
            input[pos + 1],
            input[pos + 2],
            input[pos + 3],
        ]);
        let value = input[pos + 4];
        pos += 5;
        
        let new_len = output.len() + count as usize;
        if new_len > MAX_OUTPUT_SIZE {
            return Err(io::Error::new(
                io::ErrorKind::InvalidData,
                "output size exceeds maximum allowed (1 GiB)",
            ));
        }
        
        for _ in 0..count {
            output.push(value);
        }
    }
    
    Ok(output)
}
