# encoding 编码算法集合 | Encoding Algorithms Collection

[![CI](https://github.com/LessUp/encoding/actions/workflows/ci.yml/badge.svg)](https://github.com/LessUp/encoding/actions/workflows/ci.yml)
[![Deploy Docs](https://github.com/LessUp/encoding/actions/workflows/docs.yml/badge.svg)](https://github.com/LessUp/encoding/actions/workflows/docs.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

[English](README.md) | 简体中文
![C++](https://img.shields.io/badge/C++-17-blue.svg)
![Go](https://img.shields.io/badge/Go-1.21+-00ADD8.svg)
![Rust](https://img.shields.io/badge/Rust-1.70+-orange.svg)

> 📖 **文档站 | Docs**：[https://lessup.github.io/encoding/](https://lessup.github.io/encoding/)

> 🎓 一个用多种语言实现经典压缩编码算法的学习与对比项目
>
> 🎓 A multi-language implementation of classic compression encoding algorithms for learning and comparison

---

## 🎯 Why This Project | 为什么做这个项目

**中文**：
- 📚 **学习导向**：通过阅读和对比不同语言的实现，深入理解压缩算法原理
- 🔬 **性能对比**：在相同算法下对比 C++、Go、Rust 的性能差异
- 🔄 **跨语言兼容**：所有实现使用相同的文件格式，支持交叉验证
- 🛠️ **实践友好**：简单一致的 CLI 接口，便于快速上手和实验

**English**:
- 📚 **Learning-oriented**: Understand compression algorithms by reading and comparing implementations across languages
- 🔬 **Performance comparison**: Compare C++, Go, and Rust performance for the same algorithms
- 🔄 **Cross-language compatible**: All implementations use the same file format for cross-validation
- 🛠️ **Practice-friendly**: Simple and consistent CLI interface for quick experimentation

---

## 📊 Algorithm Comparison | 算法对比

| Algorithm | Best For | Compression Ratio | Speed | Languages |
|-----------|----------|-------------------|-------|-----------|
| **Huffman** | General purpose, text | Medium | Fast | C++, Go, Rust |
| **Arithmetic** | Maximum compression | High | Medium | C++, Go, Rust |
| **Range Coder** | Balanced performance | High | Fast | C++, Go, Rust |
| **RLE** | Repetitive data | Variable* | Very Fast | C++, Go, Rust |

\* RLE compression ratio depends heavily on input data characteristics

| 算法 | 适用场景 | 压缩率 | 速度 | 支持语言 |
|------|----------|--------|------|----------|
| **Huffman** | 通用、文本 | 中等 | 快 | C++, Go, Rust |
| **Arithmetic** | 追求最大压缩 | 高 | 中等 | C++, Go, Rust |
| **Range Coder** | 平衡性能 | 高 | 快 | C++, Go, Rust |
| **RLE** | 重复数据 | 可变* | 非常快 | C++, Go, Rust |

\* RLE 压缩率高度依赖输入数据特征

---

## 🚀 Quick Start | 快速开始

### Prerequisites | 前置要求

- C++ compiler (g++ 9+ or clang++ 10+)
- Go 1.21+
- Rust 1.70+
- Python 3.8+ (for benchmarks)

### Build & Run | 构建与运行

```bash
# Clone the repository | 克隆仓库
git clone https://github.com/LessUp/encoding.git
cd encoding

# Example: Huffman encoding | 示例：Huffman 编码
cd huffman/cpp
g++ -std=c++17 -O2 main.cpp -o huffman_cpp

# Encode | 编码
./huffman_cpp encode input.bin output.huf

# Decode | 解码
./huffman_cpp decode output.huf restored.bin
```

### Run Benchmarks | 运行基准测试

```bash
# Generate test data and run all benchmarks
# 生成测试数据并运行所有基准测试
python3 scripts/run_all_bench.py
```

---

## 📁 Project Structure | 项目结构

```
encoding/
├── huffman/           # Huffman encoding | Huffman 编码
│   ├── cpp/          # C++ implementation
│   ├── go/           # Go implementation
│   ├── rust/         # Rust implementation
│   └── benchmark/    # Cross-language benchmark
├── arithmetic/        # Arithmetic coding | 算术编码
│   ├── cpp/          # C++ implementation
│   ├── go/           # Go implementation
│   ├── rust/         # Rust implementation
│   └── benchmark/    # Cross-language benchmark
├── range/            # Range coder | 区间编码
│   ├── cpp/          # C++ implementation
│   ├── go/           # Go implementation (library + CLI)
│   ├── rust/         # Rust implementation (library + CLI)
│   └── benchmark/    # Cross-language benchmark
├── rle/              # RLE encoding | 游程编码
│   ├── cpp/          # C++ implementation
│   ├── go/           # Go implementation
│   ├── rust/         # Rust implementation
│   └── benchmark/    # Cross-language benchmark
├── scripts/          # Utility scripts | 工具脚本
└── tests/            # Test data generation | 测试数据生成
```

---

## 📖 Algorithm Details | 算法详解

### Huffman Encoding | Huffman 编码

基于前缀码的无损压缩算法。实现中先扫描输入统计频率，构建 Huffman 树，再按位写入编码结果。

A lossless compression algorithm based on prefix codes. The implementation scans input to build frequency statistics, constructs a Huffman tree, and writes bit-encoded output.

**File Format | 文件格式**:
- Magic: `HFMN` (4 bytes)
- Frequency table (257 × 4 bytes, little-endian)
- Encoded data (bit stream)

### Arithmetic Coding | 算术编码

使用区间逐步细分表示整段消息的概率，压缩效率更接近信息熵上界。

Uses interval subdivision to represent message probability, achieving compression efficiency closer to the entropy limit.

### Range Coder | 区间编码

一种等价于算术编码的实现方式，但常在实践中更高效。

An implementation equivalent to arithmetic coding but often more efficient in practice.

**API (library)**:
- Go: `rangecoder.Encode(data []byte) ([]byte, error)`
- Rust: `rangecoder::encode(input: &[u8]) -> Result<Vec<u8>, RangeError>`

**CLI**:
```bash
# C++
./rangecoder_cpp encode input.bin output.rcnc
# Go
./rangecoder_go encode input.bin output.rcnc
# Rust
cargo run --bin rangecoder -- encode input.bin output.rcnc
```

### Run-Length Encoding (RLE) | 游程编码

适用于包含大量相同字节连续重复的数据。

Suitable for data with many consecutive repeated bytes.

**File Format | 文件格式**:
- Repeated `(count, value)` pairs
- `count`: 4 bytes, little-endian, unsigned integer
- `value`: 1 byte

---

## 🧪 Testing | 测试

### Correctness Verification | 正确性验证

All implementations pass encode-decode round-trip tests:

所有实现都通过编码-解码往返测试：

```bash
# Encode with C++, decode with Go (cross-language test)
./huffman_cpp encode input.bin encoded.huf
./huffman_go decode encoded.huf decoded.bin
diff input.bin decoded.bin  # Should produce no output
```

### Benchmark Results | 基准测试结果

Run `python scripts/run_all_bench.py` to generate benchmark reports in `reports/` directory.

运行 `python scripts/run_all_bench.py` 在 `reports/` 目录生成基准测试报告。

---

## 🤝 Contributing | 贡献

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

欢迎贡献！请查看 [CONTRIBUTING.md](CONTRIBUTING.md) 了解贡献指南。

### Quick Contribution Guide | 快速贡献指南

1. Fork the repository | Fork 仓库
2. Create a feature branch | 创建功能分支
3. Make your changes | 进行修改
4. Run tests | 运行测试
5. Submit a PR | 提交 PR

---

## 📜 Code of Conduct | 行为准则

This project follows the [Contributor Covenant Code of Conduct](CODE_OF_CONDUCT.md).

本项目遵循 [Contributor Covenant 行为准则](CODE_OF_CONDUCT.md)。

---

## 🔒 Security | 安全

For security issues, please see [SECURITY.md](SECURITY.md).

安全问题请查看 [SECURITY.md](SECURITY.md)。

---

## 📄 License | 许可证

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件。

---

## 🙏 Acknowledgments | 致谢

- Classic compression algorithm literature and textbooks
- The open source community for inspiration and best practices

---

## 📈 Roadmap | 路线图

- [x] All 4 algorithms implemented in C++, Go, Rust
- [x] Cross-language encode/decode compatibility for all algorithms
- [x] Unified CLI interface and cross-language benchmarks
- [ ] Add more compression algorithms (LZ77, LZ78, LZSS)
- [ ] Add Python implementations
- [ ] Add WebAssembly builds for browser demos
- [ ] Add interactive visualization of compression process
- [ ] Add more real-world test datasets

---

<p align="center">
  Made with ❤️ for learning compression algorithms
  <br>
  用 ❤️ 制作，为了学习压缩算法
</p>
