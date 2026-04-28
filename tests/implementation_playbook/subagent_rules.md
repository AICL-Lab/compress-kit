# Subagent Execution Rules

Do not parallelize language implementations until:

- frame flags are frozen
- failing conformance tests exist
- one reference implementation can pass the baseline tests

## First Recommended Parallel Split

- Main agent: spec + integration + CI wiring
- Subagent A: Go implementation
- Subagent B: Rust implementation
- Subagent C: docs + benchmark + review notes
