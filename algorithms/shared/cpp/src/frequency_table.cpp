#include "compresskit/frequency_table.hpp"

#include <array>
#include <limits>

namespace compresskit {
namespace {

bool write_u32_le(std::ostream& out, uint32_t value) {
    const std::array<char, 4> bytes = {
        static_cast<char>(value & 0xFFu),
        static_cast<char>((value >> 8) & 0xFFu),
        static_cast<char>((value >> 16) & 0xFFu),
        static_cast<char>((value >> 24) & 0xFFu),
    };
    out.write(bytes.data(), static_cast<std::streamsize>(bytes.size()));
    return static_cast<bool>(out);
}

bool read_u32_le(std::istream& in, uint32_t& value) {
    std::array<unsigned char, 4> bytes{};
    in.read(reinterpret_cast<char*>(bytes.data()), static_cast<std::streamsize>(bytes.size()));
    if (!in) {
        return false;
    }
    value = static_cast<uint32_t>(bytes[0]) | (static_cast<uint32_t>(bytes[1]) << 8) |
            (static_cast<uint32_t>(bytes[2]) << 16) | (static_cast<uint32_t>(bytes[3]) << 24);
    return true;
}

}  // namespace

bool write_frequency_table(std::ostream& out, const std::vector<uint32_t>& freq) {
    if (!write_u32_le(out, static_cast<uint32_t>(freq.size()))) {
        return false;
    }
    for (uint32_t value : freq) {
        if (!write_u32_le(out, value)) {
            return false;
        }
    }
    return true;
}

FrequencyTableReadStatus read_frequency_table(std::istream& in, std::vector<uint32_t>& freq,
                                              uint32_t expected_count, uint32_t* actual_count) {
    uint32_t count = 0;
    if (!read_u32_le(in, count)) {
        freq.clear();
        return FrequencyTableReadStatus::TRUNCATED;
    }
    if (actual_count) {
        *actual_count = count;
    }
    if (expected_count != 0 && count != expected_count) {
        freq.clear();
        return FrequencyTableReadStatus::BAD_COUNT;
    }

    freq.assign(count, 0);
    for (uint32_t& value : freq) {
        if (!read_u32_le(in, value)) {
            freq.clear();
            return FrequencyTableReadStatus::TRUNCATED;
        }
    }
    return FrequencyTableReadStatus::OK;
}

FrequencyCountStatus accumulate_frequencies(std::istream& in, std::vector<uint32_t>& freq,
                                            uint32_t* overflow_symbol) {
    std::array<unsigned char, 32 * 1024> buffer{};
    for (;;) {
        in.read(reinterpret_cast<char*>(buffer.data()), static_cast<std::streamsize>(buffer.size()));
        const std::streamsize read_count = in.gcount();
        for (std::streamsize i = 0; i < read_count; ++i) {
            const uint32_t symbol = static_cast<uint32_t>(buffer[static_cast<std::size_t>(i)]);
            if (freq[symbol] == std::numeric_limits<uint32_t>::max()) {
                if (overflow_symbol) {
                    *overflow_symbol = symbol;
                }
                return FrequencyCountStatus::OVERFLOW;
            }
            ++freq[symbol];
        }
        if (in.eof()) {
            return FrequencyCountStatus::OK;
        }
        if (!in) {
            return FrequencyCountStatus::IO_ERROR;
        }
    }
}

}  // namespace compresskit
