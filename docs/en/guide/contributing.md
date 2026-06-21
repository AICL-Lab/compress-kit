# Contributing Guide

CompressKit is a C++17 compression repository. A change is complete when the implementation works and the documented binary format contract is preserved.

## Development baseline

| Command | Purpose |
|---------|---------|
| `make build` | Compile the C++ implementations (wraps CMake) |
| `make test` | Run unit tests and streaming contract tests |
| `make lint` | Run `clang-format` checks |
| `make format` | Run `clang-format` |
| `npm run docs:build` | Build the VitePress documentation site |

`make lint` is intentionally strict. Do not hide linter failures with shell fallbacks; either fix the issue or document why a specific lint cannot apply.

## Implementation standards

| Language | Expectations |
|----------|--------------|
| C++17 | Keep the single-file algorithm CLIs compatible with the shared format; use `.clang-format` before submitting changes. |
| Python 3.8+ | Use Python for repository scripts and benchmark orchestration, not as a production algorithm target. |

## Binary compatibility rules

- Every algorithm CLI must preserve `encode|decode input output`.
- Huffman, Arithmetic, Range, and RLE formats must remain stable.
- Security limits are part of the contract: maximum input is 4 GiB and maximum decoded output is 1 GiB.
- The Range Coder large-file decode performance issue is documented and should not be treated as an incidental cleanup task.

## Pull request checklist

- `make test` passes locally.
- `make lint` and `npm run docs:build` pass when the touched files require them.
- Documentation is updated only where it helps a reader choose, use, or validate the project.
