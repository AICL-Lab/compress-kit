# Production Readiness Program Design

## Summary

This design turns the broad "make CompressKit production-ready" goal into an
ordered, contract-preserving program with four tracks:

1. Reproducible developer environment and baseline stabilization
2. CLI contract and cross-language test hardening
3. Benchmark truthfulness and release-quality reporting
4. GitHub Pages and documentation delivery hardening

The recommendation is to start with **baseline stabilization**, because the
current repository already has substantial coverage, but the green paths are not
yet trustworthy enough to support production delivery.

## Current Audit Snapshot

The repository already has strong foundations:

- Unified CLI shape is present across shared launchers and binaries
- `make test`, `make lint`, `make bench`, and `npm run docs:build` exist
- GitHub Pages deployment workflow exists
- VitePress site structure is already bilingual and reasonably mature

The main blockers discovered during audit are:

- No `.devcontainer/` exists, so local setup is not reproducible
- `make test` currently fails in `algorithms/shared/go/codec`
- `make lint` fails for the same Go regression
- `make bench` exits successfully, but `scripts/run_all_bench.py` points at
  `ROOT / "huffman" / "benchmark"`-style paths instead of
  `ROOT / "algorithms" / ...`, so benchmark drivers are silently skipped
- Benchmark execution only warns on failures instead of failing the suite
- Docs Pages workflow uses `npm install --no-package-lock`, which is less
  reproducible than a lockfile-driven install
- CLI correctness is validated indirectly through conformance, but there is no
  dedicated smoke suite for every binary's command contract, usage path, and
  basic exit semantics

## Constraints

- Preserve product name: **CompressKit**
- Preserve default branch: `master`
- Preserve docs URL: <https://lessup.github.io/compress-kit/>
- Preserve unified CLI contract:
  `<binary> <encode|decode> <input> <output>`
- Preserve security limits: 4 GiB input, 1 GiB decoded output
- Preserve cross-language binary compatibility
- Preserve documented Range Coder large-file limitation
- Keep error messages in English
- Avoid OpenSpec changes by preserving public contracts

## Approaches Considered

### Approach A: Big-bang hardening across all surfaces

Fix devcontainer, tests, CLI, benchmarks, workflows, and docs in one pass, then
run the full validation stack at the end.

**Pros**

- Fastest route on paper
- Minimizes intermediate states

**Cons**

- High risk of mixing unrelated regressions
- Hard to isolate whether failures come from tooling, algorithm code, docs, or
  release automation
- Weak reviewability in a dirty worktree

### Approach B: Ordered production-readiness program

First restore a truthful baseline, then harden CLI and tests, then make
benchmarks trustworthy, then finish delivery surfaces such as GitHub Pages and
developer environment polish.

**Pros**

- Best fault isolation
- Makes every later green build more meaningful
- Prevents false confidence from broken benchmark or CI plumbing
- Aligns with existing Makefile and workflow entry points

**Cons**

- Requires phased execution
- Some docs polish waits until the underlying automation is trustworthy

**Recommendation:** choose this approach.

### Approach C: Docs-first polish

Finish the GitHub Pages site and developer-facing docs first, then come back to
tests, benchmarks, and environment work.

**Pros**

- Visible user-facing improvement early
- Lower short-term implementation risk

**Cons**

- Produces a polished surface over an untrustworthy engineering baseline
- Leaves "all commands work" and benchmark correctness unresolved

## Recommended Design

### Track 1: Reproducible environment and truthful baseline

**Goal**

Make the repository bootstrap and validate consistently on a clean machine.

**Scope**

- Add a production-quality `.devcontainer/`
- Fix the shared Go frequency-table regression that currently breaks test/lint
- Ensure the devcontainer includes the required toolchain for:
  - C++17 build and tests
  - Go
  - Rust + clippy + rustfmt
  - Python
  - Node.js / npm for docs
  - optional wasm support only if already exercised by repo commands
- Prefer reproducible installs where lockfiles exist

**Design**

- Treat `make test`, `make lint`, and `npm run docs:build` as the baseline gate
- Make the devcontainer explicitly provision every tool needed for those gates
- Keep the fix for the Go regression surgical: resolve duplicate frequency-byte
  parsing definitions and align tests with the current shared API

### Track 2: CLI contract and test hardening

**Goal**

Guarantee that every shipped CLI binary behaves correctly, not just that
cross-language decode happens to work.

**Scope**

- Add dedicated CLI smoke coverage for all algorithms and languages
- Verify:
  - happy-path encode
  - happy-path decode
  - usage error on wrong arity
  - invalid mode handling
  - non-zero exit on failure

**Design**

- Reuse generated corpus and existing binaries
- Put the smoke suite near conformance tooling so it becomes part of the same
  production-readiness story
- Avoid changing the CLI shape itself

### Track 3: Benchmark truthfulness and reporting

**Goal**

Make `make bench` actually measure all supported benchmark drivers and fail when
the benchmark pipeline is broken.

**Scope**

- Fix benchmark discovery paths in `scripts/run_all_bench.py`
- Make missing benchmark drivers or non-zero benchmark exits fail loudly
- Verify report generation in `reports/`
- Keep Range Coder large-file safety behavior aligned with documented limits

**Design**

- Replace silent skips with explicit validation
- Treat benchmark orchestration as release tooling, not best-effort scripting
- Preserve per-algorithm benchmark entry points under `algorithms/*/benchmark/`

### Track 4: GitHub Pages and docs delivery hardening

**Goal**

Make the published documentation site reflect a trustworthy engineering state.

**Scope**

- Keep the existing VitePress information architecture
- Tighten workflow reproducibility
- Ensure benchmark/docs pages match the now-truthful benchmark pipeline
- Add devcontainer/onboarding guidance only after the environment exists

**Design**

- Prefer lockfile-based docs installation in CI when possible
- Update docs only where they directly reflect changed behavior or workflows
- Preserve bilingual navigation and product branding

## Execution Order

1. Track 1: baseline and devcontainer
2. Track 2: CLI contract hardening
3. Track 3: benchmark truthfulness
4. Track 4: GitHub Pages and docs delivery polish
5. Run full repository validation and only then treat the repo as
   production-ready for this change set

## Error Handling

- Do not add silent success-shaped fallbacks
- Missing tools, missing benchmark drivers, and broken validation steps should
  fail explicitly
- CLI smoke tests should assert current exit-code semantics instead of masking
  them
- Benchmark orchestration should stop pretending success when nothing ran

## Validation Strategy

The program is complete only when these existing commands all succeed with
truthful results:

- `make build`
- `make test`
- `make lint`
- `make bench`
- `npm run docs:build`

Additionally, the devcontainer must be able to run the same repository commands
without undocumented host prerequisites.

## Out of Scope

- Changing binary formats
- Changing public CLI semantics
- Fixing the documented Range Coder large-file performance limitation
- Adding new algorithms
