# Project Structure

## Overview

```
encoding/
├── huffman/              # Huffman coding
│   ├── cpp/              #   C++ single-file implementation
│   ├── go/               #   Go implementation (go.mod)
│   ├── rust/             #   Rust implementation
│   └── benchmark/        #   Cross-language benchmark
├── arithmetic/           # Arithmetic coding
│   ├── cpp/              #   C++ single-file implementation
│   ├── go/               #   Go implementation
│   ├── rust/             #   Rust implementation
│   └── benchmark/        #   Cross-language benchmark
├── range/                # Range coder
│   ├── cpp/              #   C++ single-file implementation
│   ├── go/               #   Go library + CLI
│   ├── rust/             #   Rust library crate + CLI
│   └── benchmark/        #   Cross-language benchmark
├── rle/                  # Run-length encoding
│   ├── cpp/              #   C++ single-file implementation
│   ├── go/               #   Go implementation
│   ├── rust/             #   Rust implementation
│   └── benchmark/        #   Cross-language benchmark
├── tests/                # Test data generation
│   ├── gen_testdata.py   #   Generate benchmark test files
│   └── data/             #   Generated test data
├── docs/                 # VitePress documentation site
│   ├── .vitepress/       #   VitePress config
│   ├── guide/            #   User guides
│   └── public/           #   Static assets
├── .github/workflows/    # CI + Pages deployment
├── Makefile              # Build/test/bench entry point
├── package.json          # npm scripts for docs
└── go.work               # Go workspace (multi-module)
```

## Language Implementation Conventions

| Language | Version | Build Method | Characteristics |
|----------|---------|--------------|-----------------|
| C++ | C++17 | `g++ -std=c++17 -O2` | Single-file, zero dependencies |
| Go | 1.21+ | Go modules (`go.mod`) | Range Coder provides library API |
| Rust | 1.70+ | Cargo / rustc | Range Coder provides library crate |

## Unified CLI Interface

All implementations follow the same CLI pattern:

```bash
<algorithm>_<lang> encode <input> <output>
<algorithm>_<lang> decode <input> <output>
```

Examples: `huffman_cpp`, `arithmetic_go`, `rangecoder_rust`, `rle_cpp`

## File Format Compatibility

All language implementations of the same algorithm use identical binary formats:

| Algorithm | Magic Header | Extension | Format |
|-----------|--------------|-----------|--------|
| Huffman | `HFMN` | `.huf` | Magic + freq table + bit stream |
| Arithmetic | `AENC` | `.aenc` | Magic + freq table + bit stream |
| Range Coder | `RCNC` | `.rcnc` | Magic + freq table + byte stream |
| RLE | None | `.rle` | (count: 4B LE, value: 1B) pairs |

### Cross-Language Verification Matrix

| Encode ↓ / Decode → | C++ | Go | Rust |
|---------------------|-----|-----|------|
| C++ | ✓ | ✓ | ✓ |
| Go | ✓ | ✓ | ✓ |
| Rust | ✓ | ✓ | ✓ |

## CI/CD Pipeline

| Workflow | File | Trigger | Purpose |
|----------|------|---------|---------|
| CI | `ci.yml` | Push / PR | Build, test, correctness |
| Pages | `pages.yml` | `docs/` change | Deploy documentation |

### CI Job Matrix

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│  build-cpp  │     │  build-go   │     │ build-rust  │
│  Ubuntu     │     │  Ubuntu     │     │  Ubuntu     │
│  macOS      │     │             │     │             │
└──────┬──────┘     └──────┬──────┘     └──────┬──────┘
       │                   │                   │
       └───────────────────┼───────────────────┘
                           │
                    ┌──────▼──────┐
                    │ correctness │
                    │   tests     │
                    └─────────────┘
```

## Security Limits

All implementations enforce:

| Limit | Value | Purpose |
|-------|-------|---------|
| Max input size | 4 GiB | Prevent frequency overflow |
| Max output size | 1 GiB | Prevent decompression bombs |
