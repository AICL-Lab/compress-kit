#pragma once

#include <string>

#include "compresskit/buffer_api.hpp"

bool huffman_encode_file(const std::string& input_path, const std::string& output_path);
bool huffman_decode_file(const std::string& input_path, const std::string& output_path);
bool arithmetic_encode_file(const std::string& input_path, const std::string& output_path);
bool arithmetic_decode_file(const std::string& input_path, const std::string& output_path);
bool rangecoder_encode_file(const std::string& input_path, const std::string& output_path);
bool rangecoder_decode_file(const std::string& input_path, const std::string& output_path);
bool rle_encode_file(const std::string& input_path, const std::string& output_path);
bool rle_decode_file(const std::string& input_path, const std::string& output_path);

namespace compresskit {

inline BufferEncoder make_huffman_encoder() {
    return BufferEncoder(::huffman_encode_file);
}

inline BufferDecoder make_huffman_decoder() {
    return BufferDecoder(::huffman_decode_file);
}

inline BufferEncoder make_arithmetic_encoder() {
    return BufferEncoder(::arithmetic_encode_file);
}

inline BufferDecoder make_arithmetic_decoder() {
    return BufferDecoder(::arithmetic_decode_file);
}

inline BufferEncoder make_range_encoder() {
    return BufferEncoder(::rangecoder_encode_file);
}

inline BufferDecoder make_range_decoder() {
    return BufferDecoder(::rangecoder_decode_file);
}

inline BufferEncoder make_rle_encoder() {
    return BufferEncoder(::rle_encode_file);
}

inline BufferDecoder make_rle_decoder() {
    return BufferDecoder(::rle_decode_file);
}

}  // namespace compresskit
