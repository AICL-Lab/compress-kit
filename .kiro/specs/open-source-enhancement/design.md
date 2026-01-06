# Design Document: Open Source Enhancement

## Overview

本设计文档描述如何将编码算法学习项目完善为一个优秀的开源项目。设计遵循开源社区最佳实践，确保项目具备良好的可发现性、贡献者友好度和代码质量保障。

## Architecture

项目结构将保持现有的算法目录组织方式，新增以下顶层文件和目录：

```
encoding/
├── .github/
│   ├── workflows/
│   │   └── ci.yml                 # CI/CD 流水线配置
│   ├── ISSUE_TEMPLATE/
│   │   ├── bug_report.md          # Bug 报告模板
│   │   └── feature_request.md     # 功能请求模板
│   └── PULL_REQUEST_TEMPLATE.md   # PR 模板
├── LICENSE                        # 开源许可证
├── CONTRIBUTING.md                # 贡献指南
├── CODE_OF_CONDUCT.md             # 行为准则
├── CHANGELOG.md                   # 变更日志
├── SECURITY.md                    # 安全策略
├── README.md                      # 增强后的 README（双语）
├── README_EN.md                   # 英文 README（可选）
└── [existing directories...]      # 现有算法目录
```

## Components and Interfaces

### Component 1: LICENSE File

采用 MIT 许可证，这是最宽松且被广泛认可的开源许可证之一。

**接口**: 静态文本文件，无编程接口。

### Component 2: CONTRIBUTING.md

结构化的贡献指南，包含以下章节：

1. **Prerequisites** - 列出所需工具和版本
   - C++ 编译器 (g++ 9+ 或 clang++ 10+)
   - Go 1.19+
   - Rust 1.70+
   - Python 3.8+

2. **Development Setup** - 克隆仓库和环境配置步骤

3. **Code Style** - 各语言的代码风格要求
   - C++: 遵循 Google C++ Style Guide
   - Go: 使用 gofmt 和 go vet
   - Rust: 使用 rustfmt 和 clippy
   - Python: 遵循 PEP 8

4. **Testing** - 如何运行测试和基准测试

5. **Pull Request Process** - PR 提交流程和检查清单

### Component 3: CODE_OF_CONDUCT.md

采用 Contributor Covenant v2.1，这是最广泛使用的开源社区行为准则。

### Component 4: Issue Templates

**Bug Report Template** 包含字段：
- 环境信息（OS、编译器版本）
- 重现步骤
- 预期行为
- 实际行为
- 相关日志或截图

**Feature Request Template** 包含字段：
- 问题描述
- 建议的解决方案
- 替代方案
- 附加上下文

### Component 5: PR Template

包含检查清单：
- [ ] 代码已通过本地测试
- [ ] 代码风格符合项目规范
- [ ] 已更新相关文档
- [ ] 已添加必要的测试
- [ ] CHANGELOG 已更新（如适用）

### Component 6: CI/CD Pipeline (GitHub Actions)

**Workflow 触发条件**:
- push 到 main 分支
- pull_request 到 main 分支

**Jobs**:

1. **build-cpp**: 编译所有 C++ 实现
   - Matrix: ubuntu-latest, macos-latest
   - 编译 huffman/cpp, arithmetic/cpp, range/cpp, Run-Length/cpp

2. **build-go**: 构建和检查 Go 实现
   - 运行 go build, go vet, go fmt -d
   - 构建 huffman/go, range/go, Run-Length/go

3. **build-rust**: 构建和检查 Rust 实现
   - 运行 cargo build, cargo clippy, cargo fmt --check
   - 构建 huffman/rust, range/rust, Run-Length/rust

4. **test-correctness**: 验证编解码正确性
   - 生成测试数据
   - 对每个算法执行 encode → decode → diff 验证

### Component 7: Enhanced README

**结构**:

```markdown
# encoding 编码算法集合 | Encoding Algorithms Collection

[Badges: CI Status, License, Languages]

[中文描述]

[English Description]

## 🎯 Why This Project / 为什么做这个项目

## 📊 Algorithm Comparison / 算法对比

[表格: 算法名称 | 压缩率 | 速度 | 适用场景]

## 🚀 Quick Start / 快速开始

## 📖 Documentation / 文档

## 🤝 Contributing / 贡献

## 📄 License / 许可证
```

### Component 8: CHANGELOG.md

遵循 [Keep a Changelog](https://keepachangelog.com/) 格式：

```markdown
# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

### Added
- ...

### Changed
- ...

### Fixed
- ...
```

### Component 9: SECURITY.md

包含：
- 支持的版本
- 报告漏洞的方式（建议使用 GitHub Security Advisories）
- 预期响应时间

## Data Models

本项目主要涉及静态文档文件，无复杂数据模型。CI 配置使用 YAML 格式。

### CI Configuration Schema (ci.yml)

```yaml
name: CI
on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  build-cpp:
    runs-on: ${{ matrix.os }}
    strategy:
      matrix:
        os: [ubuntu-latest, macos-latest]
    steps:
      - uses: actions/checkout@v4
      - name: Build C++ implementations
        run: |
          # Build commands

  build-go:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.21'
      - name: Build and check Go
        run: |
          # Go commands

  build-rust:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: dtolnay/rust-toolchain@stable
        with:
          components: clippy, rustfmt
      - name: Build and check Rust
        run: |
          # Rust commands

  test-correctness:
    needs: [build-cpp, build-go, build-rust]
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Verify encode/decode correctness
        run: |
          # Correctness verification
```

## Correctness Properties

*A property is a characteristic or behavior that should hold true across all valid executions of a system—essentially, a formal statement about what the system should do. Properties serve as the bridge between human-readable specifications and machine-verifiable correctness guarantees.*

本项目主要涉及文档和配置文件的创建，属于静态内容而非算法实现。因此，大多数验收标准是文件存在性检查或内容验证，适合通过示例测试而非属性测试来验证。

**可测试的示例验证**（通过 CI 脚本验证）：

1. 文件存在性检查：LICENSE, CONTRIBUTING.md, CODE_OF_CONDUCT.md, CHANGELOG.md, SECURITY.md
2. 目录结构检查：.github/ISSUE_TEMPLATE/, .github/workflows/
3. 内容关键词检查：各文档包含必要的章节和关键信息

**不适用属性测试的原因**：
- 文档内容是静态的，不涉及输入/输出转换
- 没有需要验证的算法逻辑
- 验收标准都是具体的存在性或内容检查

## Error Handling

### CI Pipeline Errors

1. **编译失败**: CI 将报告具体的编译错误信息，包括文件名和行号
2. **代码风格检查失败**: 报告不符合规范的文件列表和具体问题
3. **正确性验证失败**: 报告哪个算法的 encode/decode 结果不一致

### 文档验证错误

如果必要文件缺失或内容不完整，CI 可以通过简单的 shell 脚本检查并报告。

## Testing Strategy

### 验证方法

由于本项目主要是文档和配置文件，采用以下验证策略：

1. **CI 自动验证**
   - 文件存在性检查（通过 shell 脚本）
   - 代码编译验证（各语言编译器）
   - 代码风格检查（gofmt, rustfmt, clang-format）
   - 编解码正确性验证（encode → decode → diff）

2. **人工审查**
   - 文档内容的准确性和完整性
   - 链接的有效性
   - 格式的一致性

### CI 验证脚本示例

```bash
#!/bin/bash
# 检查必要文件是否存在
required_files=(
  "LICENSE"
  "CONTRIBUTING.md"
  "CODE_OF_CONDUCT.md"
  "CHANGELOG.md"
  "SECURITY.md"
  ".github/workflows/ci.yml"
  ".github/ISSUE_TEMPLATE/bug_report.md"
  ".github/ISSUE_TEMPLATE/feature_request.md"
  ".github/PULL_REQUEST_TEMPLATE.md"
)

for file in "${required_files[@]}"; do
  if [ ! -f "$file" ]; then
    echo "ERROR: Missing required file: $file"
    exit 1
  fi
done

echo "All required files present."
```

### 编解码正确性验证

对于现有的编码算法，CI 将执行以下验证流程：

```bash
# 对每个算法和语言实现
for algo in huffman rle; do
  for lang in cpp go rust; do
    # 1. 编码测试文件
    ./${algo}_${lang} encode test_input.bin encoded.bin
    # 2. 解码
    ./${algo}_${lang} decode encoded.bin decoded.bin
    # 3. 验证一致性
    diff test_input.bin decoded.bin || exit 1
  done
done
```

