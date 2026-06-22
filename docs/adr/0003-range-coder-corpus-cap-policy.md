# ADR 0003: Range Coder corpus cap policy

Date: 2026-05-31  
Status: Accepted

## Context

The Range Coder has a known large-file decode performance problem. The conformance matrix currently skips corpus files over **100 KiB** with reason code `range_coder_corpus_cap_100_kib`. This cap is hardcoded in `tests/conformance/run_decode_matrix.py` and documented in:

- `tests/conformance/README.md` (workaround note)
- `docs/en/algorithms/range.md` (known issue section)
- `algorithms/range/benchmark/bench.py` (benchmark default input note)

The cap value and its rationale have never been formally recorded as a policy decision.

## Decision

The 100 KiB conformance skip cap is **retained as-is** and recorded here as a deliberate policy rather than an incidental constant.

Rationale:
- The underlying performance issue is a known limitation documented in OpenSpec (`openspec/specs/cross-language-testing/spec.md`).
- Fixing the performance issue requires a non-trivial algorithm change that may affect binary format or cross-language behavior, triggering a new OpenSpec change.
- Raising the cap without fixing the root cause would cause CI to hang or time out, which is worse than an explicit skip.

The cap is **not permanent**. If the performance issue is resolved, the cap can be raised or removed, but doing so changes test-gate semantics and therefore requires an OpenSpec change before implementation.

## Constraints

- Do not silently raise, lower, or remove the cap without an OpenSpec change.
- Do not treat the cap value (100 KiB) as canonical in user-facing docs; document it as "current workaround" to preserve flexibility.
- When the metadata module (ADR 0001) is implemented, the cap must appear as a named constant with this ADR reference.

## Consequences

**Positive**: Policy is now explicit and findable; reviewers know the cap is intentional; future removal has a clear trigger condition.

**Negative**: The known limitation persists; Range Coder is under-tested on real-world corpus sizes.

**Deferred**: Root-cause fix for Range Coder large-file performance. That work requires a separate OpenSpec proposal and is out of scope for this architecture deepening cycle.
