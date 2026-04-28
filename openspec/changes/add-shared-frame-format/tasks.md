# Tasks: add-shared-frame-format

## Legend

- [ ] not started
- [x] done

---

## Prerequisites

- [ ] P1. Confirm the `add-streaming-api-foundation` error catalogue is approved and its shared error constants (`ERR_CORRUPT`, `ERR_VERSION_UNSUPPORTED`, `ERR_TRUNCATED`, `ERR_UNKNOWN_ALGO`) are frozen before implementing any frame parser tasks below.

---

## Phase A0 — Parser Contract Planning

- [ ] A0.1. Create `tests/shared_frame_contract/frame_examples.md` and write parser-only tests before encoder/decoder integration.
- [ ] A0.2. Create `tests/shared_frame_contract/README.md` describing fixture scope and negative-frame coverage.
- [ ] A0.3. Keep frame parser independent from algorithm payload parsing while the contract tests are being introduced.

## Phase A — Spec & Design Validation

- [ ] A1. Verify magic `0x434B5A4D` does not conflict with other known compression formats (gzip `1f 8b`, zstd `fd 2f b5 28`, LZ4 `04 22 4d 18`).
- [ ] A2. Confirm xxHash-64 library availability: `github.com/cespare/xxhash` (Go), `xxhash` crate (Rust), and xxhash-cpp (C++). Add to dependency manifests if absent.
- [ ] A3. Confirm `FLAG_*` bit positions are stable; document rationale for bit 3–15 reserved policy.

## Phase B — Shared Header Library (C++17)

- [ ] B1. Create `algorithms/shared/cpp/include/compresskit/frame.hpp` — `FrameHeader`, `FrameTrailer`, `ExtensionBlock` structs.
- [ ] B2. Implement `frame_write()` and `frame_read()` in `algorithms/shared/cpp/src/frame.cpp`.
- [ ] B3. Add unit tests in `algorithms/shared/cpp/tests/test_frame.cpp` covering:
  - round-trip header write/read
  - wrong magic → error
  - unknown flags → error
  - version != 1 → error
  - unknown `algo_id` → `ERR_UNKNOWN_ALGO`
  - truncated header / payload / trailer → `ERR_TRUNCATED`
  - trailer CRC mismatch → `ERR_CORRUPT`
  - dict extension block
  - skippable block skip

## Phase C — Shared Header Library (Go)

- [ ] C1. Reuse the existing `algorithms/shared/go` module/workspace scaffold owned by `add-streaming-api-foundation`; this change only extends that shared host with frame helpers.
- [ ] C2. Create `algorithms/shared/go/frame/frame.go` — `Header`, `Trailer`, and extension types.
- [ ] C3. Implement `WriteHeader`, `ReadHeader`, `WriteTrailer`, `ReadTrailer`.
- [ ] C4. Add unit tests in `algorithms/shared/go/frame/frame_test.go` covering:
  - round-trip header/trailer write/read
  - wrong magic → error
  - unknown flags → error
  - version != 1 → error
  - unknown `algo_id` → `ERR_UNKNOWN_ALGO`
  - truncated header / payload / trailer → `ERR_TRUNCATED`
  - trailer CRC mismatch → `ERR_CORRUPT`
  - dict extension block
  - skippable block skip

## Phase D — Shared Header Library (Rust)

- [ ] D1. Reuse the existing `algorithms/shared/rust` crate scaffold owned by `add-streaming-api-foundation`; this change only extends that shared host with frame modules.
- [ ] D2. Create `algorithms/shared/rust/src/frame/mod.rs` — `FrameHeader`, `FrameTrailer`, and extension types.
- [ ] D3. Implement `write_header`, `read_header`, `write_trailer`, `read_trailer`.
- [ ] D4. Add unit tests in `algorithms/shared/rust/tests/frame.rs` covering:
  - round-trip header/trailer write/read
  - wrong magic → error
  - unknown flags → error
  - version != 1 → error
  - unknown `algo_id` → `ERR_UNKNOWN_ALGO`
  - truncated header / payload / trailer → `ERR_TRUNCATED`
  - trailer CRC mismatch → `ERR_CORRUPT`
  - dict extension block
  - skippable block skip

## Phase E — Algorithm Integration

- [ ] E1. Update Huffman encoder/decoder (all 3 languages) to write/read frame header.
- [ ] E2. Update Arithmetic encoder/decoder (all 3 languages).
- [ ] E3. Update Range encoder/decoder (all 3 languages).
- [ ] E4. Update RLE encoder/decoder (all 3 languages).
- [ ] E5. Verify `content_size` is written correctly; use `0` when streaming (size unknown at start).
- [ ] E6. Verify `checksum` is computed over uncompressed content before encoding.

## Phase F — Migration Tooling

- [ ] F1. Implement `compresskit-migrate` CLI sub-command (or script) that reads a legacy (pre-frame) file and re-wraps it with a frame header.
- [ ] F2. Add a test in the CI pipeline that encodes a file with the migration tool and decodes it successfully.

## Phase G — Documentation

- [ ] G1. Add `docs/en/format/frame.md` with the full frame layout table and field descriptions.
- [ ] G2. Add `docs/zh/format/frame.md` (Chinese translation).
- [ ] G3. Update architecture overview to reference the new frame format.
