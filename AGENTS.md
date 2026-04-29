# CompressKit Agent Guide

This file is the project-level operating guide for AI agents working in
`LessUp/compress-kit`.

## Project identity

- Product name: **CompressKit**
- Repository: `LessUp/compress-kit`
- Documentation: <https://lessup.github.io/compress-kit/>
- Default branch: `master`

Do not use the old external branding `encoding` or `LessUp/encoding`.

## What this repository is

CompressKit is a multi-language educational compression laboratory. It keeps
four classic lossless algorithms implemented in three languages:

| Algorithm | Languages | Stable CLI |
|-----------|-----------|------------|
| Huffman Coding | C++17, Go, Rust | `<binary> encode|decode <input> <output>` |
| Arithmetic Coding | C++17, Go, Rust | `<binary> encode|decode <input> <output>` |
| Range Coder | C++17, Go, Rust | `<binary> encode|decode <input> <output>` |
| RLE | C++17, Go, Rust | `<binary> encode|decode <input> <output>` |

The core engineering promise is cross-language binary compatibility for each
algorithm's current stable format.

## Source of truth

OpenSpec requirements in `openspec/specs/` are normative:

- `encoding-project`: product scope, algorithms, security limits, branding
- `core-architecture`: directory layout, CLI shape, binary-format architecture
- `cross-language-testing`: conformance, benchmark expectations, known issues

Archived proposals in `openspec/changes/archive/` are historical or deferred
design context. Do not treat deferred archive entries as active implementation
requirements unless a new OpenSpec change is explicitly opened.

## Change policy

Create or update an OpenSpec change before implementing:

- new algorithms
- binary format changes
- public API or CLI contract changes
- cross-language conformance semantics
- CI behavior that changes required quality gates

Small documentation fixes, internal refactors, and bug fixes that preserve the
existing contract may be implemented directly, but still check the relevant spec
first.

## Validation commands

Use the smallest command that proves the change, then run the broader baseline
before completing significant work.

| Command | Purpose |
|---------|---------|
| `make build` | Build every C++/Go/Rust CLI |
| `make test` | Main repository baseline: unit, streaming, conformance |
| `make test-conformance` | Cross-language encode/decode matrix |
| `make lint` | Existing lint path |
| `make format` | Existing formatter path |
| `npm run docs:build` | VitePress documentation build |
| `openspec validate --all` | Validate specs and archived changes |

## Known limitation

Range Coder decode performance is known to degrade on files larger than 500 KiB.
This is documented in the Range Coder docs and `cross-language-testing` spec.
Do not "fix" it unless the requested scope explicitly targets Range Coder
performance. Conformance and benchmark code should cap Range-heavy sweeps at
small inputs.

## Generated artifacts

Do not commit generated outputs:

- `tests/data/*.bin`
- algorithm binaries such as `huffman_cpp`, `huffman_go`, `huffman_rust`
- Rust `target/`
- `reports/`
- `docs/.vitepress/dist/`
- temporary conformance directories under `tests/.conformance-*`

Run `make clean` when you need to inspect a source-only tree.

## Style expectations

- Error messages in code must be English.
- C++17: follow the existing `.clang-format` / Google-style layout.
- Go: use `gofmt` and `go vet`.
- Rust: use `rustfmt` and `cargo clippy` where the existing workflow calls it.
- Python scripts are tooling only; keep them deterministic and dependency-light.
- Documentation is bilingual when user-facing. README files should remain concise
  repository gateways; detailed explanations belong in `docs/`.

## Review checklist for agents

Before finishing a code-affecting change:

1. Check whether `openspec/specs/` needed updates.
2. Verify the relevant language tests.
3. If encoded bytes or decoding behavior changed, run `make test-conformance`.
4. If docs changed, run `npm run docs:build`.
5. Confirm no generated artifacts are staged.
