# CompressKit

<p align="center">
  <img src="docs/public/logo.svg" width="120" alt="CompressKit Logo">
</p>

<p align="center">
  <strong>Classic lossless compression algorithms in C++17.</strong>
</p>

<p align="center">
  <a href="https://github.com/LessUp/compress-kit/actions/workflows/ci.yml"><img src="https://github.com/LessUp/compress-kit/actions/workflows/ci.yml/badge.svg" alt="CI Status"></a>
  <a href="https://lessup.github.io/compress-kit/"><img src="https://img.shields.io/badge/Docs-Online-blue?logo=readthedocs&logoColor=white" alt="Documentation"></a>
  <a href="https://opensource.org/licenses/MIT"><img src="https://img.shields.io/badge/License-MIT-green.svg" alt="License"></a>
</p>

<p align="center">
  <b>English</b> | <a href="README.zh-CN.md">简体中文</a> | <a href="https://lessup.github.io/compress-kit/">Documentation</a>
</p>

CompressKit is an educational, verification-focused repository for studying four
classic compression algorithms in a single C++17 implementation. It is not
a black-box package: the point is to read the implementations, run the same
inputs through each algorithm, and verify round-trip correctness.

## What is included

| Algorithm | Best fit |
|-----------|----------|
| Huffman Coding | General text/data and prefix-code learning |
| Arithmetic Coding | Entropy-coding concepts and ratio comparison |
| Range Coder | Arithmetic-coder style implementation comparison |
| Run-Length Encoding | Highly repetitive data and simple format study |

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

For a pinned toolchain, open the checked-in `.devcontainer/` in VS Code or
Codespaces. For local docs work:

```bash
npm ci
(cd docs && npm ci)
```

Quick round-trip check:

```bash
printf "Hello CompressKit\n" > input.txt
./build/huffman_cpp encode input.txt output.huf
./build/huffman_cpp decode output.huf restored.txt
diff input.txt restored.txt
```

## Documentation

| Need | Link |
|------|------|
| Full documentation portal | <https://lessup.github.io/compress-kit/> |
| Setup and first run | <https://lessup.github.io/compress-kit/en/guide/getting-started> |
| Algorithm comparison | <https://lessup.github.io/compress-kit/en/guide/algorithms> |
| API references | <https://lessup.github.io/compress-kit/en/api/streaming> |

## Repository shape

```text
algorithms/   # huffman, arithmetic, range, rle; each has cpp/
tests/        # generated corpus, CLI smoke tests
docs/         # VitePress documentation site
```

## Engineering baseline

| Command | Purpose |
|---------|---------|
| `make build` | Build all C++ CLI tools (CMake) |
| `make test` | Run unit tests and CLI smoke tests |
| `make lint` | clang-format dry-run |
| `npm run docs:build` | Build the documentation site |

## License

[MIT License](LICENSE) · Copyright © 2025-2026 LessUp
