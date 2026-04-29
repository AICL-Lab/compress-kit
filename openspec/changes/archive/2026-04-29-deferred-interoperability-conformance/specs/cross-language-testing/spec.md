## ADDED Requirements

### Requirement: REQ-TEST-008 Header Parsing Conformance

All implementations SHALL reject malformed frame headers with appropriate errors.

#### Scenario: Bad magic rejection
- **GIVEN** a byte stream where the first 4 bytes do not equal `0x434B5A4D`
- **WHEN** any implementation attempts to decode it
- **THEN** decoder SHALL return an error and produce no output

#### Scenario: Unsupported version rejection
- **GIVEN** a byte stream with a valid magic but version byte != `0x01`
- **WHEN** any implementation attempts to decode it
- **THEN** decoder SHALL return `ERR_VERSION_UNSUPPORTED` and produce no output

#### Scenario: Unknown flags rejection
- **GIVEN** a byte stream with reserved flag bits (bits 3–15) set to non-zero
- **WHEN** any implementation attempts to decode it
- **THEN** decoder SHALL reject the frame

#### Scenario: Header-truncated stream
- **GIVEN** a byte stream that ends before the full 24-byte header is read
- **WHEN** any implementation attempts to decode it
- **THEN** decoder SHALL return `ERR_TRUNCATED` and produce no output

### Requirement: REQ-TEST-009 Truncation and Corruption Robustness

All implementations SHALL handle truncated and corrupted inputs without crashing or producing silent garbage output.

#### Scenario: Mid-payload truncation
- **GIVEN** a valid CompressKit file truncated at any byte offset within the compressed payload
- **WHEN** any implementation decodes the truncated file
- **THEN** decoder SHALL return an error (non-zero exit or error code)
- **AND** SHALL NOT hang beyond a 10-second timeout
- **AND** SHALL NOT produce a partial output file presented as successful

#### Scenario: Single-byte corruption
- **GIVEN** a valid CompressKit file with exactly one byte flipped (XOR 0xFF)
- **WHEN** any implementation decodes the corrupted file
- **THEN** decoder SHALL return an error
- **AND** SHALL NOT silently output data that differs from the original

### Requirement: REQ-TEST-010 Stream Concatenation

Implementations SHALL correctly decode concatenated CompressKit frames.

#### Scenario: Two-frame concatenation
- **GIVEN** two independently valid CompressKit frames concatenated into a single byte stream
- **WHEN** a compliant decoder processes the concatenated stream
- **THEN** output SHALL equal the concatenation of the decoded contents of each frame in order

### Requirement: REQ-TEST-011 Reproducible Corpus and Decode Matrix

The test suite SHALL use a fixed corpus and document cross-language decode coverage.

#### Scenario: Corpus reproducibility
- **GIVEN** `tests/gen_corpus.py --seed 42`
- **WHEN** run on any platform
- **THEN** it SHALL produce byte-identical corpus files matching `tests/corpus/MANIFEST.sha256`

#### Scenario: Full decode matrix coverage
- **GIVEN** the cross-language decode matrix test
- **WHEN** run against all corpus files for all algorithms
- **THEN** all 9 encoder/decoder language pairs SHALL produce PASS for each algorithm
- **EXCEPT** Range Coder corpus files exceeding 100 KiB which SHALL be skipped per REQ-TEST-004 (Range Coder has a known decode performance issue for files >500 KiB; the 100 KiB cap is a conservative safety margin)

#### Scenario: Matrix results are committed
- **GIVEN** a successful CI run
- **WHEN** the decode matrix test completes
- **THEN** results SHALL be persisted to `tests/results/` for regression comparison
