# Cross-Language Testing — Delta Spec (add-benchmark-governance)

<!-- DELTA TYPE: MODIFIED REQ-TEST-003, ADDED REQ-TEST-012, REQ-TEST-013 -->
<!-- REQ-TEST-003 is strengthened below; REQ-TEST-012 and REQ-TEST-013 are new. -->
<!-- These changes are applied to openspec/specs/cross-language-testing/spec.md when this change is archived. -->

---

### REQ-TEST-003: Benchmark Suite (MODIFIED — replaces existing REQ-TEST-003)

Performance benchmarks SHALL run across all implementations using a fixed corpus and SHALL produce a structured, machine-readable report.

#### Scenario: Benchmark corpus is fixed and reproducible
- **GIVEN** `make bench`
- **WHEN** executed in any environment
- **THEN** benchmarks SHALL use only the files listed in the benchmark corpus table (design.md § Benchmark Corpus)
- **AND** SHALL NOT use randomly generated files at benchmark time

#### Scenario: Benchmark execution
- **GIVEN** `make bench` command
- **WHEN** executing benchmarks
- **THEN** all implementations SHALL report metrics for all applicable corpus files

#### Scenario: Metric collection
- **GIVEN** a benchmark run
- **WHEN** collecting results
- **THEN** compression ratio, encode speed (MiB/s), decode speed (MiB/s), and peak memory (KiB) SHALL all be recorded per (algorithm, language, corpus_file) triple

#### Scenario: Report schema conformance
- **GIVEN** a completed benchmark run
- **WHEN** the output JSON is read
- **THEN** it SHALL conform to schema version `1.0` (defined in design.md § Report Schema)
- **AND** each result entry SHALL include `algorithm`, `language`, `corpus_file`, `ratio`, `encode_speed_mbps`, `decode_speed_mbps`, `peak_memory_kib`
- **AND** each intentionally skipped benchmark entry, if any, SHALL include `algorithm`, `language`, `corpus_file`, and `reason`

---

### REQ-TEST-012: Benchmark Regression Gating

CI SHALL fail if benchmark metrics regress beyond defined thresholds.

#### Scenario: Regression detected
- **GIVEN** a `baseline.json` exists in `tests/bench/`
- **WHEN** `make bench-check` is run after a code change
- **THEN** if any metric exceeds its threshold (ratio +5%, speeds −10%, memory +20%) vs. baseline
- **THEN** the command SHALL exit with code 1 and print the regressing metric(s)

#### Scenario: No baseline — graceful skip
- **GIVEN** `tests/bench/baseline.json` does not exist
- **WHEN** `make bench-check` is run
- **THEN** the command SHALL exit 0 with a message indicating no baseline to compare

#### Scenario: Range Coder corpus cap
- **GIVEN** benchmark run for Range Coder algorithm
- **WHEN** any corpus file exceeds 100 KiB
- **THEN** that file SHALL be skipped for Range Coder benchmarks
- **AND** the report SHALL record one skip entry for each skipped `(algorithm, language, corpus_file)` triple with reason `range_coder_corpus_cap_100_kib`

---

### REQ-TEST-013: Benchmark Baseline Management

A committed baseline SHALL track the approved performance envelope.

#### Scenario: Baseline is committed
- **GIVEN** the project repository
- **WHEN** checking `tests/bench/baseline.json`
- **THEN** a valid baseline file conforming to schema version 1.0 SHALL be present

#### Scenario: Intentional performance change updates baseline
- **GIVEN** an intentional performance improvement or architectural change
- **WHEN** the maintainer runs `make bench` and copies the result to `baseline.json`
- **THEN** the commit message SHALL document the reason for the baseline update
