# CompressKit

<p align="center">
  <img src="docs/public/logo.svg" width="120" alt="CompressKit Logo">
</p>

<p align="center">
  <strong>Classic Lossless Compression Algorithms in C++17, Go, and Rust</strong>
</p>

<p align="center">
  <a href="https://github.com/LessUp/compress-kit/actions/workflows/ci.yml">
    <img src="https://github.com/LessUp/compress-kit/actions/workflows/ci.yml/badge.svg" alt="CI Status">
  </a>
  <a href="https://lessup.github.io/compress-kit/">
    <img src="https://img.shields.io/badge/Docs-Online-blue?logo=readthedocs&logoColor=white" alt="Documentation">
  </a>
  <a href="https://github.com/LessUp/compress-kit/releases">
    <img src="https://img.shields.io/github/v/release/LessUp/compress-kit?include_prereleases&label=Release" alt="Release">
  </a>
  <a href="https://opensource.org/licenses/MIT">
    <img src="https://img.shields.io/badge/License-MIT-green.svg" alt="License">
  </a>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/C++-17-00599C.svg?logo=c%2B%2B" alt="C++17">
  <img src="https://img.shields.io/badge/Go-1.21+-00ADD8.svg?logo=go" alt="Go">
  <img src="https://img.shields.io/badge/Rust-1.70+-DEA584.svg?logo=rust" alt="Rust">
  <img src="https://img.shields.io/badge/Python-3.8+-3776AB.svg?logo=python" alt="Python">
</p>

<p align="center">
  <b>English</b> | <a href="README.zh-CN.md">简体中文</a> | <a href="https://lessup.github.io/compress-kit/">📖 Documentation</a>
</p>

---

## 📋 Prerequisites

| Language | Minimum Version | Installation |
|----------|-----------------|--------------|
| C++ | GCC 9+ / Clang 10+ | `apt install g++` or `brew install gcc` |
| Go | 1.21+ | [golang.org/dl](https://golang.org/dl) |
| Rust | 1.70+ | [rustup.rs](https://rustup.rs) |
| Python | 3.8+ | Only needed for generating test data |

Verify environment:
```bash
g++ --version    # Must support -std=c++17
go version
rustc --version
```

---

## ✨ Features

- 🔤 **Multi-Language** — Identical implementations in C++17, Go 1.21+, and Rust 1.70+
- 🔗 **Cross-Language Compatible** — Encode with one language, decode with another
- 📚 **Educational** — Clean, well-documented code for learning and comparison
- 🧪 **Well-Tested** — Unit tests and cross-language verification in CI
- 📊 **Benchmarked** — Performance comparison across languages

## 🧮 Algorithms

| Algorithm | Compression | Speed | Best For |
|-----------|-------------|-------|----------|
| [**Huffman**](https://lessup.github.io/compress-kit/en/guide/algorithms#huffman-coding) | Medium | Fast | General text/data |
| [**Arithmetic**](https://lessup.github.io/compress-kit/en/guide/algorithms#arithmetic-coding) | Highest | Medium | Maximum compression |
| [**Range Coder**](https://lessup.github.io/compress-kit/en/guide/algorithms#range-coder) | High | Fast | Balanced performance |
| [**RLE**](https://lessup.github.io/compress-kit/en/guide/algorithms#run-length-encoding-rle) | Variable | Fastest | Repetitive data (bitmaps, logs) |

### Algorithm Selection Guide

```
Is your data highly repetitive?
├── Yes → Use RLE (fastest, best for repeated patterns)
└── No →
    Do you need maximum compression?
    ├── Yes → Use Arithmetic Coding (closest to entropy limit)
    └── No →
        Is speed critical?
        ├── Yes → Use Range Coder (fast + good compression)
        └── No → Use Huffman (simple, general purpose)
```

## 🚀 Quick Start

```bash
git clone https://github.com/LessUp/compress-kit.git
cd compress-kit

# 1. Build all implementations
make build

# 2. Generate test data (requires Python 3.8+)
make test-data

# 3. Run tests
make test
```

### Quick Verification

```bash
# Create a test file
echo "Hello, World! Hello, World!" > input.txt

# Encode with C++
./algorithms/huffman/cpp/huffman_cpp encode input.txt output.huf

# Decode with Go
./algorithms/huffman/go/huffman_go decode output.huf restored.txt

# Verify
diff input.txt restored.txt && echo "✓ Cross-language verification passed"
```

### Cross-Language Verification (Alternative)

```bash
# Encode with C++
./algorithms/huffman/cpp/huffman_cpp encode input.txt output.huf

# Decode with Go
./algorithms/huffman/go/huffman_go decode output.huf restored.txt
diff input.txt restored.txt  # No output = identical
```

**C++ ↔ Go ↔ Rust** — all implementations share identical binary formats.

## Project Structure

```
compresskit/
├── algorithms/           # Compression algorithm implementations
│   ├── huffman/         # Prefix-code compression
│   ├── arithmetic/      # Arithmetic coding
│   ├── range/           # Range coder (byte-level arithmetic)
│   └── rle/             # Run-length encoding
│       ├── cpp/         # C++17: single file, zero deps
│       ├── go/          # Go 1.21+: library + CLI
│       ├── rust/        # Rust 1.70+: rustc or cargo
│       └── benchmark/   # Performance scripts
├── docs/                # VitePress site (en + zh)
├── specs/               # Spec-driven development docs
├── tests/               # Test data generation
└── Makefile             # Build entry point
```

## Build & Test

| Command | Description |
|---------|-------------|
| `make build` | Build all implementations |
| `make test` | Run unit tests |
| `make bench` | Run benchmarks |
| `make clean` | Remove build artifacts |

## 💻 Usage

All implementations follow the unified CLI interface:

```bash
<binary> <encode|decode> <input> <output>
```

### CLI Examples

```bash
# Huffman - C++
./algorithms/huffman/cpp/huffman_cpp encode input.txt output.huf
./algorithms/huffman/cpp/huffman_cpp decode output.huf restored.txt

# Huffman - Go
./algorithms/huffman/go/huffman_go encode input.txt output.huf
./algorithms/huffman/go/huffman_go decode output.huf restored.txt

# Huffman - Rust
./algorithms/huffman/rust/huffman_rust encode input.txt output.huf
./algorithms/huffman/rust/huffman_rust decode output.huf restored.txt

# All tools support --help for detailed options
./algorithms/huffman/go/huffman_go --help
```

### Go Library

```go
import "github.com/LessUp/compress-kit/algorithms/huffman/go"

err := huffman.EncodeFile("input.bin", "output.huf")
err = huffman.DecodeFile("output.huf", "decoded.bin")
```

Note: When using as a library, import the package and call functions directly. For standalone CLI usage, build with `go build -o huffman_go ./cmd`.

## 📚 Documentation

| Resource | Link |
|----------|------|
| 📖 Full Documentation | [lessup.github.io/compress-kit](https://lessup.github.io/compress-kit/) |
| 🔧 API Reference | [Go](https://lessup.github.io/compress-kit/en/api/go) · [Rust](https://lessup.github.io/compress-kit/en/api/rust) · [C++](https://lessup.github.io/compress-kit/en/api/cpp) |
| 📊 Benchmark Results | [Performance](https://lessup.github.io/compress-kit/en/benchmarks/results) |
| 🤝 Contributing Guide | [How to Contribute](https://lessup.github.io/compress-kit/en/guide/contributing) |
| 📋 Technical Specs | [specs/](specs/) |

## 🎯 Why This Project

- **📖 Learn** — Compare clean implementations across C++, Go, and Rust
- **✅ Verify** — Cross-language tests guarantee format compatibility
- **📐 SDD** — Built with Spec-Driven Development methodology

## 🤝 Contributing

This project follows **Spec-Driven Development (SDD)**:

1. Read specs first — `/specs/` is the single source of truth
2. Update specs before code — if interfaces change, specs change first
3. Test across languages — verify C++ ↔ Go ↔ Rust compatibility

See [Contributing Guide](https://lessup.github.io/compress-kit/en/guide/contributing) for details.

## ⚠️ Security Notes

- **Maximum input size:** 4 GiB per file
- **Maximum output size:** 1 GiB (protection against decompression bombs)
- All binary formats include integrity validation
- File formats are stable and backward compatible within major versions

## 📜 Changelog

See [CHANGELOG.md](CHANGELOG.md) for version history and migration guides.

## License

[MIT License](LICENSE) · Copyright © 2025-2026 LessUp
