# 更新日志

此页面记录 CompressKit 每个版本的变更。

All notable user-facing changes to CompressKit are tracked here.

The project follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/)
style categories and uses semantic versioning for releases.

## Unreleased

- **BREAKING**: Project refactored to C++17-only. Go and Rust implementations removed.
- **BREAKING**: Build system migrated from raw g++ Makefile to CMake.
- Removed OpenSpec, Cursor, and Claude skill meta-tooling directories.
- Removed cross-language conformance matrix and streaming API contract tests.
- **Simplification (Occam's razor)**: Removed VitePress documentation site, bilingual docs (en/zh), and root `package.json`/`docs/` tree. README is now single-language (Chinese).
- **Docs restoration**: Restored a slimmed VitePress docs site (`docs/` + `docs-pages.yml` workflow) keeping only algorithms, architecture, and benchmarks content. Removed stale streaming API, state-machine academy, ADRs, contributing, project-structure, and bibliography pages. Updated API/architecture pages to reflect the `BufferTransform`-based buffer layer and 3-value `StatusCode` enum; fixed broken base-path links.
- **Simplification**: Removed governance docs (`CODE_OF_CONDUCT.md`, `SECURITY.md`, `CONTEXT.md`, `CONTRIBUTING.md`) and `.devcontainer/` — unnecessary for a hobby/learning project.
- **Simplification**: Removed Streaming state-machine layer (`encoder.hpp` with `State`/`Encoder`/`Decoder` abstract classes, `BufferEncoder`/`BufferDecoder` wrappers). `encode_buffer`/`decode_buffer` now take a `BufferTransform` function pointer directly. Algorithms never used the streaming interface — the state machine was dead abstraction.
- **Simplification**: Merged CI workflows into a single `ci.yml` (removed `ci-docs.yml`, `codeql.yml`; `docs-pages.yml` later restored with the slimmed docs site).
- **Simplification**: Pruned `StatusCode` enum to only used values (`OK`/`ERR_CORRUPT`/`ERR_SIZE_LIMIT`); removed `INITIAL_DECODE_OVERHEAD` constant.
- Buffer layer rewritten to use in-memory transforms instead of temporary files.
- Huffman encoding now uses `uint64_t` code words instead of `std::string` for 8x density.
- Huffman decoding now uses 8-bit lookup table per internal node (~8x faster than bit-by-bit tree walk).
- Range Coder now uses shared `frequency_table` utilities (was duplicating scale/cumulative/header logic).
- `scale_frequencies` rewritten from O(N×M) decrement loop to O(N log N) proportional reduction.
- `build_cumulative` now signals all-zero tables via empty result instead of silent fallback.
- Centralized constants (`SYMBOL_LIMIT`, `EOF_SYMBOL`, magic bytes, size limits) in `constants.hpp`.
- VitePress documentation pruned of Go/Rust content; bilingual (en/zh) structure retained.
- Issue templates pruned of Go/Rust/OpenSpec/cross-language references; Language scope reduced to C++17 / Python scripts / Docs (feature template adds CI).

- `CMakeLists.txt` with static library target, 4 algorithm executables, and CTest integration.
- `algorithms/shared/cpp/include/compresskit/constants.hpp` for shared constants.
- Shared `count_frequencies`, `scale_frequencies`, `build_cumulative` utilities in `frequency_table.hpp`.
- Architecture Decision Records 0001, 0003, 0004, 0005 under `docs/adr/` (validation metadata module shape, range coder corpus cap policy, RLE buffered streaming stance, range C++ bench mode migration path).
- `tests/metadata.py` validation metadata module (C++17-only: `LANGUAGE_ORDER = ("cpp",)`).

### Changed (merge integration)

- Merged remote `master` (PR #9 architecture-deepening) into the C++17-only tree.
- Dropped ADR 0002 (cross-language semantic error alignment) and `docs/architecture/contract-inventory.md`: both depend on the removed Go/Rust implementations.
- `CONTEXT.md` "参考资料" section rebuilt: keeps ADR 0001/0003/0004/0005 references, drops the 0002 link and the deleted `openspec/` spec links.
- `docs/en/architecture/index.md` CLI section kept on the C++17-only single-binary contract (`./build/<algo>_cpp`); dropped multi-language `--lang` and `<algo>_<lang>` invocations.

### Changed (Clean Code: extract shared utilities)

- Extracted in-memory little-endian serialization helpers (`write_u32_le`, `write_magic`, `write_frequency_header`, `read_frequency_header`) to new `compresskit/serialization.hpp`. Eliminates duplicated `push_u32` lambdas and `read_frequencies` functions across huffman/arithmetic/range (~90 lines removed).
- Extracted `BitWriter` and `BitReader` to new `compresskit/bit_io.hpp`. Eliminates duplicated `BitWriter` class across huffman/arithmetic (~28 lines removed).
- RLE encode now uses shared `write_magic` and `write_u32_le` instead of inline copies.
- Removed unused `name` parameter from `compresskit::cli::run` (and all 4 algorithm call sites).
- Renamed `kInitialEncodeOverhead` (Google-style) to `INITIAL_ENCODE_OVERHEAD` (matches codebase UPPER_CASE convention for local constants).

### Changed (Clean Code: replace magic numbers with named constants)

- Added shared binary-format constants to `constants.hpp`: `MAGIC_SIZE`, `U32_SIZE`, `BITS_PER_BYTE`, `BYTE_VALUES`, `RLE_PAIR_SIZE`, `STREAM_READ_BUFFER_SIZE`, `INITIAL_ENCODE_OVERHEAD`, `INITIAL_DECODE_OVERHEAD`.
- Replaced bare `4` (magic size) with `MAGIC_SIZE` across huffman/arithmetic/range/rle decode paths and `serialization.hpp`.
- Replaced bare `4` (uint32 LE size) with `U32_SIZE` in `serialization.hpp`, `frequency_table.cpp`, and rle inline count decode.
- Replaced bare `8` / `7` (bits per byte) with `BITS_PER_BYTE` in `bit_io.hpp`, huffman decode table, and `buffer_api.cpp` encode-limit calculation.
- Replaced bare `256` (byte value count) with `BYTE_VALUES` in huffman 8-bit decode table.
- Replaced bare `5` (RLE count+value pair size) with `RLE_PAIR_SIZE` in rle decode.
- Replaced bare `32 * 1024` (stream read buffer) with `STREAM_READ_BUFFER_SIZE` in `frequency_table.cpp`.
- Promoted `INITIAL_ENCODE_OVERHEAD` from `buffer_api.cpp` anonymous namespace to `constants.hpp`; replaced bare `2048` reserve overhead in huffman/range encode.
- Added `INITIAL_DECODE_OVERHEAD` constant; replaced bare `1024` decode buffer overhead in `buffer_api.cpp`.
- Replaced `0xFFFFFFFFu` with `UINT32_MAX` in range coder encoder/decoder state.
- Named `STATE_BYTES = 4` (range coder 32-bit state width) in range main.
- Named `MAX_TREE_NODES = 2 * SYMBOL_LIMIT` (huffman worst-case node count) in huffman main.
- Named `EXPECTED_ARGC = 4` (program + mode + input + output) in `cli_launcher.cpp`.
- Refactored `write_u32_le` / `write_magic` / `read_frequency_header` and rle count decode from unrolled byte shifts to `U32_SIZE`/`MAGIC_SIZE` loops.

### Changed (Clean Code: finish magic-number sweep)

- Converted stream-based `write_u32_le` / `read_u32_le` in `frequency_table.cpp` from unrolled byte shifts to `U32_SIZE` loops (matches `serialization.hpp` style).
- Replaced bare `<<8` / `>>24` in range coder byte-renormalisation with `<< BITS_PER_BYTE` / `>> TOP_BYTE_SHIFT`; added local `TOP_BYTE_SHIFT = (STATE_BYTES - 1) * BITS_PER_BYTE` constant.

## 1.0.0 (2026-01-07)

- Huffman Coding, Arithmetic Coding, Range Coder, and Run-Length Encoding implementations.
- C++17, Go, and Rust command-line tools for all four algorithms.
- Unified CLI shape: `<binary> <encode|decode> <input> <output>`.
- Cross-language file compatibility goals for educational verification.
- Test data generation scripts and benchmark scripts.
- VitePress documentation site with English and Chinese content.
- MIT license, contribution guide, code of conduct, security policy, issue templates, and pull request template.

### Security

- Documented maximum input size of 4 GiB.
- Documented maximum decoded output size of 1 GiB for decompression-bomb protection.

[Unreleased]: https://github.com/LessUp/compress-kit/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/LessUp/compress-kit/releases/tag/v1.0.0
