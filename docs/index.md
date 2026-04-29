---
layout: home

hero:
  name: CompressKit
  text: Classic compression, verified across languages
  tagline: A compact learning and verification lab for Huffman, Arithmetic Coding, Range Coder, and RLE in C++17, Go, and Rust.
  image:
    src: /logo.svg
    alt: CompressKit Logo
  actions:
    - theme: brand
      text: English Docs
      link: /en/
    - theme: alt
      text: 中文文档
      link: /zh/
    - theme: alt
      text: GitHub
      link: https://github.com/LessUp/compress-kit

features:
  - icon: 🔁
    title: Encode here, decode there
    details: The repository continuously checks C++17, Go, and Rust implementations against the same binary formats.
  - icon: 🧠
    title: Algorithms you can study
    details: Each implementation is intentionally small enough to read while still covering real file IO, errors, limits, and tests.
  - icon: 🧪
    title: Verification-first engineering
    details: Unit tests, streaming contracts, and cross-language conformance are part of the project shape rather than afterthoughts.
---

## Choose your entry point

| I want to... | Go to |
|-------------|-------|
| Learn the project quickly | [English guide](/en/guide/getting-started) · [中文指南](/zh/guide/getting-started) |
| Compare algorithms | [Algorithm guide](/en/guide/algorithms) · [算法详解](/zh/guide/algorithms) |
| Use the library APIs | [Go](/en/api/go) · [Rust](/en/api/rust) · [C++](/en/api/cpp) |
| Verify compatibility | [Cross-language testing](/en/testing/cross-language) · [跨语言测试](/zh/testing/cross-language) |

## What makes this repository different

CompressKit is not a generic compression package with a thin README. It is a
multi-language compression laboratory: the same four classic algorithms are
implemented three times, then checked against the same command-line contract and
binary compatibility expectations.

The site focuses on the parts that matter for readers: how the algorithms work,
how the implementations line up across languages, where known limits exist, and
how to reproduce the verification locally.
