# Changelog

This page documents the changes in each CompressKit release.

All notable user-facing changes to CompressKit are tracked here.

The project follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/)
style categories and uses semantic versioning for releases.

## Unreleased

- **BREAKING**: Project refactored to C++17-only. Go and Rust implementations removed.
- **BREAKING**: Build system migrated from raw g++ Makefile to CMake.
- Removed OpenSpec, Cursor, and Claude skill meta-tooling directories.
- Removed cross-language conformance matrix and streaming API contract tests.
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
