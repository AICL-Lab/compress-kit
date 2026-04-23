# 算术编码

算术编码是一种无损数据压缩算法，将**整个消息编码为区间 [0, 1) 中的单个数字**。这种方法比霍夫曼编码更接近熵极限。

## 工作原理

与霍夫曼编码为每个符号分配整数位长的编码不同，算术编码可以分配**分数位长的编码**，使其更接近理论熵极限。

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
    
    // 输出 [low, high) 中的数字
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

## 算术编码 vs 霍夫曼

| 方面 | 霍夫曼 | 算术编码 |
|------|--------|----------|
| 编码长度 | 整数位 | 分数位 |
| 效率 | H ≤ L < H + 1 | L ≈ H + ε |
| 速度 | 更快 | 较慢 |
| 复杂度 | 简单 | 精度管理 |
| 输出 | 位 | 单个数字 |

## 复杂度

| 方面 | 复杂度 | 说明 |
|------|--------|------|
| 时间（编码） | O(n) | 单次遍历，区间更新 |
| 时间（解码） | O(n) | 逆过程 |
| 空间 | O(σ) | 概率表 |
| 精度 | 固定 | 需要小心处理下溢 |

## 精度考虑

算术编码需要管理精度以避免下溢：

1. **重归一化**：当区间变得太小时周期性输出位
2. **整数运算**：生产实现使用带缩放的整数运算
3. **流结束标记**：解码时需要知道何时完成

## 适用场景

- ✅ **最大压缩** — 当每个位都很重要时
- ✅ **统计数据** — 已有概率模型的数据
- ✅ **学术研究** — 理解熵编码
- ❌ **实时流处理** — 精度管理复杂
- ❌ **嵌入式系统** — 浮点运算需求

## 性能

| 输入类型 | 压缩率 | 速度 |
|----------|--------|------|
| 文本 | 1.8-2.0× | 中等 |
| 随机 | 0.99× | 中等 |
| 重复 | 50-100× | 中等 |

## 延伸阅读

- [区间编码](/zh/algorithms/range) — 用于生产环境的整数实现
- [霍夫曼编码](/zh/algorithms/huffman) — 更简单的替代方案
- [基准测试](/zh/benchmarks/results)
