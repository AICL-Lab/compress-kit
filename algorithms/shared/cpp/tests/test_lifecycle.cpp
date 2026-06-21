#include <algorithm>
#include <cassert>
#include <cstdint>
#include <sstream>
#include <string>
#include <vector>

#include "compresskit/algorithms.hpp"
#include "compresskit/frequency_table.hpp"

namespace {

struct AlgorithmCase {
    const char* name;
    compresskit::BufferEncoder (*make_encoder)();
    compresskit::BufferDecoder (*make_decoder)();
};

struct ScriptedEncoder final : compresskit::Encoder {
    int finish_calls = 0;

    compresskit::Result<std::size_t> process(compresskit::ByteView,
                                             compresskit::MutableByteView) override {
        return {compresskit::StatusCode::OK, 0};
    }

    compresskit::Result<std::size_t> flush(compresskit::MutableByteView) override {
        // encode_buffer does not call flush; this only satisfies the interface in the scripted test
        // double.
        return {compresskit::StatusCode::OK, 0};
    }

    compresskit::Result<std::size_t> finish(compresskit::MutableByteView out) override {
        ++finish_calls;
        if (finish_calls == 1) {
            std::copy_n("abc", 3, out.data);
            return {compresskit::StatusCode::BUF_TOO_SMALL, 3};
        }
        std::copy_n("def", 3, out.data);
        return {compresskit::StatusCode::OK, 3};
    }

    void reset() noexcept override {}
    compresskit::State state() const noexcept override { return compresskit::State::STREAMING; }
};

struct ScriptedDecoder final : compresskit::Decoder {
    int finish_calls = 0;

    compresskit::Result<std::size_t> process(compresskit::ByteView,
                                             compresskit::MutableByteView) override {
        return {compresskit::StatusCode::OK, 0};
    }

    compresskit::Result<std::size_t> flush(compresskit::MutableByteView) override {
        // decode_buffer does not call flush; this only satisfies the interface in the scripted test
        // double.
        return {compresskit::StatusCode::OK, 0};
    }

    compresskit::Result<std::size_t> finish(compresskit::MutableByteView out) override {
        ++finish_calls;
        if (finish_calls == 1) {
            std::copy_n("uvw", 3, out.data);
            return {compresskit::StatusCode::BUF_TOO_SMALL, 3};
        }
        std::copy_n("xyz", 3, out.data);
        return {compresskit::StatusCode::OK, 3};
    }

    void reset() noexcept override {}
    compresskit::State state() const noexcept override { return compresskit::State::STREAMING; }
};

void test_roundtrip_and_lifecycle(const AlgorithmCase& algorithm) {
    std::vector<uint8_t> input = {'s', 't', 'r', 'e', 'a', 'm', '-', 'a', 'p', 'i'};

    compresskit::BufferEncoder encoder = algorithm.make_encoder();
    assert(encoder.state() == compresskit::State::READY);

    auto process_one = encoder.process({input.data(), 4}, {nullptr, 0});
    assert(process_one.status == compresskit::StatusCode::OK);
    assert(process_one.value == 0);
    assert(encoder.state() == compresskit::State::STREAMING);

    auto flush = encoder.flush({nullptr, 0});
    assert(flush.status == compresskit::StatusCode::OK);
    assert(flush.value == 0);
    assert(encoder.state() == compresskit::State::FLUSHING);

    auto process_two = encoder.process({input.data() + 4, input.size() - 4}, {nullptr, 0});
    assert(process_two.status == compresskit::StatusCode::OK);
    assert(process_two.value == 0);
    assert(encoder.state() == compresskit::State::STREAMING);

    std::vector<uint8_t> tiny(1);
    auto finish_small = encoder.finish({tiny.data(), tiny.size()});
    assert(finish_small.status == compresskit::StatusCode::BUF_TOO_SMALL);
    assert(encoder.state() == compresskit::State::STREAMING);

    std::vector<uint8_t> encoded(4096);
    auto finish = encoder.finish({encoded.data(), encoded.size()});
    assert(finish.status == compresskit::StatusCode::OK);
    assert(finish.value > 0);
    assert(encoder.state() == compresskit::State::FINISHED);
    encoded.resize(finish.value);

    auto invalid = encoder.process({input.data(), input.size()}, {nullptr, 0});
    assert(invalid.status == compresskit::StatusCode::ERR_INVALID_STATE);
    assert(invalid.value == 0);
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

void test_encode_buffer_preserves_finish_retry_prefix() {
    ScriptedEncoder encoder;
    auto encoded = compresskit::encode_buffer(encoder, std::vector<uint8_t>{'x'});
    assert(encoded.status == compresskit::StatusCode::OK);
    assert(encoder.finish_calls == 2);
    assert(std::string(encoded.value.begin(), encoded.value.end()) == "abcdef");
}

void test_decode_buffer_preserves_finish_retry_prefix() {
    ScriptedDecoder decoder;
    auto decoded = compresskit::decode_buffer(decoder, std::vector<uint8_t>{'x'});
    assert(decoded.status == compresskit::StatusCode::OK);
    assert(decoder.finish_calls == 2);
    assert(std::string(decoded.value.begin(), decoded.value.end()) == "uvwxyz");
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
        {"Huffman", compresskit::make_huffman_encoder, compresskit::make_huffman_decoder},
        {"Arithmetic", compresskit::make_arithmetic_encoder, compresskit::make_arithmetic_decoder},
        {"Range", compresskit::make_range_encoder, compresskit::make_range_decoder},
        {"RLE", compresskit::make_rle_encoder, compresskit::make_rle_decoder},
    };

    for (const AlgorithmCase& algorithm : algorithms) {
        test_roundtrip_and_lifecycle(algorithm);
    }

    test_encode_buffer_preserves_finish_retry_prefix();
    test_decode_buffer_preserves_finish_retry_prefix();
    test_write_frequency_table_uses_little_endian_layout();
    test_read_frequency_table_decodes_little_endian_values();
    test_read_frequency_table_reports_bad_count();
    test_accumulate_frequencies_reports_overflow();

    return 0;
}
