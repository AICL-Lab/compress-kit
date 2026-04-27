# CompressKit Project Instructions

This file contains project-specific instructions for Claude Code when working on the CompressKit (formerly Encoding) project.

## Project Overview

CompressKit is a **multi-language lossless compression algorithm library** for educational purposes. It implements 4 classic compression algorithms (Huffman, Arithmetic Coding, Range Coder, RLE) in 3 languages (C++17, Go 1.21+, Rust 1.70+) with guaranteed cross-language binary format compatibility.

## Brand Guidelines

- **Project Name**: CompressKit
- **Repository**: `LessUp/compress-kit`
- **Documentation URL**: `https://lessup.github.io/compress-kit/`

Always use `CompressKit` and `LessUp/compress-kit` in documentation and code comments.

## Development Workflow

This project follows **OpenSpec Spec-Driven Development**. See [AGENTS.md](AGENTS.md) for the complete workflow.

### Quick Commands

```bash
make build        # Build all implementations
make test         # Run all tests
make bench        # Run benchmarks
make format       # Format all code
make lint         # Lint all code
make clean        # Clean build artifacts
```

## Known Issues

### Range Coder Performance

The Range Coder has a known decode performance issue for files >500KB. This is documented in:
- `docs/en/algorithms/range.md`
- `docs/zh/algorithms/range.md`
- `openspec/specs/cross-language-testing/spec.md`

Do not attempt to fix this without explicit user request, as it's documented for future improvement.

## Code Style

| Language | Tool | Config |
|----------|------|--------|
| C++ | clang-format | `.clang-format` (Google Style) |
| Go | gofmt | Built-in |
| Rust | rustfmt | Built-in |
| Python | black | PEP 8 |

## Cross-Language Compatibility

All implementations MUST share identical binary file formats. When modifying any algorithm:

1. Update the format spec in `openspec/specs/`
2. Ensure all 3 language implementations match
3. Run cross-language verification tests
4. Update documentation if format changes

## Security Constraints

- Maximum input size: 4 GiB
- Maximum output size: 1 GiB (decompression bomb protection)
- All implementations must validate file sizes

## Documentation

- VitePress site in `docs/`
- Bilingual: English (`docs/en/`) and Chinese (`docs/zh/`)
- Run `npm run docs:dev` for local preview
