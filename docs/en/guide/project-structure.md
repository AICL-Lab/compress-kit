# Project Structure

CompressKit is organized around algorithms first, then languages. This keeps the
same algorithm easy to compare across C++17, Go, and Rust without hiding
language-specific conventions.

## Source layout

```text
algorithms/
├── shared/        # streaming and buffer API foundations
├── huffman/       # cpp/, go/, rust/
├── arithmetic/    # cpp/, go/, rust/
├── range/         # cpp/, go/, rust/
└── rle/           # cpp/, go/, rust/

tests/
├── gen_testdata.py
├── streaming_api_contract/
└── conformance/

docs/              # VitePress site: root portal + en/ + zh/
openspec/          # stable specs and archived design changes
```

## Responsibility boundaries

| Area | Owns | Does not own |
|------|------|--------------|
| `algorithms/<algo>/<lang>/` | Algorithm implementation, CLI entrypoint, language tests | Global docs or cross-language orchestration |
| `algorithms/shared/` | Streaming lifecycle, buffer convenience APIs, shared contract tests | Algorithm-specific file formats |
| `tests/conformance/` | Cross-language encode/decode matrix for stable formats | Future shared-frame validation |
| `tests/data/` | Generated local corpus from `make test-data` | Source-controlled fixtures |
| `docs/` | User-facing guide, API notes, known limitations | OpenSpec change tracking |
| `openspec/` | Normative requirements and archived proposal history | Marketing copy |

## Binary formats

The current terminal baseline keeps per-algorithm formats stable:

| Algorithm | Magic/header | Extension | Payload |
|-----------|--------------|-----------|---------|
| Huffman | `HFMN` + frequency table | `.huf` | Bit stream |
| Arithmetic | `AENC` + frequency table | `.aenc` | Bit stream |
| Range Coder | `RCNC` + frequency table | `.rcnc` | Byte stream |
| RLE | No magic header | `.rle` | `(count: uint32 LE, value: byte)` pairs |

Future shared-frame proposals are archived under `openspec/changes/archive/` and
are not part of the active file format contract.

## Generated artifacts

Build outputs and generated data are intentionally ignored:

- algorithm binaries such as `huffman_cpp`, `huffman_go`, `huffman_rust`
- Rust `target/` directories
- `tests/data/*.bin`
- benchmark reports and temporary conformance directories
- `docs/.vitepress/dist/`

Use `make clean` before packaging or reviewing repository shape.

## Related pages

- [Getting Started](/en/guide/getting-started)
- [Streaming API](/en/api/streaming)
- [Cross-Language Testing](/en/testing/cross-language)
