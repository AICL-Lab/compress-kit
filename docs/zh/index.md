---
layout: home

hero:
  name: CompressKit
  text: 值得信赖的压缩算法
  tagline: Huffman、Arithmetic、Range、RLE —— 支持 C++17、Go、Rust
  actions:
    - theme: brand
      text: 快速开始
      link: /zh/guide/getting-started
    - theme: alt
      text: 查看 GitHub
      link: https://github.com/LessUp/compress-kit
---

<script setup>
import { onBeforeMount } from 'vue'

onBeforeMount(() => {
  if (typeof window !== 'undefined') {
    localStorage.setItem('docs-lang-preference', '/zh/')
  }
})
</script>

## 快速开始

```bash
git clone https://github.com/LessUp/compress-kit.git
cd compress-kit
make build && make test
```

## 算法

| 算法 | 适用场景 | 速度 |
|------|----------|------|
| [Huffman](/zh/algorithms/huffman) | 通用压缩、文本 | 快 |
| [Arithmetic](/zh/algorithms/arithmetic) | 最大压缩率 | 中 |
| [Range Coder](/zh/algorithms/range) | 生产系统 | 快 |
| [RLE](/zh/algorithms/rle) | 重复数据 | 极快 |

## 跨语言兼容

C++ 编码，Go 解码。Rust 编码，C++ 解码。所有实现产生完全相同的二进制输出。

```bash
# C++ 编码，Go 解码
./cpp/huffman encode input.txt output.huf
./go/huffman decode output.huf decoded.txt
# 完美运行 —— 相同字节，不同语言
```

## 下一步

- [构建说明](/zh/guide/getting-started) — 本地运行
- [算法指南](/zh/guide/algorithms) — 选择合适的算法
- [API 参考](/zh/api/streaming) — 作为库使用
