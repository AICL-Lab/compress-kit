# Changelog

All notable user-facing changes to CompressKit are tracked here.

The project follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/)
style categories and uses semantic versioning for releases.

## [Unreleased]

### Changed

- **BREAKING**: Project refactored to C++17-only. Go and Rust implementations removed.
- **BREAKING**: Build system migrated from raw g++ Makefile to CMake (`CMakeLists.txt` with static library target, 4 algorithm executables, and CTest integration).
- Buffer layer rewritten to use in-memory transforms instead of temporary files. `encode_buffer` / `decode_buffer` now take a `BufferTransform` function pointer directly.
- **Simplification**: Removed Streaming state-machine layer (`encoder.hpp` with `State`/`Encoder`/`Decoder` abstract classes, `BufferEncoder`/`BufferDecoder` wrappers). The state machine was a dead abstraction - algorithms never used the streaming interface.
- **Simplification**: Pruned `StatusCode` enum to only used values (`OK` / `ERR_CORRUPT` / `ERR_SIZE_LIMIT`); removed `INITIAL_DECODE_OVERHEAD` constant.
- **Simplification**: Merged CI workflows into a single `ci.yml` (removed `ci-docs.yml`, `codeql.yml`); `docs-pages.yml` later restored for the slimmed docs site.
- Huffman encoding now uses `uint64_t` code words instead of `std::string` for 8x density.
- Huffman decoding now uses 8-bit lookup table per internal node (~8x faster than bit-by-bit tree walk).
- Range Coder now uses shared `frequency_table` utilities (was duplicating scale/cumulative/header logic).
- `scale_frequencies` rewritten from O(N×M) decrement loop to O(N log N) proportional reduction.
- `build_cumulative` now signals all-zero tables via empty result instead of silent fallback.
- Centralized constants (`SYMBOL_LIMIT`, `EOF_SYMBOL`, magic bytes, size limits) in `constants.hpp`.
- Issue templates pruned of Go/Rust/OpenSpec/cross-language references; language scope reduced to C++17 / Python scripts / Docs (feature template adds CI).

### Added

- `CMakeLists.txt` with static library target, 4 algorithm executables, and CTest integration.
- `algorithms/shared/cpp/include/compresskit/constants.hpp` for shared constants.
- Shared `count_frequencies`, `scale_frequencies`, `build_cumulative` utilities in `frequency_table.hpp`.
- `tests/metadata.py` validation metadata module (C++17-only: `LANGUAGE_ORDER = ("cpp",)`).
- `algorithms/shared/cpp/include/compresskit/serialization.hpp` - shared in-memory little-endian serialization helpers (`write_u32_le`, `write_magic`, `write_frequency_header`, `read_frequency_header`). Eliminates duplicated `push_u32` lambdas and `read_frequencies` across huffman/arithmetic/range (~90 lines removed).
- `algorithms/shared/cpp/include/compresskit/bit_io.hpp` - shared `BitWriter` / `BitReader`. Eliminates duplicated `BitWriter` class across huffman/arithmetic (~28 lines removed).
- Slimmed VitePress docs site (`docs/` + `docs-pages.yml` workflow) keeping only algorithms, architecture, and benchmarks content. API/architecture pages reflect the `BufferTransform`-based buffer layer and 3-value `StatusCode` enum.

### Removed

- Go and Rust implementations and their cross-language conformance matrix / streaming API contract tests.
- OpenSpec, Cursor, and Claude skill meta-tooling directories.
- Governance docs (`CODE_OF_CONDUCT.md`, `SECURITY.md`, `CONTEXT.md`, `CONTRIBUTING.md`) and `.devcontainer/` - unnecessary for a hobby/learning project.
- Stale streaming API, state-machine academy, ADRs, contributing, project-structure, and bibliography doc pages.

### Clean Code: shared utilities & named constants

- RLE encode now uses shared `write_magic` and `write_u32_le` instead of inline copies.
- Removed unused `name` parameter from `compresskit::cli::run` (and all 4 algorithm call sites).
- Renamed `kInitialEncodeOverhead` (Google-style) to `INITIAL_ENCODE_OVERHEAD` (matches codebase UPPER_CASE convention for local constants).
- Added shared binary-format constants to `constants.hpp`: `MAGIC_SIZE`, `U32_SIZE`, `BITS_PER_BYTE`, `BYTE_VALUES`, `RLE_PAIR_SIZE`, `STREAM_READ_BUFFER_SIZE`, `INITIAL_ENCODE_OVERHEAD`, `INITIAL_DECODE_OVERHEAD`.
- Replaced bare magic numbers (`4`, `8`, `7`, `256`, `5`, `32 * 1024`, `0xFFFFFFFFu`) with named constants across huffman/arithmetic/range/rle decode paths, `serialization.hpp`, `frequency_table.cpp`, `bit_io.hpp`, and `buffer_api.cpp`.
- Named local constants: `STATE_BYTES = 4` (range coder 32-bit state width), `MAX_TREE_NODES = 2 * SYMBOL_LIMIT` (huffman worst-case node count), `EXPECTED_ARGC = 4` (program + mode + input + output), `TOP_BYTE_SHIFT = (STATE_BYTES - 1) * BITS_PER_BYTE`.
- Refactored `write_u32_le` / `write_magic` / `read_frequency_header` and rle count decode from unrolled byte shifts to `U32_SIZE` / `MAGIC_SIZE` loops.

## [1.0.0] - 2026-01-07

### Added

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
