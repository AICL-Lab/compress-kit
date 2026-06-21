# 参与贡献指南

CompressKit 是一个 C++17 压缩算法仓库。一个变更不是“实现能跑”就结束，而是实现正确且文档化的二进制格式契约保持一致后才算完成。

## 开发基线

| 命令 | 用途 |
|------|------|
| `make build` | 编译 C++ 实现（封装 CMake） |
| `make test` | 运行单元测试和 streaming 契约测试 |
| `make lint` | 运行 `clang-format` 检查 |
| `make format` | 运行 `clang-format` |
| `npm run docs:build` | 构建 VitePress 文档站 |

`make lint` 必须是真实门禁。不要用 shell fallback 吞掉 lint 失败；要么修复问题，要么说明某条 lint 为什么不适用于本项目。

## 实现标准

| 语言 | 要求 |
|------|------|
| C++17 | 保持单文件算法 CLI 与共享格式兼容，提交前使用 `.clang-format`。 |
| Python 3.8+ | 仅作为仓库脚本和基准编排语言，不作为生产算法目标。 |

## 二进制兼容规则

- 每个算法 CLI 都必须保持 `encode|decode input output`。
- Huffman、Arithmetic、Range、RLE 格式必须保持稳定。
- 安全限制属于契约：最大输入 4 GiB，最大解码输出 1 GiB。
- Range Coder 大文件解码性能问题是已知限制，不应作为顺手清理项处理。

## Pull Request 检查清单

- 本地 `make test` 通过。
- 触及相关文件时，`make lint` 与 `npm run docs:build` 通过。
- 文档更新只保留能帮助读者选择、使用或验证项目的信息。
