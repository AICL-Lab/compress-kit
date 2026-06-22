#include "compresskit/frequency_table.hpp"

#include <algorithm>
#include <array>
#include <limits>

#include "compresskit/constants.hpp"

namespace compresskit {
namespace {

bool write_u32_le(std::ostream& out, uint32_t value) {
    char bytes[U32_SIZE];
    for (std::size_t i = 0; i < U32_SIZE; ++i) {
        bytes[i] = static_cast<char>((value >> (i * BITS_PER_BYTE)) & 0xFFu);
    }
    out.write(bytes, static_cast<std::streamsize>(U32_SIZE));
    return static_cast<bool>(out);
}

bool read_u32_le(std::istream& in, uint32_t& value) {
    unsigned char bytes[U32_SIZE]{};
    in.read(reinterpret_cast<char*>(bytes), static_cast<std::streamsize>(U32_SIZE));
    if (!in) {
        return false;
    }
    value = 0;
    for (std::size_t i = 0; i < U32_SIZE; ++i) {
        value |= static_cast<uint32_t>(bytes[i]) << (i * BITS_PER_BYTE);
    }
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
    std::array<unsigned char, STREAM_READ_BUFFER_SIZE> buffer{};
    for (;;) {
        in.read(reinterpret_cast<char*>(buffer.data()),
                static_cast<std::streamsize>(buffer.size()));
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

std::vector<uint32_t> count_frequencies(const std::vector<uint8_t>& data) {
    std::vector<uint32_t> freq(SYMBOL_LIMIT, 0);
    for (uint8_t b : data) {
        ++freq[b];
    }
    return freq;
}

void scale_frequencies(std::vector<uint32_t>& freq, uint32_t max_total) {
    uint64_t total = 0;
    for (uint32_t f : freq) {
        total += f;
    }
    if (total == 0) {
        for (uint32_t& v : freq) {
            v = 1;
        }
        return;
    }
    if (total <= max_total) {
        return;
    }

    // Proportional reduction in a single pass.
    uint64_t new_total = 0;
    for (uint32_t& v : freq) {
        if (v == 0) {
            continue;
        }
        uint64_t scaled = (static_cast<uint64_t>(v) * max_total) / total;
        if (scaled == 0) {
            scaled = 1;
        }
        v = static_cast<uint32_t>(scaled);
        new_total += scaled;
    }

    if (new_total == 0) {
        // All non-zero entries collapsed to 0 above (impossible since we floor to 1).
        uint32_t base = max_total / static_cast<uint32_t>(freq.size());
        if (base == 0) {
            base = 1;
        }
        for (uint32_t& v : freq) {
            v = base;
        }
        return;
    }

    // One correction sweep: distribute the excess by decrementing the largest entries.
    // O(N log N) via sort of indices by frequency; avoids the previous O(N*M) loop.
    if (new_total > max_total) {
        std::vector<uint32_t> idx;
        idx.reserve(freq.size());
        for (uint32_t i = 0; i < freq.size(); ++i) {
            if (freq[i] > 1) {
                idx.push_back(i);
            }
        }
        std::sort(idx.begin(), idx.end(),
                  [&](uint32_t a, uint32_t b) { return freq[a] > freq[b]; });
        uint64_t excess = new_total - max_total;
        for (uint32_t i : idx) {
            uint64_t can_take = freq[i] - 1;
            uint64_t take = std::min(excess, can_take);
            freq[i] -= static_cast<uint32_t>(take);
            excess -= take;
            if (excess == 0) {
                break;
            }
        }
    }
}

std::vector<uint32_t> build_cumulative(const std::vector<uint32_t>& freq) {
    std::vector<uint32_t> cumulative(freq.size() + 1, 0);
    for (std::size_t i = 0; i < freq.size(); ++i) {
        cumulative[i + 1] = cumulative[i] + freq[i];
    }
    if (cumulative.back() == 0) {
        // All-zero table: signal error to caller via empty result.
        cumulative.clear();
    }
    return cumulative;
}

}  // namespace compresskit
