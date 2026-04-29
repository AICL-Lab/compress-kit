---
name: Bug report
about: Report a CompressKit build, codec, or cross-language compatibility bug
title: "[BUG] "
labels: bug
assignees: ""
---

## Summary

Describe the failure in one or two sentences.

## Scope

- Algorithm: Huffman / Arithmetic / Range Coder / RLE / Shared
- Language: C++17 / Go / Rust / Python scripts / Docs
- Surface: CLI / library API / streaming API / conformance / docs / CI

## Reproduction

```bash
# Paste the smallest command sequence that reproduces the issue.
make test-conformance
```

If the issue depends on an input file, include:

- file size:
- file type or pattern:
- whether the file can be shared:

## Expected behavior

What should have happened?

## Actual behavior

What happened instead? Include stderr/stdout when relevant.

```text
paste output here
```

## Environment

- OS:
- C++ compiler:
- Go version:
- Rust version:
- Python version, if scripts are involved:

## Notes

If this affects a binary format or cross-language behavior, mention which
encoder/decoder pair failed, for example `arithmetic cpp -> rust`.
