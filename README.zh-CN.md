# Encoding

<p align="center">
  <img src="docs/public/logo.svg" width="120" alt="Encoding Logo">
</p>

<p align="center">
  <a href="https://github.com/LessUp/encoding/actions/workflows/ci.yml">
    <img src="https://github.com/LessUp/encoding/actions/workflows/ci.yml/badge.svg" alt="CI">
  </a>
  <a href="https://github.com/LessUp/encoding/actions/workflows/pages.yml">
    <img src="https://github.com/LessUp/encoding/actions/workflows/pages.yml/badge.svg" alt="文档">
  </a>
  <a href="https://github.com/LessUp/encoding/releases">
    <img src="https://img.shields.io/github/v/release/LessUp/encoding?include_prereleases" alt="发布版本">
  </a>
  <a href="https://opensource.org/licenses/MIT">
    <img src="https://img.shields.io/badge/License-MIT-yellow.svg" alt="许可证: MIT">
  </a>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/C++-17-blue.svg" alt="C++17">
  <img src="https://img.shields.io/badge/Go-1.21+-00ADD8.svg" alt="Go 1.21+">
  <img src="https://img.shields.io/badge/Rust-1.70+-orange.svg" alt="Rust 1.70+">
  <img src="https://img.shields.io/badge/Python-3.8+-3776AB.svg" alt="Python 3.8+">
</p>

<p align="center">
  <a href="README.md">English</a> | <b>简体中文</b> | <a href="https://lessup.github.io/encoding/">文档站点</a>
</p>

---

## 🚀 快速开始

```bash
# 克隆仓库
git clone https://github.com/LessUp/encoding.git
cd encoding

# 构建所有实现
make build

# 运行测试
make test

# 运行基准测试
make bench
```

## ✨ 特性

- **4 种经典算法**：Huffman、算术编码、区间编码和 RLE
- **3 种语言**：C++17、Go 1.21+、Rust 1.70+
- **跨语言兼容**：所有实现使用相同的二进制格式
- **面向学习**：文档侧重算法原理和对比
- **生产就绪**：完整的 CI/CD 与自动化测试验证

## 📊 算法对比

| 算法 | 压缩率 | 速度 | 复杂度 | 适用场景 |
|------|--------|------|--------|----------|
| Huffman | 中等 | 快 | O(n log σ) | 通用场景 |
| 算术编码 | 高 | 中等 | O(n) | 最大压缩 |
| 区间编码 | 高 | 快 | O(n) | 平衡性能 |
| RLE | 可变 | 极快 | O(n) | 重复数据 |

## 📖 文档

| 资源 | 描述 | 链接 |
|------|------|------|
| **文档站点** | 完整的中英文双语文档 | [lessup.github.io/encoding](https://lessup.github.io/encoding/) |
| **项目规范** | 产品需求、RFC 和技术规范（单一真实来源） | [specs/](specs/) |
| **快速开始** | 环境配置、构建和基本用法 | [指南 →](https://lessup.github.io/encoding/zh/guide/getting-started) |
| **算法详解** | 算法说明和对比 | [指南 →](https://lessup.github.io/encoding/zh/guide/algorithms) |
| **项目结构** | 目录结构和约定 | [指南 →](https://lessup.github.io/encoding/zh/guide/project-structure) |
| **更新日志** | 版本历史和发布说明 | [查看 →](CHANGELOG.md) |

## 💡 使用示例

```bash
# 使用 Huffman (C++) 编码
./algorithms/huffman/cpp/huffman_cpp encode input.txt output.huf

# 使用另一种语言 (Go) 解码
./algorithms/huffman/go/huffman_go decode output.huf restored.txt

# 验证正确性
diff input.txt restored.txt  # 无输出 = 相同
```

## 🛠️ 构建选项

| 命令 | 描述 |
|------|------|
| `make build` | 构建所有实现 |
| `make build-huffman` | 仅构建 Huffman |
| `make build-arithmetic` | 仅构建算术编码 |
| `make build-range` | 仅构建区间编码 |
| `make build-rle` | 仅构建 RLE |
| `make test` | 运行所有测试 |
| `make bench` | 运行性能基准测试 |
| `make clean` | 清理构建产物 |

## 🤝 参与贡献

我们欢迎贡献！请阅读我们的[贡献指南](CONTRIBUTING.md)了解：

- C++、Go 和 Rust 的代码风格指南
- 测试要求
- Pull Request 流程

## 🙏 致谢

本项目受压缩算法教育资源的启发，旨在提供：

- 用于学习的清晰、可读实现
- 公平的跨语言性能对比
- 通过广泛测试验证的正确实现

## 📄 许可证

本项目基于 [MIT 许可证](LICENSE) 开源。

版权所有 © 2025-2026 LessUp

---

<p align="center">
  <sub>用 ❤️ 为开源社区构建</sub>
</p>
