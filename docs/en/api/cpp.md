# C++ Implementation Reference

See also: [Streaming API](/en/api/streaming)

All C++ implementations keep single-file algorithm cores, but now depend on a shared streaming/buffer facade under `algorithms/shared/cpp/include/compresskit/`.

## Compilation

The project is built with CMake. `make build` is a thin wrapper around the CMake invocation:

```bash
cmake -S . -B build && cmake --build build
```

Binaries are produced under `build/` (e.g. `build/huffman_cpp`).

### Recommended Flags

The CMake build enables C++17 with optimization and warnings. For debug builds, AddressSanitizer and UBSan can be enabled through the CMake configuration.

| Flag | Purpose |
|------|---------|
| `-std=c++17` | Enable C++17 features |
| `-O2` | Optimization level |
| `-Wall -Wextra` | Warnings |
| `-fsanitize=address` | AddressSanitizer (debug builds) |
| `-fsanitize=undefined` | UBSan (debug builds) |

## Huffman (`algorithms/huffman/cpp/main.cpp`)

### Usage

```bash
./build/huffman_cpp encode input.bin output.huf
./build/huffman_cpp decode output.huf decoded.bin
```

### Internal Structure

- `BitWriter` / `BitReader` — bit-level I/O classes
- `Node` — Huffman tree node with custom deleter
- `compress_file()` / `decompress_file()` — main encode/decode functions

### File Format

| Offset | Size | Field |
|--------|------|-------|
| 0 | 4B | Magic: `HFMN` |
| 4 | 4B | Frequency table size (always 257) |
| 8 | 1028B | Frequency table (257 × uint32 LE) |
| 1036+ | Variable | Encoded bit stream |

---

## Arithmetic (`algorithms/arithmetic/cpp/main.cpp`)

### Usage

```bash
./build/arithmetic_cpp encode input.bin output.aenc
./build/arithmetic_cpp decode output.aenc decoded.bin
```

### Key Classes

- `ArithmeticEncoder` — state machine with `low`, `high`, `pendingBits`
- `ArithmeticDecoder` — decoder with `code` initialization

---

## Range Coder (`algorithms/range/cpp/main.cpp`)

### Usage

```bash
./build/rangecoder_cpp encode input.bin output.rcnc
./build/rangecoder_cpp decode output.rcnc decoded.bin
```

### File Format

| Offset | Size | Field |
|--------|------|-------|
| 0 | 4B | Magic: `RCNC` |
| 4 | 4B | Frequency table size |
| 8 | Variable | Frequency table |
| ... | Variable | Byte stream (renormalized intervals) |

---

## RLE (`algorithms/rle/cpp/main.cpp`)

### Usage

```bash
./build/rle_cpp encode input.bin output.rle
./build/rle_cpp decode output.rle decoded.bin
```

### File Format

| Offset | Size | Field |
|--------|------|-------|
| 0 | 4B | Magic: `RLE\x00` |
| 4+ | Variable | Repeated `(count: uint32 LE, value: byte)` pairs |

---

## Common Patterns

| Pattern | Description |
|---------|-------------|
| Single-file core | Each algorithm core in one `main.cpp` |
| Shared dependencies | Uses common code from `algorithms/shared/cpp/` |
| Error handling | `fprintf(stderr, ...)` + `exit(1)` |
| Memory management | `std::unique_ptr` with custom deleters |

## Shared Streaming Facade

The streaming and buffer helpers live in:

- `compresskit/result.hpp`
- `compresskit/encoder.hpp`
- `compresskit/buffer_api.hpp`
- `compresskit/algorithms.hpp`

These headers provide a common `Encoder` / `Decoder` lifecycle and algorithm factories such as `compresskit::make_huffman_encoder()`.
