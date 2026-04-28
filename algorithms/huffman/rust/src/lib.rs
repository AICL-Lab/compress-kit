use compresskit_codec::codec::{CodecError, Decoder, Encoder, State};
use std::cmp::Ordering;
use std::collections::BinaryHeap;
use std::io::{self, Read, Write};

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
        let node = Box::new(Node {
            symbol: s as u32,
            freq: f as u64,
            left: None,
            right: None,
        });
        heap.push(HeapItem {
            freq: node.freq,
            symbol: node.symbol,
            node,
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
        let item = heap.pop().expect("heap should have one element");
        let only = item.node;
        let parent = Box::new(Node {
            symbol: only.symbol,
            freq: only.freq,
            left: Some(only),
            right: None,
        });
        heap.push(HeapItem {
            freq: parent.freq,
            symbol: parent.symbol,
            node: parent,
        });
    }
    while heap.len() > 1 {
        let a = heap.pop().expect("heap should have element").node;
        let b = heap.pop().expect("heap should have element").node;
        let freq_sum = a.freq + b.freq;
        let parent = Box::new(Node {
            symbol: a.symbol.min(b.symbol),
            freq: freq_sum,
            left: Some(a),
            right: Some(b),
        });
        let symbol = parent.symbol;
        heap.push(HeapItem {
            freq: parent.freq,
            symbol,
            node: parent,
        });
    }
    heap.pop().expect("heap should have final element").node
}

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

    fn read_bit(&mut self) -> Option<u8> {
        if self.bits_remaining == 0 {
            if self.pos >= self.data.len() {
                return None;
            }
            self.current_byte = self.data[self.pos];
            self.pos += 1;
            self.bits_remaining = 8;
        }
        self.bits_remaining -= 1;
        Some((self.current_byte >> self.bits_remaining) & 1)
    }
}

fn build_codes(node: &Node, codes: &mut [String], prefix: &mut String) {
    if is_leaf(node) {
        if prefix.is_empty() {
            codes[node.symbol as usize] = "0".to_string();
        } else {
            codes[node.symbol as usize] = prefix.clone();
        }
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

fn build_frequencies(data: &[u8]) -> Vec<u32> {
    let mut freq = vec![0u32; SYMBOL_LIMIT];
    for &b in data {
        freq[b as usize] += 1;
    }
    freq[EOF_SYMBOL as usize] = 1;
    freq
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

// In-memory encode
pub fn encode_memory(input: &[u8]) -> Result<Vec<u8>, CodecError> {
    let freq = build_frequencies(input);
    let root = build_tree(&freq);
    let mut codes = vec![String::new(); SYMBOL_LIMIT];
    let mut prefix = String::new();
    build_codes(&root, &mut codes, &mut prefix);

    let mut out = Vec::new();
    out.extend_from_slice(b"HFMN");
    write_u32_le(&mut out, freq.len() as u32);
    for &f in &freq {
        write_u32_le(&mut out, f);
    }

    let mut bit_writer = BitWriter::new();
    for &b in input {
        let code = &codes[b as usize];
        for ch in code.as_bytes() {
            let bit = if *ch == b'1' { 1 } else { 0 };
            bit_writer.write_bit(bit);
        }
    }
    let eof_code = &codes[EOF_SYMBOL as usize];
    for ch in eof_code.as_bytes() {
        let bit = if *ch == b'1' { 1 } else { 0 };
        bit_writer.write_bit(bit);
    }

    out.extend_from_slice(&bit_writer.into_bytes());
    Ok(out)
}

// In-memory decode
pub fn decode_memory(input: &[u8]) -> Result<Vec<u8>, CodecError> {
    if input.len() < 4 {
        return Err(CodecError::Truncated);
    }
    if &input[0..4] != b"HFMN" {
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

    let root = build_tree(&freq);
    let mut bit_reader = BitReader::new(&input[pos..]);
    let mut output = Vec::new();
    let mut node_ref: &Node = &root;
    let mut saw_eof = false;

    loop {
        let bit = bit_reader.read_bit().ok_or(CodecError::Truncated)?;
        if bit == 0 {
            match node_ref.left {
                Some(ref left) => {
                    node_ref = left;
                }
                None => return Err(CodecError::Corrupt),
            }
        } else {
            match node_ref.right {
                Some(ref right) => {
                    node_ref = right;
                }
                None => return Err(CodecError::Corrupt),
            }
        }

        if is_leaf(node_ref) {
            if node_ref.symbol == EOF_SYMBOL {
                saw_eof = true;
                break;
            }
            output.push(node_ref.symbol as u8);
            node_ref = &root;
        }
    }

    if !saw_eof {
        return Err(CodecError::Truncated);
    }

    Ok(output)
}

// Streaming Encoder Adapter
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

// Streaming Decoder Adapter
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
