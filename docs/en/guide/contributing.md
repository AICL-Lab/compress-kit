# Contributing Guide

CompressKit is a spec-driven, multi-language compression repository. A change is not complete when one implementation works; it is complete when the relevant specification, all affected language implementations, and the cross-language binary contract agree.

## Start from OpenSpec

Use `openspec/specs/` as the source of truth:

| Spec | When it matters |
|------|-----------------|
| `encoding-project` | Algorithm scope, quality gates, security limits, public positioning |
| `core-architecture` | Directory layout, CLI shape, binary format boundaries |
| `cross-language-testing` | Compatibility matrix, fixtures, benchmark expectations, known limitations |

Create an OpenSpec change before adding a new algorithm, changing a binary format, changing CLI behavior, or widening the compatibility contract. Small documentation fixes and implementation-only bug fixes can update the existing spec directly if they only restore documented behavior.

## Development baseline

| Command | Purpose |
|---------|---------|
| `make build` | Compile C++, Go, and Rust implementations |
| `make test` | Run unit tests, shell tests, and the cross-language conformance matrix |
| `make test-conformance` | Run only the encode/decode compatibility matrix |
| `make lint` | Run `go vet` and strict Rust `clippy` with warnings denied |
| `make format` | Run `gofmt`, `cargo fmt`, and `clang-format` |
| `npm run docs:build` | Build the VitePress documentation site |

`make lint` is intentionally strict. Do not hide linter failures with shell fallbacks; either fix the issue or document why a specific lint cannot apply.

## Implementation standards

| Language | Expectations |
|----------|--------------|
| C++17 | Keep the single-file algorithm CLIs compatible with the shared format; use `.clang-format` before submitting changes. |
| Go 1.21+ | Use `gofmt`, `go vet`, and idiomatic package-level tests. |
| Rust 1.70+ | Keep `cargo test`, `cargo fmt`, and `cargo clippy --all-targets -- -D warnings` clean for each crate. |
| Python 3.8+ | Use Python for repository scripts and conformance orchestration, not as a production algorithm target. |

## Binary compatibility rules

- Every algorithm CLI must preserve `encode|decode input output`.
- Huffman, Arithmetic, Range, and RLE formats must remain compatible across C++, Go, and Rust unless an approved OpenSpec change says otherwise.
- Security limits are part of the contract: maximum input is 4 GiB and maximum decoded output is 1 GiB.
- The Range Coder large-file decode performance issue is documented and should not be treated as an incidental cleanup task.

## Pull request checklist

- The relevant OpenSpec requirement is still true, or the PR includes the spec change.
- `make test` passes locally.
- `make lint` and `npm run docs:build` pass when the touched files require them.
- Cross-language fixtures cover any binary-format or streaming-adapter behavior change.
- Documentation is updated only where it helps a reader choose, use, or validate the project.
