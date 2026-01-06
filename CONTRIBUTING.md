# Contributing to encoding | 贡献指南

Thank you for your interest in contributing to this project! 感谢你对本项目的关注！

## Table of Contents | 目录

- [Prerequisites | 前置要求](#prerequisites--前置要求)
- [Development Setup | 开发环境配置](#development-setup--开发环境配置)
- [Code Style | 代码风格](#code-style--代码风格)
- [Testing | 测试](#testing--测试)
- [Pull Request Process | PR 流程](#pull-request-process--pr-流程)

## Prerequisites | 前置要求

Before you begin, ensure you have the following tools installed:

开始之前，请确保已安装以下工具：

### Required | 必需

- **C++ Compiler**: g++ 9+ or clang++ 10+ (with C++17 support)
- **Go**: 1.19 or later
- **Rust**: 1.70 or later (with rustfmt and clippy)
- **Python**: 3.8 or later (for benchmark scripts)

### Optional | 可选

- **Make**: For simplified build commands
- **Docker**: For consistent build environment

### Installation Examples | 安装示例

**Ubuntu/Debian:**
```bash
sudo apt update
sudo apt install g++ golang rustc python3
```

**macOS (with Homebrew):**
```bash
brew install gcc go rust python3
```

**Windows:**
- Install [MSYS2](https://www.msys2.org/) for g++
- Install [Go](https://go.dev/dl/)
- Install [Rust](https://rustup.rs/)
- Install [Python](https://www.python.org/downloads/)

## Development Setup | 开发环境配置

1. **Fork and clone the repository | Fork 并克隆仓库**

```bash
git clone https://github.com/YOUR_USERNAME/encoding.git
cd encoding
```

2. **Verify your environment | 验证环境**

```bash
# Check C++ compiler
g++ --version

# Check Go
go version

# Check Rust
rustc --version
cargo --version

# Check Python
python3 --version
```

3. **Build all implementations | 构建所有实现**

```bash
# C++ (example: Huffman)
cd huffman/cpp && g++ -std=c++17 -O2 main.cpp -o huffman_cpp && cd ../..

# Go (example: Huffman)
cd huffman/go && go build -o huffman_go . && cd ../..

# Rust (example: Huffman)
cd huffman/rust && rustc -O main.rs -o huffman_rust && cd ../..
```

## Code Style | 代码风格

### C++

- Follow [Google C++ Style Guide](https://google.github.io/styleguide/cppguide.html)
- Use 4 spaces for indentation
- Use `snake_case` for functions and variables
- Use `PascalCase` for classes and structs

```bash
# Format check (if clang-format is available)
clang-format --dry-run --Werror *.cpp
```

### Go

- Use `gofmt` for formatting (mandatory)
- Use `go vet` for static analysis
- Follow [Effective Go](https://go.dev/doc/effective_go)

```bash
# Format code
gofmt -w .

# Check for issues
go vet ./...
```

### Rust

- Use `rustfmt` for formatting (mandatory)
- Use `clippy` for linting
- Follow [Rust API Guidelines](https://rust-lang.github.io/api-guidelines/)

```bash
# Format code
cargo fmt

# Lint code
cargo clippy -- -D warnings
```

### Python

- Follow [PEP 8](https://peps.python.org/pep-0008/)
- Use 4 spaces for indentation

## Testing | 测试

### Running Benchmarks | 运行基准测试

```bash
# Generate test data
python3 tests/gen_testdata.py

# Run all benchmarks
python3 scripts/run_all_bench.py

# Run specific algorithm benchmark
cd huffman/benchmark && python3 bench.py
cd Run-Length/benchmark && python3 bench.py
```

### Verifying Correctness | 验证正确性

Each implementation should pass the encode-decode round-trip test:

每个实现都应通过编码-解码往返测试：

```bash
# Example: Huffman C++
./huffman_cpp encode input.bin encoded.huf
./huffman_cpp decode encoded.huf decoded.bin
diff input.bin decoded.bin  # Should produce no output
```

### Cross-Language Verification | 跨语言验证

Files encoded by one language implementation should be decodable by others:

一种语言编码的文件应能被其他语言正确解码：

```bash
# Encode with C++, decode with Go
./huffman_cpp encode input.bin encoded.huf
./huffman_go decode encoded.huf decoded.bin
diff input.bin decoded.bin
```

## Pull Request Process | PR 流程

1. **Create a feature branch | 创建功能分支**

```bash
git checkout -b feature/your-feature-name
```

2. **Make your changes | 进行修改**

- Write clean, readable code
- Add comments for complex logic
- Update documentation if needed

3. **Test your changes | 测试修改**

```bash
# Build and test
# Run benchmarks to verify correctness
```

4. **Commit with clear messages | 提交清晰的 commit 信息**

```bash
git commit -m "feat: add XXX feature"
# or
git commit -m "fix: resolve XXX issue"
```

Commit message format | 提交信息格式:
- `feat:` New feature | 新功能
- `fix:` Bug fix | 修复 bug
- `docs:` Documentation | 文档
- `style:` Code style | 代码风格
- `refactor:` Refactoring | 重构
- `test:` Tests | 测试
- `chore:` Maintenance | 维护

5. **Push and create PR | 推送并创建 PR**

```bash
git push origin feature/your-feature-name
```

Then create a Pull Request on GitHub.

6. **PR Checklist | PR 检查清单**

Before submitting, ensure:

提交前请确保：

- [ ] Code compiles without errors | 代码编译无错误
- [ ] All tests pass | 所有测试通过
- [ ] Code follows the style guide | 代码符合风格指南
- [ ] Documentation is updated | 文档已更新
- [ ] CHANGELOG.md is updated (if applicable) | CHANGELOG 已更新（如适用）

## Questions? | 有问题？

Feel free to open an issue if you have any questions or suggestions!

如有任何问题或建议，欢迎提交 Issue！
