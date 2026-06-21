# Streaming API

CompressKit exposes a two-layer in-memory API in C++17:

- A stateful streaming layer with `process`, `flush`, `finish`, and `reset`
- A stateless buffer layer that wraps the full lifecycle for one-shot byte-slice operations

The streaming layer operates as an in-memory transform: input bytes are processed through the algorithm core and produced as output bytes in a single pass. The buffer layer is a thin convenience wrapper over the same streaming lifecycle.

## Lifecycle

The streaming layer follows a shared state machine:

| State | Allowed calls | Notes |
|-------|---------------|-------|
| `READY` | `process`, `flush`, `finish`, `reset` | `flush` is a no-op here |
| `STREAMING` | `process`, `flush`, `finish`, `reset` | algorithms may buffer input |
| `FLUSHING` | `process`, `flush`, `finish`, `reset` | `process` moves back to `STREAMING` |
| `FINISHED` | `reset` | all other calls return `ERR_INVALID_STATE` |
| `ERROR` | `reset` | error state is terminal until reset |

`finish()` always performs the final flush implicitly.

## Error Codes

| Code | Meaning |
|------|---------|
| `OK` | Success |
| `BUF_TOO_SMALL` | Caller output buffer was too small and state is unchanged |
| `ERR_TRUNCATED` | Input ended before a complete stream could be decoded |
| `ERR_CORRUPT` | Encoded data failed structural validation |
| `ERR_INVALID_STATE` | Call order violated the lifecycle contract |
| `ERR_SIZE_LIMIT` | Input exceeded 4 GiB or decoded output exceeded 1 GiB |
| `ERR_VERSION_UNSUPPORTED` | Reserved for future frame-layer validation |
| `ERR_UNKNOWN_ALGO` | Reserved for future frame-layer validation |

## Buffer Layer

The buffer layer is equivalent to:

```text
new encoder -> process(input) -> finish()
new decoder -> process(input) -> finish()
```

It exists to keep file-to-file paths and in-memory callers on the same implementation path.

## C++ Example

```cpp
#include <vector>

#include "compresskit/algorithms.hpp"

std::vector<uint8_t> encode(const std::vector<uint8_t>& input) {
    auto encoder = compresskit::make_huffman_encoder();
    auto result = compresskit::encode_buffer(encoder, input);
    return result.value;
}
```

## Verification

`make test` runs the shared streaming-layer tests for C++ before the algorithm-specific suites.
