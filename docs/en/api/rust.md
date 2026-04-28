# Rust Crate API Reference

See also: [Streaming API](/en/api/streaming)

Each Rust implementation provides a reusable library crate, and all four algorithms now expose shared streaming adapters plus the `compresskit-codec` buffer helpers.

## Huffman (`algorithms/huffman/rust`)

### Add as Dependency

```toml
[dependencies]
huffman = { path = "../algorithms/huffman/rust" }
```

### API

#### `huffman_encode_file(input: &str, output: &str) -> io::Result<()>`

Encodes a file using Huffman coding.

```rust
use huffman;

fn main() -> std::io::Result<()> {
    huffman::huffman_encode_file("input.bin", "output.huf")?;
    huffman::huffman_decode_file("output.huf", "decoded.bin")?;
    Ok(())
}
```

### Internal Functions

- `compress_file(input: &str) -> io::Result<Vec<u8>>` — returns encoded bytes
- `decompress_file(input: &str) -> io::Result<Vec<u8>>` — returns decoded bytes

---

## Arithmetic (`algorithms/arithmetic/rust`)

### API

```rust
use arithmetic;

fn main() -> std::io::Result<()> {
    arithmetic::arithmetic_encode_file("input.bin", "output.aenc")?;
    arithmetic::arithmetic_decode_file("output.aenc", "decoded.bin")?;
    Ok(())
}
```

---

## Range Coder (`algorithms/range/rust`)

### Add as Dependency

```toml
[dependencies]
rangecoder = { path = "../algorithms/range/rust" }
```

### API

```rust
use rangecoder;

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let encoded = rangecoder::encode(&input_data)?;
    let decoded = rangecoder::decode(&encoded)?;
    Ok(())
}
```

#### `encode(data: &[u8]) -> Result<Vec<u8>, RangeError>`

Encodes data using range coding.

#### `decode(data: &[u8]) -> Result<Vec<u8>, RangeError>`

Decodes range-coded data.

### Error Type

```rust
pub enum RangeError {
    InvalidHeader(String),
    CorruptedData(String),
    IoError(std::io::Error),
}
```

### CLI Binary

```bash
cargo build --bin rangecoder --release
./rangecoder encode input.bin output.rcnc
./rangecoder decode output.rcnc decoded.bin
```

---

## RLE (`algorithms/rle/rust`)

### API

```rust
use rle;

fn main() -> std::io::Result<()> {
    rle::rle_encode_file("input.bin", "output.rle")?;
    rle::rle_decode_file("output.rle", "decoded.bin")?;
    Ok(())
}
```

---

## Common Patterns

All Rust implementations follow these conventions:

| Pattern | Description |
|---------|-------------|
| `*_encode_file` | File-based encoding, returns `io::Result<()>` |
| `*_decode_file` | File-based decoding, returns `io::Result<()>` |
| Error propagation | Uses `?` operator throughout |
| Buffered I/O | All file operations use `BufReader`/`BufWriter` |

### Error Handling

- `io::Error` — file I/O failures
- `Box<dyn Error>` — algorithm-specific errors (Range Coder)
- All errors propagate up to `main()` and print to stderr
