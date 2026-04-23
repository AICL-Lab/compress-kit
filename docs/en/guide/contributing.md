# Contributing Guide

Thank you for your interest in contributing to **Encoding**! This project follows the **Spec-Driven Development (SDD)** paradigm. All contributions must be based on specification documents.

## How to Contribute

### 1. Read the Specs First

Before writing any code, read the relevant documents in `/specs/`:
- `/specs/product/` — Product requirements
- `/specs/rfc/` — Technical design documents
- `/specs/testing/` — Test specifications

If your request conflicts with existing specs, **update the spec first**.

### 2. Implementation Standards

Each language has specific requirements:

| Language | Build | Test | Format |
|----------|-------|------|--------|
| **C++17** | `g++ -std=c++17 -O2 -Wall -Wextra` | Add `#ifdef TEST` or separate test file | `clang-format` |
| **Go 1.21+** | Go modules (`go.mod`) | `go test ./...` | `gofmt` |
| **Rust 1.70+** | `rustc` or `cargo` | `cargo test` or `rustc --test` | `rustfmt` + `clippy` |

### 3. Submit a Pull Request

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/my-feature`
3. Make changes following the standards above
4. Ensure all tests pass: `make test`
5. Ensure the build works: `make build`
6. Push and open a PR against `master`

### 4. PR Checklist

- [ ] Code follows language-specific conventions
- [ ] Unit tests added (or CI shell tests updated)
- [ ] Cross-language encode/decode verified
- [ ] Documentation updated (if behavior changed)
- [ ] Specs updated (if interface/behavior changed)

## Adding a New Algorithm

1. **Create a spec** in `/specs/rfc/` with:
   - Algorithm description
   - File format specification (magic bytes, field layout)
   - Acceptance criteria

2. **Create directory structure**:
   ```
   algorithms/<name>/
   ├── cpp/main.cpp
   ├── go/go.mod, main.go (or library + cmd/)
   ├── rust/main.rs (or Cargo.toml + src/)
   └── benchmark/bench.py
   ```

3. **Implement** in all three languages

4. **Add tests**:
   - Go: `*_test.go`
   - Rust: `#[cfg(test)]` module
   - C++: CI shell test in `ci.yml`

5. **Update**:
   - `Makefile` — add build targets
   - `.github/workflows/ci.yml` — add build/test jobs
   - `docs/en/guide/algorithms.md` — algorithm documentation
   - `docs/zh/guide/algorithms.md` — Chinese translation

## Adding a New Language

If you want to add support for another language (e.g., Python, Zig):

1. Discuss in an issue first
2. Create an RFC in `/specs/rfc/`
3. Implement all algorithms in the new language
4. Add to CI workflow

## Code of Conduct

Please read our [Code of Conduct](https://github.com/LessUp/compresskit/blob/master/CODE_OF_CONDUCT.md).

## License

By contributing, you agree that your contributions will be licensed under the [MIT License](https://github.com/LessUp/compresskit/blob/master/LICENSE).
