## ADDED Requirements

### Requirement: REQ-ARCH-012 Shared Frame Header

All CompressKit compressed files SHALL begin with a 24-byte frame header.

#### Scenario: Magic byte validation
- **GIVEN** any CompressKit compressed file
- **WHEN** the first 4 bytes are read
- **THEN** they SHALL equal `0x43 0x4B 0x5A 0x4D` ("CKZM" in ASCII)
- **AND** a file not starting with this magic SHALL be rejected with an error

#### Scenario: Version field
- **GIVEN** a CompressKit file with version byte `0x01`
- **WHEN** a decoder reads the header
- **THEN** decoding SHALL proceed normally
- **GIVEN** a file with any version byte other than `0x01`
- **WHEN** a decoder reads the header
- **THEN** decoder SHALL reject with `ERR_VERSION_UNSUPPORTED`

#### Scenario: Algorithm ID
- **GIVEN** a frame header
- **WHEN** `algo_id` field is read
- **THEN** decoder SHALL select the corresponding algorithm from the registry
- **AND** an unrecognised `algo_id` SHALL cause rejection with `ERR_UNKNOWN_ALGO`

#### Scenario: Flags — unknown bits
- **GIVEN** a frame header with bits 3–15 of the flags field set to any non-zero value
- **WHEN** a decoder reads the header
- **THEN** decoder SHALL reject the frame (forward-compatibility guard)

### Requirement: REQ-ARCH-013 Content Size and Checksum Fields

The frame header SHALL carry optional content size and content integrity fields.

#### Scenario: content_size pre-allocation hint
- **GIVEN** a frame where `content_size` is non-zero
- **WHEN** a decoder reads the header
- **THEN** decoder MAY pre-allocate `content_size` bytes for output
- **AND** if decoded output length differs from `content_size`, decoder SHALL return `ERR_CORRUPT`

#### Scenario: content_size unknown
- **GIVEN** a frame where `content_size` is zero
- **WHEN** a decoder reads the header
- **THEN** decoder SHALL not pre-allocate and SHALL stream output without a size assumption

#### Scenario: Checksum validation
- **GIVEN** a frame with `FLAG_NO_CHECKSUM` NOT set
- **WHEN** decoding is complete
- **THEN** decoder SHALL compute xxHash-64 of the decoded output and compare to `checksum`
- **AND** mismatch SHALL return `ERR_CORRUPT`

### Requirement: REQ-ARCH-014 Extension Blocks and Frame Trailer

The frame format SHALL support extensibility via typed extension blocks and a frame trailer.

#### Scenario: Dictionary ID extension block
- **GIVEN** a frame with `FLAG_HAS_DICT` set
- **WHEN** the frame is parsed
- **THEN** a DICT extension block with an 8-byte `dict_id` SHALL follow the header

#### Scenario: Skippable metadata block
- **GIVEN** a frame with `FLAG_HAS_META` set
- **WHEN** the frame is parsed
- **THEN** one or more "Skip" extension blocks SHALL follow
- **AND** a decoder that does not understand a skippable block SHALL skip its payload and continue

#### Scenario: Frame trailer integrity
- **GIVEN** any CompressKit compressed file
- **WHEN** the last 8 bytes are read
- **THEN** bytes 0–3 SHALL equal `0x45 0x4E 0x44 0x00` ("END\0")
- **AND** bytes 4–7 SHALL be the CRC-32/ISO-HDLC of the compressed payload
- **AND** a CRC mismatch SHALL return `ERR_CORRUPT`
