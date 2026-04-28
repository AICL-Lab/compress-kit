## ADDED Requirements

### Requirement: REQ-TEST-007 Frame Format Cross-Language Parsing Tests

The test suite SHALL verify that all three language implementations produce and parse identical frame headers.

#### Scenario: Header round-trip consistency
- **GIVEN** a frame header encoded by C++ implementation
- **WHEN** parsed by Go and Rust implementations
- **THEN** all fields (`magic`, `version`, `algo_id`, `flags`, `content_size`, `checksum`) SHALL be identical to the original

#### Scenario: Cross-language checksum agreement
- **GIVEN** the same uncompressed test corpus encoded by C++, Go, and Rust
- **WHEN** the `checksum` field is extracted from each output file
- **THEN** all three values SHALL be equal (same xxHash-64 of original content)

#### Scenario: Extension block interoperability
- **GIVEN** a file with a skippable metadata block encoded by one implementation
- **WHEN** decoded by a different language implementation
- **THEN** decoding SHALL succeed (skippable block skipped)
- **AND** decoded content SHALL be identical to original

#### Scenario: Corrupt magic rejection
- **GIVEN** a file with the magic bytes altered to `0x00000000`
- **WHEN** any implementation attempts to decode it
- **THEN** all three implementations SHALL reject the file with an error

#### Scenario: Trailer CRC validation
- **GIVEN** a valid CompressKit file with one byte of the compressed payload flipped
- **WHEN** any implementation decodes it
- **THEN** all three implementations SHALL return an error (trailer CRC or content checksum mismatch)
