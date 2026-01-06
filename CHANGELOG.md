# Changelog | 变更日志

All notable changes to this project will be documented in this file.

本文件记录本项目的所有重要变更。

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

格式基于 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.1.0/)，
本项目遵循 [语义化版本](https://semver.org/lang/zh-CN/)。

## [Unreleased]

### Added | 新增
- Nothing yet

### Changed | 变更
- Nothing yet

### Fixed | 修复
- Nothing yet

## [1.0.0] - 2026-01-07

### Added | 新增

#### Algorithms | 算法
- Huffman encoding implementation in C++, Go, and Rust
- Arithmetic coding implementation in C++
- Range coder implementation in C++, Go, and Rust
- Run-Length Encoding (RLE) implementation in C++, Go, and Rust

#### Features | 功能
- Unified CLI interface across all implementations (`encode`/`decode` commands)
- Cross-language file format compatibility
- Benchmark scripts for performance comparison
- Test data generation scripts

#### Documentation | 文档
- Comprehensive README with bilingual (Chinese/English) content
- Algorithm comparison table
- Quick start guide
- Project structure documentation

#### Open Source Infrastructure | 开源基础设施
- MIT License
- Contributing guidelines (CONTRIBUTING.md)
- Code of Conduct (CODE_OF_CONDUCT.md)
- Security policy (SECURITY.md)
- GitHub Issue templates (bug report, feature request)
- GitHub Pull Request template
- GitHub Actions CI/CD pipeline
  - Multi-platform C++ builds (Ubuntu, macOS)
  - Go build and lint checks
  - Rust build and clippy checks
  - Encode/decode correctness verification
  - Required files check

### Technical Details | 技术细节

#### File Formats | 文件格式
- Huffman: Magic header `HFMN`, frequency table, bit-encoded data
- RLE: `(count, value)` pairs with 4-byte little-endian count

#### Build Requirements | 构建要求
- C++17 compatible compiler
- Go 1.19+
- Rust 1.70+
- Python 3.8+ (for scripts)

---

## Version History Summary | 版本历史摘要

| Version | Date | Highlights |
|---------|------|------------|
| 1.0.0 | 2026-01-07 | Initial release with 4 algorithms, 3 languages, full open source setup |

---

[Unreleased]: https://github.com/YOUR_USERNAME/encoding/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/YOUR_USERNAME/encoding/releases/tag/v1.0.0
