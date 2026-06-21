#pragma once

#include <cstdint>
#include <string>
#include <vector>

#include "compresskit/encoder.hpp"

namespace compresskit {

// BufferTransform maps an in-memory byte buffer to its encoded/decoded form.
// Throws std::runtime_error on failure. Used by the Buffer layer and CLI.
using BufferTransform = std::vector<uint8_t> (*)(const std::vector<uint8_t>&);

// BufferEncoder implements the Encoder interface over a single-shot BufferTransform.
// Input is accumulated in memory; finish() invokes the transform once.
class BufferEncoder final : public Encoder {
public:
    explicit BufferEncoder(BufferTransform transform);

    Result<std::size_t> process(ByteView in, MutableByteView out) override;
    Result<std::size_t> flush(MutableByteView out) override;
    Result<std::size_t> finish(MutableByteView out) override;
    void reset() noexcept override;
    State state() const noexcept override;

private:
    BufferTransform transform_;
    State state_;
    std::vector<uint8_t> input_;
    std::uint64_t total_input_;
};

class BufferDecoder final : public Decoder {
public:
    explicit BufferDecoder(BufferTransform transform);

    Result<std::size_t> process(ByteView in, MutableByteView out) override;
    Result<std::size_t> flush(MutableByteView out) override;
    Result<std::size_t> finish(MutableByteView out) override;
    void reset() noexcept override;
    State state() const noexcept override;

private:
    BufferTransform transform_;
    State state_;
    std::vector<uint8_t> input_;
    std::uint64_t total_input_;
};

// Convenience: run encoder/decoder over a complete input buffer with auto-growing output.
Result<std::vector<uint8_t>> encode_buffer(Encoder& encoder, const std::vector<uint8_t>& input);
Result<std::vector<uint8_t>> decode_buffer(Decoder& decoder, const std::vector<uint8_t>& input);

// CLI helpers: read input file, apply transform, write output file.
bool encode_file_via_buffer(BufferTransform transform, const std::string& input_path,
                            const std::string& output_path);
bool decode_file_via_buffer(BufferTransform transform, const std::string& input_path,
                            const std::string& output_path);

}  // namespace compresskit
