# Encoding

<p align="center">
  <img src="docs/public/logo.svg" width="120" alt="Encoding Logo">
</p>

<p align="center">
  <strong>经典无损压缩算法 · C++17、Go、Rust 实现</strong>
</p>

<p align="center">
  <a href="https://github.com/LessUp/encoding/actions/workflows/ci.yml">
    <img src="https://github.com/LessUp/encoding/actions/workflows/ci.yml/badge.svg" alt="CI Status">
  </a>
  <a href="https://lessup.github.io/encoding/">
    <img src="https://img.shields.io/badge/Docs-在线文档-blue?logo=readthedocs&logoColor=white" alt="Documentation">
  </a>
  <a href="https://github.com/LessUp/encoding/releases">
    <img src="https://img.shields.io/github/v/release/LessUp/encoding?include_prereleases&label=Release" alt="Release">
  </a>
  <a href="https://opensource.org/licenses/MIT">
    <img src="https://img.shields.io/badge/许可证-MIT-green.svg" alt="License">
  </a>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/C++-17-00599C.svg?logo=c%2B%2B" alt="C++17">
  <img src="https://img.shields.io/badge/Go-1.21+-00ADD8.svg?logo=go" alt="Go">
  <img src="https://img.shields.io/badge/Rust-1.70+-DEA584.svg?logo=rust" alt="Rust">
  <img src="https://img.shields.io/badge/Python-3.8+-3776AB.svg?logo=python" alt="Python">
</p>

<p align="center">
  <a href="README.md">English</a> | <b>简体中文</b> | <a href="https://lessup.github.io/encoding/">📖 文档站点</a>
</p>

---

## 📋 前置要求

| 语言 | 最低版本 | 安装方式 |
|------|----------|----------|
| C++ | GCC 9+ / Clang 10+ | `apt install g++` 或 `brew install gcc` |
| Go | 1.21+ | [golang.org/dl](https://golang.org/dl) |
| Rust | 1.70+ | [rustup.rs](https://rustup.rs) |
| Python | 3.8+ | 仅需用于生成测试数据 |

验证环境：
```bash
g++ --version    # 必须支持 -std=c++17
go version
rustc --version
```

---

## ✨ 特性

- 🔤 **多语言实现** — C++17、Go 1.21+、Rust 1.70+ 三种语言完整实现
- 🔗 **跨语言兼容** — 用一种语言编码，另一种语言解码
- 📚 **面向学习** — 清晰、文档完善的代码，便于学习和对比
- 🧪 **完善测试** — CI 包含单元测试和跨语言验证
- 📊 **性能基准** — 跨语言性能对比

## 🧮 算法

| 算法 | 压缩率 | 速度 | 适用场景 |
|------|--------|------|----------|
| [**Huffman**](https://lessup.github.io/encoding/zh/guide/algorithms#huffman-编码) | 中等 | 快 | 通用文本/数据 |
| [**算术编码**](https://lessup.github.io/encoding/zh/guide/algorithms#算术编码) | 最高 | 中等 | 最大压缩需求 |
| [**区间编码**](https://lessup.github.io/encoding/zh/guide/algorithms#区间编码) | 高 | 快 | 平衡性能 |
| [**RLE**](https://lessup.github.io/encoding/zh/guide/algorithms#行程长度编码-rle) | 可变 | 最快 | 重复数据（位图、日志） |

### 算法选择指南

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

## 🚀 快速开始

```bash
git clone https://github.com/LessUp/encoding.git
cd encoding

# 1. 构建所有实现
make build

# 2. 生成测试数据（需要 Python 3.8+）
make test-data

# 3. 运行测试
make test
```

### 快速验证

```bash
# 创建测试文件
echo "Hello, World! Hello, World!" > input.txt

# 使用 C++ 编码
./algorithms/huffman/cpp/huffman_cpp encode input.txt output.huf

# 使用 Go 解码
./algorithms/huffman/go/huffman_go decode output.huf restored.txt

# 验证
diff input.txt restored.txt && echo "✓ 跨语言验证通过"
```

### 跨语言验证（替代方式）

```bash
# 使用 C++ 编码
./algorithms/huffman/cpp/huffman_cpp encode input.txt output.huf

# 使用 Go 解码
./algorithms/huffman/go/huffman_go decode output.huf restored.txt
diff input.txt restored.txt  # 无输出 = 完全相同
```

**C++ ↔ Go ↔ Rust** — 所有实现共享相同的二进制格式。

## 项目结构

```
encoding/
├── algorithms/           # 压缩算法实现
│   ├── huffman/         # 前缀码压缩
│   ├── arithmetic/      # 算术编码
│   ├── range/           # 区间编码
│   └── rle/             # 行程编码
│       ├── cpp/         # C++17: 单文件，零依赖
│       ├── go/          # Go 1.21+: 库 API + CLI
│       ├── rust/        # Rust 1.70+: rustc 或 cargo
│       └── benchmark/   # 性能脚本
├── docs/                # VitePress 文档站点（en + zh）
├── specs/               # 规范驱动开发文档
├── tests/               # 测试数据生成
└── Makefile             # 构建入口
```

## 构建与测试

| 命令 | 描述 |
|------|------|
| `make build` | 构建所有实现 |
| `make test` | 运行单元测试 |
| `make bench` | 运行基准测试 |
| `make clean` | 清理构建产物 |

## 💻 使用方法

所有实现遵循统一的 CLI 接口：

```bash
<可执行文件> <encode|decode> <输入> <输出>
```

### CLI 示例

```bash
# Huffman - C++
./algorithms/huffman/cpp/huffman_cpp encode input.txt output.huf
./algorithms/huffman/cpp/huffman_cpp decode output.huf restored.txt

# Huffman - Go
./algorithms/huffman/go/huffman_go encode input.txt output.huf
./algorithms/huffman/go/huffman_go decode output.huf restored.txt

# Huffman - Rust
./algorithms/huffman/rust/huffman_rust encode input.txt output.huf
./algorithms/huffman/rust/huffman_rust decode output.huf restored.txt

# 所有工具都支持 --help 查看详细选项
./algorithms/huffman/go/huffman_go --help
```

### Go 库使用

```go
import "github.com/LessUp/encoding/algorithms/huffman/go"

err := huffman.EncodeFile("input.bin", "output.huf")
err = huffman.DecodeFile("output.huf", "decoded.bin")
```

注意：作为库使用时直接导入包调用函数。独立 CLI 使用请用 `go build -o huffman_go ./cmd` 构建。

## 📚 文档

| 资源 | 链接 |
|------|------|
| 📖 完整文档 | [lessup.github.io/encoding](https://lessup.github.io/encoding/) |
| 🔧 API 参考 | [Go](https://lessup.github.io/encoding/zh/api/go) · [Rust](https://lessup.github.io/encoding/zh/api/rust) · [C++](https://lessup.github.io/encoding/zh/api/cpp) |
| 📊 基准测试 | [性能对比](https://lessup.github.io/encoding/zh/benchmarks/results) |
| 🤝 贡献指南 | [如何参与](https://lessup.github.io/encoding/zh/guide/contributing) |
| 📋 技术规范 | [specs/](specs/) |

## 🎯 项目特点

- **📖 学习** — 对比 C++、Go、Rust 的清晰实现
- **✅ 验证** — 跨语言测试保证格式兼容
- **📐 SDD** — 采用规范驱动开发方法

## 🤝 参与贡献

本项目遵循**规范驱动开发 (SDD)**：

1. 先阅读规范 — `/specs/` 是单一真实来源
2. 代码前先更新规范 — 接口变更时规范优先
3. 跨语言测试 — 验证 C++ ↔ Go ↔ Rust 兼容性

详见 [贡献指南](https://lessup.github.io/encoding/zh/guide/contributing)。

## ⚠️ 安全说明

- **最大输入文件大小：** 4 GiB
- **最大输出文件大小：** 1 GiB（防止解压炸弹攻击）
- 所有二进制格式包含完整性校验
- 文件格式在主要版本内稳定且向后兼容

## 📜 更新日志

查看 [CHANGELOG.md](CHANGELOG.md) 了解版本历史和迁移指南。

## 许可证

[MIT 许可证](LICENSE) · 版权所有 © 2025-2026 LessUp
