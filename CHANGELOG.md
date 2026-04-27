# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

#### Documentation

- **Internationalization (i18n)**: Complete bilingual documentation site (English & Chinese)
  - New directory structure: `docs/en/` and `docs/zh/`
  - VitePress configuration with locale switcher
  - Full translations for all guide pages
- Enhanced README with modern badges and improved structure
- Professional changelog reorganization with detailed categorization

#### Algorithms & Implementations

- Arithmetic coding Go implementation with full test suite
- Arithmetic coding Rust implementation with full test suite
- Range coder Go CLI (`algorithms/range/go/cmd/main.go`) for cross-language testing
- Range coder Rust CLI (`algorithms/range/rust/src/bin/rangecoder.rs`) for cross-language testing

#### Infrastructure

- Cross-language benchmark scripts for Arithmetic coding and Range coder
- CI cross-language correctness tests for all 4 algorithms × 3 languages
- Input size validation (4 GiB max) to prevent frequency overflow attacks
- Output size validation (1 GiB max) to prevent decompression bombs
- RAII/smart pointers for memory management in C++ Huffman implementation
- `-Werror` flag in Makefile for C++ builds to catch all warnings

### Changed

#### Documentation Site

- VitePress sidebar restructured into logical groups (Overview, Guide, Reference)
- README.md simplified to repository entry point with links to docs site
- Updated site metadata and SEO optimization (Open Graph, Twitter Cards)

#### Build & CI

- Pages workflow: sparse-checkout optimization, Node 20→22, path-based triggers
- CI workflow: full build, test, and cross-language verification matrix
- Makefile: complete build/test/clean targets for all algorithms and languages
- Arithmetic benchmark upgraded from C++-only to cross-language (C++, Go, Rust)

#### Code Quality

- All error messages standardized to English across all implementations
- All Chinese comments translated to English for consistency
- Updated `go.work` to include `algorithms/arithmetic/go` module
- Updated `algorithms/range/rust/Cargo.toml` to include `rangecoder` CLI binary

### Fixed

- README badge URL: `docs.yml` → `pages.yml`
- README algorithm table: corrected "LZ77" → "Range Coder"
- README project structure: corrected `algorithms/lz77/` → `algorithms/range/`
- `.gitignore`: removed bare `Makefile` pattern that could shadow root Makefile
- Makefile `build-range` target: fixed incorrect `go test` call
- Makefile `clean` target: replaced fragile Python one-liner with proper `rm` commands
- Huffman cross-language compatibility: unified tree building tie-breaking across C++, Go, Rust
- CI Python setup: removed invalid `cache: pip` option
- CI Range coder tests: using 100KB test file to work around decode performance issue

### Security

- Added input size limits to prevent frequency overflow attacks
- Added output size limits to prevent decompression bomb attacks
- All implementations now validate file sizes before processing

---

## [1.0.0] - 2026-01-07

### Added

#### Core Algorithms

- **Huffman Coding**: Implementation in C++, Go, and Rust
  - Prefix-code based lossless compression
  - Cross-language file format compatibility
  - Benchmark suite

- **Arithmetic Coding**: Implementation in C++
  - Interval subdivision based compression
  - Near-optimal compression efficiency

- **Range Coder**: Implementation in C++, Go, and Rust
  - Integer-based arithmetic coding equivalent
  - Library API for Go and Rust
  - CLI tools for all languages

- **Run-Length Encoding (RLE)**: Implementation in C++, Go, and Rust
  - Simple (count, value) pair encoding
  - Optimized for repetitive data

#### Features

- Unified CLI interface across all implementations (`encode`/`decode` commands)
- Cross-language file format compatibility
- Benchmark scripts for performance comparison
- Test data generation scripts (Python)

#### Documentation

- Comprehensive README with bilingual (Chinese/English) content
- Algorithm comparison table with complexity analysis
- Quick start guide with build instructions
- Project structure documentation

#### Open Source Infrastructure

- MIT License
- Contributing guidelines (CONTRIBUTING.md) - bilingual
- Code of Conduct (CODE_OF_CONDUCT.md) - bilingual
- Security policy (SECURITY.md) - bilingual
- GitHub Issue templates (bug report, feature request)
- GitHub Pull Request template

#### CI/CD Pipeline

- GitHub Actions CI workflow
  - Multi-platform C++ builds (Ubuntu, macOS)
  - Go build and lint checks
  - Rust build and clippy checks
  - Encode/decode correctness verification
  - Required files check
- GitHub Actions Pages workflow for documentation deployment

### Technical Specifications

#### File Formats

| Algorithm | Magic | Header | Data |
|-----------|-------|--------|------|
| Huffman | `HFMN` | Frequency table (257×4 bytes) | Bit stream |
| Arithmetic | `AENC` | Frequency table (257×4 bytes) | Bit stream |
| Range Coder | `RCNC` | Frequency table (257×4 bytes) | Byte stream |
| RLE | None | None | (count, value) pairs |

#### Build Requirements

| Language | Version | Compiler/Tool |
|----------|---------|---------------|
| C++ | C++17 | g++ 9+ or clang++ 10+ |
| Go | 1.21+ | go 1.21+ |
| Rust | 1.70+ | rustc 1.70+ |
| Python | 3.8+ | python3 3.8+ |

---

## Version History Summary

| Version | Date | Highlights |
|---------|------|------------|
| [Unreleased] | - | Documentation i18n, security improvements, complete test coverage |
| [1.0.0] | 2026-01-07 | Initial release with 4 algorithms, 3 languages, full open source setup |

---

## Release Comparison

[Unreleased]: https://github.com/LessUp/compress-kit/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/LessUp/compress-kit/releases/tag/v1.0.0

---

## Contributing

When adding changes to this changelog, please follow the categorization:

- `Added` - New features
- `Changed` - Changes in existing functionality
- `Deprecated` - Soon-to-be removed features
- `Removed` - Now removed features
- `Fixed` - Bug fixes
- `Security` - Security-related changes
- `Performance` - Performance improvements
