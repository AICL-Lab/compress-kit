#include <cassert>
#include <cstdint>
#include <string>
#include <vector>

#include "compresskit/algorithms.hpp"

namespace {

struct AlgorithmCase {
    const char* name;
    compresskit::BufferEncoder (*make_encoder)();
    compresskit::BufferDecoder (*make_decoder)();
};

void test_roundtrip_and_lifecycle(const AlgorithmCase& algorithm) {
    std::vector<uint8_t> input = {'s', 't', 'r', 'e', 'a', 'm', '-', 'a', 'p', 'i'};

    compresskit::BufferEncoder encoder = algorithm.make_encoder();
    assert(encoder.state() == compresskit::State::READY);

    auto process_one = encoder.process({input.data(), 4}, {nullptr, 0});
    assert(process_one.status == compresskit::StatusCode::OK);
    assert(encoder.state() == compresskit::State::STREAMING);

    auto flush = encoder.flush({nullptr, 0});
    assert(flush.status == compresskit::StatusCode::OK);
    assert(encoder.state() == compresskit::State::FLUSHING);

    auto process_two = encoder.process({input.data() + 4, input.size() - 4}, {nullptr, 0});
    assert(process_two.status == compresskit::StatusCode::OK);
    assert(encoder.state() == compresskit::State::STREAMING);

    std::vector<uint8_t> tiny(1);
    auto finish_small = encoder.finish({tiny.data(), tiny.size()});
    assert(finish_small.status == compresskit::StatusCode::BUF_TOO_SMALL);
    assert(encoder.state() == compresskit::State::STREAMING);

    std::vector<uint8_t> encoded(4096);
    auto finish = encoder.finish({encoded.data(), encoded.size()});
    assert(finish.status == compresskit::StatusCode::OK);
    assert(encoder.state() == compresskit::State::FINISHED);
    encoded.resize(finish.value);

    auto invalid = encoder.process({input.data(), input.size()}, {nullptr, 0});
    assert(invalid.status == compresskit::StatusCode::ERR_INVALID_STATE);
    assert(encoder.state() == compresskit::State::ERROR);

    encoder.reset();
    assert(encoder.state() == compresskit::State::READY);

    compresskit::BufferDecoder decoder = algorithm.make_decoder();
    auto decoded = compresskit::decode_buffer(decoder, encoded);
    assert(decoded.status == compresskit::StatusCode::OK);
    assert(decoded.value == input);

    compresskit::BufferEncoder buffer_encoder = algorithm.make_encoder();
    auto buffer_encoded = compresskit::encode_buffer(buffer_encoder, input);
    assert(buffer_encoded.status == compresskit::StatusCode::OK);
    assert(buffer_encoded.value == encoded);
}

}  // namespace

int main() {
    const AlgorithmCase algorithms[] = {
        {"Huffman", compresskit::make_huffman_encoder, compresskit::make_huffman_decoder},
        {"Arithmetic", compresskit::make_arithmetic_encoder, compresskit::make_arithmetic_decoder},
        {"Range", compresskit::make_range_encoder, compresskit::make_range_decoder},
        {"RLE", compresskit::make_rle_encoder, compresskit::make_rle_decoder},
    };

    for (const AlgorithmCase& algorithm : algorithms) {
        test_roundtrip_and_lifecycle(algorithm);
    }

    return 0;
}
