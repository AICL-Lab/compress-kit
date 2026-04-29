# CompressKit

<p align="center">
  <img src="docs/public/logo.svg" width="120" alt="CompressKit Logo">
</p>

<p align="center">
  <strong>Classic lossless compression algorithms in C++17, Go, and Rust.</strong>
</p>

<p align="center">
  <a href="https://github.com/LessUp/compress-kit/actions/workflows/ci.yml"><img src="https://github.com/LessUp/compress-kit/actions/workflows/ci.yml/badge.svg" alt="CI Status"></a>
  <a href="https://lessup.github.io/compress-kit/"><img src="https://img.shields.io/badge/Docs-Online-blue?logo=readthedocs&logoColor=white" alt="Documentation"></a>
  <a href="https://opensource.org/licenses/MIT"><img src="https://img.shields.io/badge/License-MIT-green.svg" alt="License"></a>
</p>

<p align="center">
  <b>English</b> | <a href="README.zh-CN.md">简体中文</a> | <a href="https://lessup.github.io/compress-kit/">Documentation</a>
</p>

CompressKit is an educational, verification-focused repository for comparing four
classic compression algorithms across three implementation languages. It is not
a black-box package: the point is to read the implementations, run the same
inputs through each language, and verify that compatible formats stay compatible.

## What is included

| Algorithm | C++17 | Go | Rust | Best fit |
|-----------|------:|---:|-----:|----------|
| Huffman Coding | ✓ | ✓ | ✓ | General text/data and prefix-code learning |
| Arithmetic Coding | ✓ | ✓ | ✓ | Entropy-coding concepts and ratio comparison |
| Range Coder | ✓ | ✓ | ✓ | Arithmetic-coder style implementation comparison |
| Run-Length Encoding | ✓ | ✓ | ✓ | Highly repetitive data and simple format study |

All command-line tools use:

```bash
<binary> <encode|decode> <input> <output>
```

## Quick start

```bash
git clone https://github.com/LessUp/compress-kit.git
cd compress-kit

make build
make test
```

Minimal cross-language check:

```bash
printf "Hello CompressKit\n" > input.txt
./algorithms/huffman/cpp/huffman_cpp encode input.txt output.huf
./algorithms/huffman/go/huffman_go decode output.huf restored.txt
diff input.txt restored.txt
```

## Documentation

| Need | Link |
|------|------|
| Full documentation portal | <https://lessup.github.io/compress-kit/> |
| Setup and first run | <https://lessup.github.io/compress-kit/en/guide/getting-started> |
| Algorithm comparison | <https://lessup.github.io/compress-kit/en/guide/algorithms> |
| API references | <https://lessup.github.io/compress-kit/en/api/streaming> |
| Cross-language testing | <https://lessup.github.io/compress-kit/en/testing/cross-language> |

## Repository shape

```text
algorithms/   # huffman, arithmetic, range, rle; each has cpp/go/rust
tests/        # generated corpus, streaming contracts, conformance matrix
docs/         # VitePress documentation site
openspec/     # project specifications and archived design changes
```

## Engineering baseline

| Command | Purpose |
|---------|---------|
| `make build` | Build all C++/Go/Rust CLI tools |
| `make test` | Run unit, streaming, and cross-language conformance tests |
| `make test-conformance` | Run the executable decode matrix |
| `make bench` | Run benchmark scripts |
| `npm run docs:build` | Build the documentation site |

Known limitation: the Range Coder has documented decode performance issues on
large files; local conformance and benchmark paths cap Range-heavy sweeps
accordingly. See the Range Coder docs before using it for large inputs.

## License

[MIT License](LICENSE) · Copyright © 2025-2026 LessUp
