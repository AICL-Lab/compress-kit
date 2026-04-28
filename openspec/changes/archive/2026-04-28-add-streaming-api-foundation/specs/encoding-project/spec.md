## ADDED Requirements

### Requirement: REQ-PROD-009 Streaming and Buffer API Availability

All algorithm implementations SHALL provide both a streaming and a buffer-oriented public API.

#### Scenario: Streaming API availability
- **GIVEN** any CompressKit language implementation (C++17, Go, Rust)
- **WHEN** a caller imports the library
- **THEN** streaming `Encoder` and `Decoder` interfaces SHALL be available for all four algorithms

#### Scenario: Buffer API convenience
- **GIVEN** a caller that wants to compress an in-memory byte slice
- **WHEN** the caller invokes `encode_buffer(algo, data)`
- **THEN** the function SHALL return the compressed output without requiring file I/O

#### Scenario: Lifecycle error isolation
- **GIVEN** a caller that misuses the API (e.g. calls `process` after `finish`)
- **WHEN** the invalid call is made
- **THEN** the implementation SHALL return `ERR_INVALID_STATE`
- **AND** SHALL NOT crash, corrupt memory, or produce silent data loss
