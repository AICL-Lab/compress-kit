#include <cstdint>
#include <vector>
#include <iostream>
#include <fstream>

static const uint32_t SYMBOL_LIMIT = 257;
static const uint32_t EOF_SYMBOL = SYMBOL_LIMIT - 1;
static const uint32_t MAX_TOTAL = 1u << 24;
static const uint32_t RENORM_THRESHOLD = 1u << 24;

static std::vector<uint32_t> build_cumulative(const std::vector<uint32_t>& freq) {
    std::vector<uint32_t> cumulative(freq.size() + 1, 0);
    for (size_t i = 0; i < freq.size(); i++) {
        cumulative[i + 1] = cumulative[i] + freq[i];
    }
    return cumulative;
}

static bool read_u32_le(const std::vector<uint8_t>& in, size_t& pos, uint32_t& v) {
    if (pos + 4 > in.size()) return false;
    v = static_cast<uint32_t>(in[pos]) |
        (static_cast<uint32_t>(in[pos + 1]) << 8) |
        (static_cast<uint32_t>(in[pos + 2]) << 16) |
        (static_cast<uint32_t>(in[pos + 3]) << 24);
    pos += 4;
    return true;
}

static bool read_header(const std::vector<uint8_t>& in, size_t& pos, std::vector<uint32_t>& freq) {
    if (in.size() < 8) return false;
    if (in[0] != 'R' || in[1] != 'C' || in[2] != 'N' || in[3] != 'C') return false;
    pos = 4;
    uint32_t count = 0;
    if (!read_u32_le(in, pos, count)) return false;
    if (count == 0 || count > 1024) return false;
    freq.assign(count, 0);
    for (uint32_t i = 0; i < count; i++) {
        uint32_t v = 0;
        if (!read_u32_le(in, pos, v)) return false;
        freq[i] = v;
    }
    return true;
}

class RangeDecoder {
public:
    RangeDecoder(const uint8_t* data, size_t size)
        : data_(data), size_(size), pos_(0), low_(0), high_(0xFFFFFFFFu), code_(0) {
        for (int i = 0; i < 4; ++i) {
            code_ = (code_ << 8) | read_byte();
        }
    }

    uint32_t decode_symbol(const std::vector<uint32_t>& cumulative) {
        uint64_t range = static_cast<uint64_t>(high_) - low_ + 1;
        uint64_t total = cumulative.back();
        uint64_t offset = code_ - low_;
        uint64_t value = ((offset + 1) * total - 1) / range;

        uint32_t lo = 0;
        uint32_t hi = static_cast<uint32_t>(cumulative.size() - 1);
        while (lo + 1 < hi) {
            uint32_t mid = lo + (hi - lo) / 2;
            if (cumulative[mid] > value) {
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
            low_ <<= 8;
            high_ = (high_ << 8) | 0xFFu;
            code_ = (code_ << 8) | read_byte();
        }

        return symbol;
    }

private:
    const uint8_t* data_;
    size_t size_;
    size_t pos_;
    uint32_t low_;
    uint32_t high_;
    uint32_t code_;

    uint32_t read_byte() {
        if (pos_ < size_) {
            return static_cast<uint32_t>(data_[pos_++]);
        }
        return 0;
    }
};

int main(int argc, char** argv) {
    if (argc != 2) {
        std::cerr << "Usage: " << argv[0] << " <encoded_file>\n";
        return 1;
    }
    
    std::ifstream in(argv[1], std::ios::binary);
    if (!in) {
        std::cerr << "Cannot open file\n";
        return 1;
    }
    in.seekg(0, std::ios::end);
    auto size = in.tellg();
    in.seekg(0, std::ios::beg);
    std::vector<uint8_t> encoded(size);
    in.read(reinterpret_cast<char*>(encoded.data()), size);
    
    size_t pos = 0;
    std::vector<uint32_t> freq;
    if (!read_header(encoded, pos, freq)) {
        std::cerr << "Invalid header\n";
        return 1;
    }
    
    std::cout << "Header size: " << pos << " bytes\n";
    std::cout << "Symbol count: " << freq.size() << "\n";
    std::cout << "Encoded data size: " << (encoded.size() - pos) << " bytes\n";
    
    std::vector<uint32_t> cumulative = build_cumulative(freq);
    std::cout << "Total frequency: " << cumulative.back() << "\n";
    
    RangeDecoder decoder(encoded.data() + pos, encoded.size() - pos);
    
    size_t count = 0;
    for (;;) {
        uint32_t sym = decoder.decode_symbol(cumulative);
        count++;
        if (sym == EOF_SYMBOL) {
            std::cout << "Decoded EOF at symbol count: " << count << "\n";
            break;
        }
        if (count > 10000000) {
            std::cout << "Too many symbols decoded without EOF!\n";
            break;
        }
    }
    
    return 0;
}
