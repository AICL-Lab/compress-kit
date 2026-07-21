# Arithmetic

Arithmetic 是一种无损数据压缩算法，将**整个消息编码为区间 [0, 1) 中的单个数字**。这种方法比 Huffman 更接近熵极限。

## 工作原理

与 Huffman 为每个符号分配整数位长的编码不同，Arithmetic 可以分配**分数位长的编码**，使其更接近理论熵极限。

```cpp
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

## Arithmetic vs Huffman

| 方面 | Huffman | Arithmetic |
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

Arithmetic 需要管理精度以避免下溢：

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

- [Range Coder](/algorithms/range) — 用于生产环境的整数实现
- [Huffman](/algorithms/huffman) — 更简单的替代方案
- [基准测试](/benchmarks/results)
