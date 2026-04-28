# Conformance Expansion Plan

This directory tracks the interoperability conformance surface that must be locked before new mainstream algorithms are added.

## Planned Corpus Entries

- random_1MiB.bin
- random_10MiB.bin
- repetitive_10MiB.bin
- textlike_10MiB.bin
- empty.bin
- single_byte.bin
- alternating.bin
- small_dictionary_like.bin

## Required Scenarios

- encode A -> decode B
- parse header without full decode
- truncated payload fails with actionable error
- corrupted checksum fails with actionable error
- concatenated frames obey documented semantics
