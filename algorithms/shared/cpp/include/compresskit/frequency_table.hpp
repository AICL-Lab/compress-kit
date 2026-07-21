#pragma once

#include <cstdint>
#include <vector>

namespace compresskit {

// In-memory frequency helpers (used by algorithms that operate on byte buffers).

// Counts byte frequencies of `data` into a SYMBOL_LIMIT-sized vector (EOF left at 0).
std::vector<uint32_t> count_frequencies(const std::vector<uint8_t>& data);

// Scales frequencies so the total does not exceed max_total.
// O(N) single pass: proportional reduction followed by at most one correction sweep.
void scale_frequencies(std::vector<uint32_t>& freq, uint32_t max_total);

// Builds cumulative frequency table of size freq.size()+1.
// Returns empty vector on all-zero input (caller must reject).
std::vector<uint32_t> build_cumulative(const std::vector<uint32_t>& freq);

// Builds the entropy-coder frequency table: byte counts + EOF marker, scaled
// to fit `max_total`. Shared by arithmetic and range coders.
std::vector<uint32_t> build_entropy_frequencies(const std::vector<uint8_t>& data,
                                                uint32_t max_total);

}  // namespace compresskit
