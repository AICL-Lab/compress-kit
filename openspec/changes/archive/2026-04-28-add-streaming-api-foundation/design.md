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

## Algorithm Lifecycle Review (Task A1)

All four algorithms (Huffman, Arithmetic Coding, Range Coder, RLE) have been reviewed against the lifecycle state machine. Findings:

### Huffman
- **Current implementation**: Buffers entire input, builds frequency table, then encodes in one pass
- **Lifecycle compatibility**: ✓ Compatible. Can buffer in STREAMING state, emit on finish()
- **State machine gaps**: None. Frequency table must be built before encoding begins; buffering all input is acceptable
- **End-of-stream marker**: EOF symbol (index 256) already present in current implementation

### Arithmetic Coding
- **Current implementation**: Buffers entire input for frequency analysis, then encodes with range updates
- **Lifecycle compatibility**: ✓ Compatible. Internal state (low, high, pendingBits) already tracked
- **State machine gaps**: None. `Finish()` method already exists and handles final bit emission
- **End-of-stream marker**: EOF symbol (index 256) already encoded

### Range Coder
- **Current implementation**: []byte-based, buffers entire input and output
- **Lifecycle compatibility**: ✓ Compatible. Encoder already maintains `low`, `high` state
- **State machine gaps**: None. Renormalization logic already supports incremental emission
- **End-of-stream marker**: EOF symbol already present

### RLE (Run-Length Encoding)
- **Current implementation**: Streaming-friendly — processes byte-by-byte, emits (count, value) pairs
- **Lifecycle compatibility**: ✓ Fully compatible. Already incremental
- **State machine gaps**: None. Current run must be tracked across process() calls
- **End-of-stream marker**: Not needed — decoder stops at EOF naturally

**Conclusion**: All four algorithms are lifecycle-compatible. No design.md changes required for state machine correctness.

## Maximum Output Expansion Formulas (Task A2)

Each algorithm documents its worst-case output expansion to support `BUF_TOO_SMALL` contract and buffer sizing.

### Huffman
**Formula**: `max_output = 4 + 257*4 + input_len*8 + ceil(input_len/8)`
- Header: 4 bytes (magic "HFMN")
- Frequency table: 257 symbols × 4 bytes = 1028 bytes
- Encoded data: worst case = 1 bit per input byte × input_len, padded to byte boundary
- **Conservative upper bound**: `1032 + input_len*8 + 8 bytes`
- **Simplified**: `max_output = 1040 + input_len*8`

### Arithmetic Coding
**Formula**: `max_output = 4 + 257*4 + input_len*4 + 32`
- Header: 4 bytes (magic "AENC")
- Frequency table: 1028 bytes
- Encoded data: worst case ≈ 4 bytes per input byte (when range narrows maximally)
- Final flush: 32 bits (4 bytes)
- **Simplified**: `max_output = 1064 + input_len*4`

### Range Coder
**Formula**: `max_output = 4 + 257*4 + input_len*4 + 4`
- Header: 4 bytes (magic "RCNC")
- Frequency table: 1028 bytes
- Encoded data: similar to arithmetic coding, ≈ 4 bytes per input byte worst case
- Final flush: 4 bytes
- **Simplified**: `max_output = 1036 + input_len*4`

### RLE (Run-Length Encoding)
**Formula**: `max_output = input_len*5`
- Each input byte becomes: 4 bytes (count=1 as uint32) + 1 byte (value)
- No header overhead
- **Worst case**: Every byte is different, no runs
- **Simplified**: `max_output = input_len*5`

**Note**: These formulas represent absolute worst-case scenarios. Typical compression achieves much better ratios, but callers MUST allocate using these bounds to avoid `BUF_TOO_SMALL` errors.

## Go Interface io.Reader/io.Writer Compatibility Review (Task A3)

The Go interface design has been reviewed for standard library compatibility:

### Encoder Interface
```go
type Encoder interface {
    Process(in []byte, out []byte) (written int, err error)
    Flush(out []byte) (written int, err error)
    Finish(out []byte) (written int, err error)
    Reset()
}
```

### io.Writer Wrapping Strategy

A `WriterEncoder` adapter can implement `io.Writer` by:
1. Accepting a backing `Encoder` instance
2. Maintaining an internal output buffer
3. Mapping `Write(p []byte)` → `encoder.Process(p, outBuf)` + flush to underlying io.Writer
4. Exposing `Close()` to call `Finish()` and flush final output

**Example usage**:
```go
encoder := huffman.NewEncoder()
w := codec.NewWriterEncoder(encoder, outputWriter)
io.Copy(w, inputReader)  // Standard io.Copy works
w.Close()                 // Calls Finish() and flushes
```

### io.Reader Wrapping Strategy

A `ReaderDecoder` can wrap `Decoder`:
1. Read from source via `Read(p []byte)`
2. Call `decoder.Process(chunk, outBuf)` as input arrives
3. Return decoded bytes from internal buffer via `Read(p []byte)`
4. Handle `Finish()` when source EOF is reached

**Compatibility**: ✓ The interface design is fully compatible with Go's io.Reader/io.Writer patterns. No design.md changes needed.
