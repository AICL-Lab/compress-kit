---
layout: home

hero:
  name: CompressKit
  text: Compression algorithms you can trust
  tagline: Huffman, Arithmetic, Range, and RLE in C++17, Go, and Rust
  actions:
    - theme: brand
      text: Get Started
      link: /en/guide/getting-started
    - theme: alt
      text: View on GitHub
      link: https://github.com/LessUp/compress-kit
---

<script setup>
import { onBeforeMount } from 'vue'

onBeforeMount(() => {
  if (typeof window !== 'undefined') {
    localStorage.setItem('docs-lang-preference', '/en/')
  }
})
</script>

## Quick Start

```bash
git clone https://github.com/LessUp/compress-kit.git
cd compress-kit
make build && make test
```

## Algorithms

| Algorithm | Best For | Speed |
|-----------|----------|-------|
| [Huffman](/en/algorithms/huffman) | General purpose, text | Fast |
| [Arithmetic](/en/algorithms/arithmetic) | Maximum compression | Medium |
| [Range Coder](/en/algorithms/range) | Production systems | Fast |
| [RLE](/en/algorithms/rle) | Repetitive data | Very Fast |

## Cross-Language Compatibility

Encode in C++, decode in Go. Encode in Rust, decode in C++. All implementations produce identical binary output.

```bash
# C++ encode, Go decode
./cpp/huffman encode input.txt output.huf
./go/huffman decode output.huf decoded.txt
# Works perfectly — same bytes, different languages
```

## Next

- [Build instructions](/en/guide/getting-started) — Get running locally
- [Algorithm guide](/en/guide/algorithms) — Choose the right one
- [API reference](/en/api/streaming) — Use as a library
