// Minimal in-memory Huffman encode/decode API for testing.
// Main CLI remains in main.rs. This library provides simple helpers.

use std::cmp::Ordering;
use std::collections::BinaryHeap;
use std::io;

const SYMBOL_LIMIT: usize = 257;
const EOF_SYMBOL: u32 = (SYMBOL_LIMIT - 1) as u32;

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

fn build_codes(node: &Node, codes: &mut [String], prefix: &mut String) {
    if is_leaf(node) {
        codes[node.symbol as usize] = if prefix.is_empty() {
            "0".to_string()
        } else {
            prefix.clone()
        };
        return;
    }
    if let Some(ref left) = node.left {
        prefix.push('0');
        build_codes(left, codes, prefix);
        prefix.pop();
    }
    if let Some(ref right) = node.right {
        prefix.push('1');
        build_codes(right, codes, prefix);
        prefix.pop();
    }
}

pub fn encode(input: &[u8]) -> Result<Vec<u8>, io::Error> {
    let mut freq = vec![0u32; SYMBOL_LIMIT];
    for &b in input {
        freq[b as usize] += 1;
    }
    freq[EOF_SYMBOL as usize] = 1;

    let root = build_tree(&freq);
    let mut codes = vec![String::new(); SYMBOL_LIMIT];
    let mut prefix = String::new();
    build_codes(&root, &mut codes, &mut prefix);

    let mut output = Vec::new();
    output.extend_from_slice(b"HFMN");
    output.extend_from_slice(&(SYMBOL_LIMIT as u32).to_le_bytes());
    for &f in &freq {
        output.extend_from_slice(&f.to_le_bytes());
    }

    let mut bitstring = String::new();
    for &b in input {
        bitstring.push_str(&codes[b as usize]);
    }
    bitstring.push_str(&codes[EOF_SYMBOL as usize]);

    let mut byte = 0u8;
    let mut bit_count = 0;
    for ch in bitstring.bytes() {
        byte = (byte << 1) | (if ch == b'1' { 1 } else { 0 });
        bit_count += 1;
        if bit_count == 8 {
            output.push(byte);
            byte = 0;
            bit_count = 0;
        }
    }
    if bit_count > 0 {
        byte <<= 8 - bit_count;
        output.push(byte);
    }
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

    let count = u32::from_le_bytes([input[4], input[5], input[6], input[7]]) as usize;
    if count != SYMBOL_LIMIT {
        return Err(io::Error::new(
            io::ErrorKind::InvalidData,
            "invalid freq count",
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
    CodecError, Decoder as CodecDecoder, Encoder as CodecEncoder, State, MAX_INPUT_SIZE,
    MAX_OUTPUT_SIZE,
};

fn io_error_to_codec_error(e: io::Error) -> CodecError {
    match e.kind() {
        io::ErrorKind::UnexpectedEof => CodecError::Truncated,
        io::ErrorKind::InvalidData => {
            let msg = e.to_string();
            if msg.contains("truncated") || msg.contains("too short") {
                CodecError::Truncated
            } else if msg.contains("invalid") || msg.contains("bad") {
                CodecError::Corrupt
            } else {
                CodecError::Other(msg)
            }
        }
        _ => CodecError::Other(e.to_string()),
    }
}

pub struct StreamingEncoder {
    state: State,
    buffer: Vec<u8>,
    total_input: usize,
}

impl StreamingEncoder {
    pub fn new() -> Self {
        StreamingEncoder {
            state: State::Ready,
            buffer: Vec::new(),
            total_input: 0,
        }
    }
}

impl Default for StreamingEncoder {
    fn default() -> Self {
        Self::new()
    }
}

impl CodecEncoder for StreamingEncoder {
    fn process(&mut self, input: &[u8], _output: &mut [u8]) -> Result<usize, CodecError> {
        match self.state {
            State::Ready | State::Flushing => {
                if self.total_input > MAX_INPUT_SIZE.saturating_sub(input.len()) {
                    self.state = State::Error;
                    return Err(CodecError::SizeLimit);
                }
                self.state = State::Streaming;
                self.buffer.extend_from_slice(input);
                self.total_input += input.len();
                Ok(0)
            }
            State::Streaming => {
                if self.total_input > MAX_INPUT_SIZE.saturating_sub(input.len()) {
                    self.state = State::Error;
                    return Err(CodecError::SizeLimit);
                }
                self.buffer.extend_from_slice(input);
                self.total_input += input.len();
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
                let encoded = encode(&self.buffer).map_err(io_error_to_codec_error)?;
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
        self.total_input = 0;
    }

    fn state(&self) -> State {
        self.state
    }
}

pub struct StreamingDecoder {
    state: State,
    buffer: Vec<u8>,
    total_input: usize,
}

impl StreamingDecoder {
    pub fn new() -> Self {
        StreamingDecoder {
            state: State::Ready,
            buffer: Vec::new(),
            total_input: 0,
        }
    }
}

impl Default for StreamingDecoder {
    fn default() -> Self {
        Self::new()
    }
}

impl CodecDecoder for StreamingDecoder {
    fn process(&mut self, input: &[u8], _output: &mut [u8]) -> Result<usize, CodecError> {
        match self.state {
            State::Ready | State::Flushing => {
                if self.total_input > MAX_INPUT_SIZE.saturating_sub(input.len()) {
                    self.state = State::Error;
                    return Err(CodecError::SizeLimit);
                }
                self.state = State::Streaming;
                self.buffer.extend_from_slice(input);
                self.total_input += input.len();
                Ok(0)
            }
            State::Streaming => {
                if self.total_input > MAX_INPUT_SIZE.saturating_sub(input.len()) {
                    self.state = State::Error;
                    return Err(CodecError::SizeLimit);
                }
                self.buffer.extend_from_slice(input);
                self.total_input += input.len();
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
                let decoded = decode(&self.buffer).map_err(io_error_to_codec_error)?;
                if decoded.len() > MAX_OUTPUT_SIZE {
                    self.state = State::Error;
                    return Err(CodecError::SizeLimit);
                }
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
        self.total_input = 0;
    }

    fn state(&self) -> State {
        self.state
    }
}
