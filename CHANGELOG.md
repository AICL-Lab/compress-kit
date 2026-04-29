# Changelog

All notable user-facing changes to CompressKit are tracked here.

The project follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/)
style categories and uses semantic versioning for releases.

## [Unreleased]

### Added

- RLE format now includes 4-byte magic number `RLE\x00` for file type identification.
- Added `.golangci.yml` for Go linting configuration.
- Added `.clippy.toml` for Rust linting configuration.
- Executable cross-language conformance matrix via `make test-conformance`.
- Streaming API lifecycle and buffer contract coverage across shared C++/Go/Rust layers.

### Fixed

- C++ buffer API now uses platform-appropriate temp directory instead of hardcoded `/tmp/`.
- Rust RLE decode now validates count=0 to prevent invalid data.
- Rust Huffman encoding performance improved by using `Vec<u8>` instead of `String` for bitstream.
- Fixed Rust Arithmetic Coding streaming decode compatibility for short bitstreams.
- Fixed Rust Arithmetic Coding treatment of `0x00` input bytes so they are not confused with the EOF symbol.

### Changed

- **BREAKING**: RLE format now includes 4-byte magic header. Old RLE files without magic are incompatible.
- Archived future shared-frame, extended-conformance, and benchmark-governance proposals as deferred OpenSpec design context.
- Refined README and documentation entry points so the GitHub README stays a concise repository gateway.
- Removed 41 unused BMAD skills from `.claude/skills/` directory (~2MB reduction).
- Simplified `AGENTS.md` and `CLAUDE.md` for better AI agent guidance.

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
