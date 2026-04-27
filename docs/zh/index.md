---
layout: home

hero:
  name: CompressKit
  text: 压缩算法集
  tagline: 使用 C++17、Go 和 Rust 实现的生产级压缩算法。在多语言环境下学习、对比和验证。
  image:
    src: /logo.svg
    alt: CompressKit Logo
  actions:
    - theme: brand
      text: 快速开始
      link: /zh/guide/getting-started
    - theme: alt
      text: 查看 GitHub
      link: https://github.com/LessUp/compress-kit
    - theme: alt
      text: English
      link: /en/

features:
  - icon: 🌐
    title: 多语言对比
    details: 每种算法都有 C++17、Go 和 Rust 实现，便于对比性能、代码风格和工程实践。
  - icon: 📦
    title: 跨语言兼容
    details: 所有语言实现共享相同的二进制格式。用 C++ 编码，用 Go 解码，用 Rust 验证——完全互通。
  - icon: 📚
    title: 面向学习
    details: 通过清晰的解释和三种语言的工作代码示例，深入理解每个算法背后的原理。
  - icon: ✅
    title: 生产级验证
    details: 完整的 CI/CD 流水线，自动化构建、跨语言正确性测试和持续基准测试。
---

<StatsBar />

## 选择算法

<AlgorithmGrid />

## 快速对比

| 算法 | 压缩率 | 速度 | 适用场景 |
|------|--------|------|----------|
| **霍夫曼编码** | 中等 | 快 | 通用文本/数据 |
| **算术编码** | 高 | 中等 | 最大压缩需求 |
| **区间编码** | 高 | 快 | 平衡性能 |
| **行程编码** | 可变 | 极快 | 高度重复数据 |

## 快速开始

```bash
# 克隆仓库
git clone https://github.com/LessUp/compress-kit.git
cd compress-kit

# 构建所有实现
make build

# 运行测试
make test

# 运行基准测试
make bench
```

## 跨语言验证

CompressKit 的核心特性——用任意语言编码，用任意其他语言解码：

```bash
# 用 C++ 编码
./algorithms/huffman/cpp/huffman_cpp encode input.bin encoded.huf

# 用 Go 解码
./algorithms/huffman/go/huffman_go decode encoded.huf restored.bin

# 验证正确性
diff input.bin restored.bin  # 无输出 = 相同
```

## 文档结构

| 章节 | 描述 |
|------|------|
| [快速开始](/zh/guide/getting-started) | 环境配置、构建说明 |
| [算法详解](/zh/guide/algorithms) | 详细说明和对比 |
| [API 参考](/zh/api/go) | 各语言的库 API |
| [基准测试](/zh/benchmarks/results) | 性能结果和方法论 |

## 社区

- 💬 [GitHub 讨论区](https://github.com/LessUp/compress-kit/discussions)
- 🐛 [GitHub Issues](https://github.com/LessUp/compress-kit/issues)
- 🤝 [贡献指南](/zh/guide/contributing)

---

**CompressKit** © 2025-2026 LessUp. MIT 许可证发布。

<style>
:root {
  --vp-home-hero-name-color: transparent;
  --vp-home-hero-name-background: linear-gradient(135deg, #2563eb 0%, #0ea5e9 50%, #10b981 100%);
}
</style>
