# Change Proposal: add-shared-frame-format

> **Disposition:** Deferred during the 2026-04-29 finalization pass. This proposal
> is preserved as future design context, but it was not applied to the main
> specs because the current terminal baseline keeps the existing per-algorithm
> binary formats stable.

## Summary

Define a unified binary frame envelope for all CompressKit compressed files, adding a versioned header with magic bytes, algorithm ID, flags, content size, checksum, and optional extension fields (dictionary ID, skippable metadata blocks).

## Why

Currently each algorithm writes its own bespoke header (or none). The existing `REQ-ARCH-003` frequency table format is algorithm-specific. There is no:
- common magic to identify CompressKit files,
- version field to enable forward-compatibility,
- checksum to detect corruption,
- `content_size` hint to allow decoders to pre-allocate,
- extension point for future features (dictionaries, metadata).

Without a shared envelope, interoperability testing (`add-interoperability-conformance`) cannot define a canonical parse path, and the streaming API (`add-streaming-api-foundation`) has nowhere to write the end-of-stream marker.

## Scope

### In scope
- Byte-level frame format: layout, endianness, field widths
- Magic bytes value and rationale
- Version field semantics (current = 1)
- Algorithm ID registry
- Flags bitmap definitions
- `content_size` field (0 = unknown)
- `checksum` field (xxHash-64 of uncompressed content)
- Optional `dictionary_id` extension block
- Optional `skippable metadata` block mechanism
- Format versioning / backward-compatibility rules

### Out of scope
- Streaming API lifecycle (→ `add-streaming-api-foundation`)
- Conformance test matrix (→ `add-interoperability-conformance`)
- Benchmark infrastructure (→ `add-benchmark-governance`)
- The existing per-algorithm frequency table format (REQ-ARCH-003 unchanged)

## Impact

| Spec | Change type |
|------|-------------|
| `core-architecture/spec.md` | ADDED requirements REQ-ARCH-012, REQ-ARCH-013, REQ-ARCH-014 |
| `cross-language-testing/spec.md` | ADDED requirement REQ-TEST-007 |

## Dependencies

This change is self-contained at the frame-format level — the byte layout, magic bytes, and checksum fields can be specified independently. However, the error codes surfaced when a frame fails to parse (e.g., bad magic, unsupported version, checksum mismatch) are canonically defined in `add-streaming-api-foundation`.

**Ordering note for implementers**: settle the `add-streaming-api-foundation` error catalogue before writing the frame-parsing code that references those error codes. The spec documents here may be authored in any order, but the implementation of frame parsing in all three languages must not assign ad-hoc error values that conflict with the streaming error registry.

## Risks

- Adding a new header breaks existing encoded files. Mitigation: version field + migration tooling task in tasks.md.
- xxHash-64 dependency may be unavailable in some embedded targets — flagged as a future concern.
- Skippable block mechanism must use a well-known marker to avoid false-positive detection.
