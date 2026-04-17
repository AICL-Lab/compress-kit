# Encoding

<p align="center">
  <img src="docs/public/logo.svg" width="120" alt="Encoding Logo">
</p>

<p align="center">
  <strong>使用 C++17、Go 和 Rust 实现的经典压缩算法</strong>
</p>

<p align="center">
  <a href="https://github.com/LessUp/encoding/actions/workflows/ci.yml">
    <img src="https://github.com/LessUp/encoding/actions/workflows/ci.yml/badge.svg" alt="CI">
  </a>
  <a href="https://github.com/LessUp/encoding/actions/workflows/pages.yml">
    <img src="https://github.com/LessUp/encoding/actions/workflows/pages.yml/badge.svg" alt="文档">
  </a>
  <a href="https://lessup.github.io/encoding/">
    <img src="https://img.shields.io/badge/Docs-在线文档-blue" alt="文档">
  </a>
  <a href="https://opensource.org/licenses/MIT">
    <img src="https://img.shields.io/badge/许可证-MIT-yellow.svg" alt="许可证: MIT">
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

## 🧭 我应该使用哪种算法？

```
你的数据是否高度重复？
├── 是 → 使用 RLE（最快，适合重复模式）
└── 否 →
    是否需要最大压缩？
    ├── 是 → 使用算术编码（最接近熵限）
    └── 否 →
        速度是否关键？
        ├── 是 → 使用区间编码（快速 + 压缩好）
        └── 否 → 使用 Huffman（简单通用）
```

## 📊 算法对比

| 算法 | 压缩率 | 速度 | 适用场景 | 使用时机 |
|------|--------|------|----------|----------|
| **Huffman** | 中等 | 快 | 通用文本/数据 | 简单可靠的压缩 |
| **算术编码** | ★ 最高 | 中等 | 最大压缩需求 | 每个字节都很重要 |
| **区间编码** | ★ 高 | 快 | 平衡性能 | 速度与压缩的最佳平衡 |
| **RLE** | 可变 | ★ 最快 | 重复数据（位图、日志）| 数据有长串相同字节 |

## 🚀 快速开始

```bash
git clone https://github.com/LessUp/encoding.git
cd encoding
make build && make test
```

### 跨语言验证

```bash
# 使用 C++ 编码
./algorithms/huffman/cpp/huffman_cpp encode input.txt output.huf

# 使用 Go 解码 — 任何组合都可以！
./algorithms/huffman/go/huffman_go decode output.huf restored.txt
diff input.txt restored.txt  # ✓ 无输出 = 完全相同
```

**C++ ↔ Go ↔ Rust** — 所有实现共享相同的二进制格式。

## 🏗️ 项目结构

```
encoding/
├── algorithms/
│   ├── huffman/          # 前缀码压缩
│   ├── arithmetic/       # 算术编码
│   ├── range/            # 区间编码（字节级算术编码）
│   └── rle/              # 行程长度编码
│       ├── cpp/          #   C++17: 单文件，零依赖
│       ├── go/           #   Go 1.21+: 库 API + cmd/ CLI
│       ├── rust/         #   Rust 1.70+: rustc 或 cargo
│       └── benchmark/    #   性能测量脚本
├── docs/                 # VitePress 文档站点（en + zh）
├── specs/                # 规范驱动开发（RFC、产品规范）
├── tests/                # 测试数据生成
└── Makefile              # 构建、测试、基准测试入口
```

每种算法都有 **3 种语言实现**，使用相同的文件格式：

| 语言 | 构建方式 | 结构 |
|------|----------|------|
| **C++17** | `g++ -std=c++17 -O2` | 单文件，零依赖 |
| **Go** | `go build ./cmd` | 库 API（`package <algo>`）+ CLI（`cmd/main.go`） |
| **Rust** | `rustc -O` 或 `cargo` | `main.rs` 含可复用函数 |

## 📖 文档

| 资源 | 链接 |
|------|------|
| 📚 **完整文档** | [lessup.github.io/encoding](https://lessup.github.io/encoding/) |
| 🔧 **API 参考**（Go / Rust / C++） | [API 文档](https://lessup.github.io/encoding/zh/api/go) |
| 📈 **基准测试结果** | [性能对比](https://lessup.github.io/encoding/zh/benchmarks/results) |
| 🤝 **贡献指南** | [如何参与](https://lessup.github.io/encoding/zh/guide/contributing) |
| 📋 **技术规范**（RFC） | [specs/](specs/) |
| 📝 **更新日志** | [CHANGELOG.md](CHANGELOG.md) |

## 💻 构建与测试命令

| 命令 | 描述 |
|------|------|
| `make build` | 构建所有实现（C++、Go、Rust） |
| `make build-huffman` | 仅构建 Huffman |
| `make test` | 运行 Go + Rust 单元测试 |
| `make bench` | 运行性能基准测试 |
| `make clean` | 清理构建产物 |

## 🔬 Go 库使用示例

所有 Go 实现都提供可复用的库 API：

```go
import "huffman" // 或 "arithmetic"、"rle"

// 编码文件
err := huffman.EncodeFile("input.bin", "output.huf")

// 解码文件
err := huffman.DecodeFile("output.huf", "decoded.bin")
```

## 🏅 为什么存在这个项目

- **🎓 学习** — 并排阅读清晰、有良好文档的实现
- **🔬 对比** — 了解 C++、Go 和 Rust 如何处理相同算法
- **✅ 验证** — 跨语言测试保证输出格式一致
- **📐 SDD** — 采用规范驱动开发：每个功能都从规范开始

## 🤝 参与贡献

我们欢迎贡献！本项目遵循**规范驱动开发 (SDD)**：

1. **先阅读规范** — `/specs/` 是单一真实来源
2. **代码前先更新规范** — 如果接口有变，规范先改
3. **跨语言测试** — 验证 C++ ↔ Go ↔ Rust 兼容性

详见 [贡献指南](https://lessup.github.io/encoding/zh/guide/contributing)。

## 📄 许可证

[MIT 许可证](LICENSE) · 版权所有 © 2025-2026 LessUp
