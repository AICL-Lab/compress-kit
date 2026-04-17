# Encoding

<p align="center">
  <img src="docs/public/logo.svg" width="120" alt="Encoding Logo">
</p>

<p align="center">
  <strong>Classic compression algorithms implemented in C++17, Go, and Rust</strong>
</p>

<p align="center">
  <a href="https://github.com/LessUp/encoding/actions/workflows/ci.yml">
    <img src="https://github.com/LessUp/encoding/actions/workflows/ci.yml/badge.svg" alt="CI">
  </a>
  <a href="https://github.com/LessUp/encoding/actions/workflows/pages.yml">
    <img src="https://github.com/LessUp/encoding/actions/workflows/pages.yml/badge.svg" alt="Docs">
  </a>
  <a href="https://lessup.github.io/encoding/">
    <img src="https://img.shields.io/badge/Docs-Online-blue" alt="Docs">
  </a>
  <a href="https://opensource.org/licenses/MIT">
    <img src="https://img.shields.io/badge/License-MIT-yellow.svg" alt="License: MIT">
  </a>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/C++-17-blue.svg" alt="C++17">
  <img src="https://img.shields.io/badge/Go-1.21+-00ADD8.svg" alt="Go 1.21+">
  <img src="https://img.shields.io/badge/Rust-1.70+-orange.svg" alt="Rust 1.70+">
  <img src="https://img.shields.io/badge/Python-3.8+-3776AB.svg" alt="Python 3.8+">
</p>

<p align="center">
  <b>English</b> | <a href="README.zh-CN.md">简体中文</a> | <a href="https://lessup.github.io/encoding/">Documentation</a>
</p>

---

## 🧭 Which Algorithm Should I Use?

```
Is your data highly repetitive?
├── Yes → Use RLE (fastest, best for repeated patterns)
└── No →
    Do you need maximum compression?
    ├── Yes → Use Arithmetic Coding (closest to entropy limit)
    └── No →
        Is speed critical?
        ├── Yes → Use Range Coder (fast + good compression)
        └── No → Use Huffman (simple, general purpose)
```

## 📊 Algorithm Comparison

| Algorithm | Compression | Speed | Best For | Use When |
|-----------|-------------|-------|----------|----------|
| **Huffman** | Medium | Fast | General text/data | You want simple, reliable compression |
| **Arithmetic** | ★ Highest | Medium | Maximum compression | Every byte matters |
| **Range Coder** | ★ High | Fast | Balanced performance | Best speed/compression tradeoff |
| **RLE** | Variable | ★ Fastest | Repeated data (bitmaps, logs) | Data has long runs of identical bytes |

## 🚀 Quick Start

```bash
git clone https://github.com/LessUp/encoding.git
cd encoding
make build && make test
```

### Cross-Language Verification

```bash
# Encode with C++
./algorithms/huffman/cpp/huffman_cpp encode input.txt output.huf

# Decode with Go — any combination works!
./algorithms/huffman/go/huffman_go decode output.huf restored.txt
diff input.txt restored.txt  # ✓ No output = identical
```

**C++ ↔ Go ↔ Rust** — all implementations share identical binary formats.

## 🏗️ Project Structure

```
encoding/
├── algorithms/
│   ├── huffman/          # Prefix-code compression
│   ├── arithmetic/       # Arithmetic coding
│   ├── range/            # Range coder (byte-level arithmetic)
│   └── rle/              # Run-length encoding
│       ├── cpp/          #   C++17: single file, zero deps
│       ├── go/           #   Go 1.21+: library + cmd/ CLI
│       ├── rust/         #   Rust 1.70+: rustc or cargo
│       └── benchmark/    #   Performance measurement scripts
├── docs/                 # VitePress site (en + zh)
├── specs/                # Spec-driven development (RFCs, product specs)
├── tests/                # Test data generation
└── Makefile              # Build, test, benchmark entry point
```

Each algorithm has **3 language implementations** with identical file formats:

| Language | Build | Structure |
|----------|-------|-----------|
| **C++17** | `g++ -std=c++17 -O2` | Single file, zero dependencies |
| **Go** | `go build ./cmd` | Library API (`package <algo>`) + CLI (`cmd/main.go`) |
| **Rust** | `rustc -O` or `cargo` | `main.rs` with reusable functions |

## 📖 Documentation

| Resource | Link |
|----------|------|
| 📚 **Full Documentation** | [lessup.github.io/encoding](https://lessup.github.io/encoding/) |
| 🔧 **API Reference** (Go / Rust / C++) | [API Docs](https://lessup.github.io/encoding/en/api/go) |
| 📈 **Benchmark Results** | [Performance](https://lessup.github.io/encoding/en/benchmarks/results) |
| 🤝 **Contributing Guide** | [How to Contribute](https://lessup.github.io/encoding/en/guide/contributing) |
| 📋 **Technical Specs** (RFCs) | [specs/](specs/) |
| 📝 **Changelog** | [CHANGELOG.md](CHANGELOG.md) |

## 💻 Build & Test Commands

| Command | Description |
|---------|-------------|
| `make build` | Build all implementations (C++, Go, Rust) |
| `make build-huffman` | Build Huffman only |
| `make test` | Run Go + Rust unit tests |
| `make bench` | Run performance benchmarks |
| `make clean` | Remove build artifacts |

## 🔬 Go Library Usage

All Go implementations expose a reusable library API:

```go
import "huffman" // or "arithmetic", "rle"

// Encode a file
err := huffman.EncodeFile("input.bin", "output.huf")

// Decode a file
err := huffman.DecodeFile("output.huf", "decoded.bin")
```

## 🏅 Why This Project Exists

- **🎓 Learn** — Read clean, well-documented implementations side by side
- **🔬 Compare** — See how C++, Go, and Rust handle the same algorithm differently
- **✅ Verify** — Cross-language tests guarantee identical output formats
- **📐 SDD** — Built with Spec-Driven Development: every feature starts with specs

## 🤝 Contributing

We welcome contributions! This project follows **Spec-Driven Development (SDD)**:

1. **Read specs first** — `/specs/` is the single source of truth
2. **Update specs before code** — if interfaces change, specs change first
3. **Test across languages** — verify C++ ↔ Go ↔ Rust compatibility

See the [Contributing Guide](https://lessup.github.io/encoding/en/guide/contributing) for details.

## 📄 License

[MIT License](LICENSE) · Copyright © 2025-2026 LessUp
