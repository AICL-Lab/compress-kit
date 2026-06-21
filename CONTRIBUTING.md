# Contributing to CompressKit | 贡献指南

感谢你对 CompressKit 项目的关注！本指南将帮助你快速参与项目开发。

## 前置要求 | Prerequisites

### 必需工具

- **C++ 编译器**: g++ 9+ 或 clang++ 10+ (支持 C++17)
- **CMake**: 3.16 或更高版本
- **Python**: 3.8 或更高版本（用于测试脚本）
- **clang-format**: 代码格式化
- **Make**: 简化构建命令

### 安装示例

**Ubuntu/Debian:**
```bash
sudo apt update && sudo apt install g++ cmake clang-format python3 make
```

**macOS (Homebrew):**
```bash
brew install cmake clang-format python3 make
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
- ✅ 二进制格式兼容性未被破坏
- ✅ 相关文档已更新

## 代码规范 | Code Style

### C++ (C++17)

- 遵循 [Google C++ Style Guide](https://google.github.io/styleguide/cppguide.html)
- 使用 clang-format 格式化
- 函数和变量使用 `snake_case`
- 类和结构体使用 `PascalCase`

```bash
# 格式化
make format

# 格式检查
make lint
```

## 核心约束 | Core Constraints

### 二进制格式兼容性

**最重要：** 维护二进制格式兼容性是项目的核心约束。

- Magic Numbers 不可更改
- Frequency Table 格式必须统一
- 编码/解码 round-trip 必须正确

### 安全限制

- 输入大小上限：4 GiB
- 解码输出上限：1 GiB
- 这些限制用于防止解压缩炸弹攻击

## 项目结构 | Project Structure

```
algorithms/           # 四种算法的 C++ 实现
├── huffman/
├── arithmetic/
├── range/
├── rle/
└── shared/           # 共享工具（频率表、Buffer 层、CLI）
tests/                # 测试套件和 CLI smoke 测试
docs/                 # VitePress 文档站点
```

## 获取帮助 | Get Help

- 查看 [文档站点](https://lessup.github.io/compress-kit/)
- 阅读 [架构指南](CONTEXT.md)
- 提交 [Issue](https://github.com/LessUp/compress-kit/issues)

---

**许可**：MIT License · 版权所有 © 2025-2026 LessUp
