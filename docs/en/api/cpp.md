# C++ Implementation Reference

All C++ implementations keep single-file algorithm cores and share a thin buffer facade under `algorithms/shared/cpp/include/compresskit/`. There is no stateful streaming layer; `encode_buffer` / `decode_buffer` take a `BufferTransform` function pointer directly.

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

## Shared Buffer Facade

The buffer helpers live in:

- `compresskit/result.hpp` — `StatusCode` enum and `Result<T>` template
- `compresskit/buffer_api.hpp` — `BufferTransform`, `encode_buffer`, `decode_buffer`, file helpers
- `compresskit/algorithms.hpp` — per-algorithm `*_encode_buffer` / `*_decode_buffer` entry points

### Status Codes

| Code | Meaning |
|------|---------|
| `OK` | Success |
| `ERR_CORRUPT` | Encoded data failed structural validation |
| `ERR_SIZE_LIMIT` | Input exceeded 4 GiB or decoded output exceeded 1 GiB |

### Buffer Transform

`BufferTransform` is a function pointer alias mapping an in-memory byte buffer to its transformed form. Each algorithm exposes one encode and one decode transform.

```cpp
#include <vector>
#include "compresskit/algorithms.hpp"
#include "compresskit/buffer_api.hpp"

std::vector<uint8_t> encode(const std::vector<uint8_t>& input) {
    auto result = compresskit::encode_buffer(huffman_encode_buffer, input);
    return result.ok() ? std::move(result.value) : std::vector<uint8_t>{};
}
```
