# Claude Instructions for CompressKit

Follow `AGENTS.md` first. This file adds Claude-specific reminders.

## Quick Validation

```bash
make test && make lint
```

## Compression Guardrails

- Do not change magic bytes, frequency table layout, endian rules, or RLE pair layout
- RLE magic is `RLE\x00`
- Keep security limits: 4 GiB max input, 1 GiB max decoded output

## Documentation Stance

- README: Short gateway
- Git Pages: Product portal
- Changelog: User-facing changes only
