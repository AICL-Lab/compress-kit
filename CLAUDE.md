# Claude Instructions for CompressKit

Follow `AGENTS.md` first. This file adds Claude-specific reminders for working
efficiently in this repository.

## Operating mode

- Work from the real project state, not from generic compression-library habits.
- Read the relevant OpenSpec requirement before changing code behavior.
- Prefer small, verifiable edits over broad rewrites.
- Keep the current `master` single-mainline flow unless the user explicitly asks
  for branch migration.

## High-value checks

Use these commands as the canonical local proof points:

```bash
openspec validate --all
make test
npm run docs:build
```

For focused work:

```bash
make test-conformance
go test ./algorithms/shared/go/... ./algorithms/huffman/go/... ./algorithms/arithmetic/go/... ./algorithms/range/go/... ./algorithms/rle/go/...
cargo test --manifest-path algorithms/arithmetic/rust/Cargo.toml
```

## Compression-specific guardrails

- Maintain cross-language compatibility inside each algorithm family.
- Do not silently change magic bytes, frequency table layout, endian rules, or
  RLE pair layout.
- Treat Range Coder large-file decode performance as a documented limitation,
  not opportunistic cleanup.
- Keep security limits visible: 4 GiB max input and 1 GiB max decoded output.

## Documentation stance

- README is a short gateway.
- Git Pages is the product/documentation portal.
- OpenSpec is the requirement source of truth.
- Changelog records user-facing changes only; do not use it as an architecture
  diary.

## AI/tooling stance

- Use OpenSpec skills for requirement-level changes.
- Use local search and targeted tests for small bug fixes.
- Use code review skills before merging large cross-language changes.
- Avoid adding new MCP or plugin dependencies unless they clearly reduce future
  token/context cost for this exact repository.
