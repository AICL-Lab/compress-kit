# Shared Frame Contract Test Plan

This directory freezes the parser-facing shared frame contract before frame writers and frame-aware decoders are implemented.

The first milestone here is parser-only coverage:

- valid minimum layout examples
- negative frames for structural errors
- explicit separation between frame parsing and algorithm payload parsing

Production encoder/decoder integration starts only after these examples have been converted into failing parser tests.
