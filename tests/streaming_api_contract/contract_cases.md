# Streaming API Contract Cases

This document defines the canonical test cases for the streaming API contract.
All language implementations must pass these cases.

**Implementation:** The Go implementation of these contract tests is located at:
`algorithms/shared/go/codec/lifecycle_test.go`

Each test case below includes a reference to its corresponding test function.

## Lifecycle Test Cases

### L1. Empty input
- **GIVEN** a new encoder in READY state
- **WHEN** `finish()` is called without any `process()` calls
- **THEN** output SHALL contain valid end-of-stream marker
- **AND** state SHALL be FINISHED
- **Test:** `TestLifecycle_EmptyInput`

### L2. Single-byte input
- **GIVEN** a new encoder
- **WHEN** `process([0x42])` then `finish()` are called
- **THEN** output SHALL be decodable back to `[0x42]`
- **Test:** `TestLifecycle_SingleByte`

### L3. Chunked input where output is delayed
- **GIVEN** a new encoder
- **WHEN** `process(chunk1)`, `process(chunk2)`, `process(chunk3)`, then `finish()` are called
- **THEN** concatenated output SHALL decode to `chunk1 + chunk2 + chunk3`
- **AND** implementations MAY buffer input (output may be delayed until flush/finish)
- **Test:** `TestLifecycle_ChunkedInput`

### L4. Flush without finish
- **GIVEN** an encoder in STREAMING state (after at least one `process()` call)
- **WHEN** `flush()` is called
- **THEN** state SHALL transition to FLUSHING
- **AND** all buffered output SHALL be written
- **AND** calling `process()` afterward SHALL transition back to STREAMING
- **Test:** `TestLifecycle_FlushWithoutFinish`

### L5. Finish after multiple process calls
- **GIVEN** an encoder that has called `process()` multiple times
- **WHEN** `finish()` is called
- **THEN** final output SHALL include end-of-stream marker
- **AND** state SHALL be FINISHED
- **AND** subsequent calls to `process()`, `flush()`, or `finish()` SHALL return `ERR_INVALID_STATE`
- **Test:** `TestLifecycle_FinishAfterMultipleProcess`

### L6. Reset after finish
- **GIVEN** an encoder in FINISHED state
- **WHEN** `reset()` is called
- **THEN** state SHALL return to READY
- **AND** encoder can be reused for new input
- **Test:** `TestLifecycle_ResetAfterFinish`

### L7. Reset after error
- **GIVEN** an encoder in ERROR state (e.g., after buffer too small error resolved)
- **WHEN** `reset()` is called
- **THEN** state SHALL return to READY with no residual data
- **Test:** `TestLifecycle_ResetAfterError`

## Buffer Contract Test Cases

### B1. BUF_TOO_SMALL is transactional
- **GIVEN** an encoder and an output buffer smaller than `max_output_expansion(input_len)`
- **WHEN** `process(input, small_buffer)` is called
- **THEN** call SHALL return `BUF_TOO_SMALL` error
- **AND** internal state SHALL be unchanged
- **AND** retrying with a larger buffer SHALL succeed
- **Test:** `TestBuffer_BufTooSmallTransactional`

### B2. Buffer encode full path
- **GIVEN** input bytes
- **WHEN** `encode_buffer(algo, input)` is called
- **THEN** output SHALL be equivalent to `new encoder → process(input) → finish()`
- **Test:** `TestBuffer_EncodeFullPath`

### B3. Buffer decode full path
- **GIVEN** valid encoded bytes
- **WHEN** `decode_buffer(algo, encoded)` is called
- **THEN** output SHALL match original input
- **Test:** `TestBuffer_DecodeFullPath`

## Error Handling Test Cases

### E1. Truncated frame on decode
- **GIVEN** an encoded stream with bytes removed from the end
- **WHEN** `decode()` then `finish()` are called
- **THEN** `finish()` SHALL return `ERR_TRUNCATED`
- **Test:** `TestError_TruncatedFrame`

### E2. Corrupt data on decode
- **GIVEN** an encoded stream with data corruption (e.g., flipped bits in payload)
- **WHEN** `decode()` processes the stream
- **THEN** decoder SHALL detect corruption and return `ERR_CORRUPT` or `ERR_TRUNCATED`
- **Test:** `TestError_TruncatedFrame` (covers both truncation and corruption)

### E3. Input size limit (4 GiB)
- **GIVEN** an encoder receiving input chunks
- **WHEN** cumulative input exceeds 4 GiB
- **THEN** encoder SHALL return `ERR_SIZE_LIMIT` and enter ERROR state
- **Test:** `TestConstants` (verifies limit constant)

### E4. Output size limit (1 GiB decode)
- **GIVEN** a decoder processing malicious input that expands unboundedly
- **WHEN** cumulative decoded output would exceed 1 GiB
- **THEN** decoder SHALL return `ERR_SIZE_LIMIT` and enter ERROR state
- **Test:** `TestConstants` (verifies limit constant)

### E5. Invalid state transitions
- **GIVEN** an encoder in FINISHED state
- **WHEN** `process()` is called without `reset()`
- **THEN** call SHALL return `ERR_INVALID_STATE`
- **Test:** Covered in `TestLifecycle_FinishAfterMultipleProcess`

## Cross-Algorithm Test Cases

### X1. All algorithms support full lifecycle
- **GIVEN** each of the four algorithms (Huffman, Arithmetic, Range, RLE)
- **WHEN** lifecycle cases L1-L7 are executed
- **THEN** all SHALL pass
- **Test:** All lifecycle tests run for all 4 algorithms

### X2. Buffer API works for all algorithms
- **GIVEN** each algorithm
- **WHEN** buffer encode/decode is used
- **THEN** results SHALL match streaming API results
- **Test:** All buffer tests run for all 4 algorithms
