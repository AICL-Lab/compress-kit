#pragma once

#include <cstdint>
#include <string>
#include <vector>

#include "compresskit/encoder.hpp"

namespace compresskit {

using FileTransform = bool (*)(const std::string&, const std::string&);

class BufferEncoder final : public Encoder {
public:
    explicit BufferEncoder(FileTransform transform);

    Result<std::size_t> process(ByteView in, MutableByteView out) override;
    Result<std::size_t> flush(MutableByteView out) override;
    Result<std::size_t> finish(MutableByteView out) override;
    void reset() noexcept override;
    State state() const noexcept override;

private:
    FileTransform transform_;
    State state_;
    std::vector<uint8_t> input_;
    std::uint64_t total_input_;
};

class BufferDecoder final : public Decoder {
public:
    explicit BufferDecoder(FileTransform transform);

    Result<std::size_t> process(ByteView in, MutableByteView out) override;
    Result<std::size_t> flush(MutableByteView out) override;
    Result<std::size_t> finish(MutableByteView out) override;
    void reset() noexcept override;
    State state() const noexcept override;

private:
    FileTransform transform_;
    State state_;
    std::vector<uint8_t> input_;
    std::uint64_t total_input_;
};

Result<std::vector<uint8_t>> encode_buffer(Encoder& encoder, const std::vector<uint8_t>& input);
Result<std::vector<uint8_t>> decode_buffer(Decoder& decoder, const std::vector<uint8_t>& input);

bool encode_file_via_buffer(FileTransform transform, const std::string& input_path, const std::string& output_path);
bool decode_file_via_buffer(FileTransform transform, const std::string& input_path, const std::string& output_path);

}  // namespace compresskit
