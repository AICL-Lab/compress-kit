---
layout: home

hero:
  name: CompressKit
  text: 可阅读、可运行、可对比的无损压缩算法
  tagline: 使用 C++17、Go、Rust 实现 Huffman、算术编码、区间编码和 RLE，并通过跨语言二进制兼容性验证。
  image:
    src: /logo.svg
    alt: CompressKit Logo
  actions:
    - theme: brand
      text: 从指南开始
      link: /zh/guide/getting-started
    - theme: alt
      text: 对比算法
      link: /zh/guide/algorithms
    - theme: alt
      text: English
      link: /en/

features:
  - icon: 🧩
    title: 四种经典算法
    details: Huffman、算术编码、区间编码、RLE 并排实现，便于学习算法与工程取舍。
  - icon: 🌐
    title: 三套语言实现
    details: C++17、Go 1.21+、Rust 1.70+ 共享统一 CLI 形态和兼容文件格式。
  - icon: ✅
    title: 验证进入主流程
    details: 测试基线包含 streaming 生命周期检查和可执行的跨语言解码矩阵。
---

<StatsBar />

## 算法地图

<AlgorithmGrid />

## 最短有效路径

```bash
git clone https://github.com/LessUp/compress-kit.git
cd compress-kit
make build
make test
```

按目标阅读文档：

| 目标 | 页面 |
|------|------|
| 环境准备与首次运行 | [快速开始](/zh/guide/getting-started) |
| 算法行为与取舍 | [算法详解](/zh/guide/algorithms) |
| 公共 API 与 streaming 门面 | [API 参考](/zh/api/streaming) |
| 兼容性验证 | [跨语言测试](/zh/testing/cross-language) |

## 当前工程边界

CompressKit 当前保持各算法既有二进制格式稳定。统一 frame 格式与 benchmark
治理方案已作为未来 OpenSpec 设计上下文归档，只有重新以聚焦变更实施时才会进入主规范。
