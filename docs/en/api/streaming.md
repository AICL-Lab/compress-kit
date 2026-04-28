# Streaming API

CompressKit now exposes the same two-layer in-memory API across C++17, Go, and Rust:

- A stateful streaming layer with `process`, `flush`, `finish`, and `reset`
- A stateless buffer layer that wraps the full lifecycle for one-shot byte-slice operations

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

## Go Example

```go
import (
    "github.com/LessUp/compress-kit/algorithms/shared/go/codec"
    "huffman"
)

func encode(data []byte) ([]byte, error) {
    return codec.EncodeBuffer(huffman.NewStreamingEncoder(), data)
}

func decode(encoded []byte) ([]byte, error) {
    return codec.DecodeBuffer(huffman.NewStreamingDecoder(), encoded)
}
```

## Rust Example

```rust
use compresskit_codec::codec::{decode_buffer, encode_buffer};
use huffman::{StreamingDecoder, StreamingEncoder};

fn roundtrip(input: &[u8]) -> Result<Vec<u8>, compresskit_codec::codec::CodecError> {
    let mut encoder = StreamingEncoder::new();
    let encoded = encode_buffer(&mut encoder, input)?;

    let mut decoder = StreamingDecoder::new();
    decode_buffer(&mut decoder, &encoded)
}
```

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

`make test` now runs shared streaming-layer tests for C++, Go, and Rust before the algorithm-specific suites.
