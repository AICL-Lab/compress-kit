---
layout: home

hero:
  name: CompressKit
  text: Compression Algorithms
  tagline: Production-ready compression algorithms in C++17, Go, and Rust. Learn, compare, and verify across languages with identical binary formats.
  image:
    src: /logo.svg
    alt: CompressKit Logo
  actions:
    - theme: brand
      text: Get Started
      link: /en/guide/getting-started
    - theme: alt
      text: View on GitHub
      link: https://github.com/LessUp/compress-kit
    - theme: alt
      text: 中文
      link: /zh/

features:
  - icon: 🌐
    title: Multi-Language Comparison
    details: Every algorithm implemented in C++17, Go, and Rust. Compare performance, code style, and engineering practices across languages.
  - icon: 📦
    title: Cross-Language Compatible
    details: All implementations share identical binary formats. Encode in C++, decode in Go, verify in Rust — seamless interoperability.
  - icon: 📚
    title: Learning-Focused Documentation
    details: Understand the theory behind each algorithm with clear explanations and working code examples in three languages.
  - icon: ✅
    title: Production-Ready Verification
    details: Complete CI/CD with automated builds, cross-language correctness tests, and continuous benchmarking pipelines.
---

<StatsBar />

## Explore Algorithms

<AlgorithmGrid />

## Quick Comparison

| Algorithm | Compression | Speed | Best For |
|-----------|-------------|-------|----------|
| **Huffman** | Medium | Fast | General text/data |
| **Arithmetic** | High | Medium | Maximum compression needs |
| **Range Coder** | High | Fast | Balanced performance |
| **RLE** | Variable | Very Fast | Highly repetitive data |

## Quick Start

```bash
# Clone the repository
git clone https://github.com/LessUp/compress-kit.git
cd compress-kit

# Build all implementations
make build

# Run tests
make test

# Run benchmarks
make bench
```

## Cross-Language Verification

A key feature of CompressKit — encode in any language, decode in any other:

```bash
# Encode with C++
./algorithms/huffman/cpp/huffman_cpp encode input.bin encoded.huf

# Decode with Go
./algorithms/huffman/go/huffman_go decode encoded.huf restored.bin

# Verify correctness
diff input.bin restored.bin  # No output = identical
```

## Documentation Structure

| Section | Description |
|---------|-------------|
| [Getting Started](/en/guide/getting-started) | Environment setup, build instructions, basic usage |
| [Algorithm Guide](/en/guide/algorithms) | Detailed explanations, comparisons, and use cases |
| [API Reference](/en/api/go) | Complete API documentation for Go, Rust, and C++ |
| [Benchmarks](/en/benchmarks/results) | Performance results and methodology |
| [Contributing](/en/guide/contributing) | How to contribute to the project |

## Community

- 💬 Ask questions in [GitHub Discussions](https://github.com/LessUp/compress-kit/discussions)
- 🐛 Report bugs in [GitHub Issues](https://github.com/LessUp/compress-kit/issues)
- 🤝 Read the [Contributing Guide](/en/guide/contributing)

---

**CompressKit** © 2025-2026 LessUp. Released under the MIT License.
