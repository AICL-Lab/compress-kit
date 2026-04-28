# Tasks: add-streaming-api-foundation

## Legend

- [ ] not started
- [x] done

---

## Phase A — Spec & Design Review

- [ ] A1. Review design.md lifecycle state machine against all four algorithms (Huffman, AC, RC, RLE) for correctness; update design.md if gaps found.
- [ ] A2. Confirm `max_output_expansion` formula for each algorithm and document in design.md § Partial Output.
- [ ] A3. Review Go interface against standard `io.Reader`/`io.Writer` wrapping compatibility; update design.md if needed.

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

- [ ] C1. Create `algorithms/shared/go/go.mod` — shared Go module manifest for streaming foundation helpers.
- [ ] C2. Update `go.work` to include `./algorithms/shared/go` so algorithm modules can import the shared package during local development.
- [ ] C3. Create `algorithms/shared/go/codec/encoder.go` — `Encoder` and `Decoder` interfaces.
- [ ] C4. Create `algorithms/shared/go/codec/buffer.go` — `EncodeBuffer` / `DecodeBuffer` helpers.
- [ ] C5. Create `algorithms/shared/go/codec/errors.go` — error sentinel values and `StatusCode` type.
- [ ] C6. Add `WriterEncoder` adapter in `algorithms/shared/go/codec/writer.go` implementing `io.Writer` via `Process`.
- [ ] C7. Adapt Huffman Go to implement the interfaces from `algorithms/shared/go/codec`.
- [ ] C8. Adapt Arithmetic Go to implement the shared interfaces.
- [ ] C9. Adapt Range Go to implement the shared interfaces.
- [ ] C10. Adapt RLE Go to implement the shared interfaces.
- [ ] C11. Add lifecycle unit tests in `algorithms/shared/go/codec/lifecycle_test.go`.

## Phase D — Rust Implementation

- [ ] D1. Create `algorithms/shared/rust/Cargo.toml` — shared Rust crate manifest for streaming foundation helpers.
- [ ] D2. Create `algorithms/shared/rust/src/lib.rs` exporting the shared streaming modules.
- [ ] D3. Create `algorithms/shared/rust/src/codec/encoder.rs` — `Encoder` and `Decoder` traits.
- [ ] D4. Create `algorithms/shared/rust/src/codec/buffer.rs` — `encode_buffer` / `decode_buffer` free functions.
- [ ] D5. Create `algorithms/shared/rust/src/codec/error.rs` — `CodecError` enum.
- [ ] D6. Implement `WriteEncoder` adapter in `algorithms/shared/rust/src/codec/write.rs` (`impl std::io::Write`).
- [ ] D7. Adapt Huffman Rust to implement the shared traits.
- [ ] D8. Adapt Arithmetic Rust to implement the shared traits.
- [ ] D9. Adapt Range Rust to implement the shared traits.
- [ ] D10. Adapt RLE Rust to implement the shared traits.
- [ ] D11. Add lifecycle unit tests in `algorithms/shared/rust/tests/lifecycle.rs`.

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
