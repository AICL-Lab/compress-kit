# Change Proposal: add-streaming-api-foundation

## Summary

Introduce a two-layer public API (buffer layer + streaming layer) across all three language implementations, defining a common `process / flush / finish` lifecycle and conformance rules for partial input, partial output, and EOF signalling.

## Why

The current implementation exposes only file-to-file encode/decode paths. Callers who want to:
- compress network streams without intermediate files,
- integrate into pipelines (pipes, HTTP bodies, in-memory transforms),
- feed data incrementally (e.g. chunked uploads),

...have no supported interface. Worse, each language currently manages buffering ad-hoc, so adding streaming later would require breaking changes. Establishing the API contract now—before the shared frame format or interoperability test matrix is built—ensures every subsequent change can assume a stable, testable surface.

## Scope

### In scope
- Buffer-layer API (synchronous, whole-buffer encode/decode)
- Streaming-layer API (incremental push: `process` → `flush` → `finish`)
- Lifecycle state machine shared by all three languages
- Partial-input and partial-output semantic rules
- EOF conformance contract
- Language-specific mapping notes (C++17 class, Go interface, Rust trait)

### Out of scope
- Wire format fields (→ `add-shared-frame-format`)
- Cross-language test matrix (→ `add-interoperability-conformance`)
- Benchmark metrics / regression thresholds (→ `add-benchmark-governance`)
- Async / non-blocking variants (future change)

## Impact

| Spec | Change type |
|------|-------------|
| `core-architecture/spec.md` | ADDED requirements REQ-ARCH-009, REQ-ARCH-010, REQ-ARCH-011 |
| `encoding-project/spec.md` | ADDED requirement REQ-PROD-009 |
| `cross-language-testing/spec.md` | ADDED requirement REQ-TEST-006 |

## Dependencies

None. This is a foundational change with no upstream OpenSpec dependencies.

## Risks

- Rust trait design must avoid `std::io::Read/Write` conflicts — see design.md.
- Go interface must be compatible with `io.Reader`/`io.Writer` wrapping patterns.
- Lifecycle state machine must be deterministic under concurrent misuse (error, not UB).
