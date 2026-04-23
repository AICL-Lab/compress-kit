.PHONY: build build-huffman build-arithmetic build-range build-rle \
       test test-huffman-go test-arithmetic-go test-range-go test-rle-go \
       test-huffman-rust test-arithmetic-rust test-range-rust test-rle-rust \
       test-data bench clean spec-init spec-list spec-status

# ── Build ──────────────────────────────────────────────────────────────────

build: build-huffman build-arithmetic build-range build-rle

build-huffman:
	g++ -std=c++17 -O2 -Wall -Wextra -Werror algorithms/huffman/cpp/main.cpp -o algorithms/huffman/cpp/huffman_cpp
	go build -o algorithms/huffman/go/huffman_go ./algorithms/huffman/go/cmd
	rustc -O algorithms/huffman/rust/main.rs -o algorithms/huffman/rust/huffman_rust

build-arithmetic:
	g++ -std=c++17 -O2 -Wall -Wextra -Werror algorithms/arithmetic/cpp/main.cpp -o algorithms/arithmetic/cpp/arithmetic_cpp
	go build -o algorithms/arithmetic/go/arithmetic_go ./algorithms/arithmetic/go/cmd
	rustc -O algorithms/arithmetic/rust/main.rs -o algorithms/arithmetic/rust/arithmetic_rust

build-range:
	g++ -std=c++17 -O2 -Wall -Wextra -Werror algorithms/range/cpp/main.cpp -o algorithms/range/cpp/rangecoder_cpp
	go build -o algorithms/range/go/rangecoder_go ./algorithms/range/go/cmd
	cargo build --manifest-path algorithms/range/rust/Cargo.toml --release

build-rle:
	g++ -std=c++17 -O2 -Wall -Wextra -Werror algorithms/rle/cpp/main.cpp -o algorithms/rle/cpp/rle_cpp
	go build -o algorithms/rle/go/rle_go ./algorithms/rle/go/cmd
	rustc -O algorithms/rle/rust/main.rs -o algorithms/rle/rust/rle_rust

# ── Test ───────────────────────────────────────────────────────────────────

test: test-data \
      test-huffman-go test-arithmetic-go test-range-go test-rle-go \
      test-huffman-rust test-arithmetic-rust test-range-rust test-rle-rust

test-huffman-go:
	go test ./algorithms/huffman/go/... ./algorithms/huffman/go/cmd/...

test-arithmetic-go:
	go test ./algorithms/arithmetic/go/... ./algorithms/arithmetic/go/cmd/...

test-range-go:
	go test ./algorithms/range/go/...

test-rle-go:
	go test ./algorithms/rle/go/... ./algorithms/rle/go/cmd/...

test-huffman-rust:
	rustc --test algorithms/huffman/rust/main.rs -o algorithms/huffman/rust/huffman_rust_test
	./algorithms/huffman/rust/huffman_rust_test

test-arithmetic-rust:
	rustc --test algorithms/arithmetic/rust/main.rs -o algorithms/arithmetic/rust/arithmetic_rust_test
	./algorithms/arithmetic/rust/arithmetic_rust_test

test-range-rust:
	cargo test --manifest-path algorithms/range/rust/Cargo.toml

test-rle-rust:
	rustc --test algorithms/rle/rust/main.rs -o algorithms/rle/rust/rle_rust_test
	./algorithms/rle/rust/rle_rust_test

# ── Data / Bench / Clean ──────────────────────────────────────────────────

test-data:
	python3 tests/gen_testdata.py

bench: test-data
	python3 scripts/run_all_bench.py

clean:
	rm -rf reports tests/data
	rm -rf algorithms/huffman/benchmark/tmp algorithms/arithmetic/benchmark/tmp algorithms/range/benchmark/tmp algorithms/rle/benchmark/tmp
	rm -f algorithms/huffman/cpp/huffman_cpp algorithms/huffman/go/huffman_go algorithms/huffman/rust/huffman_rust algorithms/huffman/rust/huffman_rust_test
	rm -f algorithms/arithmetic/cpp/arithmetic_cpp algorithms/arithmetic/go/arithmetic_go algorithms/arithmetic/rust/arithmetic_rust algorithms/arithmetic/rust/arithmetic_rust_test
	rm -f algorithms/range/cpp/rangecoder_cpp algorithms/range/go/rangecoder_go
	rm -f algorithms/rle/cpp/rle_cpp algorithms/rle/go/rle_go algorithms/rle/rust/rle_rust algorithms/rle/rust/rle_rust_test
	cargo clean --manifest-path algorithms/range/rust/Cargo.toml 2>/dev/null || true

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
	@echo "OpenSpec Commands:"
	@echo "  make spec-init      Initialize OpenSpec"
	@echo "  make spec-list      List active changes"
	@echo "  make spec-status    Show current status"
	@echo ""
	@echo "Other:"
	@echo "  make clean          Clean build artifacts"
