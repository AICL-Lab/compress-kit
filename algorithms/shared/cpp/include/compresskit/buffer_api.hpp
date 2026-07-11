#pragma once

#include <cstdint>
#include <string>
#include <vector>

#include "compresskit/result.hpp"

namespace compresskit {

// BufferTransform maps an in-memory byte buffer to its encoded/decoded form.
// Throws std::runtime_error on failure. Used by the Buffer layer and CLI.
using BufferTransform = std::vector<uint8_t> (*)(const std::vector<uint8_t>&);

// Convenience: run transform over a complete input buffer with size-limit checks.
Result<std::vector<uint8_t>> encode_buffer(BufferTransform transform,
                                           const std::vector<uint8_t>& input);
Result<std::vector<uint8_t>> decode_buffer(BufferTransform transform,
                                           const std::vector<uint8_t>& input);

// CLI helpers: read input file, apply transform, write output file.
bool encode_file_via_buffer(BufferTransform transform, const std::string& input_path,
                            const std::string& output_path);
bool decode_file_via_buffer(BufferTransform transform, const std::string& input_path,
                            const std::string& output_path);

}  // namespace compresskit
