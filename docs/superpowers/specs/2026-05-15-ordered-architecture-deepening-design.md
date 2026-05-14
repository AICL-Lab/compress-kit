# Ordered Architecture Deepening Design

## Summary

This design decomposes the current architecture work into five ordered, contract-preserving phases:

1. Shared Buffer Layer policy module
2. Streaming Layer wrapper module
3. Frequency Table module
4. Magic Number and header parsing module
5. Writer adapter module

The goal is to increase **depth** at the existing **seams** without changing CompressKit's public CLI semantics, binary formats, security limits, or cross-language behavior.

## Context

CompressKit currently exposes the right top-level layers:

- CLI Entry Point
- Buffer Layer
- Streaming Layer
- Algorithm Core
- Shared Utilities

The friction is inside several shallow modules where the public **interface** is small, but important policy is duplicated or leaked across nearby callers. The strongest examples are:

- Buffer growth and `BUF_TOO_SMALL` retry semantics split across Go and Rust Buffer Layer implementations
- Repeated Streaming Layer wrappers around `BufferedEncoder` and `BufferedDecoder`
- Frequency Table logic split between Shared Utilities and Algorithm Core modules
- Magic Number and header validation repeated across algorithms and languages
- Writer adapters that also own output-buffer policy

## Constraints

- Preserve the unified CLI shape: `<binary> <encode|decode> <input> <output>`
- Preserve current binary formats and cross-language compatibility
- Preserve Streaming Layer and Buffer Layer public behavior required by:
  - `openspec/specs/core-architecture/spec.md`
  - `openspec/specs/encoding-project/spec.md`
  - `openspec/specs/cross-language-testing/spec.md`
- Preserve security limits: 4 GiB max input, 1 GiB max decoded output
- Keep error messages in English
- Avoid OpenSpec changes by keeping all public contracts stable

## Approaches Considered

### Approach A: Big-bang refactor across all five modules

Refactor all five shallow areas in one pass, then repair tests and conformance issues at the end.

**Pros**

- Fastest path to the target architecture on paper
- Can shape the final implementation holistically

**Cons**

- High risk of cross-language drift
- Large review surface
- Poor fault isolation when conformance fails
- Hard to prove which phase introduced a behavior change

### Approach B: Ordered phase-by-phase deepening at existing seams

Deepen one module cluster at a time, starting with the Buffer Layer policy and finishing with the Writer adapter after the lower-level policy is concentrated.

**Pros**

- Best locality for review and debugging
- Each phase can preserve the existing interface while improving implementation shape
- Easier to validate against existing test and conformance gates
- Minimizes risk of accidental public API drift

**Cons**

- Requires temporary coexistence of old and new internal paths during each phase
- Takes longer than a big-bang change

**Recommendation:** choose this approach.

### Approach C: Start from algorithm-specific modules and generalize later

Refactor Huffman, Arithmetic, Range Coder, and RLE one by one, then extract common modules afterward.

**Pros**

- Keeps each algorithm change small
- Fits teams that work algorithm-by-algorithm

**Cons**

- Repeats the same design work four times
- Encourages shallow abstractions because the shared seam is postponed
- Weak leverage until the final extraction

## Recommended Design

### Architecture principle

Each phase deepens an existing module by moving policy-heavy implementation behind a smaller, more stable **interface**. The rule is to concentrate behavior where the deletion test is strongest: deleting the deepened module should make complexity reappear across many callers.

### Phase 1: Shared Buffer Layer policy module

**Scope**

- Go: `algorithms/shared/go/codec/buffer.go`, `buffer_loop.go`, nearby callers such as `files.go`
- Rust: `algorithms/shared/rust/src/codec/buffer.rs`

**Problem**

The Buffer Layer public interface is already correct, but the implementation of retry loops, growth policy, output limits, and `BUF_TOO_SMALL` transactionality is duplicated per language and partially reused by adjacent modules.

**Design**

- Concentrate full-buffer orchestration into one deeper internal module per language
- Put these concerns behind the Buffer Layer seam:
  - initial buffer sizing
  - growth policy
  - retry loop for `process()` and `finish()`
  - size-limit enforcement
  - preservation of already-written bytes across retries
- Keep the external Buffer Layer entry points unchanged

**Dependency category**

In-process only

**Why first**

Phase 5 depends on this policy. Starting here avoids re-solving the same output-buffer problem later in the Writer adapter.

### Phase 2: Streaming Layer wrapper module

**Scope**

- Go algorithm wrapper modules under `algorithms/*/go/streaming.go`
- Parallel Rust wrapper entry points built on `BufferedEncoder` and `BufferedDecoder`

**Problem**

The wrappers are shallow modules that mainly wire Algorithm Core functions into shared streaming helpers. Understanding the Streaming Layer requires bouncing between repeated wrapper files and the shared implementation.

**Design**

- Move wrapper construction into a deeper shared module per language
- Let algorithm modules provide only the Algorithm Core encode/decode hook
- Keep algorithm-facing constructor names and Streaming Layer behavior stable

**Dependency category**

In-process only

**Why second**

This phase builds directly on the stronger shared buffering/story from Phase 1 and reduces repeated wrapper code before touching shared binary-format helpers.

### Phase 3: Frequency Table module

**Scope**

- Shared Utilities for Frequency Table behavior
- Static-model algorithms: Huffman, Arithmetic, Range Coder

**Problem**

Frequency Table build, scaling, cumulative-table generation, and serialization are spread across Shared Utilities and Algorithm Core modules. The seam leaks binary-format knowledge and scaling policy.

**Design**

- Introduce one deep Frequency Table module per language for static-model algorithms
- Concentrate:
  - symbol counting
  - EOF insertion
  - scaling policy
  - cumulative-table construction
  - serialization and deserialization
- Leave algorithm-specific coding logic in Algorithm Core modules

**Dependency category**

In-process only

**Why third**

This is the first phase that materially touches shared binary-format helpers. By this point, the Buffer Layer and Streaming Layer refactors should already have reduced surrounding noise.

### Phase 4: Magic Number and header parsing module

**Scope**

- Decode-time format identification and header reads in Go, Rust, and C++

**Problem**

Magic Number validation and basic header parsing are scattered across many decode implementations with inconsistent error wording and duplicated structure.

**Design**

- Introduce one dedicated format/header module per language
- Concentrate:
  - Magic Number validation
  - truncated-header detection
  - basic header shape checks
  - consistent mapping to existing error kinds
- Keep algorithm-specific payload decoding in the Algorithm Core module

**Dependency category**

In-process only

**Why fourth**

It depends on the shared format concepts clarified during Phase 3, but should remain separate because format identification is a distinct seam from Frequency Table policy.

### Phase 5: Writer adapter module

**Scope**

- Go: `algorithms/shared/go/codec/writer.go`
- Rust: `algorithms/shared/rust/src/codec/write.rs`

**Problem**

The Writer adapter is supposed to bridge an `Encoder` to an output writer, but it also owns buffer growth and retry policy. That makes the adapter shallow at its transport role and too deep in the wrong place.

**Design**

- Slim the Writer adapter down to transport responsibilities
- Reuse the policy concentrated in Phase 1 for output-buffer handling
- Keep the Writer adapter interface stable

**Dependency category**

In-process only

**Why last**

This phase should consume the deeper Buffer Layer policy rather than invent its own policy. Doing it earlier would either duplicate design work or create a premature seam.

## Resulting Module Shape

After all five phases, the high-level call flow remains:

`CLI Entry Point -> Buffer Layer -> Streaming Layer -> Algorithm Core`

The difference is that more policy sits behind deeper shared modules:

- Buffer Layer policy becomes one concentrated module per language
- Streaming Layer wrappers become shared construction logic instead of repeated glue
- Frequency Table behavior becomes a coherent static-model module
- Magic Number and header checks become a dedicated format module
- Writer adapter becomes a thin adapter over already-deep buffer policy

## Error Handling

- Preserve existing error kinds and security-limit behavior
- Keep `BUF_TOO_SMALL` transactional
- Keep decode-side truncation/corruption reporting aligned with current contract
- Normalize error paths only when they map to existing semantics; do not invent new public error categories

## Testing Strategy

### Phase gates

- After every phase: `make test` and `make lint`
- After Phases 3 and 4: also run `make test-conformance`

### Test focus by phase

- **Phase 1:** strengthen Buffer Layer interface tests around retries, preserved bytes, growth stops, and limit enforcement
- **Phase 2:** keep Streaming Layer tests focused on constructor behavior and shared lifecycle semantics, not duplicated wrapper internals
- **Phase 3:** add focused Frequency Table tests for scaling, EOF handling, cumulative tables, and serialization invariants
- **Phase 4:** add focused header parsing tests for bad magic, truncation, and corrupt header cases
- **Phase 5:** test the Writer adapter through its own interface while reusing Phase 1 policy coverage

### Cross-language stance

- Public behavior must stay aligned across C++17, Go, and Rust
- Internal module shapes may differ by language as long as the seam and observable behavior remain consistent

## Risks and Mitigations

### Risk: accidental public interface drift

**Mitigation:** keep all public constructor names, Buffer Layer functions, and CLI behavior unchanged.

### Risk: cross-language divergence

**Mitigation:** keep each phase contract-preserving and use conformance gates when touching shared binary-format behavior.

### Risk: creating new shallow pass-through modules

**Mitigation:** apply the deletion test at every extraction. If deleting a module only removes indirection, do not keep it.

## Non-goals

- Changing binary formats
- Changing CLI semantics
- Introducing new algorithms
- Reworking Range Coder large-file performance
- Adding speculative abstraction across languages

## Decision

Proceed with **Approach B**: one contract-preserving deepening phase at a time, in the order listed above.
