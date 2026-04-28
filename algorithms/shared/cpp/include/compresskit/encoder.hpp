#pragma once

#include <cstddef>
#include <cstdint>

#include "compresskit/result.hpp"

namespace compresskit {

enum class State {
    READY = 0,
    STREAMING,
    FLUSHING,
    FINISHED,
    ERROR,
};

struct ByteView {
    const uint8_t* data = nullptr;
    std::size_t size = 0;
};

struct MutableByteView {
    uint8_t* data = nullptr;
    std::size_t size = 0;
};

class Encoder {
public:
    virtual ~Encoder() = default;

    virtual Result<std::size_t> process(ByteView in, MutableByteView out) = 0;
    virtual Result<std::size_t> flush(MutableByteView out) = 0;
    virtual Result<std::size_t> finish(MutableByteView out) = 0;
    virtual void reset() noexcept = 0;
    virtual State state() const noexcept = 0;
};

class Decoder {
public:
    virtual ~Decoder() = default;

    virtual Result<std::size_t> process(ByteView in, MutableByteView out) = 0;
    virtual Result<std::size_t> flush(MutableByteView out) = 0;
    virtual Result<std::size_t> finish(MutableByteView out) = 0;
    virtual void reset() noexcept = 0;
    virtual State state() const noexcept = 0;
};

}  // namespace compresskit
