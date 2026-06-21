# CompressKit

<p align="center">
  <img src="docs/public/logo.svg" width="120" alt="CompressKit Logo">
</p>

<p align="center">
  <strong>使用 C++17 实现的经典无损压缩算法。</strong>
</p>

<p align="center">
  <a href="https://github.com/LessUp/compress-kit/actions/workflows/ci.yml"><img src="https://github.com/LessUp/compress-kit/actions/workflows/ci.yml/badge.svg" alt="CI Status"></a>
  <a href="https://lessup.github.io/compress-kit/"><img src="https://img.shields.io/badge/Docs-在线文档-blue?logo=readthedocs&logoColor=white" alt="Documentation"></a>
  <a href="https://opensource.org/licenses/MIT"><img src="https://img.shields.io/badge/许可证-MIT-green.svg" alt="License"></a>
</p>

<p align="center">
  <a href="README.md">English</a> | <b>简体中文</b> | <a href="https://lessup.github.io/compress-kit/">文档站点</a>
</p>

CompressKit 是一个面向学习与验证的压缩算法仓库：同一组经典算法用
C++17 实现，通过统一命令行契约和 round-trip 测试验证正确性。
它不是黑盒压缩库，而是可以阅读、运行、对比和验证的算法实验室。

## 包含内容

| 算法 | 适用场景 |
|------|----------|
| Huffman 编码 | 通用文本/数据，学习前缀码 |
| 算术编码 | 理解熵编码与压缩率对比 |
| 区间编码 | 对比算术编码风格实现 |
| RLE 行程编码 | 高重复数据与简单格式学习 |

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

如需使用项目固定工具链，建议在 VS Code 或 Codespaces 中打开仓库内置的
`.devcontainer/`。若在本地维护文档：

```bash
npm ci
(cd docs && npm ci)
```

快速 round-trip 验证：

```bash
printf "Hello CompressKit\n" > input.txt
./build/huffman_cpp encode input.txt output.huf
./build/huffman_cpp decode output.huf restored.txt
diff input.txt restored.txt
```

## 文档

| 目标 | 链接 |
|------|------|
| 完整文档门户 | <https://lessup.github.io/compress-kit/> |
| 环境准备与首次运行 | <https://lessup.github.io/compress-kit/zh/guide/getting-started> |
| 算法对比 | <https://lessup.github.io/compress-kit/zh/guide/algorithms> |
| API 参考 | <https://lessup.github.io/compress-kit/zh/api/streaming> |

## 仓库结构

```text
algorithms/   # huffman、arithmetic、range、rle；每个算法含 cpp/
tests/        # 生成语料、CLI smoke 测试
docs/         # VitePress 文档站点
```

## 工程基线

| 命令 | 用途 |
|------|------|
| `make build` | 构建全部 C++ CLI 工具（CMake） |
| `make test` | 运行单元测试与 CLI smoke 测试 |
| `make lint` | clang-format dry-run |
| `npm run docs:build` | 构建文档站 |

## 许可证

[MIT 许可证](LICENSE) · 版权所有 © 2025-2026 LessUp
