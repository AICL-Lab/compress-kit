#include <cassert>
#include <cstdint>
#include <sstream>
#include <string>
#include <vector>

#include "compresskit/algorithms.hpp"
#include "compresskit/buffer_api.hpp"
#include "compresskit/frequency_table.hpp"

namespace {

struct AlgorithmCase {
    const char* name;
    compresskit::BufferTransform encode;
    compresskit::BufferTransform decode;
};

void test_roundtrip(const AlgorithmCase& algorithm) {
    std::vector<uint8_t> input = {'s', 't', 'r', 'e', 'a', 'm', '-', 'a', 'p', 'i'};

    auto encoded = compresskit::encode_buffer(algorithm.encode, input);
    assert(encoded.status == compresskit::StatusCode::OK);
    assert(!encoded.value.empty());

    auto decoded = compresskit::decode_buffer(algorithm.decode, encoded.value);
    assert(decoded.status == compresskit::StatusCode::OK);
    assert(decoded.value == input);
}

void test_write_frequency_table_uses_little_endian_layout() {
    std::ostringstream out(std::ios::binary);
    const std::vector<uint32_t> freq = {0x78563412u, 0x01020304u};

    assert(compresskit::write_frequency_table(out, freq));

    const std::string bytes = out.str();
    const std::string expected(
        "\x02\x00\x00\x00"
        "\x12\x34\x56\x78"
        "\x04\x03\x02\x01",
        12);
    assert(bytes == expected);
}

void test_read_frequency_table_decodes_little_endian_values() {
    const std::string bytes(
        "\x02\x00\x00\x00"
        "\x12\x34\x56\x78"
        "\x04\x03\x02\x01",
        12);
    std::istringstream in(bytes, std::ios::binary);
    std::vector<uint32_t> freq;
    uint32_t actual_count = 0;

    const auto status = compresskit::read_frequency_table(in, freq, 2, &actual_count);
    assert(status == compresskit::FrequencyTableReadStatus::OK);
    assert(actual_count == 2);
    assert((freq == std::vector<uint32_t>{0x78563412u, 0x01020304u}));
}

void test_read_frequency_table_reports_bad_count() {
    const std::string bytes("\x02\x00\x00\x00", 4);
    std::istringstream in(bytes, std::ios::binary);
    std::vector<uint32_t> freq;
    uint32_t actual_count = 0;

    const auto status = compresskit::read_frequency_table(in, freq, 3, &actual_count);

    assert(status == compresskit::FrequencyTableReadStatus::BAD_COUNT);
    assert(actual_count == 2);
    assert(freq.empty());
}

void test_accumulate_frequencies_reports_overflow() {
    std::vector<uint32_t> freq(257, 0);
    freq[0] = UINT32_MAX;
    std::istringstream in(std::string(1, '\0'), std::ios::binary);
    uint32_t overflow_symbol = UINT32_MAX;

    const auto status = compresskit::accumulate_frequencies(in, freq, &overflow_symbol);

    assert(status == compresskit::FrequencyCountStatus::OVERFLOW);
    assert(overflow_symbol == 0);
    assert(freq[0] == UINT32_MAX);
}

}  // namespace

int main() {
    const AlgorithmCase algorithms[] = {
        {"Huffman", huffman_encode_buffer, huffman_decode_buffer},
        {"Arithmetic", arithmetic_encode_buffer, arithmetic_decode_buffer},
        {"Range", rangecoder_encode_buffer, rangecoder_decode_buffer},
        {"RLE", rle_encode_buffer, rle_decode_buffer},
    };

    for (const AlgorithmCase& algorithm : algorithms) {
        test_roundtrip(algorithm);
    }

    test_write_frequency_table_uses_little_endian_layout();
    test_read_frequency_table_decodes_little_endian_values();
    test_read_frequency_table_reports_bad_count();
    test_accumulate_frequencies_reports_overflow();

    return 0;
}
