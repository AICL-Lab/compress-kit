#include "compresskit/buffer_api.hpp"

#include <algorithm>
#include <fstream>
#include <limits>
#include <stdexcept>
#include <utility>

#include "compresskit/constants.hpp"

namespace compresskit {
namespace {

std::size_t encode_limit_for(std::size_t input_size) {
    if (input_size > (std::numeric_limits<std::size_t>::max() - INITIAL_ENCODE_OVERHEAD) /
                         static_cast<std::size_t>(BITS_PER_BYTE)) {
        throw std::overflow_error("encode limit overflow");
    }
    return input_size * static_cast<std::size_t>(BITS_PER_BYTE) + INITIAL_ENCODE_OVERHEAD;
}

bool write_file(const std::string& path, const std::vector<uint8_t>& data) {
    std::ofstream out(path, std::ios::binary);
    if (!out) {
        return false;
    }
    if (!data.empty()) {
        out.write(reinterpret_cast<const char*>(data.data()),
                  static_cast<std::streamsize>(data.size()));
    }
    return static_cast<bool>(out);
}

std::vector<uint8_t> read_file(const std::string& path) {
    std::ifstream in(path, std::ios::binary | std::ios::ate);
    if (!in) {
        throw std::runtime_error("cannot open input file");
    }
    std::streampos size = in.tellg();
    if (size < 0) {
        throw std::runtime_error("cannot determine input file size");
    }
    std::vector<uint8_t> data(static_cast<std::size_t>(size));
    in.seekg(0, std::ios::beg);
    if (!data.empty()) {
        in.read(reinterpret_cast<char*>(data.data()), static_cast<std::streamsize>(data.size()));
        if (!in) {
            throw std::runtime_error("failed to read input file");
        }
    }
    return data;
}

Result<std::size_t> invalid_state_result() {
    return {StatusCode::ERR_INVALID_STATE, 0};
}

template <typename Step>
Result<std::size_t> run_buffer_step(std::vector<uint8_t>& out, std::size_t total_written,
                                    std::size_t limit, Step step) {
    for (;;) {
        Result<std::size_t> result = step({out.data() + total_written, out.size() - total_written});
        if (result.status != StatusCode::BUF_TOO_SMALL) {
            if (!result.ok()) {
                return result;
            }
            return {StatusCode::OK, total_written + result.value};
        }

        total_written += result.value;
        if (total_written > limit || out.size() >= limit) {
            return {StatusCode::ERR_SIZE_LIMIT, 0};
        }

        out.resize(std::min(limit, std::max<std::size_t>(out.size() * 2, out.size() + 1)));
    }
}

}  // namespace

BufferEncoder::BufferEncoder(BufferTransform transform)
    : transform_(transform), state_(State::READY), total_input_(0) {}

Result<std::size_t> BufferEncoder::process(ByteView in, MutableByteView) {
    if (state_ == State::FINISHED || state_ == State::ERROR) {
        state_ = State::ERROR;
        return invalid_state_result();
    }
    if (total_input_ > MAX_INPUT_SIZE - in.size) {
        state_ = State::ERROR;
        return {StatusCode::ERR_SIZE_LIMIT, 0};
    }
    input_.insert(input_.end(), in.data, in.data + in.size);
    total_input_ += in.size;
    state_ = State::STREAMING;
    return {StatusCode::OK, 0};
}

Result<std::size_t> BufferEncoder::flush(MutableByteView) {
    if (state_ == State::FINISHED || state_ == State::ERROR) {
        state_ = State::ERROR;
        return invalid_state_result();
    }
    if (state_ == State::STREAMING) {
        state_ = State::FLUSHING;
    }
    return {StatusCode::OK, 0};
}

Result<std::size_t> BufferEncoder::finish(MutableByteView out) {
    if (state_ == State::FINISHED || state_ == State::ERROR) {
        state_ = State::ERROR;
        return invalid_state_result();
    }

    std::vector<uint8_t> encoded;
    try {
        encoded = transform_(input_);
    } catch (const std::exception&) {
        state_ = State::ERROR;
        return {StatusCode::ERR_CORRUPT, 0};
    }

    if (encoded.size() > out.size) {
        return {StatusCode::BUF_TOO_SMALL, 0};
    }
    if (!encoded.empty()) {
        std::copy(encoded.begin(), encoded.end(), out.data);
    }
    state_ = State::FINISHED;
    return {StatusCode::OK, encoded.size()};
}

void BufferEncoder::reset() noexcept {
    state_ = State::READY;
    input_.clear();
    total_input_ = 0;
}

State BufferEncoder::state() const noexcept {
    return state_;
}

BufferDecoder::BufferDecoder(BufferTransform transform)
    : transform_(transform), state_(State::READY), total_input_(0) {}

Result<std::size_t> BufferDecoder::process(ByteView in, MutableByteView) {
    if (state_ == State::FINISHED || state_ == State::ERROR) {
        state_ = State::ERROR;
        return invalid_state_result();
    }
    if (total_input_ > MAX_INPUT_SIZE - in.size) {
        state_ = State::ERROR;
        return {StatusCode::ERR_SIZE_LIMIT, 0};
    }
    input_.insert(input_.end(), in.data, in.data + in.size);
    total_input_ += in.size;
    state_ = State::STREAMING;
    return {StatusCode::OK, 0};
}

Result<std::size_t> BufferDecoder::flush(MutableByteView) {
    if (state_ == State::FINISHED || state_ == State::ERROR) {
        state_ = State::ERROR;
        return invalid_state_result();
    }
    if (state_ == State::STREAMING) {
        state_ = State::FLUSHING;
    }
    return {StatusCode::OK, 0};
}

Result<std::size_t> BufferDecoder::finish(MutableByteView out) {
    if (state_ == State::FINISHED || state_ == State::ERROR) {
        state_ = State::ERROR;
        return invalid_state_result();
    }

    std::vector<uint8_t> decoded;
    try {
        decoded = transform_(input_);
    } catch (const std::exception&) {
        state_ = State::ERROR;
        return {StatusCode::ERR_CORRUPT, 0};
    }

    if (decoded.size() > MAX_OUTPUT_SIZE) {
        state_ = State::ERROR;
        return {StatusCode::ERR_SIZE_LIMIT, 0};
    }
    if (decoded.size() > out.size) {
        return {StatusCode::BUF_TOO_SMALL, 0};
    }
    if (!decoded.empty()) {
        std::copy(decoded.begin(), decoded.end(), out.data);
    }
    state_ = State::FINISHED;
    return {StatusCode::OK, decoded.size()};
}

void BufferDecoder::reset() noexcept {
    state_ = State::READY;
    input_.clear();
    total_input_ = 0;
}

State BufferDecoder::state() const noexcept {
    return state_;
}

Result<std::vector<uint8_t>> encode_buffer(Encoder& encoder, const std::vector<uint8_t>& input) {
    if (input.size() > MAX_INPUT_SIZE) {
        return {StatusCode::ERR_SIZE_LIMIT, {}};
    }

    std::size_t limit = 0;
    try {
        limit = encode_limit_for(input.size());
    } catch (const std::overflow_error&) {
        return {StatusCode::ERR_SIZE_LIMIT, {}};
    }

    std::size_t initial_size = std::min(limit, input.size() * 2 + INITIAL_ENCODE_OVERHEAD);
    std::vector<uint8_t> out(initial_size);
    std::size_t total_written = 0;

    Result<std::size_t> result =
        run_buffer_step(out, total_written, limit, [&](MutableByteView out_view) {
            return encoder.process({input.data(), input.size()}, out_view);
        });
    if (!result.ok()) {
        return {result.status, {}};
    }
    total_written = result.value;

    result = run_buffer_step(out, total_written, limit,
                             [&](MutableByteView out_view) { return encoder.finish(out_view); });
    if (!result.ok()) {
        return {result.status, {}};
    }
    total_written = result.value;

    out.resize(total_written);
    return {StatusCode::OK, std::move(out)};
}

Result<std::vector<uint8_t>> decode_buffer(Decoder& decoder, const std::vector<uint8_t>& input) {
    if (input.size() > MAX_INPUT_SIZE) {
        return {StatusCode::ERR_SIZE_LIMIT, {}};
    }

    std::vector<uint8_t> out(input.size() + INITIAL_DECODE_OVERHEAD);
    std::size_t total_written = 0;

    Result<std::size_t> result =
        run_buffer_step(out, total_written, MAX_OUTPUT_SIZE, [&](MutableByteView out_view) {
            return decoder.process({input.data(), input.size()}, out_view);
        });
    if (!result.ok()) {
        return {result.status, {}};
    }
    total_written = result.value;

    result = run_buffer_step(out, total_written, MAX_OUTPUT_SIZE,
                             [&](MutableByteView out_view) { return decoder.finish(out_view); });
    if (!result.ok()) {
        return {result.status, {}};
    }
    total_written = result.value;

    out.resize(total_written);
    return {StatusCode::OK, std::move(out)};
}

bool encode_file_via_buffer(BufferTransform transform, const std::string& input_path,
                            const std::string& output_path) {
    try {
        std::vector<uint8_t> input = read_file(input_path);
        if (input.size() > MAX_INPUT_SIZE) {
            return false;
        }
        BufferEncoder encoder(transform);
        Result<std::vector<uint8_t>> encoded = encode_buffer(encoder, input);
        if (!encoded.ok()) {
            return false;
        }
        return write_file(output_path, encoded.value);
    } catch (const std::exception&) {
        return false;
    }
}

bool decode_file_via_buffer(BufferTransform transform, const std::string& input_path,
                            const std::string& output_path) {
    try {
        std::vector<uint8_t> input = read_file(input_path);
        BufferDecoder decoder(transform);
        Result<std::vector<uint8_t>> decoded = decode_buffer(decoder, input);
        if (!decoded.ok()) {
            return false;
        }
        return write_file(output_path, decoded.value);
    } catch (const std::exception&) {
        return false;
    }
}

}  // namespace compresskit
