#pragma once

#include <cstddef>

namespace compresskit {

enum class StatusCode {
    OK = 0,
    BUF_TOO_SMALL,
    ERR_TRUNCATED,
    ERR_CORRUPT,
    ERR_INVALID_STATE,
    ERR_SIZE_LIMIT,
    ERR_VERSION_UNSUPPORTED,
    ERR_UNKNOWN_ALGO,
};

template <typename T>
struct Result {
    StatusCode status = StatusCode::OK;
    T value{};

    bool ok() const noexcept { return status == StatusCode::OK; }
};

}  // namespace compresskit
