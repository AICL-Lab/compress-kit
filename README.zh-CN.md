# CompressKit

<p align="center">
  <img src="docs/public/logo.svg" width="120" alt="CompressKit Logo">
</p>

<p align="center">
  <strong>使用 C++17、Go、Rust 实现的经典无损压缩算法。</strong>
</p>

<p align="center">
  <a href="https://github.com/LessUp/compress-kit/actions/workflows/ci.yml"><img src="https://github.com/LessUp/compress-kit/actions/workflows/ci.yml/badge.svg" alt="CI Status"></a>
  <a href="https://lessup.github.io/compress-kit/"><img src="https://img.shields.io/badge/Docs-在线文档-blue?logo=readthedocs&logoColor=white" alt="Documentation"></a>
  <a href="https://opensource.org/licenses/MIT"><img src="https://img.shields.io/badge/许可证-MIT-green.svg" alt="License"></a>
</p>

<p align="center">
  <a href="README.md">English</a> | <b>简体中文</b> | <a href="https://lessup.github.io/compress-kit/">文档站点</a>
</p>

CompressKit 是一个面向学习与验证的压缩算法仓库：同一组经典算法分别用
C++17、Go、Rust 实现，再通过统一命令行契约和跨语言解码矩阵验证格式兼容性。
它不是黑盒压缩库，而是可以阅读、运行、对比和验证的多语言算法实验室。

## 包含内容

| 算法 | C++17 | Go | Rust | 适用场景 |
|------|------:|---:|-----:|----------|
| Huffman 编码 | ✓ | ✓ | ✓ | 通用文本/数据，学习前缀码 |
| 算术编码 | ✓ | ✓ | ✓ | 理解熵编码与压缩率对比 |
| 区间编码 | ✓ | ✓ | ✓ | 对比算术编码风格实现 |
| RLE 行程编码 | ✓ | ✓ | ✓ | 高重复数据与简单格式学习 |

所有命令行工具都遵循：

```bash
<binary> <encode|decode> <input> <output>
```

## 快速开始

```bash
git clone https://github.com/LessUp/compress-kit.git
cd compress-kit

make build
make test
```

最小跨语言验证：

```bash
printf "Hello CompressKit\n" > input.txt
./algorithms/huffman/cpp/huffman_cpp encode input.txt output.huf
./algorithms/huffman/go/huffman_go decode output.huf restored.txt
diff input.txt restored.txt
```

## 文档

| 目标 | 链接 |
|------|------|
| 完整文档门户 | <https://lessup.github.io/compress-kit/> |
| 环境准备与首次运行 | <https://lessup.github.io/compress-kit/zh/guide/getting-started> |
| 算法对比 | <https://lessup.github.io/compress-kit/zh/guide/algorithms> |
| API 参考 | <https://lessup.github.io/compress-kit/zh/api/streaming> |
| 跨语言测试 | <https://lessup.github.io/compress-kit/zh/testing/cross-language> |

## 仓库结构

```text
algorithms/   # huffman、arithmetic、range、rle；每个算法含 cpp/go/rust
tests/        # 生成语料、streaming 契约、跨语言 conformance 矩阵
docs/         # VitePress 文档站点
openspec/     # 项目规范与已归档设计变更
```

## 工程基线

| 命令 | 用途 |
|------|------|
| `make build` | 构建全部 C++/Go/Rust CLI 工具 |
| `make test` | 运行单元、streaming、跨语言 conformance 测试 |
| `make test-conformance` | 单独运行可执行解码矩阵 |
| `make bench` | 运行基准脚本 |
| `npm run docs:build` | 构建文档站 |

已知限制：Range Coder 在大文件解码上存在已记录的性能问题；本地
conformance 和 benchmark 路径会对 Range 大样本做限制。处理大输入前请先阅读
Range Coder 文档。

## 许可证

[MIT 许可证](LICENSE) · 版权所有 © 2025-2026 LessUp
