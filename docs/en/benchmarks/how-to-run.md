# How to Run Benchmarks

This project includes a Python-based benchmark framework that measures encode/decode speed and compression ratios for all algorithms and languages.

## Prerequisites

- Python 3.8+
- All implementations built: `make build`
- Test data generated: `make test-data`

## Run Benchmarks

### All Benchmarks

```bash
make bench
```

This runs the `scripts/run_all_bench.py` script, which:
1. Generates test data (if `tests/data/` is empty)
2. Runs the validated algorithm × language × dataset matrix
3. Measures timing and compression ratio
4. Saves reports to `reports/` and refreshes `docs/.vitepress/data/benchmarks.json`

### Individual Algorithm Benchmarks

```bash
cd algorithms/huffman/benchmark
python3 bench.py

cd algorithms/arithmetic/benchmark
python3 bench.py

cd algorithms/range/benchmark
python3 bench.py

cd algorithms/rle/benchmark
python3 bench.py
```

## Benchmark Configuration

### Test Data

| File | Generation Method | Size |
|------|-------------------|------|
| `tests/data/random_1MiB.bin` | `os.urandom(1024*1024)` | 1 MiB |
| `tests/data/random_10MiB.bin` | `os.urandom(10*1024*1024)` | 10 MiB |
| `tests/data/repetitive_10MiB.bin` | Repeated 256-byte pattern | 10 MiB |
| `tests/data/textlike_10MiB.bin` | Weighted English letters | 10 MiB |
| `tests/data/small_dictionary_like.bin` | Small repeated dictionary-style sample | ~8 KiB |

To regenerate:

```bash
make test-data
# or
python3 tests/gen_testdata.py
```

### Metrics Measured

| Metric | Description |
|--------|-------------|
| Encode time | Wall-clock time to compress the input |
| Decode time | Wall-clock time to restore the original |
| Encode speed | MiB/s = input_size / encode_time |
| Decode speed | MiB/s = input_size / decode_time |
| Compression ratio | output_size / input_size (lower = better) |

### Output Format

Results are saved to `reports/`, and the docs snapshot is written to
`docs/.vitepress/data/benchmarks.json`:

```text
reports/
├── huffman_report_<timestamp>.txt
├── arithmetic_report_<timestamp>.txt
├── range_report_<timestamp>.txt
└── rle_report_<timestamp>.txt

docs/.vitepress/data/
└── benchmarks.json
```

Each report contains:

```
Algorithm: Huffman
Language: C++
Input: 10 MiB random data
Encode: 245 ms (40.8 MiB/s)
Decode: 198 ms (50.5 MiB/s)
Compression ratio: 1.23
```

## Adding New Benchmarks

To add a new test dataset:

1. Edit `tests/gen_testdata.py`
2. Add generation code to `generate_random_file()` or create new generator
3. Run `make test-data`
4. Update `scripts/run_all_bench.py` so the generated docs snapshot stays aligned

## Troubleshooting

### "Binary not found"

```bash
make build  # Rebuild all implementations
```

### "Test data not found"

```bash
make test-data  # Generate test files
```

### Slow benchmark on Range Coder

::: warning
The Range Coder decoder has a known performance issue for files above 500 KiB.
The checked-in benchmark flow therefore uses `small_dictionary_like.bin` for
Range results.
:::

```bash
# Create smaller test file
dd if=tests/data/random_10MiB.bin of=small.bin bs=1024 count=100
```
