#pragma once

#include <cstddef>

namespace compresskit {

enum class StatusCode {
    OK = 0,
    ERR_CORRUPT,
    ERR_SIZE_LIMIT,
};

template <typename T>
struct Result {
    StatusCode status = StatusCode::OK;
    T value{};

    bool ok() const noexcept { return status == StatusCode::OK; }
};

}  // namespace compresskit
