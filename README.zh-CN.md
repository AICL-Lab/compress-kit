# Encoding —— 编码算法集合

[![CI](https://github.com/LessUp/encoding/actions/workflows/ci.yml/badge.svg)](https://github.com/LessUp/encoding/actions/workflows/ci.yml)
[![Deploy Docs](https://github.com/LessUp/encoding/actions/workflows/pages.yml/badge.svg)](https://github.com/LessUp/encoding/actions/workflows/pages.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
![C++](https://img.shields.io/badge/C++-17-blue.svg)
![Go](https://img.shields.io/badge/Go-1.21+-00ADD8.svg)
![Rust](https://img.shields.io/badge/Rust-1.70+-orange.svg)

[English](README.md) | 简体中文 | [文档站](https://lessup.github.io/encoding/)

Encoding 是一个面向学习、实现对比与跨语言验证的经典压缩算法集合，使用 C++17、Go 和 Rust 提供对应实现。

## 仓库入口

- 覆盖 Huffman、算术编码、区间编码、RLE 四类经典算法
- 每种算法同时提供 C++17、Go、Rust 三套实现
- 统一 CLI 约定与二进制格式，便于跨语言互相编码/解码验证
- 详细使用说明、算法导读和目录结构统一放在文档站维护

## 快速开始

```bash
make build
make test
make bench
```

如果你想先从单个算法开始：

```bash
cd huffman/cpp
g++ -std=c++17 -O2 main.cpp -o huffman_cpp
./huffman_cpp encode input.bin output.huf
./huffman_cpp decode output.huf restored.bin
```

## 接下来读什么

- [文档首页](https://lessup.github.io/encoding/)
- [快速开始](https://lessup.github.io/encoding/guide/getting-started)
- [算法详解](https://lessup.github.io/encoding/guide/algorithms)
- [项目结构](https://lessup.github.io/encoding/guide/project-structure)

## 许可证

MIT License。
