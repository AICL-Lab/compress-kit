#include <algorithm>
#include <cstdint>
#include <cstring>
#include <stdexcept>
#include <vector>

#include "compresskit/buffer_api.hpp"
#include "compresskit/constants.hpp"
#include "compresskit/serialization.hpp"

// Run-Length encoding.
// Format:
// - Magic: 4 bytes "RLE\x00"
// - (count: uint32 LE, value: byte) pairs until EOF; count must be > 0.

std::vector<uint8_t> rle_encode_buffer(const std::vector<uint8_t>& input) {
    std::vector<uint8_t> out;
    compresskit::write_magic(out, compresskit::RLE_MAGIC);

    if (input.empty()) {
        return out;
    }

    uint8_t current = input[0];
    uint32_t count = 1;
    for (std::size_t i = 1; i < input.size(); ++i) {
        if (input[i] == current && count < 0xFFFFFFFFu) {
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
    if (input.size() < 4) {
        throw std::runtime_error("RLE: input too short for magic");
    }
    if (std::memcmp(input.data(), compresskit::RLE_MAGIC, 4) != 0) {
        throw std::runtime_error("RLE: bad magic");
    }

    std::vector<uint8_t> out;
    std::size_t pos = 4;
    while (pos < input.size()) {
        if (pos + 5 > input.size()) {
            throw std::runtime_error("RLE: truncated count+value pair");
        }
        uint32_t count = static_cast<uint32_t>(input[pos]) |
                         (static_cast<uint32_t>(input[pos + 1]) << 8) |
                         (static_cast<uint32_t>(input[pos + 2]) << 16) |
                         (static_cast<uint32_t>(input[pos + 3]) << 24);
        uint8_t value = input[pos + 4];
        pos += 5;
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

#ifndef COMPRESSKIT_NO_MAIN
#include "compresskit/cli_launcher.hpp"

int main(int argc, char** argv) {
    compresskit::cli::Algorithm algo{
        [](const std::string& in, const std::string& out) {
            return compresskit::encode_file_via_buffer(rle_encode_buffer, in, out);
        },
        [](const std::string& in, const std::string& out) {
            return compresskit::decode_file_via_buffer(rle_decode_buffer, in, out);
        }};
    return compresskit::cli::run(algo, argc, argv);
}
#endif
