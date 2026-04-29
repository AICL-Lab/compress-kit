# Getting Started

This guide will help you set up the development environment, build the implementations, and run tests.

## Prerequisites

### Required Tools

| Tool | Minimum Version | Purpose |
|------|-----------------|---------|
| g++ or clang++ | 9+ / 10+ | C++17 compilation |
| Go | 1.21+ | Go implementation |
| Rust (cargo) | 1.70+ | Rust implementation |
| Python | 3.8+ | Test orchestration and benchmark scripts |
| Make | Any | Build automation |

### Optional Tools

| Tool | Purpose |
|------|---------|
| Node.js 20.19+ | Documentation site and OpenSpec tooling |
| clang-format | C++ code formatting |

### Installation

::: code-group

```bash [Ubuntu/Debian]
sudo apt update
sudo apt install g++ golang rustc python3 make
```

```bash [macOS (Homebrew)]
brew install gcc go rust python3 make
```

```bash [Windows (Chocolatey)]
choco install mingw golang rust python3 make
```

:::

## Clone and Build

### Clone the Repository

```bash
git clone https://github.com/LessUp/compress-kit.git
cd compress-kit
```

### Build All Implementations

```bash
make build
```

This will compile all algorithm implementations in all three languages.

### Build Specific Algorithm

```bash
make build-huffman      # Huffman coding (C++, Go, Rust)
make build-arithmetic   # Arithmetic coding
make build-range        # Range coder
make build-rle          # Run-length encoding
```

### Manual Compilation

If you prefer to compile manually:

::: code-group

```bash [C++]
cd algorithms/huffman/cpp
g++ -std=c++17 -O2 main.cpp -o huffman_cpp
./huffman_cpp encode input.bin output.huf
./huffman_cpp decode output.huf restored.bin
```

```bash [Go]
cd algorithms/huffman/go
go build -o huffman_go ./cmd
./huffman_go encode input.bin output.huf
./huffman_go decode output.huf restored.bin
```

```bash [Rust]
cd algorithms/huffman/rust
cargo build --bin huffman_rust --release
./target/release/huffman_rust encode input.bin output.huf
./target/release/huffman_rust decode output.huf restored.bin
```

:::

## Cross-Language Verification

One of the key features of this project is that all implementations use identical file formats, enabling cross-language verification:

```bash
# Encode with C++
./algorithms/huffman/cpp/huffman_cpp encode input.bin encoded.huf

# Decode with Go
./algorithms/huffman/go/huffman_go decode encoded.huf decoded.bin

# Verify correctness
diff input.bin decoded.bin  # No output = identical
```

Any combination works: **C++ ↔ Go ↔ Rust**

## Running Tests

### Run All Tests

```bash
make test
```

This runs shared streaming-layer tests, Go/Rust unit tests, and the executable
cross-language conformance matrix.

### Run Individual Algorithm Tests

```bash
# Go tests
cd algorithms/huffman/go && go test ./...

# Rust tests
cd algorithms/huffman/rust && cargo test
```

## Running Benchmarks

### Run All Benchmarks

```bash
make bench
```

This generates test data and runs cross-language benchmarks.

### Benchmark Output

Reports are saved to the `reports/` directory. Example output:

```
Algorithm: Huffman
Language: C++
Input: 10 MiB random data
Encode: 245 ms (40.8 MiB/s)
Decode: 198 ms (50.5 MiB/s)
Compression ratio: 1.23
```

## Makefile Command Reference

| Command | Description |
|---------|-------------|
| `make build` | Build all algorithm implementations |
| `make build-huffman` | Build only Huffman implementations |
| `make build-arithmetic` | Build only Arithmetic implementations |
| `make build-range` | Build only Range Coder implementations |
| `make build-rle` | Build only RLE implementations |
| `make test` | Run unit, streaming, and conformance tests |
| `make test-conformance` | Run the cross-language decode matrix |
| `make bench` | Generate test data and run benchmarks |
| `make test-data` | Generate test data only |
| `make clean` | Remove all build artifacts and reports |

## Troubleshooting

### C++ Compilation Errors

```bash
# Check compiler version
g++ --version  # Should be 9+

# Try with clang if g++ fails
clang++ -std=c++17 -O2 main.cpp -o huffman_cpp
```

### Go Module Issues

```bash
# Ensure Go workspace includes all modules
go work use ./algorithms/shared/go
go work use ./algorithms/huffman/go
go work use ./algorithms/arithmetic/go
go work use ./algorithms/range/go
go work use ./algorithms/rle/go
```

### Rust Build Errors

```bash
# Update Rust toolchain
rustup update stable

# Check version
rustc --version  # Should be 1.70+
```

### Range Coder Performance

The Range Coder decoder has a known performance issue for files larger than 500KB. Use smaller test files:

```bash
# Create smaller test file (100KB)
dd if=tests/data/random_10MiB.bin of=/tmp/small.bin bs=1024 count=100

# Test with smaller file
./algorithms/range/cpp/rangecoder_cpp encode /tmp/small.bin /tmp/small.enc
./algorithms/range/cpp/rangecoder_cpp decode /tmp/small.enc /tmp/small.dec
```

## Next Steps

- Learn about the [algorithms](/en/guide/algorithms) and their differences
- Explore the [project structure](/en/guide/project-structure)
- Check the [CHANGELOG](https://github.com/LessUp/compress-kit/blob/master/CHANGELOG.md) for recent updates
