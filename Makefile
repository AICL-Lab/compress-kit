.PHONY: build build-huffman build-arithmetic build-range build-rle \
       test test-huffman-go test-arithmetic-go test-range-go test-rle-go \
       test-huffman-rust test-arithmetic-rust test-range-rust test-rle-rust \
       test-data bench clean

# ── Build ──────────────────────────────────────────────────────────────────

build: build-huffman build-arithmetic build-range build-rle

build-huffman:
	g++ -std=c++17 -O2 -Wall -Wextra -Werror huffman/cpp/main.cpp -o huffman/cpp/huffman_cpp
	go build -o huffman/go/huffman_go ./huffman/go
	rustc -O huffman/rust/main.rs -o huffman/rust/huffman_rust

build-arithmetic:
	g++ -std=c++17 -O2 -Wall -Wextra -Werror arithmetic/cpp/main.cpp -o arithmetic/cpp/arithmetic_cpp
	go build -o arithmetic/go/arithmetic_go ./arithmetic/go
	rustc -O arithmetic/rust/main.rs -o arithmetic/rust/arithmetic_rust

build-range:
	g++ -std=c++17 -O2 -Wall -Wextra -Werror range/cpp/main.cpp -o range/cpp/rangecoder_cpp
	go build -o range/go/rangecoder_go ./range/go/cmd
	cargo build --manifest-path range/rust/Cargo.toml --release

build-rle:
	g++ -std=c++17 -O2 -Wall -Wextra -Werror rle/cpp/main.cpp -o rle/cpp/rle_cpp
	go build -o rle/go/rle_go ./rle/go
	rustc -O rle/rust/main.rs -o rle/rust/rle_rust

# ── Test ───────────────────────────────────────────────────────────────────

test: test-data \
      test-huffman-go test-arithmetic-go test-range-go test-rle-go \
      test-huffman-rust test-arithmetic-rust test-range-rust test-rle-rust

test-huffman-go:
	go test ./huffman/go/...

test-arithmetic-go:
	go test ./arithmetic/go/...

test-range-go:
	go test ./range/go/...

test-rle-go:
	go test ./rle/go/...

test-huffman-rust:
	rustc --test huffman/rust/main.rs -o huffman/rust/huffman_rust_test
	./huffman/rust/huffman_rust_test

test-arithmetic-rust:
	rustc --test arithmetic/rust/main.rs -o arithmetic/rust/arithmetic_rust_test
	./arithmetic/rust/arithmetic_rust_test

test-range-rust:
	cargo test --manifest-path range/rust/Cargo.toml

test-rle-rust:
	rustc --test rle/rust/main.rs -o rle/rust/rle_rust_test
	./rle/rust/rle_rust_test

# ── Data / Bench / Clean ──────────────────────────────────────────────────

test-data:
	python3 tests/gen_testdata.py

bench: test-data
	python3 scripts/run_all_bench.py

clean:
	rm -rf reports tests/data
	rm -rf huffman/benchmark/tmp arithmetic/benchmark/tmp range/benchmark/tmp rle/benchmark/tmp
	rm -f huffman/cpp/huffman_cpp huffman/go/huffman_go huffman/rust/huffman_rust huffman/rust/huffman_rust_test
	rm -f arithmetic/cpp/arithmetic_cpp arithmetic/go/arithmetic_go arithmetic/rust/arithmetic_rust arithmetic/rust/arithmetic_rust_test
	rm -f range/cpp/rangecoder_cpp range/go/rangecoder_go
	rm -f rle/cpp/rle_cpp rle/go/rle_go rle/rust/rle_rust rle/rust/rle_rust_test
	cargo clean --manifest-path range/rust/Cargo.toml 2>/dev/null || true
