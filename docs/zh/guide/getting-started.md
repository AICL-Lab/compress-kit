# 快速开始

本指南将帮助您配置开发环境、构建实现并运行测试。

## 前置要求

### 必需工具

| 工具 | 最低版本 | 用途 |
|------|----------|------|
| g++ 或 clang++ | 9+ / 10+ | C++17 编译 |
| CMake | 3.16+ | 构建系统 |
| Python | 3.8+ | 测试编排与基准测试脚本 |
| Make | 任意版本 | 构建自动化 |

### 可选工具

| 工具 | 用途 |
|------|------|
| Node.js 20.19+ | 文档站点 |
| clang-format | C++ 代码格式化 |

### 安装方法

::: code-group

```bash [Ubuntu/Debian]
sudo apt update
sudo apt install g++ cmake python3 make
```

```bash [macOS (Homebrew)]
brew install gcc cmake python3 make
```

```bash [Windows (Chocolatey)]
choco install mingw cmake python3 make
```

:::

## 克隆与构建

### 克隆仓库

```bash
git clone https://github.com/LessUp/compress-kit.git
cd compress-kit
```

### 使用 Dev Container（推荐）

仓库内置了 `.devcontainer/`，提供项目支持的 C++ / Node
工具链。请在 VS Code 或 Codespaces 中打开仓库后选择 **Reopen in
Container**，即可获得与项目文档一致的基线环境。

容器创建时会自动执行可复现的安装路径：

```bash
npm ci
(cd docs && npm ci)
```

### 本地可复现的文档工具链

如果您仍在宿主机上工作，请在构建文档前执行同样的
锁文件安装：

```bash
npm ci
(cd docs && npm ci)
```

### 构建所有实现

```bash
make build
```

`make build` 封装了 CMake 调用（`cmake -S . -B build && cmake --build build`），将所有算法实现编译到 `build/` 目录。

### 构建特定算法

```bash
make build-huffman      # Huffman
make build-arithmetic   # Arithmetic
make build-range        # Range Coder
make build-rle          # RLE
```

### 手动编译

如果您更喜欢直接调用 CMake：

```bash
cmake -S . -B build && cmake --build build
./build/huffman_cpp encode input.bin output.huf
./build/huffman_cpp decode output.huf restored.bin
```

## 运行测试

### 运行所有测试

```bash
make test
```

这将运行共享 streaming 层测试以及各算法特定测试套件。

## 运行基准测试

### 运行所有基准测试

```bash
make bench
```

这会生成测试数据、运行已验证的基准矩阵，并把原始报告写入 `reports/`，
同时刷新交互式文档图表使用的 `docs/.vitepress/data/benchmarks.json`。

### 基准测试产物

| 产物 | 用途 |
|------|------|
| `reports/*.txt` | 每次运行的原始基准日志 |
| `docs/.vitepress/data/benchmarks.json` | 文档基准图表读取的生成快照 |
| `/zh/benchmarks/results` | 面向读者的生成 JSON 视图 |

当前仓库快照中的基准事实应以生成的 JSON 文件和图表为准。

## Makefile 命令参考

| 命令 | 描述 |
|------|------|
| `make build` | 构建所有算法实现（封装 CMake） |
| `make build-huffman` | 仅构建 Huffman 实现 |
| `make build-arithmetic` | 仅构建 Arithmetic 实现 |
| `make build-range` | 仅构建 Range Coder 实现 |
| `make build-rle` | 仅构建 RLE 实现 |
| `make test` | 运行单元和 streaming 测试 |
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

### Range Coder 性能

Range Coder 解码器对大于 500 KiB 的文件存在已知性能问题，因此基准测试
会为 Range 结果使用较小的 `small_dictionary_like` 数据集。手动检查时也请
保持在类似的小样本范围内：

```bash
# 创建较小的测试文件 (100 KiB)
dd if=tests/data/random_10MiB.bin of=small.bin bs=1024 count=100

# 使用较小文件测试
./build/rangecoder_cpp encode small.bin small.enc
./build/rangecoder_cpp decode small.enc small.dec
```

## 下一步

- 了解 [算法详解](/zh/guide/algorithms) 及其差异
- 探索 [项目结构](/zh/guide/project-structure)
- 查看 [CHANGELOG](https://github.com/LessUp/compress-kit/blob/master/CHANGELOG.md) 了解最新更新
