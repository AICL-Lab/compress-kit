#include <algorithm>
#include <cstdint>
#include <stdexcept>
#include <vector>

#include "compresskit/buffer_api.hpp"
#include "compresskit/constants.hpp"
#include "compresskit/serialization.hpp"

// Run-Length encoding.
// Format:
// - Magic: 4 bytes "RLE\x00"
// - (count: uint32 LE, value: byte) pairs until EOF; count must be > 0.

namespace compresskit {

std::vector<uint8_t> rle_encode_buffer(const std::vector<uint8_t>& input) {
    std::vector<uint8_t> out;
    out.reserve(input.size() / 8 + compresskit::MAGIC_SIZE + compresskit::RLE_PAIR_SIZE);
    compresskit::write_magic(out, compresskit::RLE_MAGIC);

    if (input.empty()) {
        return out;
    }

    uint8_t current = input[0];
    uint32_t count = 1;
    for (std::size_t i = 1; i < input.size(); ++i) {
        if (input[i] == current && count < UINT32_MAX) {
            ++count;
        } else {
            compresskit::write_u32_le(out, count);
            out.push_back(current);
            current = input[i];
            count = 1;
        }
    }
    compresskit::write_u32_le(out, count);
    out.push_back(current);
    return out;
}

std::vector<uint8_t> rle_decode_buffer(const std::vector<uint8_t>& input) {
    std::size_t pos = 0;
    compresskit::verify_magic(input, pos, compresskit::RLE_MAGIC, "RLE");

    std::vector<uint8_t> out;
    while (pos < input.size()) {
        if (pos + compresskit::RLE_PAIR_SIZE > input.size()) {
            throw std::runtime_error("RLE: truncated count+value pair");
        }
        uint32_t count = 0;
        for (std::size_t b = 0; b < compresskit::U32_SIZE; ++b) {
            count |= static_cast<uint32_t>(input[pos + b]) << (b * compresskit::BITS_PER_BYTE);
        }
        uint8_t value = input[pos + compresskit::U32_SIZE];
        pos += compresskit::RLE_PAIR_SIZE;
        if (count == 0) {
            throw std::runtime_error("RLE: count must not be 0");
        }
        if (out.size() + count > compresskit::MAX_OUTPUT_SIZE) {
            throw std::runtime_error("RLE: output size limit exceeded");
        }
        out.insert(out.end(), count, value);
    }
    return out;
}

}  // namespace compresskit

#ifndef COMPRESSKIT_NO_MAIN
#include "compresskit/cli_launcher.hpp"

int main(int argc, char** argv) {
    compresskit::cli::Algorithm algo{compresskit::rle_encode_buffer,
                                     compresskit::rle_decode_buffer};
    return compresskit::cli::run(algo, argc, argv);
}
#endif
