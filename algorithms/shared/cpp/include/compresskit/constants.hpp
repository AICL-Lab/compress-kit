#pragma once

#include <cstddef>
#include <cstdint>

namespace compresskit {

// Symbol table size: 256 byte values + 1 EOF marker.
constexpr uint32_t SYMBOL_LIMIT = 257;
constexpr uint32_t EOF_SYMBOL = SYMBOL_LIMIT - 1;

// Security limits.
constexpr uint64_t MAX_INPUT_SIZE = 4ULL * 1024 * 1024 * 1024;   // 4 GiB
constexpr uint64_t MAX_OUTPUT_SIZE = 1ULL * 1024 * 1024 * 1024;  // 1 GiB

// Binary format element sizes.
constexpr std::size_t MAGIC_SIZE = 4;  // algorithm magic prefix length
constexpr std::size_t U32_SIZE = 4;    // little-endian uint32 serialisation
constexpr int BITS_PER_BYTE = 8;
constexpr uint32_t BYTE_VALUES = 1u << BITS_PER_BYTE;  // 256
constexpr std::size_t RLE_PAIR_SIZE = U32_SIZE + 1;    // count + value byte

// I/O tuneables.
constexpr std::size_t INITIAL_ENCODE_OVERHEAD = 2048;

// Algorithm magic numbers (binary format identifiers, little-endian agnostic).
// These MUST NOT change: they are the basis of binary compatibility.
constexpr char HUFFMAN_MAGIC[4] = {'H', 'F', 'M', 'N'};
constexpr char ARITHMETIC_MAGIC[4] = {'A', 'E', 'N', 'C'};
constexpr char RANGE_MAGIC[4] = {'R', 'C', 'N', 'C'};
constexpr char RLE_MAGIC[4] = {'R', 'L', 'E', '\x00'};

}  // namespace compresskit
