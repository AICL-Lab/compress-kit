# Tasks: add-streaming-api-foundation

## Legend

- [ ] not started
- [x] done

---

## Phase A0 — Contract-First Test Planning

- [x] A0.1. Create `tests/streaming_api_contract/contract_cases.md` and write failing tests for process/flush/finish lifecycle.
- [x] A0.2. Verify tests fail against current Huffman/Range API surface before introducing shared streaming abstractions. Evidence: `tests/streaming_api_contract/red_phase_evidence.txt` (verified against base commit 90929c5).
- [x] A0.3. Create `tests/streaming_api_contract/README.md` describing scope, fixture ownership, and the rule that this stage modifies test plans only.

## Phase A — Spec & Design Review

- [x] A1. Review design.md lifecycle state machine against all four algorithms (Huffman, AC, RC, RLE) for correctness; update design.md if gaps found.
- [x] A2. Confirm `max_output_expansion` formula for each algorithm and document in design.md § Partial Output.
- [x] A3. Review Go interface against standard `io.Reader`/`io.Writer` wrapping compatibility; update design.md if needed.

## Phase B — C++17 Implementation

- [ ] B1. Create `algorithms/shared/cpp/include/compresskit/encoder.hpp` — abstract `Encoder` and `Decoder` base classes matching design.md sketches.
- [ ] B2. Create `algorithms/shared/cpp/include/compresskit/buffer_api.hpp` — `encode_buffer` / `decode_buffer` free functions.
- [ ] B3. Create `algorithms/shared/cpp/include/compresskit/result.hpp` — `Result<T>` type and error code enum.
- [ ] B4. Implement `BufferEncoder` shim in `algorithms/shared/cpp/src/buffer_api.cpp`.
- [ ] B5. Add unit tests for lifecycle state transitions in `algorithms/shared/cpp/tests/test_lifecycle.cpp`.
- [ ] B6. Adapt Huffman C++ to implement the `Encoder`/`Decoder` interfaces.
- [ ] B7. Adapt Arithmetic C++ to implement the interfaces.
- [ ] B8. Adapt Range C++ to implement the interfaces.
- [ ] B9. Adapt RLE C++ to implement the interfaces.

## Phase C — Go Implementation

- [x] C1. Create `algorithms/shared/go/go.mod` — shared Go module manifest for streaming foundation helpers.
- [x] C2. Update `go.work` to include `./algorithms/shared/go` so algorithm modules can import the shared package during local development.
- [x] C3. Create `algorithms/shared/go/codec/encoder.go` — `Encoder` and `Decoder` interfaces.
- [x] C4. Create `algorithms/shared/go/codec/buffer.go` — `EncodeBuffer` / `DecodeBuffer` helpers.
- [x] C5. Create `algorithms/shared/go/codec/errors.go` — error sentinel values and `StatusCode` type.
- [x] C6. Add `WriterEncoder` adapter in `algorithms/shared/go/codec/writer.go` implementing `io.Writer` via `Process`.
- [x] C7. Adapt Huffman Go to implement the interfaces from `algorithms/shared/go/codec`.
- [x] C8. Adapt Arithmetic Go to implement the shared interfaces.
- [x] C9. Adapt Range Go to implement the shared interfaces.
- [x] C10. Adapt RLE Go to implement the shared interfaces.
- [x] C11. Add lifecycle unit tests in `algorithms/shared/go/codec/lifecycle_test.go`.

## Phase D — Rust Implementation

- [x] D1. Create `algorithms/shared/rust/Cargo.toml` — shared Rust crate manifest for streaming foundation helpers.
- [x] D2. Create `algorithms/shared/rust/src/lib.rs` exporting the shared streaming modules.
- [x] D3. Create `algorithms/shared/rust/src/codec/encoder.rs` — `Encoder` and `Decoder` traits.
- [x] D4. Create `algorithms/shared/rust/src/codec/buffer.rs` — `encode_buffer` / `decode_buffer` free functions.
- [x] D5. Create `algorithms/shared/rust/src/codec/error.rs` — `CodecError` enum.
- [x] D6. Implement `WriteEncoder` adapter in `algorithms/shared/rust/src/codec/write.rs` (`impl std::io::Write`).
- [x] D7. Adapt Huffman Rust to implement the shared traits.
- [x] D8. Adapt Arithmetic Rust to implement the shared traits.
- [x] D9. Adapt Range Rust to implement the shared traits.
- [x] D10. Adapt RLE Rust to implement the shared traits.
- [x] D11. Add lifecycle unit tests in `algorithms/shared/rust/tests/lifecycle.rs`.

## Phase E — Integration & Verification

- [ ] E1. Update `Makefile` targets: `make test` must include streaming-layer unit tests.
- [ ] E2. Run `make test` — all tests pass.
- [ ] E3. Run `make lint` — no new warnings.
- [ ] E4. Update CLI file-to-file paths to use the new buffer-layer helpers (removes ad-hoc buffering).
- [ ] E5. Confirm security limits (4 GiB in / 1 GiB out) are enforced at the streaming layer boundary.

## Phase F — Documentation

- [ ] F1. Add `docs/en/api/streaming.md` covering lifecycle, error codes, and language examples.
- [ ] F2. Add `docs/zh/api/streaming.md` (Chinese translation).
- [ ] F3. Update `docs/en/guide/architecture.md` to reference the two-layer model.
