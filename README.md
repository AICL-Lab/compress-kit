# Encoding Algorithms Collection

[![CI](https://github.com/LessUp/encoding/actions/workflows/ci.yml/badge.svg)](https://github.com/LessUp/encoding/actions/workflows/ci.yml)
[![Deploy Docs](https://github.com/LessUp/encoding/actions/workflows/docs.yml/badge.svg)](https://github.com/LessUp/encoding/actions/workflows/docs.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
![C++](https://img.shields.io/badge/C++-17-blue.svg)
![Go](https://img.shields.io/badge/Go-1.21+-00ADD8.svg)
![Rust](https://img.shields.io/badge/Rust-1.70+-orange.svg)

English | [简体中文](README.zh-CN.md)

> 📖 **Docs**: [https://lessup.github.io/encoding/](https://lessup.github.io/encoding/)

A multi-language implementation of classic compression encoding algorithms for learning and comparison.

## Algorithms

| Algorithm | C++ | Go | Rust | Description |
|-----------|-----|-----|------|-------------|
| **Arithmetic Coding** | ✅ | ✅ | ✅ | Entropy-optimal encoding with fractional bit precision |
| **Huffman Coding** | ✅ | ✅ | ✅ | Classic prefix-free variable-length encoding |
| **LZ77** | ✅ | ✅ | ✅ | Sliding window dictionary compression |
| **Run-Length Encoding** | ✅ | ✅ | ✅ | Simple repeated symbol compression |

## Features

- **Multi-Language**: Same algorithms in C++17, Go, and Rust for cross-language comparison
- **Cross-Platform**: Verified on Linux, macOS, and Windows
- **Benchmarked**: Performance comparison across implementations
- **Well-Tested**: Unit tests, property tests, and cross-language verification
- **Educational**: Clear, documented implementations focused on learning

## Quick Start

### C++ (Arithmetic Coding)

```bash
cd arithmetic/cpp
mkdir build && cd build
cmake .. && make
./arithmetic_codec encode input.txt output.bin
./arithmetic_codec decode output.bin restored.txt
```

### Go

```bash
cd arithmetic/go
go build ./...
go test ./...
```

### Rust

```bash
cd arithmetic/rust
cargo build --release
cargo test
```

## Project Structure

```
encoding/
├── arithmetic/         # Arithmetic coding (C++, Go, Rust)
├── huffman/            # Huffman coding (C++, Go, Rust)
├── lz77/              # LZ77 compression (C++, Go, Rust)
├── rle/               # Run-length encoding (C++, Go, Rust)
├── docs/              # Documentation
└── .github/workflows/ # CI
```

## Documentation

- [Online Docs](https://lessup.github.io/encoding/)
- [Algorithm Deep Dives](https://lessup.github.io/encoding/algorithms)
- [Performance Comparison](https://lessup.github.io/encoding/benchmarks)

## License

MIT License
