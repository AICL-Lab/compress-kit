# CompressKit architecture deepening design

## Context

CompressKit must preserve cross-language binary compatibility, the unified CLI
shape (`<binary> <encode|decode> <input> <output>`), the Streaming Layer state
machine, and the 4 GiB input / 1 GiB decoded-output safety limits.

The current codebase has several shallow modules around the Buffer Layer,
Streaming Layer adapters, and Go CLI entry points. The main friction is that
buffer growth, `BUF_TOO_SMALL` retry, size-limit enforcement, and finish-loop
semantics are spread across multiple modules instead of living behind one deeper
seam per language.

## Goal

Deepen the Buffer Layer and its surrounding adapters in five sequential slices
so that:

1. the interface for callers stays stable,
2. leverage increases because shared policy lives behind smaller seams, and
3. locality improves because lifecycle and size-limit fixes land in one place.

## Chosen approach

Use **reference-first slices**.

Start with the cross-language Buffer Layer as the reference slice. Once that
seam is deeper and well-tested, refactor dependent adapters in order so they
delegate policy inward instead of reimplementing it.

Rejected alternatives:

- **Layer-first sweep across the whole repository:** more consistent at first,
  but failures would be harder to localize and rollback.
- **Language-first sweep:** faster in Go, but it would leave CompressKit with an
  uneven architecture longer and weaken cross-language leverage.

## Design

### Slice 1: Shared Buffer Layer orchestration

**Files in scope**

- `algorithms/shared/go/codec/buffer.go`
- `algorithms/shared/cpp/src/buffer_api.cpp`
- `algorithms/shared/cpp/include/compresskit/buffer_api.hpp`
- `algorithms/shared/rust/src/codec/buffer.rs`

**Design**

Keep the external Buffer Layer interface unchanged (`encode_buffer` /
`decode_buffer` and language-idiomatic equivalents). Deepen the implementation
so one internal module per language owns:

- output buffer growth,
- `BUF_TOO_SMALL` retry loops,
- encode/decode size-limit policy,
- finish-loop orchestration, and
- transactional behavior when output is too small.

The Buffer Layer remains the seam that callers and tests cross. Internal helper
modules may be introduced, but only when they concentrate policy rather than
create new pass-through seams.

### Slice 2: Go writer seam

**Files in scope**

- `algorithms/shared/go/codec/writer.go`
- any new internal helper introduced by slice 1

**Design**

Make the writer module a true adapter over the Buffer Layer seam. Its interface
should remain “accept bytes, emit encoded bytes to `io.Writer`, finish on
close.” Buffer growth, retry, and finish policy should live behind the deeper
Buffer Layer seam rather than being reimplemented in the writer module.

### Slice 3: Algorithm Streaming Layer wrappers

**Files in scope**

- `algorithms/huffman/go/streaming.go`
- `algorithms/arithmetic/go/streaming.go`
- `algorithms/range/go/streaming.go`
- `algorithms/rle/go/streaming.go`

**Design**

Collapse repetitive wrapper logic into a shared Streaming Layer adapter shape.
Each algorithm module should contribute only the algorithm-specific
implementation hook: its single-shot encode/decode behavior. Lifecycle policy
stays in the shared buffered implementation.

### Slice 4: Go CLI entry-point modules

**Files in scope**

- `algorithms/huffman/go/cmd/main.go`
- `algorithms/arithmetic/go/cmd/main.go`
- `algorithms/range/go/cmd/main.go`
- `algorithms/rle/go/cmd/main.go`

**Design**

Create one deeper launcher module for the unified CLI seam. Per-algorithm `main`
files become tiny adapters that bind the algorithm implementation to the common
launcher. Usage text, mode validation, and error reporting become uniform.

### Slice 5: C++ Buffer Layer temp-file seam

**Files in scope**

- `algorithms/shared/cpp/src/buffer_api.cpp`
- `algorithms/shared/cpp/include/compresskit/buffer_api.hpp`

**Design**

Demote temp-file transforms to an adapter behind the Buffer Layer seam. The
Buffer Layer module should own lifecycle rules and size-limit behavior; file I/O
compatibility should not be entangled with the module’s core implementation.

This slice may remain file-backed internally if needed, but the design target is
to make the temp-file path an adapter rather than a defining part of the Buffer
Layer module.

## Data flow

After the refactor, the common shape is:

1. caller crosses a small interface at the Buffer Layer or launcher seam,
2. the deep module owns lifecycle and policy,
3. adapters translate to file I/O, `io.Writer`, or algorithm-specific hooks,
4. tests exercise the same seam that production callers use.

## Error handling

The refactor must preserve:

- Streaming Layer state transitions,
- `ERR_INVALID_STATE` semantics,
- transactional `BUF_TOO_SMALL` behavior,
- 4 GiB input limit enforcement, and
- 1 GiB decoded-output limit enforcement.

No slice may change binary format, magic bytes, CLI contract shape, or
cross-language conformance semantics.

## Testing strategy

The interface is the test surface.

- **Slice 1:** strengthen Buffer Layer contract tests around retry, size limits,
  and finish behavior.
- **Slice 2:** test writer-specific adapter behavior (`Write`, `Close`, output
  forwarding) rather than re-testing generic growth policy.
- **Slice 3:** keep only tests that prove algorithm-specific hooks are wired into
  the shared Streaming Layer adapter correctly.
- **Slice 4:** add launcher-level tests for usage text, mode validation, and
  encode/decode dispatch.
- **Slice 5:** test the Buffer Layer seam while swapping the C++ file-backed
  adapter beneath it.

Where wrapper tests only duplicate behavior now covered at the deeper seam, they
should be deleted or reduced.

## Success criteria

- Callers still use the same public interfaces.
- Shared policy lives in fewer places per language.
- Later slices become smaller because they delegate inward.
- Tests move upward to the seam instead of reaching through shallow wrappers.
