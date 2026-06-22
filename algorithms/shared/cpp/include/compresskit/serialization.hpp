#pragma once

// In-memory little-endian serialization helpers shared by the buffer-based
// algorithm implementations. Stream-based variants live in frequency_table.

#include <cstdint>
#include <stdexcept>
#include <string>
#include <vector>

#include "compresskit/constants.hpp"

namespace compresskit {

// Appends a 32-bit unsigned integer in little-endian to `out`.
inline void write_u32_le(std::vector<uint8_t>& out, uint32_t value) {
    for (std::size_t i = 0; i < U32_SIZE; ++i) {
        out.push_back(static_cast<uint8_t>((value >> (i * BITS_PER_BYTE)) & 0xFFu));
    }
}

// Appends a 4-byte magic prefix.
inline void write_magic(std::vector<uint8_t>& out, const char* magic) {
    for (std::size_t i = 0; i < MAGIC_SIZE; ++i) {
        out.push_back(static_cast<uint8_t>(magic[i]));
    }
}

// Writes a static-model file header: magic + frequency table
// (count: u32 LE, then count x u32 LE entries).
inline void write_frequency_header(std::vector<uint8_t>& out, const char* magic,
                                   const std::vector<uint32_t>& freq) {
    write_magic(out, magic);
    write_u32_le(out, static_cast<uint32_t>(freq.size()));
    for (uint32_t v : freq) {
        write_u32_le(out, v);
    }
}

// Reads a frequency table following the magic prefix.
// `pos` starts after the magic; `algo_name` prefixes error messages.
inline std::vector<uint32_t> read_frequency_header(const std::vector<uint8_t>& input,
                                                   std::size_t& pos, const char* algo_name) {
    if (pos + U32_SIZE > input.size()) {
        throw std::runtime_error(std::string(algo_name) + ": truncated frequency count");
    }
    uint32_t count = 0;
    for (std::size_t i = 0; i < U32_SIZE; ++i) {
        count |= static_cast<uint32_t>(input[pos + i]) << (i * BITS_PER_BYTE);
    }
    pos += U32_SIZE;
    if (count != SYMBOL_LIMIT) {
        throw std::runtime_error(std::string(algo_name) + ": bad frequency count");
    }
    std::vector<uint32_t> freq(count, 0);
    for (uint32_t i = 0; i < count; ++i) {
        if (pos + U32_SIZE > input.size()) {
            throw std::runtime_error(std::string(algo_name) + ": truncated frequency entry");
        }
        uint32_t v = 0;
        for (std::size_t b = 0; b < U32_SIZE; ++b) {
            v |= static_cast<uint32_t>(input[pos + b]) << (b * BITS_PER_BYTE);
        }
        freq[i] = v;
        pos += U32_SIZE;
    }
    return freq;
}

}  // namespace compresskit
