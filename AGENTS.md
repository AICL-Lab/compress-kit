# CompressKit Agent Guide

## Project Identity

- Product: **CompressKit** - Multi-language compression laboratory
- Repository: `LessUp/compress-kit`
- Documentation: <https://lessup.github.io/compress-kit/>
- Default branch: `master`

## Core Contract

Four algorithms (Huffman, Arithmetic, Range, RLE) × three languages (C++17, Go, Rust).
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
| `make build` | Build all CLIs |
| `make test` | Full test suite |
| `make test-conformance` | Cross-language matrix |
| `make lint` | All linters |
| `npm run docs:build` | Build documentation |
| `openspec validate --all` | Validate specs |

## Key Constraints

- Maintain cross-language binary format compatibility
- Security limits: 4 GiB input, 1 GiB output
- Range Coder performance limitation on large files is documented
- Error messages in code must be English

## Change Policy

OpenSpec change required for:
- Binary format changes
- New algorithms
- Public API changes
- Cross-language conformance semantics

Small documentation fixes, internal refactors, and bug fixes that preserve existing contract may be implemented directly.

## Reference

See `openspec/specs/` for full requirements.
