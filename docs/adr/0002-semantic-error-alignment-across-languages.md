# ADR 0002: Semantic error alignment across languages

Date: 2026-05-31  
Status: Accepted

## Context

The three language implementations expose the same seven semantic error kinds (`BUF_TOO_SMALL`, `ERR_TRUNCATED`, `ERR_CORRUPT`, `ERR_INVALID_STATE`, `ERR_SIZE_LIMIT`, `ERR_VERSION_UNSUPPORTED`, `ERR_UNKNOWN_ALGO`) but through structurally different APIs:

- **Go**: `ErrorKind` enum + `CodecError` struct + sentinel errors + `errors.Is` mapping (plus `KindIO` for wrapped I/O)
- **Rust**: `CodecError` enum with `Display` + helper mapping from `io::Error`
- **C++**: `StatusCode` inside `Result<T>`; CLI returns boolean success/failure

Some algorithm-level direct APIs still map errors through local messages before the shared codec adapter normalizes them (Range Rust, RLE Rust, Range Go).

The `contract-inventory.md` notes that structural API identity across languages is explicitly **not** the goal.

## Decision

Align at the **semantic** level only:

1. The seven canonical error kinds listed in `CONTEXT.md` are the normative contract. Each language must be able to produce and identify each kind.
2. Language-idiomatic representations (enum, struct, sentinel, Result type) are preserved. No forced structural uniformity.
3. Algorithm-level error messages flowing to CLI output must be English and actionable (OpenSpec requirement). Internal mapping paths may remain language-local provided the final kind is correct.
4. A shared cross-language contract test (see ADR phase-2 deliverables) verifies that encode/decode round-trips produce the correct kind on corruption, truncation, size-limit, and invalid-state inputs.

## Constraints

- Changing the public streaming/buffer API error kind semantics or lifecycle requires an OpenSpec change.
- Go `KindIO` is a Go-local extension, not a cross-language contract; it must not appear in conformance test assertions.

## Consequences

**Positive**: Each language stays idiomatic; shared tests pin semantic equivalence; algorithm-level mapping quirks are explicitly tolerated while remaining testable.

**Negative**: Developers must understand that identical code structure is not the goal; code review must check semantic equivalence rather than structural similarity.

**Deferred**: Whether `KindIO` should be formally excluded from the CONTEXT vocabulary, or documented as Go-only, is a minor cleanup deferred to documentation alignment (Todo 8 / ADR 0002 amendment).
