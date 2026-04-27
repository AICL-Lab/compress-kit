# 参与贡献指南

感谢您对 **Encoding** 的贡献兴趣！本项目严格遵循 **规范驱动开发 (SDD)** 范式。所有贡献必须基于规范文档。

## 如何贡献

### 1. 先阅读规范

在编写任何代码之前，请先阅读 `/specs/` 中的相关文档：
- `/specs/product/` — 产品需求
- `/specs/rfc/` — 技术设计文档
- `/specs/testing/` — 测试规范

如果您的需求与现有规范冲突，**请先更新规范**。

### 2. 实现标准

每种语言有特定要求：

| 语言 | 构建 | 测试 | 格式 |
|------|------|------|------|
| **C++17** | `g++ -std=c++17 -O2 -Wall -Wextra` | 添加 `#ifdef TEST` 或独立测试文件 | `clang-format` |
| **Go 1.21+** | Go 模块 (`go.mod`) | `go test ./...` | `gofmt` |
| **Rust 1.70+** | `rustc` 或 `cargo` | `cargo test` 或 `rustc --test` | `rustfmt` + `clippy` |

### 3. 提交 Pull Request

1. Fork 仓库
2. 创建特性分支：`git checkout -b feature/my-feature`
3. 按照上述标准进行更改
4. 确保所有测试通过：`make test`
5. 确保构建通过：`make build`
6. 推送并对 `master` 开启 PR

### 4. PR 检查清单

- [ ] 代码遵循语言约定
- [ ] 添加了单元测试（或更新了 CI 测试）
- [ ] 跨语言编码/解码已验证
- [ ] 文档已更新（如行为有变）
- [ ] 规范已更新（如接口/行为有变）

## 添加新算法

1. **在 `/specs/rfc/` 创建规范**，包含：
   - 算法描述
   - 文件格式规范（魔数、字段布局）
   - 验收标准

2. **创建目录结构**：
   ```
   algorithms/<name>/
   ├── cpp/main.cpp
   ├── go/go.mod, main.go（或 library + cmd/）
   ├── rust/main.rs（或 Cargo.toml + src/）
   └── benchmark/bench.py
   ```

3. **用三种语言实现**

4. **添加测试**：
   - Go: `*_test.go`
   - Rust: `#[cfg(test)]` 模块
   - C++: `ci.yml` 中的 CI shell 测试

5. **更新**：
   - `Makefile` — 添加构建目标
   - `.github/workflows/ci.yml` — 添加构建/测试任务
   - `docs/en/guide/algorithms.md` — 算法文档
   - `docs/zh/guide/algorithms.md` — 中文翻译

## 添加新语言

如果您想添加另一种语言（如 Python、Zig）：

1. 先在 issue 中讨论
2. 在 `/specs/rfc/` 创建 RFC
3. 用新语言实现所有算法
4. 添加到 CI 工作流

## 行为准则

请阅读我们的 [行为准则](https://github.com/LessUp/compress-kit/blob/master/CODE_OF_CONDUCT.md)。

## 许可证

通过贡献，您同意您的贡献将在 [MIT 许可证](https://github.com/LessUp/compress-kit/blob/master/LICENSE) 下授权。
