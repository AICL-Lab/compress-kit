#include <cstdint>
#include <cstring>
#include <stdexcept>
#include <vector>

#include "compresskit/buffer_api.hpp"
#include "compresskit/constants.hpp"
#include "compresskit/frequency_table.hpp"

// Arithmetic coding (32-bit state, static model).
// Format: "AENC" + frequency table (LE) + bitstream (MSB-first).

namespace {

constexpr uint32_t MAX_TOTAL = 1u << 24;
constexpr uint64_t STATE_BITS = 32;
constexpr uint64_t FULL_RANGE = 1ull << STATE_BITS;
constexpr uint64_t HALF_RANGE = FULL_RANGE >> 1;
constexpr uint64_t FIRST_QUARTER = HALF_RANGE >> 1;
constexpr uint64_t THIRD_QUARTER = FIRST_QUARTER * 3;

class BitWriter {
public:
    void write_bit(int bit) {
        buffer_ = static_cast<uint8_t>((buffer_ << 1) | (bit & 1));
        if (++bits_ == 8) {
            out_.push_back(buffer_);
            bits_ = 0;
            buffer_ = 0;
        }
    }
    void flush() {
        if (bits_ > 0) {
            buffer_ = static_cast<uint8_t>(buffer_ << (8 - bits_));
            out_.push_back(buffer_);
            bits_ = 0;
            buffer_ = 0;
        }
    }
    std::vector<uint8_t> finish() {
        flush();
        return std::move(out_);
    }

private:
    std::vector<uint8_t> out_;
    uint8_t buffer_ = 0;
    int bits_ = 0;
};

class BitReader {
public:
    explicit BitReader(const std::vector<uint8_t>& data) : data_(data) {}
    int read_bit() {
        if (byte_pos_ >= data_.size()) {
            return 0;
        }
        int bit = (data_[byte_pos_] >> (7 - bit_pos_)) & 1;
        if (++bit_pos_ == 8) {
            ++byte_pos_;
            bit_pos_ = 0;
        }
        return bit;
    }
    bool eof() const { return byte_pos_ >= data_.size(); }

private:
    const std::vector<uint8_t>& data_;
    std::size_t byte_pos_ = 0;
    int bit_pos_ = 0;
};

class ArithmeticEncoder {
public:
    explicit ArithmeticEncoder(BitWriter& w)
        : writer_(w), low_(0), high_(FULL_RANGE - 1), pending_(0) {}

    void encode_symbol(uint32_t symbol, const std::vector<uint32_t>& cumulative) {
        uint64_t range = high_ - low_ + 1;
        uint64_t total = cumulative.back();
        uint64_t sym_low = cumulative[symbol];
        uint64_t sym_high = cumulative[symbol + 1];
        high_ = low_ + (range * sym_high) / total - 1;
        low_ = low_ + (range * sym_low) / total;
        for (;;) {
            if (high_ < HALF_RANGE) {
                output_bit(0);
            } else if (low_ >= HALF_RANGE) {
                output_bit(1);
                low_ -= HALF_RANGE;
                high_ -= HALF_RANGE;
            } else if (low_ >= FIRST_QUARTER && high_ < THIRD_QUARTER) {
                ++pending_;
                low_ -= FIRST_QUARTER;
                high_ -= FIRST_QUARTER;
            } else {
                break;
            }
            low_ <<= 1;
            high_ = (high_ << 1) | 1;
        }
    }

    void finish() {
        ++pending_;
        output_bit(low_ < FIRST_QUARTER ? 0 : 1);
        writer_.flush();
    }

private:
    void output_bit(int bit) {
        writer_.write_bit(bit);
        int complement = bit ^ 1;
        while (pending_ > 0) {
            writer_.write_bit(complement);
            --pending_;
        }
    }
    BitWriter& writer_;
    uint64_t low_, high_, pending_;
};

class ArithmeticDecoder {
public:
    explicit ArithmeticDecoder(BitReader& r)
        : reader_(r), low_(0), high_(FULL_RANGE - 1), code_(0) {
        for (uint64_t i = 0; i < STATE_BITS; ++i) {
            code_ = (code_ << 1) | static_cast<uint64_t>(reader_.read_bit());
        }
    }

    uint32_t decode_symbol(const std::vector<uint32_t>& cumulative) {
        uint64_t range = high_ - low_ + 1;
        uint64_t total = cumulative.back();
        uint64_t offset = code_ - low_;
        uint64_t value = ((offset + 1) * total - 1) / range;

        // Binary search for symbol (O(log N)).
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
        high_ = low_ + (range * sym_high) / total - 1;
        low_ = low_ + (range * sym_low) / total;
        for (;;) {
            if (high_ < HALF_RANGE) {
            } else if (low_ >= HALF_RANGE) {
                low_ -= HALF_RANGE;
                high_ -= HALF_RANGE;
                code_ -= HALF_RANGE;
            } else if (low_ >= FIRST_QUARTER && high_ < THIRD_QUARTER) {
                low_ -= FIRST_QUARTER;
                high_ -= FIRST_QUARTER;
                code_ -= FIRST_QUARTER;
            } else {
                break;
            }
            low_ <<= 1;
            high_ = (high_ << 1) | 1;
            code_ = (code_ << 1) | static_cast<uint64_t>(reader_.read_bit());
        }
        return symbol;
    }

private:
    BitReader& reader_;
    uint64_t low_, high_, code_;
};

std::vector<uint32_t> build_frequencies(const std::vector<uint8_t>& data) {
    std::vector<uint32_t> freq = compresskit::count_frequencies(data);
    freq[compresskit::EOF_SYMBOL] = 1;
    compresskit::scale_frequencies(freq, MAX_TOTAL);
    return freq;
}

void write_header(std::vector<uint8_t>& out, const std::vector<uint32_t>& freq) {
    out.push_back(static_cast<uint8_t>(compresskit::ARITHMETIC_MAGIC[0]));
    out.push_back(static_cast<uint8_t>(compresskit::ARITHMETIC_MAGIC[1]));
    out.push_back(static_cast<uint8_t>(compresskit::ARITHMETIC_MAGIC[2]));
    out.push_back(static_cast<uint8_t>(compresskit::ARITHMETIC_MAGIC[3]));
    // Frequency table: count (u32 LE) + count × u32 LE.
    auto push_u32 = [&](uint32_t v) {
        out.push_back(static_cast<uint8_t>(v & 0xFFu));
        out.push_back(static_cast<uint8_t>((v >> 8) & 0xFFu));
        out.push_back(static_cast<uint8_t>((v >> 16) & 0xFFu));
        out.push_back(static_cast<uint8_t>((v >> 24) & 0xFFu));
    };
    push_u32(static_cast<uint32_t>(freq.size()));
    for (uint32_t v : freq) {
        push_u32(v);
    }
}

std::vector<uint32_t> read_frequencies(const std::vector<uint8_t>& input, std::size_t& pos) {
    if (pos + 4 > input.size()) {
        throw std::runtime_error("arithmetic: truncated frequency count");
    }
    uint32_t count = static_cast<uint32_t>(input[pos]) |
                     (static_cast<uint32_t>(input[pos + 1]) << 8) |
                     (static_cast<uint32_t>(input[pos + 2]) << 16) |
                     (static_cast<uint32_t>(input[pos + 3]) << 24);
    pos += 4;
    if (count != compresskit::SYMBOL_LIMIT) {
        throw std::runtime_error("arithmetic: bad frequency count");
    }
    std::vector<uint32_t> freq(count, 0);
    for (uint32_t i = 0; i < count; ++i) {
        if (pos + 4 > input.size()) {
            throw std::runtime_error("arithmetic: truncated frequency entry");
        }
        freq[i] = static_cast<uint32_t>(input[pos]) | (static_cast<uint32_t>(input[pos + 1]) << 8) |
                  (static_cast<uint32_t>(input[pos + 2]) << 16) |
                  (static_cast<uint32_t>(input[pos + 3]) << 24);
        pos += 4;
    }
    return freq;
}

}  // namespace

std::vector<uint8_t> arithmetic_encode_buffer(const std::vector<uint8_t>& input) {
    if (input.size() > compresskit::MAX_INPUT_SIZE) {
        throw std::runtime_error("arithmetic: input too large");
    }
    std::vector<uint32_t> freq = build_frequencies(input);
    std::vector<uint32_t> cumulative = compresskit::build_cumulative(freq);
    if (cumulative.empty()) {
        throw std::runtime_error("arithmetic: empty frequency table");
    }

    std::vector<uint8_t> out;
    write_header(out, freq);

    BitWriter writer;
    ArithmeticEncoder encoder(writer);
    for (uint8_t b : input) {
        encoder.encode_symbol(static_cast<uint32_t>(b), cumulative);
    }
    encoder.encode_symbol(compresskit::EOF_SYMBOL, cumulative);
    encoder.finish();

    std::vector<uint8_t> bits = writer.finish();
    out.insert(out.end(), bits.begin(), bits.end());
    return out;
}

std::vector<uint8_t> arithmetic_decode_buffer(const std::vector<uint8_t>& input) {
    if (input.size() < 4) {
        throw std::runtime_error("arithmetic: input too short");
    }
    if (std::memcmp(input.data(), compresskit::ARITHMETIC_MAGIC, 4) != 0) {
        throw std::runtime_error("arithmetic: bad magic");
    }
    std::size_t pos = 4;
    std::vector<uint32_t> freq = read_frequencies(input, pos);
    std::vector<uint32_t> cumulative = compresskit::build_cumulative(freq);
    if (cumulative.empty()) {
        throw std::runtime_error("arithmetic: invalid frequency table");
    }

    std::vector<uint8_t> payload(input.begin() + pos, input.end());
    BitReader reader(payload);
    ArithmeticDecoder decoder(reader);

    std::vector<uint8_t> out;
    for (;;) {
        uint32_t sym = decoder.decode_symbol(cumulative);
        if (sym == compresskit::EOF_SYMBOL) {
            break;
        }
        if (out.size() >= compresskit::MAX_OUTPUT_SIZE) {
            throw std::runtime_error("arithmetic: output size limit exceeded");
        }
        out.push_back(static_cast<uint8_t>(sym));
    }
    return out;
}

#ifndef COMPRESSKIT_NO_MAIN
#include "compresskit/cli_launcher.hpp"

int main(int argc, char** argv) {
    compresskit::cli::Algorithm algo{
        [](const std::string& in, const std::string& out) {
            return compresskit::encode_file_via_buffer(arithmetic_encode_buffer, in, out);
        },
        [](const std::string& in, const std::string& out) {
            return compresskit::decode_file_via_buffer(arithmetic_decode_buffer, in, out);
        }};
    return compresskit::cli::run("arithmetic", algo, argc, argv);
}
#endif
