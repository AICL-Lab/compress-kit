---
layout: home

hero:
  name: CompressKit
  text: 值得信赖的压缩算法
  tagline: 生产级 Huffman、Arithmetic、Range 和 RLE 实现，支持 C++、Go 和 Rust。跨语言验证，文档完善，适合学习。
  image:
    src: /logo.svg
    alt: CompressKit Logo
  actions:
    - theme: brand
      text: 快速开始
      link: /zh/guide/getting-started
    - theme: alt
      text: 算法对比
      link: /zh/guide/algorithms
    - theme: alt
      text: English Docs
      link: /en/

features:
  - icon: 🔄
    title: 跨语言验证
    details: 每种算法均在 C++17、Go 和 Rust 中实现，二进制格式完全一致。一种语言编码，另一种语言解码。
  - icon: 📚
    title: 阅读即学习
    details: 代码清晰、注释完善，专为教育设计。每个实现都浓缩在单文件中，真正可读。
  - icon: ⚡
    title: 生产级就绪
    details: 安全限制（4 GiB 输入，1 GiB 输出）、流式 API、全面测试、清晰文档。
  - icon: 🧪
    title: 测试驱动质量
    details: 144 项跨语言一致性测试确保二进制兼容。每次发布前均经过完整验证。
---

<StatsBar />

## 为什么选择 CompressKit？

<div class="ck-value-props">

| 你需要 | 我们提供 |
|--------|----------|
| 学习压缩算法 | 阅读单文件实现，代码清晰易懂 |
| 跨语言兼容性 | C++/Go/Rust 验证的二进制格式 |
| 生产环境使用 | 流式 API、安全限制、完善的错误处理 |
| 性能对比基准 | 运行 `make bench` 获取真实数据 |

</div>

## 算法选择指南

<AlgorithmGrid />

## 30 秒快速上手

```bash
git clone https://github.com/LessUp/compress-kit.git
cd compress-kit
make build && make test
```

就这么简单。12 个实现（4 种算法 × 3 种语言）构建并测试完成。

## 独特之处

**不是"万能压缩库"。** CompressKit 是一个压缩算法实验室：

- **透明的格式** — 每个字节都有文档，没有黑盒魔法
- **同构实现** — 相同算法，相同输出，不同语言
- **教育优先** — 代码可读，测试可追踪
- **验证为先** — 跨语言一致性是测试要求，而非事后补充

## Magic Numbers

每个算法都有 4 字节的魔数头，用于即时文件识别：

| 算法 | Magic | 描述 |
|------|-------|------|
| Huffman | `HFMN` | 基于前缀码的压缩 |
| Arithmetic | `AENC` | 熵最优编码 |
| Range Coder | `RCNC` | 快速整数区间编码 |
| RLE | `RLE\x00` | 行程长度压缩 |

## 下一步

| 目标 | 页面 |
|------|------|
| 本地构建运行 | [快速开始](/zh/guide/getting-started) |
| 选择合适的算法 | [算法指南](/zh/guide/algorithms) |
| 作为库使用 | [Streaming API](/zh/api/streaming) |
| 验证兼容性 | [跨语言测试](/zh/testing/cross-language) |
