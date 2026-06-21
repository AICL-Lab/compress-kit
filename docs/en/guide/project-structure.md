# Project Structure

CompressKit is organized around algorithms first. Each algorithm keeps a
single-file C++17 core with a shared streaming/buffer facade, so the same
algorithm is easy to compare and maintain without language-specific forks.

## Source layout

```text
algorithms/
├── shared/        # streaming and buffer API foundations
├── huffman/       # cpp/
├── arithmetic/    # cpp/
├── range/         # cpp/
└── rle/           # cpp/

tests/
├── gen_testdata.py
└── streaming_api_contract/

docs/              # VitePress site: root portal + en/ + zh/
```

## Responsibility boundaries

| Area | Owns | Does not own |
|------|------|--------------|
| `algorithms/<algo>/cpp/` | Algorithm implementation, CLI entrypoint, algorithm tests | Global docs orchestration |
| `algorithms/shared/` | Streaming lifecycle, buffer convenience APIs, shared contract tests | Algorithm-specific file formats |
| `tests/data/` | Generated local corpus from `make test-data` | Source-controlled fixtures |
| `docs/` | User-facing guide, API notes, known limitations | Marketing copy |

## Binary formats

The current terminal baseline keeps per-algorithm formats stable:

| Algorithm | Magic/header | Extension | Payload |
|-----------|--------------|-----------|---------|
| Huffman | `HFMN` + frequency table | `.huf` | Bit stream |
| Arithmetic | `AENC` + frequency table | `.aenc` | Bit stream |
| Range Coder | `RCNC` + frequency table | `.rcnc` | Byte stream |
| RLE | `RLE\x00` | `.rle` | Magic + `(count: uint32 LE, value: byte)` pairs |

## Generated artifacts

Build outputs and generated data are intentionally ignored:

- algorithm binaries under `build/` (e.g. `build/huffman_cpp`)
- `tests/data/*.bin`
- benchmark reports and temporary directories
- `docs/.vitepress/dist/`

Use `make clean` before packaging or reviewing repository shape.

## Related pages

- [Getting Started](/en/guide/getting-started)
- [Streaming API](/en/api/streaming)
