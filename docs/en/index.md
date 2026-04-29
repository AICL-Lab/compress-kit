---
layout: home

hero:
  name: CompressKit
  text: Lossless compression algorithms you can read, run, and compare
  tagline: Huffman, Arithmetic Coding, Range Coder, and RLE implemented in C++17, Go, and Rust with cross-language binary verification.
  image:
    src: /logo.svg
    alt: CompressKit Logo
  actions:
    - theme: brand
      text: Start with the guide
      link: /en/guide/getting-started
    - theme: alt
      text: Compare algorithms
      link: /en/guide/algorithms
    - theme: alt
      text: 中文
      link: /zh/

features:
  - icon: 🧩
    title: Four classic algorithms
    details: Huffman, Arithmetic Coding, Range Coder, and RLE are implemented side by side for direct comparison.
  - icon: 🌐
    title: Three language stacks
    details: C++17, Go 1.21+, and Rust 1.70+ implementations expose a shared CLI shape and matching file formats.
  - icon: ✅
    title: Conformance in the loop
    details: The test baseline includes streaming lifecycle checks and an executable cross-language decode matrix.
---

<StatsBar />

## Algorithm map

<AlgorithmGrid />

## The shortest useful path

```bash
git clone https://github.com/LessUp/compress-kit.git
cd compress-kit
make build
make test
```

After that, use the docs by intent:

| Need | Page |
|------|------|
| Build prerequisites and first run | [Getting Started](/en/guide/getting-started) |
| Algorithm behavior and trade-offs | [Algorithms](/en/guide/algorithms) |
| Public APIs and streaming facade | [API Reference](/en/api/streaming) |
| Compatibility verification | [Cross-Language Testing](/en/testing/cross-language) |

## Current engineering stance

CompressKit keeps the existing per-algorithm binary formats stable. Future frame
format and benchmark-governance proposals are archived as design context until
they are implemented as focused OpenSpec changes.
