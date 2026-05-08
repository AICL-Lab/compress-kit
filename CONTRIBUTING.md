# Contributing to CompressKit | 贡献指南

感谢你对 CompressKit 项目的关注！本指南将帮助你快速参与项目开发。

## 前置要求 | Prerequisites

### 必需工具

- **C++ 编译器**: g++ 9+ 或 clang++ 10+ (支持 C++17)
- **Go**: 1.19 或更高版本
- **Rust**: 1.70 或更高版本（包含 rustfmt 和 clippy）
- **Python**: 3.8 或更高版本（用于测试脚本）
- **Make**: 简化构建命令

### 安装示例

**Ubuntu/Debian:**
```bash
sudo apt update && sudo apt install g++ golang rustc cargo python3 make
```

**macOS (Homebrew):**
```bash
brew install gcc go rust python3 make
```

## 快速开始 | Quick Start

```bash
# 克隆仓库
git clone https://github.com/LessUp/compress-kit.git
cd compress-kit

# 构建所有实现
make build

# 运行测试套件
make test

# 运行跨语言一致性测试
make test-conformance
```

## 开发工作流 | Development Workflow

### 1. 创建功能分支

```bash
git checkout -b feature/your-feature-name
```

### 2. 开发与测试

```bash
# 开发过程中持续运行测试
make test

# 代码格式检查
make lint

# 跨语言验证
make test-conformance
```

### 3. 提交变更

使用规范的 commit 信息格式：

- `feat:` 新功能
- `fix:` Bug 修复
- `docs:` 文档更新
- `refactor:` 代码重构
- `test:` 测试相关
- `chore:` 维护性工作

```bash
git commit -m "feat: add XXX feature"
```

### 4. 创建 Pull Request

推送分支后，在 GitHub 上创建 PR。确保：

- ✅ 所有测试通过
- ✅ 代码格式符合规范
- ✅ 跨语言二进制兼容性未被破坏
- ✅ 相关文档已更新

## 代码规范 | Code Style

### C++ (C++17)

- 遵循 [Google C++ Style Guide](https://google.github.io/styleguide/cppguide.html)
- 使用 4 空格缩进
- 函数和变量使用 `snake_case`
- 类和结构体使用 `PascalCase`

```bash
# 格式检查
clang-format --dry-run --Werror algorithms/**/*.cpp
```

### Go

- 使用 `gofmt` 格式化（强制）
- 使用 `go vet` 静态分析
- 遵循 [Effective Go](https://go.dev/doc/effective_go)

```bash
# 格式化
gofmt -w algorithms/

# 静态检查
go vet ./algorithms/...
```

### Rust

- 使用 `rustfmt` 格式化（强制）
- 使用 `clippy` 进行 lint
- 遵循 [Rust API Guidelines](https://rust-lang.github.io/api-guidelines/)

```bash
# 格式化
cargo fmt --all

# Lint 检查
cargo clippy --all-targets -- -D warnings
```

## 核心约束 | Core Constraints

### 二进制格式兼容性

**最重要：** 维护跨语言二进制格式兼容性是项目的核心约束。

- Magic Numbers 不可更改
- Frequency Table 格式必须统一
- 编码在一种语言，解码在另一种语言必须正常工作

### 安全限制

- 输入大小上限：4 GiB
- 解码输出上限：1 GiB
- 这些限制用于防止解压缩炸弹攻击

### 已知限制

Range Coder 在大文件（>500KB）上性能下降，这是已知问题，不建议修复。

## 规范驱动开发 | Spec-Driven Development

本项目遵循 **OpenSpec 工作流**：

1. **变更前先提案**：对于二进制格式、公共 API、跨语言语义等变更，需要先创建 OpenSpec change
2. **审查规范**：在 `/openspec/specs/` 中查看相关规范文档
3. **按规范实现**：代码实现必须遵守规范定义
4. **验证合规**：测试必须验证规范合规性

小型的文档修复、内部重构、保持现有契约的 bug 修复可以直接实施。

## 项目结构 | Project Structure

```
algorithms/           # 四种算法的 C++/Go/Rust 实现
├── huffman/
├── arithmetic/
├── range/
└── rle/
tests/                # 测试套件和跨语言验证
docs/                 # VitePress 文档站点
openspec/             # OpenSpec 规范和变更提案
```

## 获取帮助 | Get Help

- 查看 [文档站点](https://lessup.github.io/compress-kit/)
- 阅读 [架构指南](CONTEXT.md)
- 提交 [Issue](https://github.com/LessUp/compress-kit/issues)

---

**许可**：MIT License · 版权所有 © 2025-2026 LessUp
