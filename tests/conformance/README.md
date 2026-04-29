# Cross-Language Conformance

This directory contains the executable interoperability checks for CompressKit's
current file formats.

## Executable matrix

Run the local decode matrix with:

```bash
make test-conformance
```

The target builds all implementations, generates `tests/data/`, then runs
`run_decode_matrix.py`.

The default matrix verifies every encoder language against every decoder
language for Huffman, Arithmetic Coding, Range Coder, and RLE using the small
generated corpus:

- `empty.bin`
- `single_byte.bin`
- `alternating.bin`
- `small_dictionary_like.bin`

That produces 144 round-trip checks:

```text
4 algorithms x 4 corpus files x 3 encoders x 3 decoders
```

Use `python3 tests/conformance/run_decode_matrix.py --include-large` for a
larger local sweep. Range Coder automatically skips files over 100 KiB to avoid
the documented large-file decode performance issue.

## Deferred frame-era tests

Header parsing, truncation, checksum corruption, and concatenated-frame tests
belong to the proposed shared frame format. They are intentionally not active
in the final baseline until the frame envelope is implemented across C++17, Go,
and Rust.
