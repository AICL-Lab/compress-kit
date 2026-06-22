# ADR 0005: Range C++ bench mode migration path

Date: 2026-05-31  
Status: Accepted

## Context

The public CLI contract for all production binaries is:

```
<binary> <encode|decode> <input> <output>
```

Range C++ is the only exception: its production binary also accepts:

```
<binary> bench [size_bytes] [iterations]
```

This is documented in `contract-inventory.md` (§ CLI contract and exceptions) and recorded in `algorithms/range/cpp/main.cpp:378-415`.

The `bench` mode is used by `algorithms/range/benchmark/bench.py` for micro-benchmarking. It is a **shipped CLI behavior** even though it conflicts with the unified CLI adapter shape.

## Decision

The `bench` mode is **retained without change** in the current release cycle. This ADR records the migration path for when it is eventually removed.

### Migration path (to be executed in a future cycle with an OpenSpec change)

1. Open an OpenSpec change to officially deprecate the `bench` subcommand from the production Range C++ binary.
2. Move benchmark behavior to a separate binary or script (e.g., `algorithms/range/benchmark/bench_driver` or extend `algorithms/range/benchmark/bench.py` to drive encoding/decoding directly).
3. Update `algorithms/range/benchmark/bench.py` to use the new driver.
4. Update CLI smoke tests to assert that `bench` is no longer a valid mode in the production binary.
5. Remove the `bench` branch from `algorithms/range/cpp/main.cpp`.

### Current constraints

- Do not silently remove `bench` mode before the above steps; it is a shipped behavior even if undocumented in the primary contract.
- CLI smoke tests currently test the `<encode|decode|invalid_mode>` contract; they do not assert that `bench` is **absent**. This is intentional until the migration is approved.

## Consequences

**Positive**: Migration path is explicit; no accidental removal; future reviewer knows the intent.

**Negative**: Range C++ production binary continues to violate the unified CLI contract until the migration is executed.

**Deferred**: Actual migration work. Requires a new OpenSpec change proposal before any code change.
