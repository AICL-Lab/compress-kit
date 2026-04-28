# Tasks: add-interoperability-conformance

## Legend

- [ ] not started
- [x] done

---

## Prerequisites

- [ ] P1. Confirm `add-shared-frame-format` is archived (frame header fields must be stable).
- [ ] P2. Confirm `add-streaming-api-foundation` is archived (error codes `ERR_VERSION_UNSUPPORTED`, `ERR_UNKNOWN_ALGO`, `ERR_TRUNCATED`, `ERR_CORRUPT` must be stable).

## Phase A — Corpus Generation

- [ ] A1. Create `tests/gen_corpus.py` — deterministic generator using `--seed 42`; produces the 9 synthetic corpus files listed in design.md § Corpus Matrix (everything except `binary_elf_sample.bin`).
- [ ] A2. Run `python tests/gen_corpus.py --seed 42` and commit the generated synthetic files to `tests/corpus/`.
- [ ] A3. Add committed fixture `tests/corpus/binary_elf_sample.bin` verbatim from the curated repository sample described in design.md; do not source it from `gen_corpus.py` or the CI host.
- [ ] A4. Add SHA-256 manifest `tests/corpus/MANIFEST.sha256` covering both generated corpus files and the committed binary sample for integrity verification.

## Phase B — Test Vector Generation

- [ ] B1. Create `tests/gen_vectors.py` — generates header parsing test vectors (design.md § Header Parsing Tests) as binary files + `tests/vectors/header_vectors.json` manifest.
- [ ] B2. Generate and commit vectors: `python tests/gen_vectors.py`.
- [ ] B3. Add truncation helper function to `tests/conformance/helpers.py` that truncates any file at given byte offsets.
- [ ] B4. Add corruption helper function that flips a single byte at given offset.

## Phase C — Header Parsing Tests

- [ ] C1. Implement `tests/conformance/test_header_parsing.py` — runs each vector through all three language decoders and checks expected error codes.
- [ ] C2. Integrate `test_header_parsing.py` into `make test` via `Makefile`.

## Phase D — Truncation Tests

- [ ] D1. Implement `tests/conformance/test_truncation.py` — for each corpus file × 6 truncation points, confirm `ERR_TRUNCATED` (non-zero exit, no hang, no output file).
- [ ] D2. Add timeout guard (10 seconds per test case) to detect hangs.
- [ ] D3. Integrate into `make test`.

## Phase E — Corruption Tests

- [ ] E1. Implement `tests/conformance/test_corruption.py` — for each corpus file × 6 corruption offsets, confirm non-zero exit code (no silent success).
- [ ] E2. Integrate into `make test`.

## Phase F — Concatenation Tests

- [ ] F1. Implement `tests/conformance/test_concatenation.py` — concatenates two frames and decodes them; verifies output = decoded_A ++ decoded_B.
- [ ] F2. Test same-algorithm / same-content and same-algorithm / different-content cases.
- [ ] F3. Integrate into `make test`.

## Phase G — Cross-Language Decode Matrix

- [ ] G1. Implement `tests/conformance/test_decode_matrix.py` — runs all 9 (encoder × decoder) pairs for all 4 algorithms × 10 corpus files.
- [ ] G2. Output results to `tests/results/decode_matrix_<timestamp>.json`.
- [ ] G3. Add `tests/results/baseline_matrix.json` committing the first known-good run result.
- [ ] G4. Integrate into `make test` (non-blocking for Range Coder files >100 KB per known issue).
- [ ] G5. Document Range Coder ≤100 KB limit in test runner output and in `tests/README.md`.

## Phase H — CI Integration

- [ ] H1. Update `.github/workflows/test.yml` to run all conformance tests on every PR.
- [ ] H2. Add corpus integrity check step: validate `MANIFEST.sha256` before running tests.
- [ ] H3. Set test timeout to 5 minutes for the full matrix run.

## Phase I — Documentation

- [ ] I1. Add `docs/en/testing/conformance.md` describing the test categories, corpus, and how to add new test vectors.
- [ ] I2. Add `docs/zh/testing/conformance.md` (Chinese translation).
