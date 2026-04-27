# Change Proposal: add-benchmark-governance

## Summary

Define fixed benchmark corpora, a structured report schema, four canonical metrics (compression ratio, encode speed, decode speed, memory peak), and per-algorithm regression thresholds that trigger CI failures when performance degrades beyond acceptable bounds.

## Motivation

The current `REQ-TEST-003` says "benchmarks SHALL run and report metrics" but does not define:
- Which files form the benchmark corpora (different each run = non-comparable results)
- What the report format looks like (makes automation and historical comparison impossible)
- What constitutes a regression (without a threshold, any slowdown is invisible)
- How memory usage is measured (no agreed method across C++, Go, Rust)

This change closes those gaps. It intentionally uses the same corpus files defined by `add-interoperability-conformance` to avoid duplication, but governs only the performance measurement and reporting aspects.

## Scope

### In scope
- Fixed benchmark corpus selection (subset of conformance corpus)
- JSON report schema with versioning
- Four canonical metrics: ratio, encode_speed_MBps, decode_speed_MBps, peak_memory_KiB
- Regression threshold table per algorithm (initial baseline values TBD from first run)
- CI gate: fail if any metric exceeds its per-metric threshold (ratio +5%, speed −10%, memory +20%)
- Historical report storage convention

### Out of scope
- Corpus file content (→ `add-interoperability-conformance` defines and generates them)
- Cross-language correctness testing (→ `add-interoperability-conformance`)
- API lifecycle (→ `add-streaming-api-foundation`)
- Frame format (→ `add-shared-frame-format`)

## Impact

| Spec | Change type |
|------|-------------|
| `cross-language-testing/spec.md` | MODIFIED REQ-TEST-003 (strengthened), ADDED REQ-TEST-012, REQ-TEST-013 |

## Dependencies

- `add-interoperability-conformance` must be archived first (corpus files must exist at `tests/corpus/`).

## Risks

- Initial baseline thresholds are set from first CI run — they may be too loose. Plan: tighten after 3 stable runs.
- Memory measurement differs by language (`/usr/bin/time -v` for C++ and Rust, `runtime.ReadMemStats` for Go). Methodology must be documented and consistent.
- Range Coder decode performance for >500 KB files is a known issue; benchmark corpus for Range Coder is capped at 100 KB.
