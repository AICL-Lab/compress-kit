# ADR 0001: Validation metadata module shape

Date: 2026-05-31  
Status: Accepted

## Context

Conformance, CLI smoke, test-data generation, and benchmark orchestration scripts each hardcode overlapping facts: algorithm names, binary paths, file extensions, corpus names, Range corpus caps, benchmark dataset choices, language ordering, and report schema fields. This duplication causes drift when adding algorithms, adjusting corpus policy, or changing binary paths.

Current consumers and their metadata shapes are recorded in `docs/architecture/contract-inventory.md` (§ Validation and test metadata consumers, § Benchmark metadata consumers).

## Decision

Introduce a single versioned Python metadata module (`tests/metadata.py` or similar) that exposes:

- Algorithm registry: name, extension, binary paths per language
- Corpus registry: corpus names, sizes, generation parameters, eligibility flags
- Range corpus cap: explicit named constant with documented rationale
- Benchmark job registry: per-algorithm selected input file
- Report schema: field names and ordering for benchmark JSON output

Conformance, smoke, test-data, and benchmark runners become thin adapters over this module. The module is versioned with a `METADATA_VERSION` string; consumers assert the version they were tested against.

## Constraints

- This change does **not** alter algorithm/language matrix, corpus eligibility rules, Range cap value, or benchmark schema: it only consolidates existing facts.
- Any intentional change to Range cap policy, corpus eligibility, or conformance matrix requires an OpenSpec change first (see `docs/architecture/contract-inventory.md` § OpenSpec-triggering decisions).
- Python module must remain compatible with the existing `uv`/`python3` invocation patterns in the Makefile.

## Consequences

**Positive**: Single source of truth for test/bench metadata; adding an algorithm or corpus file requires one edit instead of many; Range cap policy is explicit and findable.

**Negative**: Additional indirection layer; scripts must import the module, requiring Python path management.

**Deferred**: Deciding whether the metadata module is the right place to record Range cap policy permanently, or whether that policy should move into OpenSpec. Record that decision in ADR 0003 before migrating.
