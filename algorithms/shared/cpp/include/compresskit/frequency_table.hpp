#pragma once

#include <cstdint>
#include <istream>
#include <ostream>
#include <vector>

namespace compresskit {

enum class FrequencyTableReadStatus {
    OK = 0,
    TRUNCATED,
    BAD_COUNT,
};

enum class FrequencyCountStatus {
    OK = 0,
    IO_ERROR,
    OVERFLOW,
};

bool write_frequency_table(std::ostream& out, const std::vector<uint32_t>& freq);

FrequencyTableReadStatus read_frequency_table(std::istream& in, std::vector<uint32_t>& freq,
                                              uint32_t expected_count,
                                              uint32_t* actual_count = nullptr);

FrequencyCountStatus accumulate_frequencies(std::istream& in, std::vector<uint32_t>& freq,
                                            uint32_t* overflow_symbol = nullptr);

// In-memory frequency helpers (used by algorithms that operate on byte buffers).

// Counts byte frequencies of `data` into a SYMBOL_LIMIT-sized vector (EOF left at 0).
std::vector<uint32_t> count_frequencies(const std::vector<uint8_t>& data);

// Scales frequencies so the total does not exceed max_total.
// O(N) single pass: proportional reduction followed by at most one correction sweep.
void scale_frequencies(std::vector<uint32_t>& freq, uint32_t max_total);

// Builds cumulative frequency table of size freq.size()+1.
// Returns empty vector on all-zero input (caller must reject).
std::vector<uint32_t> build_cumulative(const std::vector<uint32_t>& freq);

}  // namespace compresskit
