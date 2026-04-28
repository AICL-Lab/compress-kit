# Design: add-streaming-api-foundation

## Architecture Overview

Two public layers sit above the existing file-to-file path:

```
┌─────────────────────────────────────────────┐
│  Caller                                      │
└───────────┬─────────────────────────────────┘
            │ buffer API  (whole-slice encode/decode)
            ▼
┌─────────────────────────────────────────────┐
│  Buffer Layer  (stateless, wraps streaming) │
└───────────┬─────────────────────────────────┘
            │ streaming API (process / flush / finish)
            ▼
┌─────────────────────────────────────────────┐
│  Streaming Layer  (stateful encoder/decoder)│
└───────────┬─────────────────────────────────┘
            │ internal byte-level codec
            ▼
┌─────────────────────────────────────────────┐
│  Existing codec core  (Huffman, AC, RC, RLE)│
└─────────────────────────────────────────────┘
```

The buffer layer is a thin shim: it calls `process(input)`, `flush()`, `finish()` in sequence and collects output.

## Lifecycle State Machine

```
         ┌──────────┐
  new()  │  READY   │◄──────────────────────────────┐
         └────┬─────┘                               │
              │ process(chunk)                      │
              ▼                                     │ reset()
         ┌──────────┐                               │
         │ STREAMING│──── process(chunk) ──────────►│
         └────┬─────┘                               │
              │ flush() (optional, repeatable)      │
              ▼                                     │
         ┌──────────┐                               │
         │ FLUSHING │──── flush() ─────────────────►│
         └────┬─────┘                               │
              │ finish()                            │
              ▼                                     │
         ┌──────────┐                               │
         │ FINISHED │                               │
         └────┬─────┘                               │
              │ reset()                             │
              └─────────────────────────────────────┘

         Any state  ──── error condition ──►  ERROR (terminal unless reset())
```

### State transition rules

| Current state | Call | Next state | Output produced |
|---------------|------|------------|-----------------|
| READY | `process(chunk)` | STREAMING | 0 or more bytes |
| READY | `flush()` | READY | 0 bytes (no-op) |
| STREAMING | `process(chunk)` | STREAMING | 0 or more bytes |
| STREAMING | `flush()` | FLUSHING | pending output flushed |
| FLUSHING | `flush()` | FLUSHING | 0 bytes (idempotent) |
| FLUSHING | `process(chunk)` | STREAMING | 0 or more bytes |
| STREAMING/FLUSHING | `finish()` | FINISHED | final bytes + trailer |
| FINISHED | any call except `reset()` | ERROR | — |
| ERROR | `reset()` | READY | — |

> **Note:** The `FLUSHING → STREAMING` path (via `process()`) is intentional — callers may interleave flush and process calls freely. The lifecycle diagram above does not show this arc explicitly; the table above is authoritative for all transitions.

## Partial Input / Partial Output Contract

### Partial input
- Implementations MAY buffer input internally; callers MUST NOT assume output is produced for every `process()` call.
- After `finish()`, all internally buffered input MUST be encoded and appended to output.

### Partial output
- Output buffers MUST be sized to `max_output_expansion(input_len)` which each algorithm documents.
- If the caller-supplied output buffer is too small, the call MUST return `BUF_TOO_SMALL` and the internal state MUST be unchanged (transactional).

### EOF / finish semantics
- `finish()` writes any algorithm-specific end-of-stream marker (e.g. EOF symbol for Huffman).
- Calling `decode.finish()` on a truncated stream MUST return `ERR_TRUNCATED`, not corrupt data.

## Language-Specific Mapping

### C++17

```cpp
// core/encoder.hpp  (sketch)
class Encoder {
public:
    virtual ~Encoder() = default;
    // Returns bytes written into out[0..out_len).
    // Returns BufTooSmall if out_len is insufficient.
    virtual Result process(std::span<const uint8_t> in,
                           std::span<uint8_t> out) = 0;
    virtual Result flush(std::span<uint8_t> out) = 0;
    virtual Result finish(std::span<uint8_t> out) = 0;
    virtual void   reset() noexcept = 0;
};
```

Callers own output buffers. `Result` carries `{bytes_written, status_code}`.

### Go

```go
// pkg/codec/encoder.go  (sketch)
type Encoder interface {
    Process(in []byte, out []byte) (written int, err error)
    Flush(out []byte) (written int, err error)
    Finish(out []byte) (written int, err error)
    Reset()
}
```

Wraps naturally with `io.Writer`: a thin `WriterEncoder` struct can implement `io.Writer.Write` via `Process`.

### Rust

```rust
// src/codec/encoder.rs  (sketch)
pub trait Encoder {
    type Error: std::error::Error;
    fn process(&mut self, input: &[u8], output: &mut [u8])
        -> Result<usize, Self::Error>;
    fn flush(&mut self, output: &mut [u8])
        -> Result<usize, Self::Error>;
    fn finish(&mut self, output: &mut [u8])
        -> Result<usize, Self::Error>;
    fn reset(&mut self);
}
```

Does not conflict with `std::io::Write`; a wrapper type can implement both.

## Buffer-Layer Helper

Each language SHALL expose a convenience function matching:

```
encode_buffer(algo, input_bytes) -> Result<output_bytes>
decode_buffer(algo, input_bytes) -> Result<output_bytes>
```

Implemented as: `new encoder → process(input) → finish() → collect output`.

## Error Code Catalogue

| Code | Meaning |
|------|---------|
| `OK` | Success |
| `BUF_TOO_SMALL` | Caller output buffer too small; state unchanged |
| `ERR_TRUNCATED` | Input stream ends prematurely during decode |
| `ERR_CORRUPT` | Checksum or structural integrity check failed |
| `ERR_INVALID_STATE` | Call not valid in current lifecycle state |
| `ERR_SIZE_LIMIT` | Input or output exceeds security limits |
| `ERR_VERSION_UNSUPPORTED` | Frame version byte is not `0x01`; canonically defined here and used by frame-layer validation |
| `ERR_UNKNOWN_ALGO` | `algo_id` in frame header does not map to a known algorithm; canonically defined here and used by frame-layer validation |

## Open Questions

1. Should `flush()` be mandatory before `finish()`, or should `finish()` auto-flush? **Decision: `finish()` auto-flushes** to simplify caller code.
2. Thread safety: are encoders thread-safe? **Decision: No. Single-threaded use only; callers manage concurrency.**
