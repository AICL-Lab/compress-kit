# Architecture

CompressKit sits above the algorithm cores with a single thin buffer layer:

```text
caller
  -> buffer layer (encode_buffer / decode_buffer, BufferTransform fn pointer)
  -> algorithm core (Huffman / Arithmetic / Range / RLE)
```

## Why It Exists

- file-to-file CLIs and in-memory callers share the same buffer-layer path
- a single `BufferTransform` function-pointer alias normalizes every algorithm behind one signature
- shared size-limit checks and error codes keep the surface small

## Language Host

- C++: `algorithms/shared/cpp`

## Security Boundary

- input accepted by the buffer layer is capped at `4 GiB`
- decoded output produced by the buffer layer is capped at `1 GiB`

## Verification

`make test` runs the shared lifecycle tests followed by the algorithm-specific suites.
