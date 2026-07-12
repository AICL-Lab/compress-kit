# Getting Started

This guide will help you set up the development environment, build the implementations, and run tests.

## Prerequisites

### Required Tools

| Tool | Minimum Version | Purpose |
|------|-----------------|---------|
| g++ or clang++ | 9+ / 10+ | C++17 compilation |
| CMake | 3.16+ | Build system |
| Python | 3.8+ | Test orchestration and benchmark scripts |
| Make | Any | Build automation |

### Optional Tools

| Tool | Purpose |
|------|---------|
| Node.js 20.19+ | Documentation site |
| clang-format | C++ code formatting |

### Installation

::: code-group

```bash [Ubuntu/Debian]
sudo apt update
sudo apt install g++ cmake python3 make
```

```bash [macOS (Homebrew)]
brew install gcc cmake python3 make
```

```bash [Windows (Chocolatey)]
choco install mingw cmake python3 make
```

:::

## Clone and Build

### Clone the Repository

```bash
git clone https://github.com/LessUp/compress-kit.git
cd compress-kit
```

### Use the Dev Container (Recommended)

The repository includes a checked-in `.devcontainer/` for the supported
C++/Node toolchain. Open the folder in VS Code or Codespaces and choose
**Reopen in Container** to get the same baseline used by the project docs.

The container runs the reproducible install path on creation:

```bash
npm ci
(cd docs && npm ci)
```

### Reproducible Local Docs Tooling

If you are staying on your host machine, run the same lockfile-backed installs
before building the docs:

```bash
npm ci
(cd docs && npm ci)
```

### Build All Implementations

```bash
make build
```

`make build` wraps the CMake invocation (`cmake -S . -B build && cmake --build build`) and compiles all algorithm implementations into `build/`.

### Build Specific Algorithm

```bash
make build-huffman      # Huffman coding
make build-arithmetic   # Arithmetic coding
make build-range        # Range coder
make build-rle          # Run-length encoding
```

### Manual Compilation

If you prefer to invoke CMake directly:

```bash
cmake -S . -B build && cmake --build build
./build/huffman_cpp encode input.bin output.huf
./build/huffman_cpp decode output.huf restored.bin
```

## Running Tests

### Run All Tests

```bash
make test
```

This runs the shared lifecycle tests and the algorithm-specific test
suites.

## Running Benchmarks

### Run All Benchmarks

```bash
make bench
```

This generates test data, runs the verified benchmark matrix, writes raw reports
to `reports/`, and refreshes `docs/.vitepress/data/benchmarks.json` for the
interactive docs chart.

### Benchmark Artifacts

| Artifact | Purpose |
|----------|---------|
| `reports/*.txt` | Raw per-run benchmark logs |
| `docs/.vitepress/data/benchmarks.json` | Generated snapshot consumed by the docs benchmark chart |
| `/en/benchmarks/results` | Reader-friendly view of the generated JSON snapshot |

Treat the generated JSON file and chart as the source of truth for benchmark
numbers in this repository snapshot.

## Makefile Command Reference

| Command | Description |
|---------|-------------|
| `make build` | Build all algorithm implementations (wraps CMake) |
| `make build-huffman` | Build only Huffman implementations |
| `make build-arithmetic` | Build only Arithmetic implementations |
| `make build-range` | Build only Range Coder implementations |
| `make build-rle` | Build only RLE implementations |
| `make test` | Run unit and lifecycle tests |
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

### Range Coder Performance

The Range Coder decoder has a known performance issue for files larger than
500 KiB. Benchmarks therefore use the smaller `small_dictionary_like` dataset
for Range results. For manual checks, stay on similarly small samples:

```bash
# Create a smaller test file (100 KiB)
dd if=tests/data/random_10MiB.bin of=small.bin bs=1024 count=100

# Test with smaller file
./build/rangecoder_cpp encode small.bin small.enc
./build/rangecoder_cpp decode small.enc small.dec
```

## Next Steps

- Learn about the [algorithms](/en/guide/algorithms) and their differences
- Read the [architecture overview](/en/guide/architecture)
- Check the [CHANGELOG](https://github.com/LessUp/compress-kit/blob/master/CHANGELOG.md) for recent updates
