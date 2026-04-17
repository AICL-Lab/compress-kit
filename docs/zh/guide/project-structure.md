# 项目结构

本指南介绍项目组织方式、文件格式和所有实现中使用的约定。

## 目录结构

```
encoding/
├── algorithms/huffman/              # Huffman 编码实现
│   ├── cpp/              #   C++ 单文件实现
│   ├── go/               #   Go 模块 (go.mod)
│   ├── rust/             #   Rust 实现
│   └── benchmark/        #   跨语言基准测试脚本
├── algorithms/arithmetic/           # 算术编码实现
│   ├── cpp/              #   C++ 单文件实现
│   ├── go/               #   Go 实现
│   ├── rust/             #   Rust 实现
│   └── benchmark/        #   跨语言基准测试
├── algorithms/range/                # 区间编码实现
│   ├── cpp/              #   C++ 单文件实现
│   ├── go/               #   Go 库 + CLI
│   ├── rust/             #   Rust 库包 + CLI
│   └── benchmark/        #   跨语言基准测试
├── algorithms/rle/                  # 行程长度编码
│   ├── cpp/              #   C++ 单文件实现
│   ├── go/               #   Go 实现
│   ├── rust/             #   Rust 实现
│   └── benchmark/        #   跨语言基准测试
├── tests/                # 测试数据生成
│   ├── gen_testdata.py   #   生成基准测试文件
│   └── data/             #   生成的测试数据
├── docs/                 # 文档站点 (VitePress)
│   ├── .vitepress/       #   VitePress 配置
│   ├── en/               #   英文文档
│   ├── zh/               #   中文文档
│   └── public/           #   静态资源（logo 等）
├── .github/workflows/    # GitHub Actions CI/CD
├── Makefile              # 构建、测试和基准测试入口
├── package.json          # 文档 npm 脚本
└── go.work               # Go 工作区（多模块）
```

## 语言实现标准

| 语言 | 版本 | 构建方式 | 特性 |
|------|------|----------|------|
| **C++** | C++17 | `g++ -std=c++17 -O2` | 单文件，零依赖 |
| **Go** | 1.21+ | Go 模块 (`go.mod`) | 区间编码提供库 API |
| **Rust** | 1.70+ | Cargo / `rustc` | 区间编码提供库包 |

## 统一 CLI 接口

所有实现遵循相同的命令行模式：

```bash
<algorithm>_<lang> encode <input_file> <output_file>
<algorithm>_<lang> decode <input_file> <output_file>
```

### 二进制名称

| 算法 | C++ | Go | Rust |
|------|-----|-----|------|
| Huffman | `huffman_cpp` | `huffman_go` | `huffman_rust` |
| 算术编码 | `arithmetic_cpp` | `arithmetic_go` | `arithmetic_rust` |
| 区间编码 | `rangecoder_cpp` | `rangecoder_go` | `rangecoder` (cargo) |
| RLE | `rle_cpp` | `rle_go` | `rle_rust` |

### 使用示例

```bash
# 使用 Huffman (C++) 编码
./huffman_cpp encode document.txt document.huf

# 使用区间编码 (Go) 解码
./rangecoder_go decode data.rcnc data.bin

# 使用 RLE (Rust) 编码
./rle_rust encode bitmap.raw bitmap.rle
```

## 文件格式兼容性

相同算法的所有语言实现使用**相同的二进制格式**。

### 格式汇总

| 算法 | 魔数头 | 扩展名 | 结构 |
|------|--------|--------|------|
| Huffman | `HFMN` | `.huf` | 魔数 + 频率表 + 比特流 |
| 算术编码 | `AENC` | `.aenc` | 魔数 + 频率表 + 比特流 |
| 区间编码 | `RCNC` | `.rcnc` | 魔数 + 频率表 + 字节流 |
| RLE | 无 | `.rle` | (计数: 4B 小端, 值: 1B) 对 |

### 跨语言验证

| 编码 ↓ / 解码 → | C++ | Go | Rust |
|-----------------|-----|-----|------|
| C++ | ✓ | ✓ | ✓ |
| Go | ✓ | ✓ | ✓ |
| Rust | ✓ | ✓ | ✓ |

任何组合都有效：**C++ ↔ Go ↔ Rust**

## CI/CD 流水线

### 工作流

| 工作流 | 文件 | 触发条件 | 目的 |
|--------|------|----------|------|
| **CI** | `.github/workflows/ci.yml` | Push / PR | 构建、测试、正确性验证 |
| **Pages** | `.github/workflows/pages.yml` | `docs/` 变更 | 部署文档 |

### CI 任务矩阵

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│  build-cpp  │     │  build-go   │     │ build-rust  │
│  ├ Ubuntu   │     │  Ubuntu     │     │  Ubuntu     │
│  └ macOS    │     │             │     │             │
└──────┬──────┘     └──────┬──────┘     └──────┬──────┘
       │                   │                   │
       └───────────────────┼───────────────────┘
                           │
                    ┌──────▼──────┐
                    │ correctness │
                    │   tests     │
                    │  (Python)   │
                    └─────────────┘
```

### CI 检查

1. **构建任务**: 在 Ubuntu 和 macOS 上编译所有实现 (C++)
2. **测试任务**: 运行 Go `go test` 和 Rust `cargo test`
3. **检查任务**: Go vet、Rust clippy
4. **正确性任务**: 跨语言编码/解码验证

## 安全考虑

### 输入/输出大小限制

所有实现强制执行以下限制：

| 限制 | 值 | 目的 |
|------|-----|------|
| 最大输入大小 | 4 GiB | 防止频率溢出 |
| 最大输出大小 | 1 GiB | 防止解压缩炸弹 |

这些限制在处理开始前应用，防止：
- 整数溢出攻击
- 解压缩炸弹攻击
- 过度内存使用

## 构建系统详情

### Makefile 目标

| 目标 | 描述 |
|------|------|
| `build` | 构建所有实现 |
| `build-huffman` | 构建 Huffman (C++, Go, Rust) |
| `build-arithmetic` | 构建算术编码 (C++, Go, Rust) |
| `build-range` | 构建区间编码 (C++, Go, Rust) |
| `build-rle` | 构建 RLE (C++, Go, Rust) |
| `test` | 运行所有 Go 和 Rust 单元测试 |
| `bench` | 生成测试数据并运行所有基准测试 |
| `test-data` | 仅生成测试数据 |
| `clean` | 删除所有构建产物和报告 |

### Go 工作区

项目使用 Go 工作区管理多个模块：

```go
// go.work
go 1.21

use (
    ./algorithms/huffman/go
    ./algorithms/arithmetic/go
    ./algorithms/range/go
    ./algorithms/rle/go
)
```

这允许：
- 跨模块依赖
- 统一构建命令
- 所有 Go 模块的 IDE 支持

## 文档站点

文档使用 [VitePress](https://vitepress.dev/) 构建：

```bash
# 安装依赖
npm install

# 启动开发服务器
npm run docs:dev

# 构建生产版本
npm run docs:build
```

构建后的站点位于 `docs/.vitepress/dist/`，并部署到 GitHub Pages。

---

## 相关文档

- [快速开始](/zh/guide/getting-started) - 设置和基础用法
- [算法详解](/zh/guide/algorithms) - 算法说明
- [GitHub 仓库](https://github.com/LessUp/encoding) - 源代码
