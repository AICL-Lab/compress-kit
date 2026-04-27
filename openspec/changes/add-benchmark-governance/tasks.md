# Tasks: add-benchmark-governance

## Legend

- [ ] not started
- [x] done

---

## Prerequisites

- [ ] P1. Confirm `add-interoperability-conformance` is archived (corpus files at `tests/corpus/` must exist).

## Phase A — Benchmark Runner Infrastructure

- [ ] A1. Create `tests/bench/run_bench.py` — benchmark runner script; accepts `--algorithm`, `--language`, `--corpus-dir`, `--output` flags.
- [ ] A2. Implement timing loop: 1 warm-up + 5 timed runs; records median wall-clock time.
- [ ] A3. Implement memory measurement helpers per language (see design.md § Memory Measurement).
- [ ] A4. Implement report writer that outputs `bench_<ISO8601>.json` matching design.md schema version 1.0.
- [ ] A5. Create `tests/bench/compare.py` — reads baseline.json and latest result; prints pass/fail per metric per triple; exits 1 if any regression.
- [ ] A6. Create `tests/bench/thresholds.json` with default threshold values and commit it to the repository:
  ```json
  {
    "version": 1,
    "speed_regression_pct": 10,
    "memory_regression_pct": 20
  }
  ```

## Phase B — Makefile Targets

- [ ] B1. Add `make bench` target: runs `run_bench.py` for all (algorithm × language × corpus_file) combinations.
- [ ] B2. Add `make bench-check` target: runs `compare.py` against baseline.json; exits 0 if no baseline exists yet.
- [ ] B3. Document both targets in `Makefile` comments.

## Phase C — Baseline Capture

- [ ] C1. Run `make bench` on CI (GitHub Actions runner) to produce the first result file.
- [ ] C2. Copy result to `tests/bench/baseline.json` and commit.
- [ ] C3. Add `tests/bench/results/` to `.gitignore`.

## Phase D — CI Integration

- [ ] D1. Add `bench-check` step to `.github/workflows/test.yml` after the test stage.
- [ ] D2. Ensure `bench-check` only fails the build if `baseline.json` exists (skip gracefully if absent).
- [ ] D3. Set a 10-minute CI timeout on the bench step.

## Phase E — Threshold Validation

- [ ] E1. After 3 stable CI runs, review whether ±10% speed and +20% memory thresholds produce false positives; adjust in `tests/bench/thresholds.json` if needed.
- [ ] E2. Document threshold adjustment history in `tests/bench/README.md`.

## Phase F — Documentation

- [ ] F1. Create `tests/bench/README.md` covering: how to run benchmarks, how to interpret results, how to update the baseline.
- [ ] F2. Add `docs/en/testing/benchmarks.md` — public-facing benchmark governance documentation.
- [ ] F3. Add `docs/zh/testing/benchmarks.md` (Chinese translation).
- [ ] F4. Note Range Coder cap at 100 KiB in both README.md and benchmark docs.
