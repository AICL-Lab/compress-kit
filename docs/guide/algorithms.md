# 算法综述

CompressKit 实现四种经典无损压缩算法。本页提供快速对比与选择指南；每个算法的
原理、文件格式、CLI 用法与性能特征详见 [算法详解](/algorithms/huffman)。

## 快速对比

| 算法 | 适用场景 | 压缩率 | 速度 | 时间复杂度 | 空间复杂度 | 详解链接 |
|------|----------|--------|------|------------|------------|----------|
| **Huffman** | 通用、文本 | 中等 | 快 | O(n log σ) | O(σ) | [霍夫曼编码](/algorithms/huffman) |
| **Arithmetic** | 最大压缩 | 高 | 中等 | O(n) | O(σ) | [算术编码](/algorithms/arithmetic) |
| **Range Coder** | 平衡性能 | 高 | 快 | O(n) | O(σ) | [区间编码](/algorithms/range) |
| **RLE** | 重复数据 | 可变 | 极快 | O(n) | O(1) | [行程编码](/algorithms/rle) |

> **图例**: σ = 字母表大小（字节级为 256），n = 输入长度

## 算法选择指南

| 数据类型 | 推荐算法 | 原因 |
|----------|----------|------|
| 文本文件 | Huffman 或 Range Coder | 自然语言频率分布不均匀 |
| 最大压缩需求 | Arithmetic | 最接近理论极限 |
| 性能关键场景 | Range Coder | 速度与压缩率的最佳平衡 |
| 高度重复（位图、日志）| RLE | 简单模式压缩效果极好 |
| 未知/混合内容 | Range Coder | 速度与压缩率最佳平衡 |

### 决策流程图


```
数据是否高度重复？
├── 是 → 使用 RLE
└── 否 →
    是否需要最大压缩？
    ├── 是 → 使用 Arithmetic
    └── 否 →
        速度是否关键？
        ├── 是 → 使用 Range Coder
        └── 否 → 使用 Huffman
```

---

## 延伸阅读

- [霍夫曼编码](/algorithms/huffman) · [算术编码](/algorithms/arithmetic) · [区间编码](/algorithms/range) · [行程编码](/algorithms/rle) - 各算法详解
- [算法学院](/academy/) - 数学原理与深度解析
- [架构设计](/architecture/) - 系统架构与二进制格式规范
- [快速开始](/guide/getting-started) - 构建和测试说明
- [基准测试结果](/benchmarks/results) - 实测性能对比
