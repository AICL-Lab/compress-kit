# Claude Instructions for CompressKit

Follow `AGENTS.md` first. This file adds Claude-specific reminders.

## Quick Validation

```bash
make test && make lint
```

## Compression Guardrails

- Do not change magic bytes, frequency table layout, endian rules, or RLE pair layout
- RLE now has magic number `RLE\x00` (added 2026-04-30)
- Range Coder large-file performance is a known limitation
- Keep security limits: 4 GiB max input, 1 GiB max decoded output

## Documentation Stance

- README: Short gateway
- Git Pages: Product portal
- OpenSpec: Requirements source
- Changelog: User-facing changes only

## AI Tooling

- Use OpenSpec skills for requirement-level changes
- Use local search and targeted tests for bug fixes
- Use code review skills before merging cross-language changes
