# Cross-Language Testing Specification

## Overview

This spec defines the cross-language verification strategy for the encoding project, ensuring that implementations in C++17, Go, and Rust produce compatible output.

## Test Strategy

### 1. Correctness Tests

All implementations must pass cross-language encode/decode tests:

```bash
# Encode with language A, decode with language B
./algorithms/huffman/cpp/huffman_cpp encode input.txt output.huf
./algorithms/huffman/go/huffman_go decode output.huf restored.txt
diff input.txt restored.txt  # Must be identical
```

### 2. Test Data Generation

Test data is generated using `tests/gen_testdata.py` and includes:

- Random binary data
- Text files
- Repetitive patterns (for RLE)
- Edge cases (empty files, single byte, etc.)

### 3. Benchmark Tests

Performance benchmarks run across all implementations:

```bash
make bench
```

Results are compared across:
- Compression ratio
- Encode speed
- Decode speed
- Memory usage

## Known Issues

### Range Coder Performance

- **Issue**: Decode hangs for files >500KB
- **Workaround**: CI uses 100KB test file
- **Status**: Under investigation

## Future Test Improvements

- [ ] Fix range coder decode performance issue
- [ ] Add adaptive probability model tests
- [ ] Add LZ77/LZSS algorithm tests
- [ ] Add benchmark visualization
- [ ] Add WebAssembly builds
- [ ] Add Python bindings tests

## Acceptance Criteria

- [ ] All algorithms produce identical output across C++, Go, and Rust
- [ ] Benchmarks run successfully on all implementations
- [ ] No memory leaks or safety issues
- [ ] CI pipeline passes on all platforms
