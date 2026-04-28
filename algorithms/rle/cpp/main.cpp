#include <cstdint>
#include <cstring>
#include <fstream>
#include <iostream>
#include <string>
#include <algorithm>

#include "compresskit/buffer_api.hpp"

// Simple Run-Length encoding implementation.
// Format: repeatedly write (count, value) pairs until EOF.
// - count: 4-byte unsigned integer, little-endian, represents how many times value repeats, must be > 0.
// - value: 1 byte, represents the byte value to repeat.
// This format is simple; all three language implementations are fully consistent for cross-decoding and benchmarking.

// Maximum output size limit (1 GiB) to prevent decompression bomb attacks
static const uint64_t MAX_OUTPUT_SIZE = 1ULL * 1024 * 1024 * 1024;

static void write_u32_le(std::ostream& out, uint32_t v) {
    unsigned char buf[4];
    buf[0] = static_cast<unsigned char>(v & 0xFFu);
    buf[1] = static_cast<unsigned char>((v >> 8) & 0xFFu);
    buf[2] = static_cast<unsigned char>((v >> 16) & 0xFFu);
    buf[3] = static_cast<unsigned char>((v >> 24) & 0xFFu);
    out.write(reinterpret_cast<const char*>(buf), 4);
}

// Read a 32-bit little-endian unsigned integer from stream.
// Returns:
//   true  - Successfully read a complete value and write to out_value
//   false - Normal EOF (no bytes read)
// On truncation (partial bytes read), error message is output to stderr.
static bool read_u32_le(std::istream& in, uint32_t& out_value) {
    unsigned char buf[4];
    in.read(reinterpret_cast<char*>(buf), 4);
    std::streamsize got = in.gcount();
    if (got == 0) {
        // Normal EOF
        return false;
    }
    if (got != 4 || !in) {
        std::cerr << "RLE data truncated: cannot read complete count field\n";
        in.setstate(std::ios::badbit);
        return false;
    }
    out_value =
        (static_cast<uint32_t>(buf[0])      ) |
        (static_cast<uint32_t>(buf[1]) << 8 ) |
        (static_cast<uint32_t>(buf[2]) << 16) |
        (static_cast<uint32_t>(buf[3]) << 24);
    return true;
}

static bool rle_encode_file_checked(const std::string& input_path, const std::string& output_path);
static bool rle_decode_file_checked(const std::string& input_path, const std::string& output_path);

bool rle_encode_file(const std::string& input_path, const std::string& output_path) {
    return rle_encode_file_checked(input_path, output_path);
}

// Perform Run-Length encoding on entire file.
// input_path is the original binary file path.
// output_path is the encoded file path.
static bool rle_encode_file_checked(const std::string& input_path, const std::string& output_path) {
    std::ifstream in(input_path, std::ios::binary);
    if (!in) {
        std::cerr << "cannot open input file for reading: " << input_path << "\n";
        return false;
    }
    std::ofstream out(output_path, std::ios::binary);
    if (!out) {
        std::cerr << "cannot open output file for writing: " << output_path << "\n";
        return false;
    }

    char c;
    if (!in.get(c)) {
        // Empty file: encoding result is also empty.
        return true;
    }
    unsigned char current = static_cast<unsigned char>(c);
    uint32_t count = 1;

    while (in.get(c)) {
        unsigned char b = static_cast<unsigned char>(c);
        if (b == current && count < 0xFFFFFFFFu) {
            ++count;
        } else {
            write_u32_le(out, count);
            out.put(static_cast<char>(current));
            if (!out) {
                std::cerr << "failed to write RLE data\n";
                return false;
            }
            current = b;
            count = 1;
        }
    }

    // Write last run
    write_u32_le(out, count);
    out.put(static_cast<char>(current));
    if (!out) {
        std::cerr << "failed to write RLE data\n";
        return false;
    }

    return true;
}

// Decode RLE encoded file back to original byte stream.
// input_path is the encoded input file path.
// output_path is the decoded output file path.
static bool rle_decode_file_checked(const std::string& input_path, const std::string& output_path) {
    std::ifstream in(input_path, std::ios::binary);
    if (!in) {
        std::cerr << "cannot open input file for reading: " << input_path << "\n";
        return false;
    }
    std::ofstream out(output_path, std::ios::binary);
    if (!out) {
        std::cerr << "cannot open output file for writing: " << output_path << "\n";
        return false;
    }

    const std::size_t BUF_SIZE = 4096;
    char buffer[BUF_SIZE];
    uint64_t total_written = 0;

    while (true) {
        uint32_t count = 0;
        if (!read_u32_le(in, count)) {
            if (!in.bad()) {
                // Normal EOF
                break;
            } else {
                // read_u32_le already output error message
                return false;
            }
        }
        if (count == 0) {
            std::cerr << "invalid RLE data: count should not be 0\n";
            return false;
        }

        // Check output size limit
        if (total_written + static_cast<uint64_t>(count) > MAX_OUTPUT_SIZE) {
            std::cerr << "output size limit exceeded (max " << MAX_OUTPUT_SIZE << " bytes)\n";
            return false;
        }

        char value_char;
        if (!in.get(value_char)) {
            std::cerr << "RLE data truncated: missing value byte\n";
            return false;
        }
        unsigned char value = static_cast<unsigned char>(value_char);

        while (count > 0) {
            std::size_t chunk = static_cast<std::size_t>(std::min<uint32_t>(count, static_cast<uint32_t>(BUF_SIZE)));
            std::memset(buffer, static_cast<int>(value), chunk);
            out.write(buffer, static_cast<std::streamsize>(chunk));
            if (!out) {
                std::cerr << "failed to write decoded data\n";
                return false;
            }
            total_written += static_cast<uint64_t>(chunk);
            count -= static_cast<uint32_t>(chunk);
        }
    }

    return true;
}

bool rle_decode_file(const std::string& input_path, const std::string& output_path) {
    return rle_decode_file_checked(input_path, output_path);
}

#ifndef COMPRESSKIT_NO_MAIN
int main(int argc, char** argv) {
    if (argc != 4) {
        std::cerr << "usage: " << argv[0] << " encode|decode input output\n";
        return 1;
    }
    std::string mode = argv[1];
    std::string input_path = argv[2];
    std::string output_path = argv[3];

    bool ok = true;

    if (mode == "encode") {
        ok = compresskit::encode_file_via_buffer(rle_encode_file, input_path, output_path);
    } else if (mode == "decode") {
        ok = compresskit::decode_file_via_buffer(rle_decode_file, input_path, output_path);
    } else {
        std::cerr << "unknown mode: " << mode << ", expected encode or decode\n";
        return 1;
    }

    return ok ? 0 : 1;
}
#endif
