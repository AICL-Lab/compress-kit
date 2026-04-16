# Algorithms Guide

## Overview

| Algorithm | Best For | Compression | Speed | Time | Space |
|-----------|----------|-------------|-------|------|-------|
| **Huffman** | General, text | Medium | Fast | O(n log n) | O(σ) |
| **Arithmetic** | Maximum compression | High | Medium | O(n) | O(σ) |
| **Range Coder** | Balanced performance | High | Fast | O(n) | O(σ) |
| **RLE** | Repetitive data | Variable | Very Fast | O(n) | O(1) |

> σ = alphabet size (256 for byte-level), n = input length

---

## Huffman Coding

Prefix-code based lossless compression. The implementation scans input to build frequency table, constructs Huffman tree, then encodes bit-by-bit.

### Algorithm

1. Count byte frequencies in input
2. Build Huffman binary tree (lower frequency → deeper depth)
3. Generate prefix codes (unambiguous)
4. Encode input using code table, write bit stream

### File Format

| Field | Size | Description |
|-------|------|-------------|
| Magic | 4 bytes | `HFMN` |
| Frequency table | 257 × 4 bytes | Little-endian uint32 |
| Encoded data | Variable | Bit stream |

### Compression Efficiency

- **Theoretical lower bound**: Average code length ≥ entropy H
- **Huffman upper bound**: H ≤ L < H + 1 (at most 1 bit extra per symbol)
- Works best on data with uneven frequency distribution

### Example Usage

```bash
# C++
./huffman_cpp encode input.bin output.huf
./huffman_cpp decode output.huf restored.bin

# Go
./huffman_go encode input.bin output.huf
./huffman_go decode output.huf restored.bin

# Rust
./huffman_rust encode input.bin output.huf
./huffman_rust decode output.huf restored.bin
```

---

## Arithmetic Coding

Uses interval subdivision to represent the entire message's probability, achieving compression closer to the entropy bound.

### Algorithm

1. Initial interval `[0, 1)`
2. Subdivide interval based on symbol probabilities
3. Represent entire message with a number in the final interval
4. Output bits ≈ `-log₂(P(message))`

### Huffman vs Arithmetic Comparison

| Aspect | Huffman | Arithmetic |
|--------|---------|------------|
| Encoding unit | At least 1 bit/symbol | Fractional bits possible |
| Theoretical efficiency | H ≤ L < H + 1 | L ≈ H + ε |
| Implementation complexity | Lower | Higher (precision management) |
| Use case | General | Maximum compression |

### Characteristics

- Theoretically optimal compression (closest to entropy)
- Slower than Huffman encoding/decoding
- Higher implementation complexity (precision handling)

---

## Range Coder

An implementation equivalent to arithmetic coding but typically more efficient in practice. Uses integer interval operations instead of floating point, avoiding patent issues while achieving better real-world performance.

### Arithmetic Coding vs Range Coder

| Aspect | Arithmetic | Range Coder |
|--------|------------|-------------|
| Output unit | Bits | Bytes |
| I/O efficiency | Lower | Higher |
| Compression rate | Nearly identical | Nearly identical |
| Patent status | Historical patents | No patent restrictions |
| Engineering preference | Academic | Production |

### API Usage

**Go**:
```go
import "encoding/range/go/rangecoder"

encoded, err := rangecoder.Encode(data)
decoded, err := rangecoder.Decode(encoded)
```

**Rust**:
```rust
use rangecoder;

let encoded = rangecoder::encode(input)?;
let decoded = rangecoder::decode(&encoded)?;
```

### CLI Usage

```bash
# C++
./rangecoder_cpp encode input.bin output.rcnc
./rangecoder_cpp decode output.rcnc restored.bin

# Go
./rangecoder_go encode input.bin output.rcnc
./rangecoder_go decode output.rcnc restored.bin

# Rust
cargo run --bin rangecoder -- encode input.bin output.rcnc
cargo run --bin rangecoder -- decode output.rcnc restored.bin
```

> **Known Issue**: Range coder decoder may be slow for files >500KB. Use smaller test files for cross-language verification.

---

## Run-Length Encoding (RLE)

Simplest compression algorithm, ideal for data with many consecutive repeated bytes.

### File Format

Repeated `(count, value)` pairs:

| Field | Size | Description |
|-------|------|-------------|
| count | 4 bytes | Little-endian unsigned int |
| value | 1 byte | Original byte |

### Characteristics

- Simplest implementation, extremely fast encoding/decoding
- Excellent for repetitive data (bitmaps, logs with repeated lines)
- May expand random data (worst case: 5× expansion)
- Often used as preprocessing for other algorithms (BWT + MTF + RLE + Arithmetic)

### Example Usage

```bash
# C++
./rle_cpp encode input.bin output.rle
./rle_cpp decode output.rle restored.bin

# Go
./rle_go encode input.bin output.rle
./rle_go decode output.rle restored.bin

# Rust
./rle_rust encode input.bin output.rle
./rle_rust decode output.rle restored.bin
```

---

## Choosing an Algorithm

| Data Type | Recommended Algorithm |
|-----------|----------------------|
| Text, mixed content | Huffman or Arithmetic |
| Maximum compression needed | Arithmetic |
| Performance-critical | Range Coder |
| Highly repetitive | RLE |
| Unknown data type | Huffman (good balance) |
