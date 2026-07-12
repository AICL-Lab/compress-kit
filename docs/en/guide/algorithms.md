# Algorithms Guide

This guide explains the four compression algorithms implemented in this project, their use cases, and key differences.

## Quick Comparison

| Algorithm | Best For | Compression | Speed | Time Complexity | Space Complexity |
|-----------|----------|-------------|-------|-----------------|------------------|
| **Huffman** | General, text | Medium | Fast | O(n log σ) | O(σ) |
| **Arithmetic** | Maximum compression | High | Medium | O(n) | O(σ) |
| **Range Coder** | Balanced performance | High | Fast | O(n) | O(σ) |
| **RLE** | Repetitive data | Variable | Very Fast | O(n) | O(1) |

> **Legend**: σ = alphabet size (256 for byte-level), n = input length

---

## Huffman Coding

Prefix-code based lossless compression algorithm. Builds an optimal prefix tree based on symbol frequencies.

### How It Works

1. **Frequency Analysis**: Count byte frequencies in input
2. **Tree Construction**: Build binary tree where lower frequency symbols have deeper paths
3. **Code Generation**: Generate prefix codes (unambiguous bit sequences)
4. **Encoding**: Replace each byte with its code and write bit stream

### File Format

| Field | Size | Description |
|-------|------|-------------|
| Magic | 4 bytes | `HFMN` (0x48 0x46 0x4D 0x4E) |
| Frequency Table | 257 × 4 bytes | Little-endian uint32 array |
| Encoded Data | Variable | Bit stream |

### Compression Efficiency

- **Theoretical lower bound**: Average code length ≥ entropy H
- **Huffman upper bound**: H ≤ L < H + 1 (at most 1 bit extra per symbol)
- Most effective on data with uneven frequency distribution

### Usage Example

```bash
./build/huffman_cpp encode input.bin output.huf
./build/huffman_cpp decode output.huf restored.bin
```

---

## Arithmetic Coding

Represents the entire message as a single number in the interval [0, 1), achieving compression closer to the entropy bound than Huffman coding.

### How It Works

1. **Initialize**: Start with interval [0, 1)
2. **Subdivide**: For each symbol, subdivide current interval based on its probability
3. **Select**: The final interval contains infinite numbers; select any to represent the message
4. **Output**: Number of bits ≈ -log₂(P(message))

### Huffman vs Arithmetic Comparison

| Aspect | Huffman | Arithmetic |
|--------|---------|------------|
| Encoding unit | At least 1 bit per symbol | Fractional bits possible |
| Theoretical efficiency | H ≤ L < H + 1 | L ≈ H + ε (closer to entropy) |
| Implementation | Simpler | More complex (precision management) |
| Speed | Faster | Slower |
| Use case | General purpose | Maximum compression |

### Characteristics

- **Optimal compression**: Theoretically closest to entropy limit
- **Slower**: Encoding/decoding overhead higher than Huffman
- **Complexity**: Requires careful precision management

---

## Range Coder

An integer-based implementation equivalent to arithmetic coding but typically more efficient in practice. Uses integer interval operations instead of floating point.

### Arithmetic vs Range Coder

| Aspect | Arithmetic | Range Coder |
|--------|------------|-------------|
| Output unit | Bits | Bytes |
| I/O efficiency | Lower | Higher |
| Compression rate | Nearly identical | Nearly identical |
| Patent status | Had historical patents | No restrictions |
| Engineering use | Academic study | Production systems |

### CLI Usage

```bash
./build/rangecoder_cpp encode input.bin output.rcnc
./build/rangecoder_cpp decode output.rcnc restored.bin
```

::: warning Performance Note
The Range Coder decoder has a known performance issue for files >500 KiB. Use smaller test files for verification.
:::

---

## Run-Length Encoding (RLE)

The simplest compression algorithm, ideal for data with consecutive repeated bytes.

### File Format

| Field | Size | Description |
|-------|------|-------------|
| Magic | 4 bytes | `RLE\x00` (0x52 0x4C 0x45 0x00) |
| Count | 4 bytes | Little-endian unsigned int (run length) |
| Value | 1 byte | The repeated byte |

Each run is stored as a `(count, value)` pair after the magic header.

### Characteristics

- **Simplicity**: Easiest to understand and implement
- **Speed**: Extremely fast encoding and decoding
- **Best for**: Repetitive data (bitmaps, logs with repeated lines)
- **Worst case**: Can expand data up to 5× for random input
- **Common use**: Preprocessing for other algorithms (e.g., BWT + MTF + RLE + Arithmetic)

### Usage Example

```bash
./build/rle_cpp encode input.bin output.rle
./build/rle_cpp decode output.rle restored.bin
```

---

## Algorithm Selection Guide

| Data Type | Recommended Algorithm | Reason |
|-----------|----------------------|--------|
| Text files | Huffman or Range Coder | Natural language has uneven frequency distribution |
| Maximum compression needed | Arithmetic | Closest to theoretical limit |
| Performance-critical | Range Coder | Fast with good compression |
| Highly repetitive (bitmaps, logs) | RLE | Simple patterns compress extremely well |
| Unknown/Mixed content | Range Coder | Best balance of speed and compression |

### Decision Flowchart

```
Is the data highly repetitive?
├── Yes → Use RLE
└── No →
    Is maximum compression required?
    ├── Yes → Use Arithmetic
    └── No →
        Is speed critical?
        ├── Yes → Use Range Coder
        └── No → Use Huffman
```

---

## Further Reading

- [Architecture](/en/guide/architecture) - Layered design and shared utilities
- [Getting Started](/en/guide/getting-started) - Build and test instructions
- [GitHub Repository](https://github.com/LessUp/compress-kit) - Source code and issues
