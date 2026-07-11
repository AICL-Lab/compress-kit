.PHONY: build test test-data test-conformance test-cli-smoke clean format lint help

# ── Build ──────────────────────────────────────────────────────────────────

build:
	cmake -S . -B build
	cmake --build build

# ── Test ───────────────────────────────────────────────────────────────────

test: test-data build
	ctest --test-dir build --output-on-failure
	$(MAKE) test-cli-smoke

test-conformance: test-data build
	python3 tests/conformance/run_cli_smoke.py

test-cli-smoke: test-data build
	python3 tests/conformance/run_cli_smoke.py

# ── Data / Clean ───────────────────────────────────────────────────────────

test-data:
	python3 tests/gen_testdata.py

clean:
	rm -rf build tests/data

# ── Format & Lint ───────────────────────────────────────────────────────────

format:
	@command -v clang-format >/dev/null || { echo "clang-format is required"; exit 1; }
	@echo "Formatting C++ code..."
	@find algorithms -type f \( -name '*.cpp' -o -name '*.hpp' \) -exec clang-format -i {} +
	@echo "Done!"

lint:
	@command -v clang-format >/dev/null || { echo "clang-format is required"; exit 1; }
	@echo "Linting C++ code (clang-format dry-run)..."
	@find algorithms -type f \( -name '*.cpp' -o -name '*.hpp' \) \
		! -path '*/build/*' \
		-exec clang-format --dry-run -Werror {} + || { echo "clang-format check failed; run 'make format'"; exit 1; }
	@echo "Done!"

# ── Help ────────────────────────────────────────────────────────────────────

help:
	@echo "Build Commands:"
	@echo "  make build          Build all algorithms (CMake)"
	@echo ""
	@echo "Test Commands:"
	@echo "  make test           Run all tests + CLI smoke"
	@echo "  make test-cli-smoke Run CLI smoke tests"
	@echo ""
	@echo "Code Quality:"
	@echo "  make format         Format all C++ code"
	@echo "  make lint           Lint all C++ code (clang-format dry-run)"
	@echo ""
	@echo "Other:"
	@echo "  make clean          Clean build artifacts"
	@echo "  make test-data      Generate test corpus"
