# Design: add-benchmark-governance

## Benchmark Corpus

The benchmark corpus is a subset of the interoperability conformance corpus (`tests/corpus/`). Range Coder files are capped at 100 KiB per the known decode performance issue.

| File | Size | Used for algorithms |
|------|------|---------------------|
| `corpus/all_zeros_64k.bin` | 64 KiB | All 4 |
| `corpus/repetitive_1k.bin` | 1 KiB | All 4 |
| `corpus/rand_1k.bin` | 1 KiB | All 4 |
| `corpus/rand_10k.bin` | 10 KiB | All 4 |
| `corpus/rand_100k.bin` | 100 KiB | Huffman, AC, RLE only |
| `corpus/text_canterbury_alice.txt` | ~148 KiB | Huffman, AC, RLE only |

Range Coder benchmarks use only files ≤ 100 KiB.

## Metrics

Four canonical metrics are collected per (algorithm, language, corpus_file) triple:

| Metric | Unit | How measured |
|--------|------|--------------|
| `ratio` | float (compressed / original) | `compressed_size / original_size` |
| `encode_speed_mbps` | MiB/s | `original_size / encode_wall_time_s / (1024*1024)` |
| `decode_speed_mbps` | MiB/s | `original_size / decode_wall_time_s / (1024*1024)` |
| `peak_memory_kib` | KiB | Language-specific (see § Memory Measurement) |

### Timing methodology

- Each encode/decode operation is repeated **5 times**; the **median** wall-clock time is used.
- Warm-up: 1 dry run before timing begins (discarded).
- OS page cache is dropped before each run on Linux (`echo 3 > /proc/sys/vm/drop_caches`) when running as root; otherwise the first-run penalty is accepted.

### Memory Measurement

| Language | Tool | Metric captured |
|----------|------|-----------------|
| C++ | `/usr/bin/time -v` | "Maximum resident set size (kbytes)" |
| Go | `runtime.ReadMemStats` | `HeapInuse + StackInuse` at peak |
| Rust | `/usr/bin/time -v` | "Maximum resident set size (kbytes)" |

Note: Memory measurement methodology varies per language (see table above). The +20% regression threshold is enforced per-language by comparing a language's current run against its own baseline entry — it does not impose a cross-language absolute comparison. Comparing C++ memory vs. Go memory across languages is informational only.

## Report Schema

Each benchmark run produces a JSON report at `tests/bench/results/bench_<ISO8601>.json`.

```json
{
  "schema_version": "1.0",
  "run_at": "2025-01-15T12:34:56Z",
  "git_sha": "abc1234",
  "platform": {
    "os": "linux",
    "arch": "amd64",
    "cpu": "Intel Core i7-12700K",
    "go_version": "1.21.0",
    "rust_version": "1.70.0",
    "cpp_compiler": "g++ 13.1"
  },
  "results": [
    {
      "algorithm": "huffman",
      "language": "cpp",
      "corpus_file": "rand_10k.bin",
      "corpus_size_bytes": 10240,
      "compressed_size_bytes": 10312,
      "ratio": 1.007,
      "encode_speed_mbps": 48.3,
      "decode_speed_mbps": 52.1,
      "peak_memory_kib": 1024,
      "samples": 5
    }
  ]
}
```

Pass/fail evaluation is **not stored in the result file**. `make bench-check` reads the result file and the baseline file and computes pass/fail externally, printing any regressing metrics to stdout.

## Regression Thresholds

Thresholds are defined as maximum allowed **degradation** from baseline:

| Metric | Threshold |
|--------|-----------|
| `ratio` | +5% (compressed size may grow by at most 5% of baseline) |
| `encode_speed_mbps` | −10% (speed may drop by at most 10%) |
| `decode_speed_mbps` | −10% |
| `peak_memory_kib` | +20% |

Initial baseline values are set from the first successful CI run and stored in `tests/bench/baseline.json` (same schema as results, with single result per triple).

### Threshold rationale

- ±10% speed covers normal CI variance; tighter thresholds cause false positives.
- +5% ratio threshold catches accidental format changes that inflate output.
- +20% memory is intentionally loose for Phase 1; will be tightened in a future governance change.

## CI Integration

```
make bench        # runs benchmarks, writes results/bench_<timestamp>.json
make bench-check  # compares latest result against baseline.json, exits 1 if any regression
```

CI pipeline adds a `bench-check` step after the test stage. It does NOT block the build if `baseline.json` does not yet exist (first run).

## Historical Report Storage

```
tests/bench/
├── baseline.json               # committed baseline
├── results/
│   ├── bench_2025-01-15T12:34:56Z.json
│   └── ...
└── README.md                   # explains how to update baseline
```

`results/` is git-ignored to keep repository size bounded; `baseline.json` is committed.

## Updating the Baseline

To update the baseline after an intentional performance change:

```bash
make bench
cp tests/bench/results/bench_<latest>.json tests/bench/baseline.json
git add tests/bench/baseline.json
git commit -m "bench: update baseline after <description>"
```
