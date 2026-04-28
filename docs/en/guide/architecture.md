# Architecture

CompressKit now uses a two-layer in-memory integration model above the existing algorithm cores:

```text
caller
  -> buffer layer (encode_buffer / decode_buffer)
  -> streaming layer (process / flush / finish / reset)
  -> existing algorithm core
```

## Why It Exists

- file-to-file CLIs now share the same buffer-layer path as in-memory callers
- lifecycle semantics are consistent across C++17, Go, and Rust
- shared error handling makes later frame-format and conformance work easier to build on

## Language Hosts

- C++: `algorithms/shared/cpp`
- Go: `algorithms/shared/go`
- Rust: `algorithms/shared/rust`

## Security Boundary

- input accepted by the streaming layer is capped at `4 GiB`
- decoded output produced by the streaming layer is capped at `1 GiB`

## Verification

`make test` includes shared streaming lifecycle tests for all three languages before the algorithm-specific suites.
