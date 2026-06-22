#pragma once

// MSB-first bit-level I/O for entropy coders that emit/consume bit streams.

#include <cstddef>
#include <cstdint>
#include <vector>

#include "compresskit/constants.hpp"

namespace compresskit {

// MSB-first bit writer. Buffers bits and flushes full bytes to an internal
// buffer; partial trailing bits are left-aligned on flush().
class BitWriter {
public:
    void write_bit(int bit) {
        buffer_ = static_cast<uint8_t>((buffer_ << 1) | (bit & 1));
        if (++bits_ == BITS_PER_BYTE) {
            out_.push_back(buffer_);
            bits_ = 0;
            buffer_ = 0;
        }
    }
    void flush() {
        if (bits_ > 0) {
            buffer_ = static_cast<uint8_t>(buffer_ << (BITS_PER_BYTE - bits_));
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

// MSB-first bit reader. Returns 0 for bits read past the end of the stream.
class BitReader {
public:
    explicit BitReader(const std::vector<uint8_t>& data) : data_(data) {}
    int read_bit() {
        if (byte_pos_ >= data_.size()) {
            return 0;
        }
        int bit = (data_[byte_pos_] >> ((BITS_PER_BYTE - 1) - bit_pos_)) & 1;
        if (++bit_pos_ == BITS_PER_BYTE) {
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

}  // namespace compresskit
