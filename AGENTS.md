# CompressKit Agent Guide

## Project Identity

- Product: **CompressKit** - C++17 compression laboratory
- Repository: `LessUp/compress-kit`
- Documentation: <https://lessup.github.io/compress-kit/>
- Default branch: `master`

## Core Contract

Four algorithms (Huffman, Arithmetic, Range, RLE) implemented in C++17.
Binary format compatibility is the primary constraint.

**Magic Numbers**:
| Algorithm | Magic |
|-----------|-------|
| Huffman | `HFMN` |
| Arithmetic | `AENC` |
| Range Coder | `RCNC` |
| RLE | `RLE\x00` |

## Validation Commands

| Command | Purpose |
|---------|---------|
| `make build` | Build all C++ CLIs (CMake) |
| `make test` | Unit tests + CLI smoke tests |
| `make lint` | clang-format dry-run |
| `npm run docs:build` | Build documentation |

## Key Constraints

- Maintain binary format compatibility (magic bytes, frequency table LE layout, RLE pair layout)
- Security limits: 4 GiB input, 1 GiB output
- Error messages in code must be English

## Change Policy

Binary format changes and new algorithms require careful design review.
Internal refactors and bug fixes that preserve existing contract may be
implemented directly.
