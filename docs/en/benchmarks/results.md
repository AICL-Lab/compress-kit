# Benchmark Results

This page is a viewer for the generated benchmark snapshot committed with the
docs.

::: tip Local numbers win
The interactive chart reads `docs/.vitepress/data/benchmarks.json`, which is
rewritten by `make bench`. Treat that generated JSON file and the chart below as
the source of truth for the current repository snapshot.
:::

## Interactive Performance Chart

<BenchmarkChart />

## Generated Data Model

| JSON location | Fields | Meaning |
|---------------|--------|---------|
| top level | `generated` | Date of the benchmark snapshot used by the docs |
| each `results[]` row | `algorithm`, `dataset` | Benchmark coordinate shown in the chart |
| each `results[]` row | `encodeTime`, `decodeTime` | Wall-clock milliseconds |
| each `results[]` row | `encodeSpeed`, `decodeSpeed` | Throughput in MiB/s |
| each `results[]` row | `compressionRatio`, `throughput` | Size ratio plus the coarse throughput label used by the UI |

## Current Dataset Coverage

| Dataset key | Chart label | Why it appears |
|-------------|-------------|----------------|
| `textlike_10MiB` | Text-like (10 MiB) | Main comparison input for Huffman and Arithmetic |
| `repetitive_10MiB` | Repetitive (10 MiB) | High-repeat input where RLE is meaningful |
| `small_dictionary_like` | Small dictionary-like sample | Reduced Range workload kept below the known large-file decode limitation |

Because the chart reads the generated JSON directly, local reruns replace the
snapshot without any manual table editing on this page.

## Reading the Snapshot Safely

- Compare results on the **same dataset** before drawing conclusions.
- Range results are intentionally shown on a smaller input; do not read them as
  apples-to-apples throughput versus the 10 MiB datasets.
- RLE is most meaningful on repetitive inputs and may expand other data.

## Refresh the Source of Truth

```bash
make bench
npm run docs:build
```

- `make bench` rewrites `reports/` and `docs/.vitepress/data/benchmarks.json`
- `npm run docs:build` verifies that the docs portal renders the refreshed
  dataset
- If you commit benchmark updates, commit the generated JSON file rather than
  hand-editing narrative numbers on this page

## See Also

- [How to Run Benchmarks](/en/benchmarks/how-to-run) — Detailed instructions
- [Algorithm Guide](/en/guide/algorithms) — Comparison and selection guide
