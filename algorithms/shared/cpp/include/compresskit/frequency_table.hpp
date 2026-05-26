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

}  // namespace compresskit
