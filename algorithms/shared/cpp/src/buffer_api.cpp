#include "compresskit/buffer_api.hpp"

#include <unistd.h>

#include <algorithm>
#include <cstdio>
#include <cstdlib>
#include <fstream>
#include <limits>
#include <stdexcept>
#include <string>
#include <utility>

namespace compresskit {
namespace {

constexpr std::uint64_t kMaxInputSize = 4ULL * 1024 * 1024 * 1024;
constexpr std::uint64_t kMaxOutputSize = 1ULL * 1024 * 1024 * 1024;
constexpr std::size_t kInitialEncodeOverhead = 2048;

class ScopedTempFile {
   public:
    explicit ScopedTempFile(const char* prefix) {
        // Use platform-appropriate temp directory
        const char* tmp_dir = std::getenv("TMPDIR");
        if (!tmp_dir)
            tmp_dir = std::getenv("TMP");
        if (!tmp_dir)
            tmp_dir = std::getenv("TEMP");
        if (!tmp_dir)
            tmp_dir = "/tmp";
        std::string pattern = std::string(tmp_dir) + "/" + prefix + "-XXXXXX";
        std::vector<char> buffer(pattern.begin(), pattern.end());
        buffer.push_back('\0');
        int fd = mkstemp(buffer.data());
        if (fd < 0) {
            throw std::runtime_error("failed to create temp file");
        }
        close(fd);
        path_ = buffer.data();
    }

    ~ScopedTempFile() {
        if (!path_.empty()) {
            std::remove(path_.c_str());
        }
    }

    const std::string& path() const noexcept { return path_; }

   private:
    std::string path_;
};

std::size_t encode_limit_for(std::size_t input_size) {
    if (input_size > (std::numeric_limits<std::size_t>::max() - kInitialEncodeOverhead) / 8) {
        throw std::overflow_error("encode limit overflow");
    }
    return input_size * 8 + kInitialEncodeOverhead;
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

Result<std::vector<uint8_t>> read_file(const std::string& path, bool enforce_output_limit) {
    std::ifstream in(path, std::ios::binary | std::ios::ate);
    if (!in) {
        return {StatusCode::ERR_CORRUPT, {}};
    }
    std::streampos size = in.tellg();
    if (size < 0) {
        return {StatusCode::ERR_CORRUPT, {}};
    }
    if (enforce_output_limit && static_cast<std::uint64_t>(size) > kMaxOutputSize) {
        return {StatusCode::ERR_SIZE_LIMIT, {}};
    }
    std::vector<uint8_t> data(static_cast<std::size_t>(size));
    in.seekg(0, std::ios::beg);
    if (!data.empty()) {
        in.read(reinterpret_cast<char*>(data.data()), static_cast<std::streamsize>(data.size()));
        if (!in) {
            return {StatusCode::ERR_CORRUPT, {}};
        }
    }
    return {StatusCode::OK, std::move(data)};
}

Result<std::vector<uint8_t>> run_transform(FileTransform transform,
                                           const std::vector<uint8_t>& input,
                                           bool enforce_output_limit) {
    try {
        ScopedTempFile input_file("compresskit-in");
        ScopedTempFile output_file("compresskit-out");
        if (!write_file(input_file.path(), input)) {
            return {StatusCode::ERR_CORRUPT, {}};
        }
        if (!transform(input_file.path(), output_file.path())) {
            Result<std::vector<uint8_t>> maybe_output =
                read_file(output_file.path(), enforce_output_limit);
            if (maybe_output.status == StatusCode::ERR_SIZE_LIMIT) {
                return maybe_output;
            }
            return {StatusCode::ERR_CORRUPT, {}};
        }
        return read_file(output_file.path(), enforce_output_limit);
    } catch (const std::exception&) {
        return {StatusCode::ERR_CORRUPT, {}};
    }
}

Result<std::size_t> invalid_state_result() {
    return {StatusCode::ERR_INVALID_STATE, 0};
}

}  // namespace

BufferEncoder::BufferEncoder(FileTransform transform)
    : transform_(transform), state_(State::READY), total_input_(0) {}

Result<std::size_t> BufferEncoder::process(ByteView in, MutableByteView) {
    if (state_ == State::FINISHED) {
        state_ = State::ERROR;
        return invalid_state_result();
    }
    if (state_ == State::ERROR) {
        return invalid_state_result();
    }
    if (total_input_ > kMaxInputSize - in.size) {
        state_ = State::ERROR;
        return {StatusCode::ERR_SIZE_LIMIT, 0};
    }
    input_.insert(input_.end(), in.data, in.data + in.size);
    total_input_ += in.size;
    state_ = State::STREAMING;
    return {StatusCode::OK, 0};
}

Result<std::size_t> BufferEncoder::flush(MutableByteView) {
    if (state_ == State::FINISHED) {
        state_ = State::ERROR;
        return invalid_state_result();
    }
    if (state_ == State::ERROR) {
        return invalid_state_result();
    }
    if (state_ == State::STREAMING) {
        state_ = State::FLUSHING;
    }
    return {StatusCode::OK, 0};
}

Result<std::size_t> BufferEncoder::finish(MutableByteView out) {
    if (state_ == State::FINISHED) {
        state_ = State::ERROR;
        return invalid_state_result();
    }
    if (state_ == State::ERROR) {
        return invalid_state_result();
    }

    Result<std::vector<uint8_t>> encoded = run_transform(transform_, input_, false);
    if (!encoded.ok()) {
        state_ = State::ERROR;
        return {encoded.status, 0};
    }
    if (encoded.value.size() > out.size) {
        return {StatusCode::BUF_TOO_SMALL, 0};
    }
    if (!encoded.value.empty()) {
        std::copy(encoded.value.begin(), encoded.value.end(), out.data);
    }
    state_ = State::FINISHED;
    return {StatusCode::OK, encoded.value.size()};
}

void BufferEncoder::reset() noexcept {
    state_ = State::READY;
    input_.clear();
    total_input_ = 0;
}

State BufferEncoder::state() const noexcept {
    return state_;
}

BufferDecoder::BufferDecoder(FileTransform transform)
    : transform_(transform), state_(State::READY), total_input_(0) {}

Result<std::size_t> BufferDecoder::process(ByteView in, MutableByteView) {
    if (state_ == State::FINISHED) {
        state_ = State::ERROR;
        return invalid_state_result();
    }
    if (state_ == State::ERROR) {
        return invalid_state_result();
    }
    if (total_input_ > kMaxInputSize - in.size) {
        state_ = State::ERROR;
        return {StatusCode::ERR_SIZE_LIMIT, 0};
    }
    input_.insert(input_.end(), in.data, in.data + in.size);
    total_input_ += in.size;
    state_ = State::STREAMING;
    return {StatusCode::OK, 0};
}

Result<std::size_t> BufferDecoder::flush(MutableByteView) {
    if (state_ == State::FINISHED) {
        state_ = State::ERROR;
        return invalid_state_result();
    }
    if (state_ == State::ERROR) {
        return invalid_state_result();
    }
    if (state_ == State::STREAMING) {
        state_ = State::FLUSHING;
    }
    return {StatusCode::OK, 0};
}

Result<std::size_t> BufferDecoder::finish(MutableByteView out) {
    if (state_ == State::FINISHED) {
        state_ = State::ERROR;
        return invalid_state_result();
    }
    if (state_ == State::ERROR) {
        return invalid_state_result();
    }

    Result<std::vector<uint8_t>> decoded = run_transform(transform_, input_, true);
    if (!decoded.ok()) {
        state_ = State::ERROR;
        return {decoded.status, 0};
    }
    if (decoded.value.size() > out.size) {
        return {StatusCode::BUF_TOO_SMALL, 0};
    }
    if (!decoded.value.empty()) {
        std::copy(decoded.value.begin(), decoded.value.end(), out.data);
    }
    state_ = State::FINISHED;
    return {StatusCode::OK, decoded.value.size()};
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
    if (input.size() > kMaxInputSize) {
        return {StatusCode::ERR_SIZE_LIMIT, {}};
    }

    std::size_t limit = 0;
    try {
        limit = encode_limit_for(input.size());
    } catch (const std::overflow_error&) {
        return {StatusCode::ERR_SIZE_LIMIT, {}};
    }

    std::size_t initial_size = std::min(limit, input.size() * 2 + kInitialEncodeOverhead);
    std::vector<uint8_t> out(initial_size);
    std::size_t total_written = 0;

    for (;;) {
        Result<std::size_t> result = encoder.process(
            {input.data(), input.size()}, {out.data() + total_written, out.size() - total_written});
        if (result.status != StatusCode::BUF_TOO_SMALL) {
            if (!result.ok()) {
                return {result.status, {}};
            }
            total_written += result.value;
            break;
        }
        total_written += result.value;
        if (total_written > limit || out.size() >= limit) {
            return {StatusCode::ERR_SIZE_LIMIT, {}};
        }
        out.resize(std::min(limit, std::max<std::size_t>(out.size() * 2, out.size() + 1)));
    }

    for (;;) {
        Result<std::size_t> result =
            encoder.finish({out.data() + total_written, out.size() - total_written});
        if (result.status != StatusCode::BUF_TOO_SMALL) {
            if (!result.ok()) {
                return {result.status, {}};
            }
            total_written += result.value;
            break;
        }
        if (out.size() >= limit) {
            return {StatusCode::ERR_SIZE_LIMIT, {}};
        }
        out.resize(std::min(limit, std::max<std::size_t>(out.size() * 2, out.size() + 1)));
    }

    if (total_written > limit) {
        return {StatusCode::ERR_SIZE_LIMIT, {}};
    }
    out.resize(total_written);
    return {StatusCode::OK, std::move(out)};
}

Result<std::vector<uint8_t>> decode_buffer(Decoder& decoder, const std::vector<uint8_t>& input) {
    if (input.size() > kMaxInputSize) {
        return {StatusCode::ERR_SIZE_LIMIT, {}};
    }

    std::vector<uint8_t> out(input.size() + 1024);
    std::size_t total_written = 0;

    for (;;) {
        Result<std::size_t> result = decoder.process(
            {input.data(), input.size()}, {out.data() + total_written, out.size() - total_written});
        if (result.status != StatusCode::BUF_TOO_SMALL) {
            if (!result.ok()) {
                return {result.status, {}};
            }
            total_written += result.value;
            break;
        }
        total_written += result.value;
        if (total_written > kMaxOutputSize || out.size() >= kMaxOutputSize) {
            return {StatusCode::ERR_SIZE_LIMIT, {}};
        }
        out.resize(std::min<std::size_t>(kMaxOutputSize,
                                         std::max<std::size_t>(out.size() * 2, out.size() + 1)));
    }

    for (;;) {
        Result<std::size_t> result =
            decoder.finish({out.data() + total_written, out.size() - total_written});
        if (result.status != StatusCode::BUF_TOO_SMALL) {
            if (!result.ok()) {
                return {result.status, {}};
            }
            total_written += result.value;
            break;
        }
        if (out.size() >= kMaxOutputSize) {
            return {StatusCode::ERR_SIZE_LIMIT, {}};
        }
        out.resize(std::min<std::size_t>(kMaxOutputSize,
                                         std::max<std::size_t>(out.size() * 2, out.size() + 1)));
    }

    if (total_written > kMaxOutputSize) {
        return {StatusCode::ERR_SIZE_LIMIT, {}};
    }
    out.resize(total_written);
    return {StatusCode::OK, std::move(out)};
}

bool encode_file_via_buffer(FileTransform transform, const std::string& input_path,
                            const std::string& output_path) {
    Result<std::vector<uint8_t>> input = read_file(input_path, false);
    if (!input.ok()) {
        return false;
    }
    BufferEncoder encoder(transform);
    Result<std::vector<uint8_t>> encoded = encode_buffer(encoder, input.value);
    if (!encoded.ok()) {
        return false;
    }
    return write_file(output_path, encoded.value);
}

bool decode_file_via_buffer(FileTransform transform, const std::string& input_path,
                            const std::string& output_path) {
    Result<std::vector<uint8_t>> input = read_file(input_path, false);
    if (!input.ok()) {
        return false;
    }
    BufferDecoder decoder(transform);
    Result<std::vector<uint8_t>> decoded = decode_buffer(decoder, input.value);
    if (!decoded.ok()) {
        return false;
    }
    return write_file(output_path, decoded.value);
}

}  // namespace compresskit
