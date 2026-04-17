# RFC-0001: Core Architecture

## Status
**Status**: Accepted  
**Created**: 2024  
**Updated**: 2026  

## Architecture Overview

```
encoding/
├── .github/
│   ├── workflows/
│   │   ├── ci.yml              # Build, test, correctness verification
│   │   └── pages.yml           # VitePress docs deployment
│   ├── ISSUE_TEMPLATE/
│   │   ├── bug_report.md
│   │   └── feature_request.md
│   └── PULL_REQUEST_TEMPLATE.md
├── specs/                      # Spec-Driven Development documents
│   ├── product/                # Product requirements
│   ├── rfc/                    # Technical design documents
│   ├── api/                    # API definitions
│   ├── db/                     # Database schemas (if applicable)
│   └── testing/                # Test specifications
├── docs/
│   ├── .vitepress/config.mts   # VitePress configuration
│   ├── index.md                # Documentation landing page
│   ├── guide/
│   │   ├── getting-started.md
│   │   ├── algorithms.md
│   │   └── project-structure.md
│   └── public/
├── algorithms/huffman/
│   ├── cpp/main.cpp
│   ├── go/main.go
│   └── rust/main.rs
├── algorithms/arithmetic/
│   ├── cpp/main.cpp
│   ├── go/main.go
│   └── rust/main.rs
├── algorithms/range/
│   ├── cpp/main.cpp
│   ├── go/ (library + cmd/)
│   └── rust/ (Cargo.toml + src/)
├── algorithms/rle/
│   ├── cpp/main.cpp
│   ├── go/main.go
│   └── rust/main.rs
├── tests/
│   ├── gen_testdata.py
│   └── data/
├── Makefile
├── package.json
├── LICENSE
├── README.md
├── README.zh-CN.md
├── CHANGELOG.md
├── CONTRIBUTING.md
├── CODE_OF_CONDUCT.md
└── SECURITY.md
```

## Component Design

### Algorithm Modules

Each algorithm follows a consistent CLI pattern across all languages:

#### C++ Pattern
```cpp
int main(int argc, char** argv) {
    if (argc != 4) { /* usage */ }
    string mode = argv[1];  // "encode" or "decode"
    string input = argv[2];
    string output = argv[3];
    // Process...
}
```

#### Go Pattern
```go
func main() {
    if len(os.Args) != 4 { /* usage */ }
    mode := os.Args[1]
    inputPath := os.Args[2]
    outputPath := os.Args[3]
    // Process...
}
```

#### Rust Pattern
```rust
fn main() {
    let args: Vec<String> = env::args().collect();
    if args.len() != 4 { /* usage */ }
    let mode = &args[1];
    let input = &args[2];
    let output = &args[3];
    // Process...
}
```

## Frequency Table Format

All static-model algorithms share the same frequency table structure:

```
+------------------+------------------------+
| Field            | Format                 |
+------------------+------------------------+
| Symbol count     | 4 bytes LE (uint32)    |
| Frequency[0]     | 4 bytes LE (uint32)    |
| ...              | ...                    |
| Frequency[256]   | 4 bytes LE (uint32)    |
| Frequency[EOF]   | 4 bytes LE (uint32)    |
+------------------+------------------------+
```

## CI/CD Workflow Design

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│  build-cpp  │     │  build-go   │     │ build-rust  │
│  (matrix)   │     │             │     │             │
└──────┬──────┘     └──────┬──────┘     └──────┬──────┘
       │                   │                   │
       └───────────────────┼───────────────────┘
                           │
                    ┌──────▼──────┐
                    │ correctness │
                    │   tests     │
                    └─────────────┘
```

## Error Handling Strategy

1. **Input Validation**: Check file size before reading
2. **Decompression Limits**: Track output size to prevent bombs
3. **Error Messages**: English, descriptive, actionable
4. **Exit Codes**: 0 for success, 1 for errors

## Performance Considerations

1. **Buffered I/O**: Use bufio in Go, ifstream/ofstream in C++
2. **Frequency Scaling**: Scale to maxTotal (2^24) for numerical stability
3. **Memory**: Process streams where possible, avoid loading entire files

## Documentation Site Design

- **Landing Page**: Project overview, target audience, reading paths
- **Getting Started**: Environment requirements, build commands, CLI usage
- **Algorithms**: Theory, complexity analysis, implementation differences
- **Project Structure**: Directory layout, CLI conventions, file formats

## Decisions

### Why Three Languages?
- **C++17**: Industry standard, manual memory management learning
- **Go**: Modern systems language, excellent concurrency
- **Rust**: Memory safety without garbage collector

### Why This Directory Structure?
- Algorithm-first organization for easy navigation
- Language subdirectories within each algorithm
- Shared test data and CI configuration at root level
