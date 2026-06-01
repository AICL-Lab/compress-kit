# ADR 0004: RLE buffered streaming stance

Date: 2026-05-31  
Status: Accepted

## Context

RLE is naturally streamable: each `(count, value)` pair is independent and requires no look-ahead beyond the current run. However, the current Go and Rust RLE adapters use buffered streaming wrappers (`BufferedEncoder`/`BufferedDecoder`) that collect the full input before processing, identical to the static-model algorithms (Huffman, Arithmetic, Range).

From `contract-inventory.md` (§ Buffer and streaming contract):

> Flush is a no-op for buffered static-model algorithms until `Finish`; RLE is currently also buffered in Go/Rust adapters despite being naturally streamable.

The C++ RLE implementation reads/writes file streams incrementally but is wrapped through the C++ buffer adapter for the CLI layer.

## Decision

**Retain buffered semantics for RLE** in all three languages as the documented contract for this product version.

Rationale:
- RLE shares the same `BufferedEncoder`/`BufferedDecoder` lifecycle as the other three algorithms; making it the only algorithm with true incremental streaming would create a two-tier streaming contract that complicates cross-language conformance tests.
- The performance benefit of incremental RLE streaming is not material for the current corpus sizes (up to 10 MiB).
- Implementing true incremental streaming for RLE would require a new streaming adapter path across all three languages — a significant scope increase.

The buffered semantic is **explicitly documented** rather than left implicit. API callers must not assume that `Process` writes any output before `Finish` is called.

## Constraints

- Do not change RLE's public streaming lifecycle behavior without an OpenSpec change (it is part of the public streaming/buffer API contract).
- If a future OpenSpec change introduces true incremental streaming for RLE, the conformance test suite must be updated to assert per-chunk output behavior.

## Consequences

**Positive**: Consistent lifecycle across all four algorithms; simpler conformance tests; explicit stance eliminates future confusion.

**Negative**: RLE cannot be used for low-latency or memory-constrained streaming use cases; that limitation is now documented.

**Deferred**: True incremental RLE streaming as a future feature, subject to a new OpenSpec change.
