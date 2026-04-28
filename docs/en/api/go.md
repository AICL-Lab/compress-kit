# Go Library API Reference

See also: [Streaming API](/en/api/streaming)

All Go implementations expose a library API in addition to the CLI. Each algorithm now follows the same layered pattern:

- legacy file/stream helpers such as `Encode`, `Decode`, `EncodeFile`, and `DecodeFile`
- shared streaming implementations via `NewStreamingEncoder()` / `NewStreamingDecoder()`
- shared buffer helpers via `github.com/LessUp/compress-kit/algorithms/shared/go/codec`

## Huffman (`algorithms/huffman/go`)

### Import

```go
import "huffman"
```

### Functions

#### `Encode(input io.Reader, w io.Writer) error`

Reads from `input` and writes Huffman-encoded output to `w`.

```go
in, _ := os.Open("input.bin")
defer in.Close()
out, _ := os.Create("output.huf")
defer out.Close()

if err := huffman.Encode(in, out); err != nil {
    log.Fatal(err)
}
```

#### `Decode(r io.Reader, w io.Writer) error`

Reads from `r` and writes decoded output to `w`.

```go
in, _ := os.Open("encoded.huf")
defer in.Close()
out, _ := os.Create("decoded.bin")
defer out.Close()

if err := huffman.Decode(in, out); err != nil {
    log.Fatal(err)
}
```

#### `EncodeFile(inputPath, outputPath string) error`

Convenience function for file-based encoding.

#### `DecodeFile(inputPath, outputPath string) error`

Convenience function for file-based decoding.

### Constants

| Constant | Value | Description |
|----------|-------|-------------|
| `SymbolLimit` | 257 | 256 bytes + 1 EOF symbol |
| `MaxInputSize` | 4 GiB | Maximum allowed input size |

### Error Handling

All functions return `error`. Common errors:

- `"cannot open input file"` — file not found or permission denied
- `"input too large"` — exceeds `MaxInputSize`
- `"invalid input file format"` — wrong magic bytes or corrupted file
- `"input data corrupted or truncated"` — decode detected invalid bit stream

---

## Arithmetic (`algorithms/arithmetic/go`)

### Import

```go
import "arithmetic"
```

### Functions

Same API as Huffman:

- `Encode(input io.Reader, w io.Writer) error`
- `Decode(r io.Reader, w io.Writer) error`
- `EncodeFile(inputPath, outputPath string) error`
- `DecodeFile(inputPath, outputPath string) error`

### Additional Functions

#### `ScaleFrequencies(freq []uint32)`

Normalizes frequencies to fit within `MaxTotal` (16M). Useful when building custom frequency tables.

#### `BuildCumulative(freq []uint32) []uint32`

Builds a cumulative frequency table from raw frequencies.

### Constants

| Constant | Value | Description |
|----------|-------|-------------|
| `MaxTotal` | 16,777,216 | Maximum total frequency |
| `MaxInputSize` | 4 GiB | Maximum allowed input size |

---

## Range Coder (`algorithms/range/go`)

### Import

```go
import "rangecoder"
```

### Functions

- `Encode(data []byte) ([]byte, error)` — returns encoded byte slice
- `Decode(data []byte) ([]byte, error)` — returns decoded byte slice

### Constants

| Constant | Value | Description |
|----------|-------|-------------|
| `MaxOutputSize` | 1 GiB | Maximum output size to prevent decompression bombs |

::: tip Note
The Range Coder uses byte-level I/O (not bit-level), making it more efficient than Arithmetic coding while achieving nearly identical compression ratios.
:::

---

## RLE (`algorithms/rle/go`)

### Import

```go
import "rle"
```

### Functions

- `Encode(input io.Reader, w io.Writer) error`
- `Decode(r io.Reader, w io.Writer) error`
- `EncodeFile(inputPath, outputPath string) error`
- `DecodeFile(inputPath, outputPath string) error`

### Constants

| Constant | Value | Description |
|----------|-------|-------------|
| `MaxOutputSize` | 1 GiB | Maximum output size to prevent decompression bombs |

### Error Handling

RLE decode validates that `count > 0` and rejects invalid RLE data:

- `"invalid RLE data: count should not be 0"` — corrupted data
- `"RLE data truncated"` — incomplete (count, value) pair
- `"output size limit exceeded"` — decompression bomb detected

---

## CLI Usage

Each Go module has a `cmd/main.go` entry point:

```bash
# Build
cd algorithms/huffman/go && go build -o huffman_go ./cmd

# Use
./huffman_go encode input.bin output.huf
./huffman_go decode output.huf decoded.bin
```

### Usage Pattern

All CLIs follow the same interface:

```
<algorithm> encode <input> <output>
<algorithm> decode <input> <output>
```

Exit codes:
- `0` — success
- `1` — error (message printed to stderr)
