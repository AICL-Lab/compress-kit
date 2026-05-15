// Minimal in-memory Huffman encode/decode API for testing.
// Main CLI remains in main.rs. This library provides simple helpers.

use std::cmp::Ordering;
use std::collections::BinaryHeap;
use std::io;

use compresskit_codec::codec::{
    build_frequencies, read_frequencies_exact, write_frequencies, BitWriter, EOF_SYMBOL,
    SYMBOL_LIMIT,
};
const MAX_OUTPUT_SIZE: usize = 1024 * 1024 * 1024; // 1 GiB

struct Node {
    symbol: u32,
    freq: u64,
    left: Option<Box<Node>>,
    right: Option<Box<Node>>,
}

fn is_leaf(node: &Node) -> bool {
    node.left.is_none() && node.right.is_none()
}

struct HeapItem {
    freq: u64,
    symbol: u32,
    node: Box<Node>,
}

impl Eq for HeapItem {}
impl PartialEq for HeapItem {
    fn eq(&self, other: &Self) -> bool {
        self.freq == other.freq && self.symbol == other.symbol
    }
}
impl Ord for HeapItem {
    fn cmp(&self, other: &Self) -> Ordering {
        other
            .freq
            .cmp(&self.freq)
            .then_with(|| other.symbol.cmp(&self.symbol))
    }
}
impl PartialOrd for HeapItem {
    fn partial_cmp(&self, other: &Self) -> Option<Ordering> {
        Some(self.cmp(other))
    }
}

fn build_tree(freq: &[u32]) -> Box<Node> {
    let mut heap = BinaryHeap::<HeapItem>::new();
    for (s, &f) in freq.iter().enumerate() {
        if f == 0 {
            continue;
        }
        heap.push(HeapItem {
            freq: f as u64,
            symbol: s as u32,
            node: Box::new(Node {
                symbol: s as u32,
                freq: f as u64,
                left: None,
                right: None,
            }),
        });
    }
    if heap.is_empty() {
        return Box::new(Node {
            symbol: EOF_SYMBOL,
            freq: 1,
            left: None,
            right: None,
        });
    }
    if heap.len() == 1 {
        let item = heap.pop().unwrap();
        let only = item.node;
        heap.push(HeapItem {
            freq: only.freq,
            symbol: only.symbol,
            node: Box::new(Node {
                symbol: only.symbol,
                freq: only.freq,
                left: Some(only),
                right: None,
            }),
        });
    }
    while heap.len() > 1 {
        let a = heap.pop().unwrap().node;
        let b = heap.pop().unwrap().node;
        let freq_sum = a.freq + b.freq;
        heap.push(HeapItem {
            freq: freq_sum,
            symbol: a.symbol.min(b.symbol),
            node: Box::new(Node {
                symbol: a.symbol.min(b.symbol),
                freq: freq_sum,
                left: Some(a),
                right: Some(b),
            }),
        });
    }
    heap.pop().unwrap().node
}

/// Build Huffman codes using Vec<bool> instead of String for efficiency.
fn build_codes_bitvec(node: &Node, codes: &mut [Vec<bool>], prefix: &mut Vec<bool>) {
    if is_leaf(node) {
        codes[node.symbol as usize] = if prefix.is_empty() {
            vec![false]
        } else {
            prefix.clone()
        };
        return;
    }
    if let Some(ref left) = node.left {
        prefix.push(false);
        build_codes_bitvec(left, codes, prefix);
        prefix.pop();
    }
    if let Some(ref right) = node.right {
        prefix.push(true);
        build_codes_bitvec(right, codes, prefix);
        prefix.pop();
    }
}

pub fn encode(input: &[u8]) -> Result<Vec<u8>, io::Error> {
    let freq = build_frequencies(input);

    let root = build_tree(&freq);
    let mut codes = vec![Vec::new(); SYMBOL_LIMIT];
    let mut prefix = Vec::new();
    build_codes_bitvec(&root, &mut codes, &mut prefix);

    let mut output = Vec::new();
    output.extend_from_slice(b"HFMN");
    write_frequencies(&mut output, &freq);

    // Use BitWriter for efficient bit packing
    let mut writer = BitWriter::new();
    for &b in input {
        for &bit in &codes[b as usize] {
            writer.write_bit(bit);
        }
    }
    // Write EOF code
    for &bit in &codes[EOF_SYMBOL as usize] {
        writer.write_bit(bit);
    }

    output.extend(writer.finish());
    Ok(output)
}

pub fn decode(input: &[u8]) -> Result<Vec<u8>, io::Error> {
    if input.len() < 8 {
        return Err(io::Error::new(
            io::ErrorKind::InvalidData,
            "input too short",
        ));
    }
    if &input[0..4] != b"HFMN" {
        return Err(io::Error::new(io::ErrorKind::InvalidData, "invalid magic"));
    }

    let mut pos = 4;
    let freq = read_frequencies_exact(
        input,
        &mut pos,
        SYMBOL_LIMIT,
        "truncated freq table",
        "invalid freq count",
    )
    .map_err(|err| io::Error::new(io::ErrorKind::InvalidData, err.message))?;

    let root = build_tree(&freq);
    let mut output = Vec::new();
    let mut node = &root;

    for &byte in &input[pos..] {
        for i in (0..8).rev() {
            let bit = (byte >> i) & 1;
            node = if bit == 0 {
                node.left.as_ref()
            } else {
                node.right.as_ref()
            }
            .ok_or_else(|| io::Error::new(io::ErrorKind::InvalidData, "invalid tree traversal"))?;

            if is_leaf(node) {
                if node.symbol == EOF_SYMBOL {
                    return Ok(output);
                }
                if output.len() >= MAX_OUTPUT_SIZE {
                    return Err(io::Error::new(
                        io::ErrorKind::InvalidData,
                        "output size exceeds maximum allowed (1 GiB)",
                    ));
                }
                output.push(node.symbol as u8);
                node = &root;
            }
        }
    }
    Err(io::Error::new(io::ErrorKind::InvalidData, "no EOF found"))
}

// Streaming adapters
use compresskit_codec::codec::{
    io_error_to_codec_error, streaming_decoder, streaming_encoder, CodecError, Decoder, Encoder,
};

/// Creates a new streaming Huffman encoder.
pub fn new_encoder() -> impl Encoder {
    streaming_encoder(huffman_encode)
}

/// Creates a new streaming Huffman decoder.
pub fn new_decoder() -> impl Decoder {
    streaming_decoder(huffman_decode)
}

fn huffman_encode(input: &[u8]) -> Result<Vec<u8>, CodecError> {
    encode(input).map_err(io_error_to_codec_error)
}

fn huffman_decode(input: &[u8]) -> Result<Vec<u8>, CodecError> {
    decode(input).map_err(io_error_to_codec_error)
}
