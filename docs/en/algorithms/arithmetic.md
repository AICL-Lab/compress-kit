# Arithmetic Coding

Arithmetic coding is a lossless data compression algorithm that **encodes the entire message as a single number** in the interval [0, 1). This approach achieves compression closer to the entropy limit than Huffman coding.

## How It Works

Unlike Huffman coding which assigns integer-length codes to individual symbols, arithmetic coding can assign **fractional-length codes**, allowing it to approach the theoretical entropy limit more closely.

::: code-group

```cpp [C++]
void encode(const vector<uint8_t>& data, 
            const vector<double>& probs) {
    double low = 0.0;
    double high = 1.0;
    
    for (uint8_t symbol : data) {
        double range = high - low;
        high = low + range * cumProb[symbol + 1];
        low = low + range * cumProb[symbol];
    }
    
    // Output number in [low, high)
    output_bits = ceil(-log2(high - low));
    write_value((low + high) / 2, output_bits);
}
```

```go [Go]
func Encode(data []byte, probs []float64) []byte {
    low, high := 0.0, 1.0
    
    for _, symbol := range data {
        range_ := high - low
        high = low + range_*cumProb[symbol+1]
        low = low + range_*cumProb[symbol]
    }
    
    bits := int(math.Ceil(-math.Log2(high - low)))
    value := (low + high) / 2
    
    return bitsToBytes(value, bits)
}
```

```rust [Rust]
pub fn encode(data: &[u8], probs: &[f64]) -> Vec<u8> {
    let mut low = 0.0;
    let mut high = 1.0;
    
    for &symbol in data {
        let range = high - low;
        high = low + range * cum_prob[symbol as usize + 1];
        low = low + range * cum_prob[symbol as usize];
    }
    
    let bits = ((high - low).log2().abs().ceil()) as usize;
    let value = (low + high) / 2.0;
    
    bits_to_bytes(value, bits)
}
```

:::

## Arithmetic vs Huffman

| Aspect | Huffman | Arithmetic |
|--------|---------|------------|
| Code length | Integer bits | Fractional bits |
| Efficiency | H ≤ L < H + 1 | L ≈ H + ε |
| Speed | Faster | Slower |
| Complexity | Simpler | Precision management |
| Output | Bits | Single number |

## Complexity

| Aspect | Complexity | Notes |
|--------|-----------|-------|
| Time (encode) | O(n) | Single pass with interval updates |
| Time (decode) | O(n) | Reverse process |
| Space | O(σ) | Probability table |
| Precision | Fixed | Requires careful underflow handling |

## Precision Considerations

Arithmetic coding requires managing precision to avoid underflow:

1. **Renormalization**: Periodically output bits when range becomes too small
2. **Integer arithmetic**: Production implementations use integer math with scaling
3. **End-of-stream marker**: Required to know when decoding is complete

## Use Cases

- ✅ **Maximum compression** — When every bit counts
- ✅ **Statistical data** — Already has probability models
- ✅ **Academic study** — Understanding entropy coding
- ❌ **Real-time streaming** — Precision management complexity
- ❌ **Embedded systems** — Floating point requirements

## Performance

| Input Type | Compression | Speed |
|------------|-------------|-------|
| Text | 1.8-2.0× | Medium |
| Random | 0.99× | Medium |
| Repetitive | 50-100× | Medium |

## Further Reading

- [Range Coder](/en/algorithms/range) — Integer-based equivalent for production use
- [Huffman Coding](/en/algorithms/huffman) — Simpler alternative
- [Benchmarks](/en/benchmarks/results)
