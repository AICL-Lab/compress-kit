#include <cstdint>
#include <stdexcept>
#include <vector>

#include "compresskit/buffer_api.hpp"
#include "compresskit/constants.hpp"
#include "compresskit/frequency_table.hpp"
#include "compresskit/serialization.hpp"

// Range coder (32-bit state, byte-renormalisation, static model).
// Format: "RCNC" + frequency table (LE) + range-coded byte stream.

namespace {

constexpr uint32_t MAX_TOTAL = 1u << 24;
constexpr uint32_t RENORM_THRESHOLD = 1u << 24;
constexpr int STATE_BYTES = 4;  // 32-bit coder state, emitted/loaded byte-by-byte
constexpr int TOP_BYTE_SHIFT = (STATE_BYTES - 1) * compresskit::BITS_PER_BYTE;  // 24

class RangeEncoder {
public:
    explicit RangeEncoder(std::vector<uint8_t>& out) : out_(out), low_(0), high_(UINT32_MAX) {}

    void encode_symbol(uint32_t symbol, const std::vector<uint32_t>& cumulative) {
        uint64_t range = static_cast<uint64_t>(high_) - low_ + 1;
        uint64_t total = cumulative.back();
        uint64_t sym_low = cumulative[symbol];
        uint64_t sym_high = cumulative[symbol + 1];
        high_ = low_ + static_cast<uint32_t>((range * sym_high) / total - 1);
        low_ = low_ + static_cast<uint32_t>((range * sym_low) / total);
        while ((low_ ^ high_) < RENORM_THRESHOLD) {
            out_.push_back(static_cast<uint8_t>(low_ >> TOP_BYTE_SHIFT));
            low_ <<= compresskit::BITS_PER_BYTE;
            high_ = (high_ << compresskit::BITS_PER_BYTE) | 0xFFu;
        }
    }

    void finish() {
        for (int i = 0; i < STATE_BYTES; ++i) {
            out_.push_back(static_cast<uint8_t>(low_ >> TOP_BYTE_SHIFT));
            low_ <<= compresskit::BITS_PER_BYTE;
        }
    }

private:
    std::vector<uint8_t>& out_;
    uint32_t low_, high_;
};

class RangeDecoder {
public:
    RangeDecoder(const uint8_t* data, std::size_t size)
        : data_(data), size_(size), pos_(0), low_(0), high_(UINT32_MAX), code_(0) {
        for (int i = 0; i < STATE_BYTES; ++i) {
            code_ = (code_ << compresskit::BITS_PER_BYTE) | read_byte();
        }
    }

    uint32_t decode_symbol(const std::vector<uint32_t>& cumulative) {
        uint64_t range = static_cast<uint64_t>(high_) - low_ + 1;
        uint64_t total = cumulative.back();
        uint64_t offset = static_cast<uint64_t>(code_) - low_;
        uint64_t value = ((offset + 1) * total - 1) / range;

        // Binary search (O(log N)).
        uint32_t lo = 0;
        uint32_t hi = static_cast<uint32_t>(cumulative.size() - 1);
        while (lo + 1 < hi) {
            uint32_t mid = lo + (hi - lo) / 2;
            if (static_cast<uint64_t>(cumulative[mid]) > value) {
                hi = mid;
            } else {
                lo = mid;
            }
        }
        uint32_t symbol = lo;
        uint64_t sym_low = cumulative[symbol];
        uint64_t sym_high = cumulative[symbol + 1];
        high_ = low_ + static_cast<uint32_t>((range * sym_high) / total - 1);
        low_ = low_ + static_cast<uint32_t>((range * sym_low) / total);
        while ((low_ ^ high_) < RENORM_THRESHOLD) {
            low_ <<= compresskit::BITS_PER_BYTE;
            high_ = (high_ << compresskit::BITS_PER_BYTE) | 0xFFu;
            code_ = (code_ << compresskit::BITS_PER_BYTE) | read_byte();
        }
        return symbol;
    }

private:
    uint32_t read_byte() {
        if (pos_ < size_) {
            return static_cast<uint32_t>(data_[pos_++]);
        }
        return 0;
    }
    const uint8_t* data_;
    std::size_t size_;
    std::size_t pos_;
    uint32_t low_, high_, code_;
};

}  // namespace

std::vector<uint8_t> rangecoder_encode_buffer(const std::vector<uint8_t>& input) {
    if (input.size() > compresskit::MAX_INPUT_SIZE) {
        throw std::runtime_error("range: input too large");
    }
    std::vector<uint32_t> freq = compresskit::build_entropy_frequencies(input, MAX_TOTAL);
    std::vector<uint32_t> cumulative = compresskit::build_cumulative(freq);
    if (cumulative.empty()) {
        throw std::runtime_error("range: empty frequency table");
    }

    std::vector<uint8_t> out;
    out.reserve(input.size() + compresskit::INITIAL_ENCODE_OVERHEAD);
    compresskit::write_frequency_header(out, compresskit::RANGE_MAGIC, freq);

    RangeEncoder encoder(out);
    for (uint8_t b : input) {
        encoder.encode_symbol(static_cast<uint32_t>(b), cumulative);
    }
    encoder.encode_symbol(compresskit::EOF_SYMBOL, cumulative);
    encoder.finish();
    return out;
}

std::vector<uint8_t> rangecoder_decode_buffer(const std::vector<uint8_t>& input) {
    std::size_t pos = 0;
    std::vector<uint32_t> freq =
        compresskit::read_magic_and_frequency_header(input, pos, compresskit::RANGE_MAGIC, "range");
    std::vector<uint32_t> cumulative = compresskit::build_cumulative(freq);
    if (cumulative.empty()) {
        throw std::runtime_error("range: invalid frequency table");
    }

    std::vector<uint8_t> out;
    if (pos >= input.size()) {
        return out;
    }
    RangeDecoder decoder(input.data() + pos, input.size() - pos);
    for (;;) {
        uint32_t sym = decoder.decode_symbol(cumulative);
        if (sym == compresskit::EOF_SYMBOL) {
            break;
        }
        if (out.size() >= compresskit::MAX_OUTPUT_SIZE) {
            throw std::runtime_error("range: output size limit exceeded");
        }
        out.push_back(static_cast<uint8_t>(sym));
    }
    return out;
}

#ifndef COMPRESSKIT_NO_MAIN
#include "compresskit/cli_launcher.hpp"

int main(int argc, char** argv) {
    compresskit::cli::Algorithm algo{rangecoder_encode_buffer, rangecoder_decode_buffer};
    return compresskit::cli::run(algo, argc, argv);
}
#endif
