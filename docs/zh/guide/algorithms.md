# 算法详解

本指南解释项目中实现的四种压缩算法、它们的用例和主要差异。

## 快速对比

| 算法 | 适用场景 | 压缩率 | 速度 | 时间复杂度 | 空间复杂度 |
|------|----------|--------|------|------------|------------|
| **Huffman** | 通用、文本 | 中等 | 快 | O(n log σ) | O(σ) |
| **算术编码** | 最大压缩 | 高 | 中等 | O(n) | O(σ) |
| **区间编码** | 平衡性能 | 高 | 快 | O(n) | O(σ) |
| **RLE** | 重复数据 | 可变 | 极快 | O(n) | O(1) |

> **图例**: σ = 字母表大小（字节级为 256），n = 输入长度

---

## Huffman 编码

基于前缀码的无损压缩算法。根据符号频率构建最优前缀树。

### 工作原理

1. **频率分析**: 统计输入中的字节频率
2. **树构建**: 构建二叉树，低频符号路径更深
3. **码表生成**: 生成前缀码（无歧义比特序列）
4. **编码**: 用码表替换每个字节，写入比特流

### 文件格式

| 字段 | 大小 | 描述 |
|------|------|------|
| 魔数 | 4 字节 | `HFMN` (0x48 0x46 0x4D 0x4E) |
| 频率表 | 257 × 4 字节 | 小端 uint32 数组 |
| 编码数据 | 可变 | 比特流 |

### 压缩效率

- **理论下限**: 平均码长 ≥ 熵 H
- **Huffman 上限**: H ≤ L < H + 1（每个符号最多多 1 位）
- 对频率分布不均匀的数据最有效

### 使用示例

::: code-group

```bash [C++]
./huffman_cpp encode input.bin output.huf
./huffman_cpp decode output.huf restored.bin
```

```bash [Go]
./huffman_go encode input.bin output.huf
./huffman_go decode output.huf restored.bin
```

```bash [Rust]
./huffman_rust encode input.bin output.huf
./huffman_rust decode output.huf restored.bin
```

:::

---

## 算术编码

将整条消息表示为区间 [0, 1) 中的单个数字，比 Huffman 编码更接近熵限。

### 工作原理

1. **初始化**: 从区间 [0, 1) 开始
2. **细分**: 对每个符号，根据概率细分当前区间
3. **选择**: 最终区间包含无限多个数，任选一个表示消息
4. **输出**: 位数 ≈ -log₂(P(消息))

### Huffman vs 算术编码对比

| 特性 | Huffman | 算术编码 |
|------|---------|----------|
| 编码单位 | 至少 1 位/符号 | 支持分数位 |
| 理论效率 | H ≤ L < H + 1 | L ≈ H + ε（更接近熵） |
| 实现复杂度 | 简单 | 复杂（精度管理） |
| 速度 | 快 | 慢 |
| 用例 | 通用 | 最大压缩 |

### 特性

- **最优压缩**: 理论上最接近熵限
- **较慢**: 编码/解码开销高于 Huffman
- **复杂度**: 需要仔细的精度管理

---

## 区间编码

整数运算实现，效果等价于算术编码，但实践中通常更高效。使用整数区间运算而非浮点数。

### 算术编码 vs 区间编码

| 特性 | 算术编码 | 区间编码 |
|------|----------|----------|
| 输出单位 | 位 | 字节 |
| I/O 效率 | 较低 | 较高 |
| 压缩率 | 几乎相同 | 几乎相同 |
| 专利状态 | 曾有历史专利 | 无限制 |
| 工程使用 | 学术研究 | 生产系统 |

### 库 API 使用

**Go 库**:
```go
import "github.com/LessUp/compress-kit/algorithms/range/go/rangecoder"

// 编码数据
encoded, err := rangecoder.Encode(data)
if err != nil {
    log.Fatal(err)
}

// 解码数据
decoded, err := rangecoder.Decode(encoded)
if err != nil {
    log.Fatal(err)
}
```

**Rust 包**:
```rust
use rangecoder;

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let encoded = rangecoder::encode(input)?;
    let decoded = rangecoder::decode(&encoded)?;
    Ok(())
}
```

### CLI 使用

::: code-group

```bash [C++]
./rangecoder_cpp encode input.bin output.rcnc
./rangecoder_cpp decode output.rcnc restored.bin
```

```bash [Go]
./rangecoder_go encode input.bin output.rcnc
./rangecoder_go decode output.rcnc restored.bin
```

```bash [Rust]
cargo run --bin rangecoder -- encode input.bin output.rcnc
cargo run --bin rangecoder -- decode output.rcnc restored.bin
```

:::

::: warning 性能提示
区间编码解码器对大于 500KB 的文件存在已知性能问题。跨语言验证时请使用较小的测试文件。
:::

---

## 行程长度编码 (RLE)

最简单的压缩算法，适合连续重复字节的数据。

### 文件格式

存储重复的 `(计数, 值)` 对：

| 字段 | 大小 | 描述 |
|------|------|------|
| 计数 | 4 字节 | 小端无符号整数 |
| 值 | 1 字节 | 重复的字节 |

### 特性

- **简单性**: 最容易理解和实现
- **速度**: 极快的编码和解码
- **最佳场景**: 重复数据（位图、有重复行的日志）
- **最坏情况**: 随机输入可能膨胀到 5 倍
- **常见用法**: 作为其他算法的预处理（如 BWT + MTF + RLE + 算术编码）

### 使用示例

::: code-group

```bash [C++]
./rle_cpp encode input.bin output.rle
./rle_cpp decode output.rle restored.bin
```

```bash [Go]
./rle_go encode input.bin output.rle
./rle_go decode output.rle restored.bin
```

```bash [Rust]
./rle_rust encode input.bin output.rle
./rle_rust decode output.rle restored.bin
```

:::

---

## 算法选择指南

| 数据类型 | 推荐算法 | 原因 |
|----------|----------|------|
| 文本文件 | Huffman 或区间编码 | 自然语言频率分布不均匀 |
| 最大压缩需求 | 算术编码 | 最接近理论极限 |
| 性能关键场景 | 区间编码 | 速度与压缩率的最佳平衡 |
| 高度重复（位图、日志）| RLE | 简单模式压缩效果极好 |
| 未知/混合内容 | 区间编码 | 速度与压缩率最佳平衡 |

### 决策流程图

```
数据是否高度重复？
├── 是 → 使用 RLE
└── 否 →
    是否需要最大压缩？
    ├── 是 → 使用算术编码
    └── 否 →
        速度是否关键？
        ├── 是 → 使用区间编码
        └── 否 → 使用 Huffman
```

---

## 延伸阅读

- [项目结构](/zh/guide/project-structure) - 文件格式和 CLI 规范
- [快速开始](/zh/guide/getting-started) - 构建和测试说明
- [GitHub 仓库](https://github.com/LessUp/compress-kit) - 源代码和问题
