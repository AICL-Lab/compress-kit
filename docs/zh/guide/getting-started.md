# 快速开始

本指南将帮助您配置开发环境、构建实现并运行测试。

## 前置要求

### 必需工具

| 工具 | 最低版本 | 用途 |
|------|----------|------|
| g++ 或 clang++ | 9+ / 10+ | C++17 编译 |
| Go | 1.21+ | Go 实现 |
| Rust (cargo) | 1.70+ | Rust 实现 |
| Python | 3.8+ | 基准测试脚本 |
| Make | 任意版本 | 构建自动化 |

### 可选工具

| 工具 | 用途 |
|------|------|
| Node.js 18+ | 文档站点开发 |
| clang-format | C++ 代码格式化 |

### 安装方法

::: code-group

```bash [Ubuntu/Debian]
sudo apt update
sudo apt install g++ golang rustc python3 make
```

```bash [macOS (Homebrew)]
brew install gcc go rust python3 make
```

```bash [Windows (Chocolatey)]
choco install mingw golang rust python3 make
```

:::

## 克隆与构建

### 克隆仓库

```bash
git clone https://github.com/LessUp/compresskit.git
cd encoding
```

### 构建所有实现

```bash
make build
```

这将编译所有三种语言的所有算法实现。

### 构建特定算法

```bash
make build-huffman      # Huffman 编码 (C++, Go, Rust)
make build-arithmetic   # 算术编码
make build-range        # 区间编码
make build-rle          # 行程长度编码
```

### 手动编译

如果您更喜欢手动编译：

::: code-group

```bash [C++]
cd algorithms/huffman/cpp
g++ -std=c++17 -O2 main.cpp -o huffman_cpp
./huffman_cpp encode input.bin output.huf
./huffman_cpp decode output.huf restored.bin
```

```bash [Go]
cd algorithms/huffman/go
go build -o huffman_go ./cmd
./huffman_go encode input.bin output.huf
./huffman_go decode output.huf restored.bin
```

```bash [Rust]
cd algorithms/huffman/rust
rustc -O main.rs -o huffman_rust
./huffman_rust encode input.bin output.huf
./huffman_rust decode output.huf restored.bin
```

:::

## 跨语言验证

本项目的一个关键特性是所有实现使用相同的文件格式，支持跨语言验证：

```bash
# 使用 C++ 编码
./algorithms/huffman/cpp/huffman_cpp encode input.bin encoded.huf

# 使用 Go 解码
./algorithms/huffman/go/huffman_go decode encoded.huf decoded.bin

# 验证正确性
diff input.bin decoded.bin  # 无输出 = 相同
```

任何组合都有效：**C++ ↔ Go ↔ Rust**

## 运行测试

### 运行所有测试

```bash
make test
```

这将运行所有算法的 Go 和 Rust 单元测试。

### 运行特定算法测试

```bash
# Go 测试
cd algorithms/huffman/go && go test ./...

# Rust 测试
cd algorithms/huffman/rust && rustc --test main.rs -o test && ./test
```

## 运行基准测试

### 运行所有基准测试

```bash
make bench
```

这将生成测试数据并运行跨语言基准测试。

### 基准测试输出

报告保存到 `reports/` 目录。示例输出：

```
Algorithm: Huffman
Language: C++
Input: 10 MiB random data
Encode: 245 ms (40.8 MiB/s)
Decode: 198 ms (50.5 MiB/s)
Compression ratio: 1.23
```

## Makefile 命令参考

| 命令 | 描述 |
|------|------|
| `make build` | 构建所有算法实现 |
| `make build-huffman` | 仅构建 Huffman 实现 |
| `make build-arithmetic` | 仅构建算术编码实现 |
| `make build-range` | 仅构建区间编码实现 |
| `make build-rle` | 仅构建 RLE 实现 |
| `make test` | 运行所有 Go 和 Rust 单元测试 |
| `make bench` | 生成测试数据并运行基准测试 |
| `make test-data` | 仅生成测试数据 |
| `make clean` | 删除所有构建产物和报告 |

## 故障排除

### C++ 编译错误

```bash
# 检查编译器版本
g++ --version  # 应为 9+

# 如果 g++ 失败，尝试 clang
clang++ -std=c++17 -O2 main.cpp -o huffman_cpp
```

### Go 模块问题

```bash
# 确保 Go 工作区包含所有模块
go work use ./algorithms/huffman/go
go work use ./algorithms/arithmetic/go
go work use ./algorithms/range/go
go work use ./algorithms/rle/go
```

### Rust 构建错误

```bash
# 更新 Rust 工具链
rustup update stable

# 检查版本
rustc --version  # 应为 1.70+
```

### 区间编码性能

区间编码解码器对大于 500KB 的文件存在已知性能问题。请使用较小的测试文件：

```bash
# 创建较小的测试文件 (100KB)
dd if=tests/data/random_10MiB.bin of=/tmp/small.bin bs=1024 count=100

# 使用较小文件测试
./algorithms/range/cpp/rangecoder_cpp encode /tmp/small.bin /tmp/small.enc
./algorithms/range/cpp/rangecoder_cpp decode /tmp/small.enc /tmp/small.dec
```

## 下一步

- 了解 [算法详解](/zh/guide/algorithms) 及其差异
- 探索 [项目结构](/zh/guide/project-structure)
- 查看 [CHANGELOG](https://github.com/LessUp/compresskit/blob/master/CHANGELOG.md) 了解最新更新
