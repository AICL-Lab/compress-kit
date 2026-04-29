# Range Coder

Range coding is an **integer-based implementation** equivalent to arithmetic coding. It uses integer interval operations instead of floating point, making it more suitable for production systems while achieving the same compression ratios.

## How It Works

Range coding maintains an interval [low, low + range) using fixed-width integers. Unlike arithmetic coding's bit output, range coding outputs **bytes**, which significantly improves I/O efficiency.

::: code-group

```cpp [C++]
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

```go [Go]
func (rc *RangeCoder) Encode(data []byte, cumFreq [257]uint32) error {
    const maxRange uint64 = 1 << 32
    const minRange uint64 = 1 << 24
    total := cumFreq[256]
    
    for _, symbol := range data {
        symLow := cumFreq[symbol]
        symHigh := cumFreq[symbol+1]
        
        rc.range /= uint64(total)
        rc.low += uint64(symLow) * rc.range
        rc.range *= uint64(symHigh - symLow)
        
        for rc.range < minRange {
            rc.outputByte(byte(rc.low >> 56))
            rc.low <<= 8
            rc.range <<= 8
        }
    }
    rc.flush()
    return nil
}
```

```rust [Rust]
pub fn encode(&mut self, data: &[u8], cum_freq: &[u32; 257]) {
    const MAX_RANGE: u64 = 1u64 << 32;
    const MIN_RANGE: u64 = 1u64 << 24;
    let total = cum_freq[256];
    
    for &symbol in data {
        let sym_low = cum_freq[symbol as usize];
        let sym_high = cum_freq[symbol as usize + 1];
        
        self.range /= total as u64;
        self.low += (sym_low as u64) * self.range;
        self.range *= (sym_high - sym_low) as u64;
        
        while self.range < MIN_RANGE {
            self.output((self.low >> 56) as u8);
            self.low <<= 8;
            self.range <<= 8;
        }
    }
    self.flush();
}
```

:::

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

| Language | Text Compression | Speed | Memory |
|----------|------------------|-------|--------|
| C++ | 1.90× | 58 MiB/s | Low |
| Go | 1.90× | 50 MiB/s | Low |
| Rust | 1.90× | 64 MiB/s | Low |

## Use Cases

- ✅ **Production systems** — Most widely deployed entropy coder
- ✅ **Balanced workloads** — Good speed and compression
- ✅ **Video codecs** — H.264, HEVC use range coding
- ✅ **Compression tools** — Used in modern archivers

## Library API

### Go

```go
import "github.com/LessUp/compress-kit/algorithms/range/go/rangecoder"

// Create encoder with frequency table
freq := rangecoder.BuildFrequencyTable(data)
cumFreq := rangecoder.BuildCumulative(freq)

// Encode
encoded, err := rangecoder.Encode(data, cumFreq)

// Decode  
decoded, err := rangecoder.Decode(encoded, cumFreq, len(data))
```

### Rust

```rust
use rangecoder;

let freq = rangecoder::build_frequency_table(&data);
let cum_freq = rangecoder::cumulative(&freq);

let encoded = rangecoder::encode(&data, &cum_freq)?;
let decoded = rangecoder::decode(&encoded, &cum_freq, data.len())?;
```

## Further Reading

- [Arithmetic Coding](/en/algorithms/arithmetic) — Floating point equivalent
- [Benchmarks](/en/benchmarks/results) — Performance comparison
- [OpenSpec Architecture Specs](https://github.com/LessUp/compress-kit/tree/master/openspec/specs/core-architecture)

## Known Limitations

::: warning Performance Issue with Large Files

The current Range Coder implementation has a **known decode performance issue** for files larger than **500 KB**. The decode operation may become significantly slower or appear to hang.

**Workaround**: For testing purposes, use files smaller than 100 KB. This is reflected in the CI pipeline which uses 100 KB test files for Range Coder verification.

**Status**: This is a known issue that is documented for future improvement. The encode operation works correctly for all file sizes.

:::
