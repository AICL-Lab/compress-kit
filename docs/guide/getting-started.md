# Getting Started

## Environment Requirements

| Tool | Minimum Version | Purpose |
|------|----------------|---------|
| g++ / clang++ | 9+ / 10+ | C++17 compilation |
| Go | 1.21+ | Go implementation |
| Rust (cargo) | 1.70+ | Rust implementation |
| Python | 3.8+ | Benchmark scripts |
| Node.js | 18+ | Documentation site (optional) |
| Make | Any | Build automation |

## Clone and Build

```bash
git clone https://github.com/LessUp/encoding.git
cd encoding
```

### Build All Implementations

```bash
make build
```

### Build Specific Algorithm

```bash
make build-huffman      # Huffman (C++, Go, Rust)
make build-arithmetic   # Arithmetic coding
make build-range        # Range coder
make build-rle          # Run-length encoding
```

### Manual Compilation Examples

::: code-group
```bash [C++]
cd huffman/cpp
g++ -std=c++17 -O2 main.cpp -o huffman_cpp
./huffman_cpp encode input.bin output.huf
./huffman_cpp decode output.huf restored.bin
```

```bash [Go]
cd huffman/go
go build -o huffman_go .
./huffman_go encode input.bin output.huf
```

```bash [Rust]
cd huffman/rust
rustc -O main.rs -o huffman_rust
./huffman_rust encode input.bin output.huf
```
:::

## Cross-Language Verification

All implementations use identical file formats, enabling cross-verification:

```bash
# C++ encode, Go decode
./huffman/cpp/huffman_cpp encode input.bin encoded.huf
./huffman/go/huffman_go decode encoded.huf decoded.bin
diff input.bin decoded.bin  # No output = identical
```

Supports any combination: C++ ↔ Go ↔ Rust

## Running Tests

```bash
make test          # Run all Go + Rust unit tests
```

### Individual Algorithm Tests

```bash
cd huffman/go && go test ./... && cd ../..
cd huffman/rust && rustc --test main.rs -o test && ./test && cd ../..
```

## Running Benchmarks

```bash
make bench         # Generate test data + run all benchmarks
```

Reports output to `reports/` directory.

### Benchmark Output Example

```
Algorithm: Huffman
Language: C++
Input: 10 MiB random data
Encode: 245 ms (40.8 MiB/s)
Decode: 198 ms (50.5 MiB/s)
Compression ratio: 1.23
```

## Documentation Site

### Local Preview

```bash
npm install
npm run docs:dev
```

Opens at `http://localhost:5173/encoding/`

### Build for Production

```bash
npm run docs:build
```

Output in `docs/.vitepress/dist/`

## Makefile Commands Reference

| Command | Description |
|---------|-------------|
| `make build` | Build all algorithm implementations |
| `make test` | Run all Go and Rust unit tests |
| `make bench` | Generate test data and run benchmarks |
| `make test-data` | Generate test data only |
| `make clean` | Remove all build artifacts and reports |

## Troubleshooting

### C++ Compilation Errors

```bash
# Check compiler version
g++ --version  # Need 9+

# Try with clang
clang++ -std=c++17 -O2 main.cpp -o huffman_cpp
```

### Go Module Issues

```bash
# Ensure Go workspace is active
go work use ./huffman/go
go work use ./arithmetic/go
go work use ./range/go
go work use ./rle/go
```

### Rust Build Errors

```bash
# Update toolchain
rustup update stable

# Check version
rustc --version  # Need 1.70+
```

### Range Coder Slow Decode

The range coder has a known performance issue for files >500KB. Use smaller test files:

```bash
# Create smaller test file
dd if=tests/data/random_10MiB.bin of=/tmp/small.bin bs=1024 count=100

# Test with smaller file
./range/cpp/rangecoder_cpp encode /tmp/small.bin /tmp/small.enc
./range/cpp/rangecoder_cpp decode /tmp/small.enc /tmp/small.dec
```
