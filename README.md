# Encoding

[![CI](https://github.com/LessUp/encoding/actions/workflows/ci.yml/badge.svg)](https://github.com/LessUp/encoding/actions/workflows/ci.yml)
[![Deploy Docs](https://github.com/LessUp/encoding/actions/workflows/pages.yml/badge.svg)](https://github.com/LessUp/encoding/actions/workflows/pages.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
![C++](https://img.shields.io/badge/C++-17-blue.svg)
![Go](https://img.shields.io/badge/Go-1.21+-00ADD8.svg)
![Rust](https://img.shields.io/badge/Rust-1.70+-orange.svg)

English | [简体中文](README.zh-CN.md) | [Docs](https://lessup.github.io/encoding/)

Encoding is a multi-language collection of classic compression algorithms for learning, implementation comparison, and cross-language verification.

## Repository Overview

- Four algorithms: Huffman, Arithmetic Coding, Range Coder, and RLE
- Three language tracks: C++17, Go, and Rust
- Unified CLI conventions and shared binary formats for cross-language validation
- Dedicated docs site for getting started, algorithm guides, and project structure

## Quick Start

```bash
make build
make test
make bench
```

If you want to start from a single algorithm first:

```bash
cd huffman/cpp
g++ -std=c++17 -O2 main.cpp -o huffman_cpp
./huffman_cpp encode input.bin output.huf
./huffman_cpp decode output.huf restored.bin
```

## Read Next

- [Documentation Home](https://lessup.github.io/encoding/)
- [Getting Started](https://lessup.github.io/encoding/guide/getting-started)
- [Algorithms Guide](https://lessup.github.io/encoding/guide/algorithms)
- [Project Structure](https://lessup.github.io/encoding/guide/project-structure)

## License

MIT License.
