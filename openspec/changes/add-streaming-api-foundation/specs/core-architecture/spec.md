# Core Architecture — Delta Spec (add-streaming-api-foundation)

<!-- DELTA TYPE: ADDED -->
<!-- These requirements are ADDED to openspec/specs/core-architecture/spec.md when this change is archived. -->

---

### REQ-ARCH-009: Streaming Encoder/Decoder Interface

All algorithm implementations SHALL expose a streaming interface with a defined lifecycle.

#### Scenario: Streaming encoder lifecycle
- **GIVEN** a freshly constructed encoder
- **WHEN** caller calls `process(chunk)` one or more times, then `finish()`
- **THEN** all input bytes SHALL be encoded and emitted across the `process` and `finish` calls in order
- **AND** state SHALL transition READY → STREAMING → FINISHED

#### Scenario: Flush produces stable output
- **GIVEN** an encoder in STREAMING state
- **WHEN** caller calls `flush()`
- **THEN** all currently buffered output SHALL be emitted
- **AND** state SHALL transition to FLUSHING
- **AND** a subsequent `process()` call SHALL transition back to STREAMING

#### Scenario: finish auto-flushes
- **GIVEN** an encoder in any non-FINISHED state
- **WHEN** caller calls `finish()`
- **THEN** remaining buffered bytes SHALL be flushed before the end-of-stream marker is written

#### Scenario: Reset from error
- **GIVEN** an encoder in ERROR state
- **WHEN** caller calls `reset()`
- **THEN** encoder SHALL return to READY state with no residual data

---

### REQ-ARCH-010: Buffer-Layer Convenience API

All implementations SHALL expose a stateless buffer-layer API wrapping the streaming interface.

#### Scenario: Buffer encode
- **GIVEN** a byte slice as input
- **WHEN** caller invokes `encode_buffer(algo, input)`
- **THEN** the function SHALL return the complete encoded output as a new byte slice
- **AND** the call SHALL be equivalent to `new encoder → process(input) → finish()`

#### Scenario: Buffer decode
- **GIVEN** a valid encoded byte slice
- **WHEN** caller invokes `decode_buffer(algo, input)`
- **THEN** function SHALL return fully decoded output
- **AND** an invalid or truncated input SHALL return `ERR_TRUNCATED` or `ERR_CORRUPT`

---

### REQ-ARCH-011: Output Buffer Contract

Implementations SHALL define and document maximum output expansion ratios.

#### Scenario: BUF_TOO_SMALL is transactional
- **GIVEN** an output buffer smaller than required
- **WHEN** `process()`, `flush()`, or `finish()` is called
- **THEN** the call SHALL return `BUF_TOO_SMALL`
- **AND** internal encoder state SHALL be unchanged (caller may retry with larger buffer)

#### Scenario: Security limits enforced at streaming boundary — output
- **GIVEN** a streaming decoder accumulating output
- **WHEN** cumulative decoded output would exceed 1 GiB
- **THEN** decoder SHALL return `ERR_SIZE_LIMIT` and enter ERROR state

#### Scenario: Security limits enforced at streaming boundary — input
- **GIVEN** a streaming encoder or decoder receiving input chunks
- **WHEN** cumulative input bytes would exceed 4 GiB
- **THEN** the implementation SHALL return `ERR_SIZE_LIMIT` and enter ERROR state
