# Streaming API Contract Test Plan

This directory defines and validates the streaming API contract across all CompressKit implementations.

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

## Test Development Workflow

This directory follows **contract-first test-driven development**:

1. **Phase A0**: Define test cases in `contract_cases.md` (this stage)
2. **Phase A0.2**: Verify streaming APIs do not exist before implementation — see `red_phase_evidence.txt` for verification against base commit 90929c5
3. **Phase C**: Implement shared streaming abstractions to make tests pass
4. **Phase E**: Run full test suite and verify all algorithms conform

## Fixture Ownership

Test fixtures and test data generation are owned by this directory:
- Test input patterns (empty, single-byte, repetitive, random)
- Expected error conditions (truncated streams, oversized inputs)
- State transition validation helpers

Language-specific test implementations SHALL reference this contract but MAY use 
language-idiomatic test frameworks (Go: testing package; C++: Catch2; Rust: cargo test).

## Rules for This Stage

**Current stage (A0): Test planning only**
- Modifications to `contract_cases.md` are allowed and encouraged
- No production code changes are permitted until failing tests exist
- The initial "RED" phase of TDD must verify that current implementations 
  (Huffman io.Reader/io.Writer, Range []byte-based) do NOT support the streaming contract

**Rationale:** If tests pass before implementation, the tests are invalid.
