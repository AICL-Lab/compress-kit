# Run-Length Encoding (RLE)

RLE is the **simplest compression algorithm**, replacing consecutive repeated bytes with a `(count, value)` pair. It's extremely fast and works exceptionally well on data with long runs of identical values.

## How It Works

::: code-group

```cpp [C++]
vector<uint8_t> encode(const vector<uint8_t>& data) {
    vector<uint8_t> result;
    
    for (size_t i = 0; i < data.size();) {
        uint8_t current = data[i];
        uint32_t count = 1;
        
        // Count consecutive identical bytes
        while (i + count < data.size() && 
               data[i + count] == current && 
               count < UINT32_MAX) {
            count++;
        }
        
        // Write count (4 bytes, little-endian) + value
        result.push_back(count & 0xFF);
        result.push_back((count >> 8) & 0xFF);
        result.push_back((count >> 16) & 0xFF);
        result.push_back((count >> 24) & 0xFF);
        result.push_back(current);
        
        i += count;
    }
    
    return result;
}
```

```go [Go]
func Encode(data []byte) []byte {
    var result bytes.Buffer
    
    for i := 0; i < len(data); {
        current := data[i]
        count := uint32(1)
        
        for i+int(count) < len(data) && 
            data[i+int(count)] == current && 
            count < math.MaxUint32 {
            count++
        }
        
        // Write count as uint32 (little-endian) + value
        binary.Write(&result, binary.LittleEndian, count)
        result.WriteByte(current)
        
        i += int(count)
    }
    
    return result.Bytes()
}
```

```rust [Rust]
pub fn encode(data: &[u8]) -> Vec<u8> {
    let mut result = Vec::new();
    let mut i = 0;
    
    while i < data.len() {
        let current = data[i];
        let mut count: u32 = 1;
        
        while i + count as usize < data.len() && 
              data[i + count as usize] == current && 
              count < u32::MAX {
            count += 1;
        }
        
        // Write count (little-endian) + value
        result.extend_from_slice(&count.to_le_bytes());
        result.push(current);
        
        i += count as usize;
    }
    
    result
}
```

:::

## File Format

| Field | Size | Description |
|-------|------|-------------|
| Magic | 4 bytes | `RLE\x00` (0x52 0x4C 0x45 0x00) |
| Count | 4 bytes | Little-endian unsigned int (run length) |
| Value | 1 byte | The repeated byte value |

## Complexity

| Aspect | Complexity | Notes |
|--------|-----------|-------|
| Time (encode) | O(n) | Single pass |
| Time (decode) | O(n) | Single pass, very fast |
| Space | O(1) | No auxiliary structures |

## Performance

| Data Type | Compression | Speed |
|-----------|-------------|-------|
| Repetitive (10 MB) | 25× | 300+ MiB/s |
| Text | 1.1× | Fast |
| Random | 0.2× (expands) | Fast |

::: warning Worst Case
On random data, RLE expands each byte to 5 bytes (4-byte count + 1-byte value), resulting in **5× size increase**.
:::

## Use Cases

- ✅ **Bitmap images** — Long runs of same color
- ✅ **Log files** — Repeated patterns
- ✅ **Preprocessing** — Before BWT or other transforms
- ✅ **FAX transmissions** — Standard compression method
- ❌ **Random data** — Severe expansion
- ❌ **Already compressed** — No benefit

## Common Use Cases

### As Preprocessing

RLE is often used as a preprocessing step in more complex compression pipelines:

```
Original → BWT → MTF → RLE → Arithmetic → Compressed
```

This combination (Burrows-Wheeler + Move-to-Front + RLE + Arithmetic) is the basis of **bzip2**.

### In Image Formats

- **BMP**: Simple RLE variants
- **PCX**: RLE compression
- **TIFF**: Optional RLE packbits

## Comparison with Other Algorithms

| Algorithm | Setup Cost | Compression | Speed |
|-----------|------------|-------------|-------|
| RLE | None | Variable | Fastest |
| Huffman | O(σ log σ) | Medium | Fast |
| Range | O(n) | High | Fast |

## Further Reading

- [Algorithm Comparison](/en/guide/algorithms) — Full comparison matrix
- [Benchmarks](/en/benchmarks/results) — Performance numbers
- [Burrows-Wheeler Transform](https://en.wikipedia.org/wiki/Burrows%E2%80%93Wheeler_transform) — Common RLE pairing
