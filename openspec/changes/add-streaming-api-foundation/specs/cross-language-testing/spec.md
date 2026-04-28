## ADDED Requirements

### Requirement: REQ-TEST-006 Streaming API Conformance Tests

The test suite SHALL verify streaming lifecycle compliance for all implementations.

#### Scenario: Single-chunk round-trip via streaming API
- **GIVEN** a test corpus file
- **WHEN** encoded with streaming API (single `process` + `finish`) and decoded with streaming API
- **THEN** decoded output SHALL be byte-identical to original

#### Scenario: Multi-chunk round-trip
- **GIVEN** a test corpus file split into 4 KiB chunks
- **WHEN** each chunk is fed via `process()` and `finish()` is called at end
- **THEN** decoded output SHALL be byte-identical to original

#### Scenario: Buffer API equals streaming API output
- **GIVEN** the same input bytes
- **WHEN** encoded via buffer API vs. streaming API (single process + finish)
- **THEN** both SHALL produce identical byte output

#### Scenario: ERR_INVALID_STATE on post-finish process
- **GIVEN** an encoder that has called `finish()`
- **WHEN** caller calls `process()` again without `reset()`
- **THEN** implementation SHALL return `ERR_INVALID_STATE`

#### Scenario: BUF_TOO_SMALL leaves state unchanged
- **GIVEN** an encoder in STREAMING state with a 1-byte output buffer
- **WHEN** `process()` is called with non-empty input
- **THEN** if BUF_TOO_SMALL is returned, a subsequent call with adequate buffer SHALL succeed
