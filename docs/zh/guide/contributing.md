# 参与贡献指南

CompressKit 是一个规范驱动的多语言压缩算法仓库。一个变更不是“某个实现能跑”就结束，而是相关规范、所有受影响语言实现、跨语言二进制契约三者都一致后才算完成。

## 从 OpenSpec 开始

以 `openspec/specs/` 作为事实来源：

| 规范 | 适用场景 |
|------|----------|
| `encoding-project` | 算法范围、质量门禁、安全限制、对外定位 |
| `core-architecture` | 目录结构、CLI 形态、二进制格式边界 |
| `cross-language-testing` | 兼容性矩阵、测试夹具、基准要求、已知限制 |

新增算法、修改二进制格式、调整 CLI 行为、扩大兼容性契约时，先创建 OpenSpec 变更。小型文档修正或“恢复既有规范行为”的实现 bugfix，可以直接同步现有规范。

## 开发基线

| 命令 | 用途 |
|------|------|
| `make build` | 编译 C++、Go、Rust 实现 |
| `make test` | 运行单元测试、shell 测试和跨语言 conformance 矩阵 |
| `make test-conformance` | 仅运行编码/解码兼容性矩阵 |
| `make lint` | 运行 `go vet` 和严格 Rust `clippy`，警告即失败 |
| `make format` | 运行 `gofmt`、`cargo fmt`、`clang-format` |
| `npm run docs:build` | 构建 VitePress 文档站 |

`make lint` 必须是真实门禁。不要用 shell fallback 吞掉 lint 失败；要么修复问题，要么说明某条 lint 为什么不适用于本项目。

## 实现标准

| 语言 | 要求 |
|------|------|
| C++17 | 保持单文件算法 CLI 与共享格式兼容，提交前使用 `.clang-format`。 |
| Go 1.21+ | 使用 `gofmt`、`go vet` 和符合 Go 习惯的包级测试。 |
| Rust 1.70+ | 每个 crate 保持 `cargo test`、`cargo fmt`、`cargo clippy --all-targets -- -D warnings` 干净。 |
| Python 3.8+ | 仅作为仓库脚本和 conformance 编排语言，不作为生产算法目标。 |

## 二进制兼容规则

- 每个算法 CLI 都必须保持 `encode|decode input output`。
- Huffman、Arithmetic、Range、RLE 格式必须在 C++、Go、Rust 之间兼容，除非已有获批 OpenSpec 变更。
- 安全限制属于契约：最大输入 4 GiB，最大解码输出 1 GiB。
- Range Coder 大文件解码性能问题是已知限制，不应作为顺手清理项处理。

## Pull Request 检查清单

- 相关 OpenSpec requirement 仍然成立，或 PR 包含对应规范变更。
- 本地 `make test` 通过。
- 触及相关文件时，`make lint` 与 `npm run docs:build` 通过。
- 二进制格式或 streaming adapter 行为变更有跨语言夹具覆盖。
- 文档更新只保留能帮助读者选择、使用或验证项目的信息。
