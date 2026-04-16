# Encoding

<p align="center">
  <img src="docs/public/logo.svg" width="120" alt="Encoding Logo">
</p>

<p align="center">
  <a href="https://github.com/LessUp/encoding/actions/workflows/ci.yml">
    <img src="https://github.com/LessUp/encoding/actions/workflows/ci.yml/badge.svg" alt="CI">
  </a>
  <a href="https://github.com/LessUp/encoding/actions/workflows/pages.yml">
    <img src="https://github.com/LessUp/encoding/actions/workflows/pages.yml/badge.svg" alt="Docs">
  </a>
  <a href="https://github.com/LessUp/encoding/releases">
    <img src="https://img.shields.io/github/v/release/LessUp/encoding?include_prereleases" alt="Release">
  </a>
  <a href="https://opensource.org/licenses/MIT">
    <img src="https://img.shields.io/badge/License-MIT-yellow.svg" alt="License: MIT">
  </a>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/C++-17-blue.svg" alt="C++17">
  <img src="https://img.shields.io/badge/Go-1.21+-00ADD8.svg" alt="Go 1.21+">
  <img src="https://img.shields.io/badge/Rust-1.70+-orange.svg" alt="Rust 1.70+">
  <img src="https://img.shields.io/badge/Python-3.8+-3776AB.svg" alt="Python 3.8+">
</p>

<p align="center">
  <b>English</b> | <a href="README.zh-CN.md">简体中文</a> | <a href="https://lessup.github.io/encoding/">Documentation</a>
</p>

---

## 🚀 Quick Start

```bash
# Clone the repository
git clone https://github.com/LessUp/encoding.git
cd encoding

# Build all implementations
make build

# Run tests
make test

# Run benchmarks
make bench
```

## ✨ Features

- **4 Classic Algorithms**: Huffman, Arithmetic Coding, Range Coder, and RLE
- **3 Languages**: C++17, Go 1.21+, Rust 1.70+
- **Cross-Language Compatible**: Same binary format across all implementations
- **Learning-Focused**: Documentation emphasizes algorithm principles and comparisons
- **Production-Ready**: Complete CI/CD with automated testing and verification

## 📊 Algorithm Comparison

| Algorithm | Compression | Speed | Complexity | Best For |
|-----------|-------------|-------|------------|----------|
| Huffman | Medium | Fast | O(n log σ) | General purpose |
| Arithmetic | High | Medium | O(n) | Maximum compression |
| Range Coder | High | Fast | O(n) | Balanced performance |
| RLE | Variable | Very Fast | O(n) | Repetitive data |

## 📖 Documentation

| Resource | Description | Link |
|----------|-------------|------|
| **Documentation Site** | Full documentation with bilingual support | [lessup.github.io/encoding](https://lessup.github.io/encoding/) |
| **Project Specs** | Product requirements, RFCs, and technical specs | [specs/](specs/) |
| **Getting Started** | Setup, build, and basic usage | [Guide →](https://lessup.github.io/encoding/guide/getting-started) |
| **Algorithms** | Algorithm explanations and comparisons | [Guide →](https://lessup.github.io/encoding/guide/algorithms) |
| **Project Structure** | Directory layout and conventions | [Guide →](https://lessup.github.io/encoding/guide/project-structure) |
| **Changelog** | Version history and release notes | [View →](CHANGELOG.md) |

## 💡 Example Usage

```bash
# Encode with Huffman (C++)
./huffman/cpp/huffman_cpp encode input.txt output.huf

# Decode with a different language (Go)
./huffman/go/huffman_go decode output.huf restored.txt

# Verify correctness
diff input.txt restored.txt  # No output = identical
```

## 🛠️ Build Options

| Command | Description |
|---------|-------------|
| `make build` | Build all implementations |
| `make build-huffman` | Build Huffman only |
| `make build-arithmetic` | Build Arithmetic only |
| `make build-range` | Build Range Coder only |
| `make build-rle` | Build RLE only |
| `make test` | Run all tests |
| `make bench` | Run performance benchmarks |
| `make clean` | Clean build artifacts |

## 🤝 Contributing

We welcome contributions! Please read our [Contributing Guide](CONTRIBUTING.md) for details on:

- Code style guidelines for C++, Go, and Rust
- Testing requirements
- Pull request process

## 🙏 Acknowledgments

This project is inspired by educational resources on compression algorithms and aims to provide:

- Clean, readable implementations for learning
- Fair cross-language performance comparisons
- Verified correct implementations through extensive testing

## 📄 License

This project is licensed under the [MIT License](LICENSE).

Copyright © 2025-2026 LessUp

---

<p align="center">
  <sub>Built with ❤️ for the open source community</sub>
</p>
