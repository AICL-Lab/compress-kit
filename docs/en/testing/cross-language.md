# Cross-Language Testing

One of CompressKit's key features is **identical binary formats** across all language implementations. This enables seamless cross-language encoding and decoding.

## Test Methodology

Our CI pipeline verifies correctness across all 12 possible encoding/decoding combinations:

```
4 algorithms × 3 languages = 12 implementations
12 encoders × 3 decoders = 36 cross-language pairs
```

## Test Matrix

| Encode ↓ / Decode → | C++ | Go | Rust |
|---------------------|-----|-----|------|
| **C++** | ✅ | ✅ | ✅ |
| **Go** | ✅ | ✅ | ✅ |
| **Rust** | ✅ | ✅ | ✅ |

## How It Works

Each algorithm uses a **strictly defined binary format**:

### Huffman Example

```
+----------+-------------------+---------------+
|  Magic   | Frequency Table   | Encoded Data  |
| (4 bytes)| (257 × 4 bytes)   | (variable)    |
+----------+-------------------+---------------+
| "HFMN"   | uint32[257]       | bit stream    |
+----------+-------------------+---------------+
```

All three implementations write this exact structure:
- Same magic bytes (`HFMN`)
- Same endianness (little-endian)
- Same table layout
- Same bit stream encoding

## Running Cross-Language Tests

### Manual Verification

```bash
# Generate test data
dd if=/dev/urandom of=test.bin bs=1M count=1

# C++ encode → Go decode
./algorithms/huffman/cpp/huffman_cpp encode test.bin encoded.huf
./algorithms/huffman/go/huffman_go decode encoded.huf restored.bin
diff test.bin restored.bin  # Should produce no output

# Go encode → Rust decode
./algorithms/huffman/go/huffman_go encode test.bin encoded.huf
./algorithms/huffman/rust/huffman_rust decode encoded.huf restored.bin
diff test.bin restored.bin

# Try any combination:
# C++ ↔ Go, C++ ↔ Rust, Go ↔ Rust
# Works for all 4 algorithms
```

### Automated Testing

```bash
make test
```

This runs the full cross-language test suite across all algorithms.

## File Format Specifications

Detailed format specifications are available in the specs directory:

- [Huffman Format](https://github.com/LessUp/compress-kit/tree/master/specs/rfc)
- [Arithmetic Format](https://github.com/LessUp/compress-kit/tree/master/specs/rfc)
- [Range Coder Format](https://github.com/LessUp/compress-kit/tree/master/specs/rfc)
- [RLE Format](https://github.com/LessUp/compress-kit/tree/master/specs/rfc)

## Why Cross-Language Matters

1. **Data Portability**: Encode data in one environment, decode in another
2. **Incremental Migration**: Gradually move from one language to another
3. **Testing Verification**: Multiple independent implementations catch bugs
4. **Learning**: Compare how the same algorithm is implemented differently

## Known Limitations

| Algorithm | Limitation | Workaround |
|-----------|-----------|------------|
| Range Coder | Decode performance degrades >500KB | Use smaller chunks |

## Reporting Issues

If you find cross-language incompatibilities:

1. Test with `diff` to verify binary mismatch
2. Check file headers with `xxd -l 20 encoded.huf`
3. Report to [GitHub Issues](https://github.com/LessUp/compress-kit/issues) with:
   - Algorithm name
   - Encode language
   - Decode language
   - Input file type
