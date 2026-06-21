#pragma once

#include <string>
#include <vector>

#include "compresskit/buffer_api.hpp"

// Memory-based encode/decode entry points for each algorithm.
// Each takes a byte buffer and returns the transformed buffer, throwing on error.
std::vector<uint8_t> huffman_encode_buffer(const std::vector<uint8_t>& input);
std::vector<uint8_t> huffman_decode_buffer(const std::vector<uint8_t>& input);
std::vector<uint8_t> arithmetic_encode_buffer(const std::vector<uint8_t>& input);
std::vector<uint8_t> arithmetic_decode_buffer(const std::vector<uint8_t>& input);
std::vector<uint8_t> rangecoder_encode_buffer(const std::vector<uint8_t>& input);
std::vector<uint8_t> rangecoder_decode_buffer(const std::vector<uint8_t>& input);
std::vector<uint8_t> rle_encode_buffer(const std::vector<uint8_t>& input);
std::vector<uint8_t> rle_decode_buffer(const std::vector<uint8_t>& input);

namespace compresskit {

inline BufferEncoder make_huffman_encoder() {
    return BufferEncoder(::huffman_encode_buffer);
}

inline BufferDecoder make_huffman_decoder() {
    return BufferDecoder(::huffman_decode_buffer);
}

inline BufferEncoder make_arithmetic_encoder() {
    return BufferEncoder(::arithmetic_encode_buffer);
}

inline BufferDecoder make_arithmetic_decoder() {
    return BufferDecoder(::arithmetic_decode_buffer);
}

inline BufferEncoder make_range_encoder() {
    return BufferEncoder(::rangecoder_encode_buffer);
}

inline BufferDecoder make_range_decoder() {
    return BufferDecoder(::rangecoder_decode_buffer);
}

inline BufferEncoder make_rle_encoder() {
    return BufferEncoder(::rle_encode_buffer);
}

inline BufferDecoder make_rle_decoder() {
    return BufferDecoder(::rle_decode_buffer);
}

}  // namespace compresskit
