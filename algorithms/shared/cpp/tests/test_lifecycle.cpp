#include <cassert>
#include <cstdint>
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

}  // namespace

int main() {
    const AlgorithmCase algorithms[] = {
        {"Huffman", compresskit::huffman_encode_buffer, compresskit::huffman_decode_buffer},
        {"Arithmetic", compresskit::arithmetic_encode_buffer,
         compresskit::arithmetic_decode_buffer},
        {"Range", compresskit::rangecoder_encode_buffer, compresskit::rangecoder_decode_buffer},
        {"RLE", compresskit::rle_encode_buffer, compresskit::rle_decode_buffer},
    };

    for (const AlgorithmCase& algorithm : algorithms) {
        test_roundtrip(algorithm);
    }

    return 0;
}
