.PHONY: build build-huffman build-arithmetic build-range build-rle \
       test test-conformance test-huffman-go test-arithmetic-go test-range-go test-rle-go \
       test-shared-cpp test-shared-go test-shared-rust \
       test-huffman-rust test-arithmetic-rust test-range-rust test-rle-rust \
       test-data bench clean format lint spec-init spec-list spec-status

# ── Build ──────────────────────────────────────────────────────────────────

build: build-huffman build-arithmetic build-range build-rle

build-huffman:
	g++ -std=c++17 -O2 -Wall -Wextra -Werror -Ialgorithms/shared/cpp/include algorithms/shared/cpp/src/buffer_api.cpp algorithms/huffman/cpp/main.cpp -o algorithms/huffman/cpp/huffman_cpp
	go build -o algorithms/huffman/go/huffman_go ./algorithms/huffman/go/cmd
	cargo build --manifest-path algorithms/huffman/rust/Cargo.toml --bin huffman_rust --release
	cp algorithms/huffman/rust/target/release/huffman_rust algorithms/huffman/rust/huffman_rust

build-arithmetic:
	g++ -std=c++17 -O2 -Wall -Wextra -Werror -Ialgorithms/shared/cpp/include algorithms/shared/cpp/src/buffer_api.cpp algorithms/arithmetic/cpp/main.cpp -o algorithms/arithmetic/cpp/arithmetic_cpp
	go build -o algorithms/arithmetic/go/arithmetic_go ./algorithms/arithmetic/go/cmd
	cargo build --manifest-path algorithms/arithmetic/rust/Cargo.toml --bin arithmetic_rust --release
	cp algorithms/arithmetic/rust/target/release/arithmetic_rust algorithms/arithmetic/rust/arithmetic_rust

build-range:
	g++ -std=c++17 -O2 -Wall -Wextra -Werror -Ialgorithms/shared/cpp/include algorithms/shared/cpp/src/buffer_api.cpp algorithms/range/cpp/main.cpp -o algorithms/range/cpp/rangecoder_cpp
	go build -o algorithms/range/go/rangecoder_go ./algorithms/range/go/cmd
	cargo build --manifest-path algorithms/range/rust/Cargo.toml --release

build-rle:
	g++ -std=c++17 -O2 -Wall -Wextra -Werror -Ialgorithms/shared/cpp/include algorithms/shared/cpp/src/buffer_api.cpp algorithms/rle/cpp/main.cpp -o algorithms/rle/cpp/rle_cpp
	go build -o algorithms/rle/go/rle_go ./algorithms/rle/go/cmd
	cargo build --manifest-path algorithms/rle/rust/Cargo.toml --bin rle_rust --release
	cp algorithms/rle/rust/target/release/rle_rust algorithms/rle/rust/rle_rust

# ── Test ───────────────────────────────────────────────────────────────────

test: test-data \
       test-shared-cpp test-shared-go test-shared-rust \
       test-huffman-go test-arithmetic-go test-range-go test-rle-go \
       test-huffman-rust test-arithmetic-rust test-range-rust test-rle-rust \
       test-conformance

test-shared-cpp:
	g++ -std=c++17 -O2 -Wall -Wextra -Werror -DCOMPRESSKIT_NO_MAIN -Ialgorithms/shared/cpp/include algorithms/shared/cpp/src/buffer_api.cpp algorithms/huffman/cpp/main.cpp algorithms/arithmetic/cpp/main.cpp algorithms/range/cpp/main.cpp algorithms/rle/cpp/main.cpp algorithms/shared/cpp/tests/test_lifecycle.cpp -o algorithms/shared/cpp/tests/test_lifecycle
	./algorithms/shared/cpp/tests/test_lifecycle

test-shared-go:
	go test ./algorithms/shared/go/...

test-shared-rust:
	cargo test --manifest-path algorithms/shared/rust/Cargo.toml

test-huffman-go:
	go test ./algorithms/huffman/go/... ./algorithms/huffman/go/cmd/...

test-arithmetic-go:
	go test ./algorithms/arithmetic/go/... ./algorithms/arithmetic/go/cmd/...

test-range-go:
	go test ./algorithms/range/go/...

test-rle-go:
	go test ./algorithms/rle/go/... ./algorithms/rle/go/cmd/...

test-huffman-rust:
	cargo test --manifest-path algorithms/huffman/rust/Cargo.toml

test-arithmetic-rust:
	cargo test --manifest-path algorithms/arithmetic/rust/Cargo.toml

test-range-rust:
	cargo test --manifest-path algorithms/range/rust/Cargo.toml

test-rle-rust:
	cargo test --manifest-path algorithms/rle/rust/Cargo.toml

test-conformance: build test-data
	python3 tests/conformance/run_decode_matrix.py

# ── Data / Bench / Clean ──────────────────────────────────────────────────

test-data:
	python3 tests/gen_testdata.py

bench: test-data
	python3 scripts/run_all_bench.py

clean:
	rm -rf reports tests/data
	rm -rf docs/.vitepress/dist
	rm -rf algorithms/huffman/benchmark/tmp algorithms/arithmetic/benchmark/tmp algorithms/range/benchmark/tmp algorithms/rle/benchmark/tmp
	rm -f algorithms/huffman/cpp/huffman_cpp algorithms/huffman/go/huffman_go algorithms/huffman/rust/huffman_rust algorithms/huffman/rust/huffman_rust_test
	rm -f algorithms/arithmetic/cpp/arithmetic_cpp algorithms/arithmetic/go/arithmetic_go algorithms/arithmetic/rust/arithmetic_rust algorithms/arithmetic/rust/arithmetic_rust_test
	rm -f algorithms/range/cpp/rangecoder_cpp algorithms/range/go/rangecoder_go
	rm -rf algorithms/huffman/rust/target algorithms/arithmetic/rust/target
	rm -rf algorithms/range/rust/target
	rm -f algorithms/rle/cpp/rle_cpp algorithms/rle/go/rle_go algorithms/rle/rust/rle_rust algorithms/rle/rust/rle_rust_test
	rm -f algorithms/shared/cpp/tests/test_lifecycle
	rm -rf algorithms/rle/rust/target algorithms/shared/rust/target

# ── Format & Lint ───────────────────────────────────────────────────────────

format:
	@echo "Formatting Go code..."
	gofmt -w algorithms/*/go
	@echo "Formatting Rust code..."
	cargo fmt --manifest-path algorithms/huffman/rust/Cargo.toml
	cargo fmt --manifest-path algorithms/arithmetic/rust/Cargo.toml
	cargo fmt --manifest-path algorithms/range/rust/Cargo.toml
	cargo fmt --manifest-path algorithms/rle/rust/Cargo.toml
	cargo fmt --manifest-path algorithms/shared/rust/Cargo.toml
	@command -v clang-format >/dev/null || { echo "clang-format is required for C++ formatting"; exit 1; }
	@echo "Formatting C++ code..."
	@for f in algorithms/*/cpp/main.cpp algorithms/shared/cpp/src/buffer_api.cpp algorithms/shared/cpp/tests/test_lifecycle.cpp; do \
		clang-format -i "$$f"; \
	done
	@echo "Done!"

lint:
	@echo "Linting Go code..."
	go vet ./algorithms/shared/go/... ./algorithms/huffman/go/... ./algorithms/arithmetic/go/... ./algorithms/range/go/... ./algorithms/rle/go/...
	@echo "Linting Rust code..."
	cargo clippy --manifest-path algorithms/huffman/rust/Cargo.toml --all-targets -- -D warnings
	cargo clippy --manifest-path algorithms/arithmetic/rust/Cargo.toml --all-targets -- -D warnings
	cargo clippy --manifest-path algorithms/range/rust/Cargo.toml --all-targets -- -D warnings
	cargo clippy --manifest-path algorithms/rle/rust/Cargo.toml --all-targets -- -D warnings
	cargo clippy --manifest-path algorithms/shared/rust/Cargo.toml --all-targets -- -D warnings
	@echo "Done!"

# ── OpenSpec ────────────────────────────────────────────────────────────────

spec-init:
	@openspec init --tools claude,cursor

spec-list:
	@openspec list

spec-status:
	@openspec status

# ── Help ────────────────────────────────────────────────────────────────────

help:
	@echo "Build Commands:"
	@echo "  make build          Build all algorithms"
	@echo "  make build-<algo>   Build specific algorithm"
	@echo ""
	@echo "Test Commands:"
	@echo "  make test           Run all tests"
	@echo "  make bench          Run benchmarks"
	@echo ""
	@echo "Code Quality:"
	@echo "  make format         Format all code"
	@echo "  make lint           Lint all code"
	@echo ""
	@echo "OpenSpec Commands:"
	@echo "  make spec-init      Initialize OpenSpec"
	@echo "  make spec-list      List active changes"
	@echo "  make spec-status    Show current status"
	@echo ""
	@echo "Other:"
	@echo "  make clean          Clean build artifacts"
