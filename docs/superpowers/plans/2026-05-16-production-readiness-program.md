# Production Readiness Program Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Restore a truthful engineering baseline, add a reproducible devcontainer, harden every CLI command contract, make `make bench` produce real benchmark data, and align the GitHub Pages site with the now-trustworthy workflows.

**Architecture:** Keep CompressKit's public contracts unchanged and work from the inside out. First fix the broken baseline and create a reproducible local environment, then add explicit CLI smoke coverage, then harden benchmark orchestration so green status means real work ran, and finally update docs workflows and pages to reflect the validated pipeline.

**Tech Stack:** Make, C++17, Go 1.21, Rust 2021, Python 3, Node.js 20, VitePress, GitHub Actions

---

## Scope note

This plan implements the ordered tracks from `docs/superpowers/specs/2026-05-16-production-readiness-program-design.md` in one sequence because they share one release gate:

1. `make build`
2. `make test`
3. `make lint`
4. `make bench`
5. `npm run docs:build`

If work stalls in Track 3 or Track 4, stop and split those tracks into dedicated follow-up plans rather than broadening Track 1 changes.

## File structure

**Create**

- `.devcontainer/devcontainer.json` - container definition for a reproducible CompressKit workspace
- `.devcontainer/Dockerfile` - installs the repo's required C++/Go/Rust/Python/Node toolchain
- `tests/conformance/run_cli_smoke.py` - repository-level CLI smoke suite for every algorithm/language binary

**Modify**

- `Makefile` - wire the CLI smoke suite into `make test`; keep existing repository commands intact
- `algorithms/shared/go/codec/frequency.go` - keep the canonical frequency-table byte helpers and exact-count validation
- `algorithms/shared/go/codec/frequency_bytes.go` - stop redeclaring byte-read helpers and delegate to the canonical shared functions
- `algorithms/shared/go/codec/frequency_bytes_test.go` - align tests with the shared API and keep exact-count coverage
- `scripts/run_all_bench.py` - fix benchmark discovery, fail loudly on missing drivers, and emit trustworthy summary artifacts
- `docs/.vitepress/data/benchmarks.json` - refresh published benchmark data from the benchmark pipeline output
- `.github/workflows/ci.yml` - keep the main validation gate aligned with the strengthened test suite
- `.github/workflows/ci-docs.yml` - use reproducible docs dependency installation
- `.github/workflows/docs-pages.yml` - use reproducible docs dependency installation and keep Pages deployment stable
- `README.md` - document the devcontainer and truthful benchmark workflow at the gateway level
- `README.zh-CN.md` - same gateway-level updates in Chinese
- `docs/en/guide/getting-started.md` - add devcontainer setup guidance and clarify benchmark workflow
- `docs/zh/guide/getting-started.md` - same guide updates in Chinese
- `docs/en/benchmarks/results.md` - explain the benchmark data source and make the chart/file the source of truth
- `docs/zh/benchmarks/results.md` - same benchmark-page updates in Chinese

**Test**

- `algorithms/shared/go/codec/frequency_bytes_test.go`
- `tests/conformance/run_cli_smoke.py`
- existing repository gates via `make test`, `make lint`, `make bench`, and `npm run docs:build`

### Task 1: Restore the broken baseline and add a reproducible devcontainer

**Files:**
- Create: `.devcontainer/devcontainer.json`
- Create: `.devcontainer/Dockerfile`
- Modify: `algorithms/shared/go/codec/frequency.go`
- Modify: `algorithms/shared/go/codec/frequency_bytes.go`
- Modify: `algorithms/shared/go/codec/frequency_bytes_test.go`
- Test: `algorithms/shared/go/codec/frequency_bytes_test.go`

- [ ] **Step 1: Write the failing targeted Go tests**

```go
func TestReadFrequenciesFromBytesExactRejectsUnexpectedCount(t *testing.T) {
	out := []byte{}
	AppendFrequencies(&out, []uint32{1, 2, 3})

	pos := 0
	_, err := ReadFrequenciesFromBytesExact(out, &pos, 4)
	if err == nil {
		t.Fatal("expected corrupt error for mismatched count")
	}
	if err.Error() != "invalid frequency table size: 3" {
		t.Fatalf("err = %q, want %q", err.Error(), "invalid frequency table size: 3")
	}
}
```

- [ ] **Step 2: Run the Go tests to capture the current failure**

Run: `go test ./algorithms/shared/go/...`
Expected: FAIL with duplicate `ReadFrequenciesFromBytes` declarations and outdated test calls.

- [ ] **Step 3: Implement the minimal Go fix and add the devcontainer files**

```go
// algorithms/shared/go/codec/frequency_bytes.go
package codec

// AppendFrequencies appends a frequency table using the shared LE wire format.
func AppendFrequencies(out *[]byte, freq []uint32) {
	WriteFrequenciesToBytes(out, freq)
}

// ReadFrequenciesFromBytesExact keeps the exact-count helper used by focused tests.
func ReadFrequenciesFromBytesExact(in []byte, pos *int, expectedCount int) ([]uint32, error) {
	return ReadFrequenciesFromBytes(in, pos, uint32(expectedCount))
}
```

```json
// .devcontainer/devcontainer.json
{
  "name": "CompressKit",
  "build": {
    "dockerfile": "Dockerfile"
  },
  "remoteUser": "vscode",
  "customizations": {
    "vscode": {
      "settings": {
        "terminal.integrated.defaultProfile.linux": "bash"
      }
    }
  },
  "postCreateCommand": "npm install && (cd docs && npm ci)"
}
```

```dockerfile
# .devcontainer/Dockerfile
FROM mcr.microsoft.com/devcontainers/base:ubuntu-24.04

RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get install -y --no-install-recommends build-essential clang-format python3 python3-pip curl ca-certificates \
    && rm -rf /var/lib/apt/lists/*

ARG GO_VERSION=1.21.13
RUN curl -fsSL "https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz" | tar -C /usr/local -xz
ENV PATH="/usr/local/go/bin:${PATH}"

USER vscode
RUN curl https://sh.rustup.rs -sSf | sh -s -- -y --profile minimal --component clippy rustfmt
ENV PATH="/home/vscode/.cargo/bin:${PATH}"
```

- [ ] **Step 4: Re-run the baseline gates affected by this task**

Run: `make test && make lint`
Expected: PASS with the shared Go codec package compiling cleanly and no duplicate-symbol error.

- [ ] **Step 5: Commit**

```bash
git add .devcontainer/devcontainer.json .devcontainer/Dockerfile algorithms/shared/go/codec/frequency.go algorithms/shared/go/codec/frequency_bytes.go algorithms/shared/go/codec/frequency_bytes_test.go
git commit -m "fix: restore baseline validation and add devcontainer" -m "Co-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>"
```

### Task 2: Add explicit CLI smoke coverage for every shipped binary

**Files:**
- Create: `tests/conformance/run_cli_smoke.py`
- Modify: `Makefile`
- Test: `tests/conformance/run_cli_smoke.py`

- [ ] **Step 1: Write the failing CLI smoke runner**

```python
ALGORITHMS = {
    "huffman": [
        ROOT / "algorithms/huffman/cpp/huffman_cpp",
        ROOT / "algorithms/huffman/go/huffman_go",
        ROOT / "algorithms/huffman/rust/huffman_rust",
    ],
    # arithmetic, range, rle follow the same pattern
}

def assert_usage(binary: Path) -> None:
    proc = subprocess.run([str(binary)], text=True, capture_output=True)
    if proc.returncode == 0:
        raise RuntimeError(f"{binary} unexpectedly succeeded without args")
    if "Usage:" not in (proc.stdout + proc.stderr):
        raise RuntimeError(f"{binary} did not print usage")
```

- [ ] **Step 2: Run the smoke runner before wiring it into `make test`**

Run: `python3 tests/conformance/run_cli_smoke.py`
Expected: FAIL until the script is finished and the repo has fresh binaries to inspect.

- [ ] **Step 3: Finish the runner and wire it into the existing test target**

```make
test: test-data \
       test-shared-cpp test-shared-go test-shared-rust \
       test-huffman-go test-arithmetic-go test-range-go test-rle-go \
       test-huffman-rust test-arithmetic-rust test-range-rust test-rle-rust \
       test-conformance test-cli-smoke

test-cli-smoke: build test-data
	python3 tests/conformance/run_cli_smoke.py
```

```python
def assert_round_trip(binary: Path, source: Path, encoded: Path, decoded: Path) -> None:
    run_checked([str(binary), "encode", str(source), str(encoded)])
    run_checked([str(binary), "decode", str(encoded), str(decoded)])
    if source.read_bytes() != decoded.read_bytes():
        raise RuntimeError(f"{binary} round-trip mismatch")
```

- [ ] **Step 4: Re-run the repository test gate**

Run: `make test`
Expected: PASS with the new `test-cli-smoke` target exercising wrong-arity, invalid-mode, encode, and decode flows.

- [ ] **Step 5: Commit**

```bash
git add Makefile tests/conformance/run_cli_smoke.py
git commit -m "test: add cli smoke coverage for all binaries" -m "Co-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>"
```

### Task 3: Make `make bench` truthful and publish real benchmark data

**Files:**
- Modify: `scripts/run_all_bench.py`
- Modify: `docs/.vitepress/data/benchmarks.json`
- Test: `scripts/run_all_bench.py`

- [ ] **Step 1: Make the benchmark orchestrator fail on false-green conditions**

```python
benchmarks = {
    "huffman": ROOT / "algorithms" / "huffman" / "benchmark" / "bench.py",
    "arithmetic": ROOT / "algorithms" / "arithmetic" / "benchmark" / "bench.py",
    "range": ROOT / "algorithms" / "range" / "benchmark" / "bench.py",
    "rle": ROOT / "algorithms" / "rle" / "benchmark" / "bench.py",
}

missing = [name for name, path in benchmarks.items() if not path.is_file()]
if missing:
    raise SystemExit(f"missing benchmark driver(s): {', '.join(missing)}")
```

- [ ] **Step 2: Run the benchmark command to confirm the current false-green behavior**

Run: `make bench`
Expected: PASS even when no per-algorithm benchmark script actually ran, proving the current orchestration is not trustworthy.

- [ ] **Step 3: Implement real discovery, hard failure, and docs data export**

```python
def emit_docs_benchmark_data(results: list[dict[str, object]]) -> None:
    target = ROOT / "docs" / ".vitepress" / "data" / "benchmarks.json"
    payload = {
        "generated": datetime.date.today().isoformat(),
        "version": "1.0.0",
        "results": results,
    }
    target.write_text(json.dumps(payload, indent=2) + "\n", encoding="utf-8")
```

```python
if proc.returncode != 0:
    raise SystemExit(f"{title} failed with exit code {proc.returncode}; see {report_path}")
```

- [ ] **Step 4: Re-run the benchmark gate**

Run: `make bench`
Expected: PASS only after all benchmark drivers run, reports are written to `reports/`, and `docs/.vitepress/data/benchmarks.json` is refreshed.

- [ ] **Step 5: Commit**

```bash
git add scripts/run_all_bench.py docs/.vitepress/data/benchmarks.json
git commit -m "fix: make benchmark orchestration truthful" -m "Co-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>"
```

### Task 4: Harden GitHub Pages and docs around the validated workflow

**Files:**
- Modify: `.github/workflows/ci.yml`
- Modify: `.github/workflows/ci-docs.yml`
- Modify: `.github/workflows/docs-pages.yml`
- Modify: `README.md`
- Modify: `README.zh-CN.md`
- Modify: `docs/en/guide/getting-started.md`
- Modify: `docs/zh/guide/getting-started.md`
- Modify: `docs/en/benchmarks/results.md`
- Modify: `docs/zh/benchmarks/results.md`

- [ ] **Step 1: Update the docs workflows to use reproducible installs**

```yaml
- name: Install docs dependencies
  working-directory: docs
  run: npm ci
```

- [ ] **Step 2: Refresh the user-facing docs around devcontainer and benchmark truth**

```md
## Dev Container

If you want a reproducible local environment, open the repository in the bundled
`.devcontainer/` and run the standard validation commands from inside the
container:

~~~bash
make test
make lint
npm run docs:build
~~~
```

- [ ] **Step 3: Align the benchmark pages with the generated data model**

```md
::: warning Source of truth
The interactive chart reads `docs/.vitepress/data/benchmarks.json`, which is
refreshed by `make bench`. The narrative below summarizes the latest generated
run, but hardware-specific numbers should always be reproduced locally.
:::
```

- [ ] **Step 4: Re-run the final delivery gates**

Run: `make build && make test && make lint && make bench && npm run docs:build`
Expected: PASS, with docs build succeeding after workflow-aligned dependency installation and benchmark pages reflecting the generated data contract.

- [ ] **Step 5: Commit**

```bash
git add .github/workflows/ci.yml .github/workflows/ci-docs.yml .github/workflows/docs-pages.yml README.md README.zh-CN.md docs/en/guide/getting-started.md docs/zh/guide/getting-started.md docs/en/benchmarks/results.md docs/zh/benchmarks/results.md
git commit -m "docs: harden release workflows and delivery docs" -m "Co-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>"
```

## Self-review checklist

- Track 1 covers the failing Go baseline and missing devcontainer
- Track 2 covers explicit CLI command validation for every binary
- Track 3 covers the false-green benchmark path and published data refresh
- Track 4 covers Pages/docs reproducibility and documentation updates
- No task changes binary formats, CLI shape, or documented Range Coder limits
