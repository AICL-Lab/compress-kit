# Change Proposal: add-interoperability-conformance

## Summary

Establish a comprehensive conformance test suite for cross-language interoperability, covering header parsing, truncated/corrupted input, stream concatenation, a fixed corpus matrix, and a full cross-language encode/decode matrix.

## Motivation

The existing `REQ-TEST-001` checks only the happy path: encode in language A, decode in language B. It does not cover:
- Malformed or partially-written files (truncation, bit-flip corruption)
- Concatenated streams (two compressed files back-to-back)
- A reproducible, version-stable test corpus that all future benchmarks and regressions can share
- A structured matrix documenting which (language × algorithm) pairs have been verified

Without these tests, silent interoperability regressions can go undetected. This change creates the test infrastructure; it does not define benchmark thresholds (→ `add-benchmark-governance`).

## Scope

### In scope
- Header parsing test vectors (valid, invalid magic, bad version, unknown flags)
- Truncation test vectors (file cut at byte N for various N values)
- Corruption test vectors (single bit/byte flipped at various offsets)
- Concatenation test: two valid frames decoded sequentially
- Reproducible corpus matrix: file list, sizes, entropy categories
- Cross-language decode matrix: N×M table (encoder language × decoder language) for each algorithm

### Out of scope
- API lifecycle tests (→ `add-streaming-api-foundation`)
- Frame field semantics (→ `add-shared-frame-format`)
- Benchmark metrics / regression thresholds (→ `add-benchmark-governance`)

## Impact

| Spec | Change type |
|------|-------------|
| `cross-language-testing/spec.md` | ADDED requirements REQ-TEST-008, REQ-TEST-009, REQ-TEST-010, REQ-TEST-011 |

## Dependencies

- Assumes `add-shared-frame-format` is archived first (header fields must be stable before header parsing tests can be written).
- Assumes `add-streaming-api-foundation` is archived first (error codes such as `ERR_VERSION_UNSUPPORTED`, `ERR_UNKNOWN_ALGO`, `ERR_TRUNCATED`, and `ERR_CORRUPT` are defined in that change's error catalogue and must be stable before conformance tests can reference them).

## Risks

- Corpus files must be checked into the repository or generated deterministically; random generation makes regression reproduction harder.
- Truncation/corruption tests may depend on internal format details — keep them at the frame boundary level to remain algorithm-agnostic.
