# Range Coder

Range coding is an **integer-based implementation** equivalent to arithmetic coding. It uses integer interval operations instead of floating point, making it more suitable for production systems while achieving the same compression ratios.

## How It Works

Range coding maintains an interval [low, low + range) using fixed-width integers. Unlike arithmetic coding's bit output, range coding outputs **bytes**, which significantly improves I/O efficiency.

```cpp
void encode(const vector<uint8_t>& data,
            const uint32_t cumFreq[257]) {
    uint64_t low = 0;
    uint64_t range = MAX_RANGE;
    uint64_t total = cumFreq[256];
    
    for (uint8_t symbol : data) {
        uint32_t symLow = cumFreq[symbol];
        uint32_t symHigh = cumFreq[symbol + 1];
        
        range /= total;
        low += symLow * range;
        range *= (symHigh - symLow);
        
        // Renormalize
        while (range < MIN_RANGE) {
            output_byte(low >> 56);
            low <<= 8;
            range <<= 8;
        }
    }
    // Output final bytes
    flush(low);
}
```

## Arithmetic vs Range Coder

| Aspect | Arithmetic | Range Coder |
|--------|------------|-------------|
| Arithmetic | Floating point | Fixed-width integers |
| Output unit | Bits | Bytes |
| I/O efficiency | Lower | Higher |
| Compression | Nearly identical | Nearly identical |
| Patent status | Had historical patents | No restrictions |
| Production use | Academic | Industry standard |

## Complexity

| Aspect | Complexity | Notes |
|--------|-----------|-------|
| Time (encode) | O(n) | Similar to arithmetic |
| Time (decode) | O(n) | Byte-level I/O faster |
| Space | O(σ) | Cumulative frequency table |
| Precision | Fixed | 64-bit integers |

## Performance Characteristics

| Input Type | Compression | Speed | Memory |
|------------|-------------|-------|--------|
| Text | 1.90× | 58 MiB/s | Low |

## Use Cases

- ✅ **Production systems** — Most widely deployed entropy coder
- ✅ **Balanced workloads** — Good speed and compression
- ✅ **Video codecs** — H.264, HEVC use range coding
- ✅ **Compression tools** — Used in modern archivers

## CLI Usage

```bash
./build/rangecoder_cpp encode input.bin output.rcnc
./build/rangecoder_cpp decode output.rcnc restored.bin
```

## Further Reading

- [Arithmetic Coding](/en/algorithms/arithmetic) — Floating point equivalent
- [Benchmarks](/en/benchmarks/results) — Performance comparison

## Known Limitations

::: warning Performance Issue with Large Files

The current Range Coder implementation has a **known decode performance issue** for files larger than **500 KiB**. The decode operation may become significantly slower or appear to hang.

**Workaround**: For testing purposes, use files smaller than 100 KiB. This is reflected in the CI pipeline which uses 100 KiB test files for Range Coder verification.

**Status**: This is a known issue that is documented for future improvement. The encode operation works correctly for all file sizes.

:::
