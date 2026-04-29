# Streaming API Contract

This directory documents the streaming API contract that is exercised by the
language-level test suites.

## Purpose

The streaming API contract establishes:
- A standardized `process() / flush() / finish()` lifecycle for all algorithms
- State machine guarantees (READY → STREAMING → FLUSHING → FINISHED)
- Error handling semantics (BUF_TOO_SMALL, ERR_TRUNCATED, ERR_CORRUPT, etc.)
- Security constraints (4 GiB input limit, 1 GiB output limit)
- Buffer-layer convenience API guarantees

## Scope

**In scope:**
- Lifecycle state transition tests
- Buffer contract tests (BUF_TOO_SMALL transactional behavior)
- Error handling tests (truncation, corruption, size limits)
- Cross-algorithm conformance (all four algorithms must pass identical test suite)

**Out of scope:**
- Wire format testing (covered by cross-language conformance tests)
- Performance benchmarks (covered by separate benchmark suite)
- Thread safety (single-threaded use is the contract)

## Executable coverage

The contract cases in `contract_cases.md` are covered by:

- `algorithms/shared/cpp/tests/test_lifecycle.cpp`
- `algorithms/shared/go/codec/lifecycle_test.go`
- `algorithms/shared/rust/tests/lifecycle.rs`
- algorithm-specific Go/Rust streaming tests under `algorithms/*/{go,rust}/`

Run them through the repository baseline:

```bash
make test
```

## Fixture Ownership

Test fixtures and test data generation are owned by this directory:
- Test input patterns (empty, single-byte, repetitive, random)
- Expected error conditions (truncated streams, oversized inputs)
- State transition validation helpers

Language-specific test implementations SHALL reference this contract but MAY use 
language-idiomatic test frameworks (Go: testing package; C++: Catch2; Rust: cargo test).

## Historical note

`red_phase_evidence.txt` records the original red phase from the streaming API
foundation change. It is kept as provenance, not as a current implementation
task list.
