# Shared Buffer Layer Deepening Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Deepen the shared Buffer Layer in Go, C++, and Rust so buffer growth, `BUF_TOO_SMALL` retry, and size-limit policy live behind one smaller seam per language without changing caller-facing behavior.

**Architecture:** Implement this in one reference-first slice. Start with Go as the reference implementation, then port the same module shape to C++ and Rust once the reference tests are green. Each language keeps the existing public Buffer Layer interface, but moves retry/growth policy into one internal helper so callers and adapters cross a smaller seam.

**Tech Stack:** Go, C++17, Rust, existing CompressKit shared codec modules, existing `make test` / `make lint` validation, Go `testing`, Rust `cargo test`, C++ shared test binary

---

## Scope note

This plan covers **slice 1 only** from the approved design:

- Go / C++ / Rust shared Buffer Layer orchestration

The later slices (Go writer seam, algorithm Streaming Layer wrappers, Go CLI launcher, C++ temp-file adapter) should each get their own follow-on plan after this slice lands cleanly.

## File structure

### Go

- Create: `algorithms/shared/go/codec/buffer_loop.go`
  - internal helper for retry/growth/limit orchestration
- Modify: `algorithms/shared/go/codec/buffer.go`
  - delegate `EncodeBuffer` / `DecodeBuffer` retry loops into the helper
- Modify: `algorithms/shared/go/codec/buffer_internal_test.go`
  - add focused regressions for partial-write retry preservation and cap handling

### C++

- Modify: `algorithms/shared/cpp/src/buffer_api.cpp`
  - extract internal retry/growth helper in the anonymous namespace and reuse it from encode/decode Buffer Layer entry points
- Modify: `algorithms/shared/cpp/tests/test_lifecycle.cpp`
  - add focused retry-preservation regression using fake `Encoder` / `Decoder` implementations

### Rust

- Modify: `algorithms/shared/rust/src/codec/buffer.rs`
  - extract internal retry/growth helper and reuse it from `encode_buffer` / `decode_buffer`
- Create: `algorithms/shared/rust/tests/buffer.rs`
  - add focused retry-preservation regression tests for the Buffer Layer seam

## Task 1: Go reference slice

**Files:**
- Create: `algorithms/shared/go/codec/buffer_loop.go`
- Modify: `algorithms/shared/go/codec/buffer.go`
- Modify: `algorithms/shared/go/codec/buffer_internal_test.go`
- Test: `algorithms/shared/go/codec/buffer_internal_test.go`

- [ ] **Step 1: Write the failing Go regression tests**

Add focused tests that describe the seam we want to preserve while the implementation moves inward:

```go
func TestEncodeBuffer_PreservesOutputAcrossFinishRetry(t *testing.T) {
	stub := &scriptedEncoder{
		process: []scriptedCall{{written: 0, err: nil}},
		finish: []scriptedCall{
			{written: 3, err: ErrBufTooSmall, payload: []byte("abc")},
			{written: 3, err: nil, payload: []byte("def")},
		},
	}

	out, err := EncodeBuffer(stub, []byte("ignored"))
	if err != nil {
		t.Fatalf("EncodeBuffer() error = %v", err)
	}
	if string(out) != "abcdef" {
		t.Fatalf("EncodeBuffer() = %q, want %q", out, "abcdef")
	}
}

func TestDecodeBuffer_ReturnsSizeLimitWhenGrowthStops(t *testing.T) {
	stub := &scriptedDecoder{
		process: []scriptedCall{{written: 0, err: ErrBufTooSmall}},
	}

	_, err := decodeBufferWithLimit(stub, []byte("ignored"), 1, 1)
	if err != ErrSizeLimit {
		t.Fatalf("decodeBufferWithLimit() error = %v, want %v", err, ErrSizeLimit)
	}
}
```

- [ ] **Step 2: Run the Go package tests to verify the new cases fail**

Run: `go test ./algorithms/shared/go/codec`

Expected: FAIL in the new retry helper tests because `decodeBufferWithLimit` and the shared retry helper do not exist yet.

- [ ] **Step 3: Write the minimal Go implementation**

Create `algorithms/shared/go/codec/buffer_loop.go` and move the repeated retry/growth policy there:

```go
package codec

type bufferStep func(out []byte) (int, error)

func runBufferStep(outBuf []byte, totalWritten int, limit int, step bufferStep) ([]byte, int, error) {
	for {
		n, err := step(outBuf[totalWritten:])
		if err != ErrBufTooSmall {
			if err != nil {
				return nil, totalWritten, err
			}
			return outBuf, totalWritten + n, nil
		}

		totalWritten += n
		if totalWritten > limit || len(outBuf) >= limit {
			return nil, totalWritten, ErrSizeLimit
		}

		newSize := growBuffer(len(outBuf), limit)
		if newSize <= len(outBuf) {
			return nil, totalWritten, ErrSizeLimit
		}

		newBuf := make([]byte, newSize)
		copy(newBuf, outBuf[:totalWritten])
		outBuf = newBuf
	}
}
```

Then refactor `buffer.go` so `EncodeBuffer` and `DecodeBuffer` call the helper instead of keeping two nearly identical retry loops:

```go
func encodeBufferWithLimit(encoder Encoder, input []byte, initialSize int, limit int) ([]byte, error) {
	outBuf := make([]byte, initialSize)
	totalWritten := 0

	var err error
	outBuf, totalWritten, err = runBufferStep(outBuf, totalWritten, limit, func(out []byte) (int, error) {
		return encoder.Process(input, out)
	})
	if err != nil {
		return nil, err
	}

	outBuf, totalWritten, err = runBufferStep(outBuf, totalWritten, limit, func(out []byte) (int, error) {
		return encoder.Finish(out)
	})
	if err != nil {
		return nil, err
	}

	return outBuf[:totalWritten], nil
}
```

- [ ] **Step 4: Run the Go package tests to verify the reference slice is green**

Run: `go test ./algorithms/shared/go/codec`

Expected: PASS

- [ ] **Step 5: Commit the Go reference slice**

```bash
git add algorithms/shared/go/codec/buffer.go algorithms/shared/go/codec/buffer_loop.go algorithms/shared/go/codec/buffer_internal_test.go
git commit -m "refactor(go): deepen shared buffer retry loop"
```

## Task 2: Port the deeper seam to C++

**Files:**
- Modify: `algorithms/shared/cpp/src/buffer_api.cpp`
- Modify: `algorithms/shared/cpp/tests/test_lifecycle.cpp`
- Test: `algorithms/shared/cpp/tests/test_lifecycle.cpp`

- [ ] **Step 1: Add failing C++ regressions for retry preservation**

Extend `test_lifecycle.cpp` with fake `Encoder` / `Decoder` implementations that return `BUF_TOO_SMALL` after writing a prefix:

```cpp
struct ScriptedEncoder final : compresskit::Encoder {
    int finish_calls = 0;

    compresskit::Result<std::size_t> process(compresskit::ByteView, compresskit::MutableByteView) override {
        return {compresskit::StatusCode::OK, 0};
    }

    compresskit::Result<std::size_t> flush(compresskit::MutableByteView) override {
        return {compresskit::StatusCode::OK, 0};
    }

    compresskit::Result<std::size_t> finish(compresskit::MutableByteView out) override {
        ++finish_calls;
        if (finish_calls == 1) {
            std::copy_n("abc", 3, out.data);
            return {compresskit::StatusCode::BUF_TOO_SMALL, 3};
        }
        std::copy_n("def", 3, out.data);
        return {compresskit::StatusCode::OK, 3};
    }

    void reset() noexcept override {}
    compresskit::State state() const noexcept override { return compresskit::State::STREAMING; }
};
```

Use it in a new assertion around the Buffer Layer helper seam:

```cpp
ScriptedEncoder encoder;
auto encoded = encode_buffer(encoder, std::vector<uint8_t>{'x'});
assert(encoded.status == compresskit::StatusCode::OK);
assert(std::string(encoded.value.begin(), encoded.value.end()) == "abcdef");
```

- [ ] **Step 2: Run the C++ shared tests to verify the regression fails**

Run: `make test`

Expected: FAIL in the new C++ retry-preservation assertion because the shared retry helper has not been extracted yet.

- [ ] **Step 3: Write the minimal C++ implementation**

Inside the anonymous namespace in `buffer_api.cpp`, extract the repeated resize/retry logic into one helper and reuse it from both `encode_buffer` and `decode_buffer`:

```cpp
template <typename Step>
Result<std::size_t> run_buffer_step(std::vector<uint8_t>& out,
                                    std::size_t total_written,
                                    std::size_t limit,
                                    Step step) {
    for (;;) {
        Result<std::size_t> result =
            step({out.data() + total_written, out.size() - total_written});
        if (result.status != StatusCode::BUF_TOO_SMALL) {
            if (!result.ok()) {
                return result;
            }
            return {StatusCode::OK, total_written + result.value};
        }

        total_written += result.value;
        if (total_written > limit || out.size() >= limit) {
            return {StatusCode::ERR_SIZE_LIMIT, 0};
        }

        out.resize(std::min(limit, std::max<std::size_t>(out.size() * 2, out.size() + 1)));
    }
}
```

Then replace the duplicated `for (;;)` loops in `encode_buffer` and `decode_buffer` with calls to `run_buffer_step`.

- [ ] **Step 4: Run the shared test baseline again**

Run: `make test`

Expected: PASS

- [ ] **Step 5: Commit the C++ port**

```bash
git add algorithms/shared/cpp/src/buffer_api.cpp algorithms/shared/cpp/tests/test_lifecycle.cpp
git commit -m "refactor(cpp): deepen shared buffer retry loop"
```

## Task 3: Port the deeper seam to Rust

**Files:**
- Modify: `algorithms/shared/rust/src/codec/buffer.rs`
- Create: `algorithms/shared/rust/tests/buffer.rs`
- Test: `algorithms/shared/rust/tests/buffer.rs`

Rust keeps `BufTooSmall` handling transactional at this seam. Unlike the Go port, the Rust streaming traits return `Result<usize, CodecError>`, so an error cannot also carry a partial byte count for accumulation.

- [ ] **Step 1: Add failing Rust regressions for retry preservation**

Create `algorithms/shared/rust/tests/buffer.rs` with fake `Encoder` / `Decoder` implementations that model partial writes on `BufTooSmall`:

```rust
use compresskit_codec::codec::{encode_buffer, decode_buffer, CodecError, Decoder, Encoder, State};

struct ScriptedEncoder {
    finish_calls: usize,
}

impl Encoder for ScriptedEncoder {
    fn process(&mut self, _: &[u8], _: &mut [u8]) -> Result<usize, CodecError> {
        Ok(0)
    }

    fn flush(&mut self, _: &mut [u8]) -> Result<usize, CodecError> {
        Ok(0)
    }

    fn finish(&mut self, out: &mut [u8]) -> Result<usize, CodecError> {
        self.finish_calls += 1;
        if self.finish_calls == 1 {
            out[..3].copy_from_slice(b"abc");
            return Err(CodecError::BufTooSmall);
        }
        out[..3].copy_from_slice(b"def");
        Ok(3)
    }

    fn reset(&mut self) {}
    fn state(&self) -> State { State::Streaming }
}

#[test]
fn encode_buffer_preserves_partial_output_across_finish_retry() {
    let mut encoder = ScriptedEncoder { finish_calls: 0 };
    let out = encode_buffer(&mut encoder, b"x").unwrap();
    assert_eq!(out, b"abcdef");
}
```

- [ ] **Step 2: Run the Rust shared crate tests to verify the regression fails**

Run: `cargo test --manifest-path algorithms/shared/rust/Cargo.toml`

Expected: FAIL in the new Rust buffer retry test because the shared helper has not been extracted yet.

- [ ] **Step 3: Write the minimal Rust implementation**

Refactor `buffer.rs` to centralize retry/growth policy in one helper:

```rust
type BufferStep<'a> = dyn FnMut(&mut [u8]) -> Result<usize, CodecError> + 'a;

fn run_buffer_step(
    out_buf: &mut Vec<u8>,
    total_written: usize,
    limit: usize,
    step: &mut BufferStep<'_>,
) -> Result<usize, CodecError> {
    let mut total_written = total_written;

    loop {
        match step(&mut out_buf[total_written..]) {
            Ok(n) => return Ok(total_written + n),
            Err(CodecError::BufTooSmall) => {
                if total_written > limit || out_buf.len() >= limit {
                    return Err(CodecError::SizeLimit);
                }

                let new_size = grow_buffer(out_buf.len(), limit);
                if new_size <= out_buf.len() {
                    return Err(CodecError::SizeLimit);
                }

                out_buf.resize(new_size, 0);
            }
            Err(err) => return Err(err),
        }
    }
}
```

Use the helper from both `encode_buffer` and `decode_buffer` so only one internal implementation owns the retry/growth policy.

- [ ] **Step 4: Run the Rust shared crate tests again**

Run: `cargo test --manifest-path algorithms/shared/rust/Cargo.toml`

Expected: PASS

- [ ] **Step 5: Commit the Rust port**

```bash
git add algorithms/shared/rust/src/codec/buffer.rs algorithms/shared/rust/tests/buffer.rs
git commit -m "refactor(rust): deepen shared buffer retry loop"
```

## Task 4: Validate the cross-language slice and leave the tree clean

**Files:**
- Modify: any files changed in Tasks 1-3
- Test: repository baseline commands

- [ ] **Step 1: Run the repository test baseline**

Run: `make test`

Expected: PASS

- [ ] **Step 2: Run the repository lint baseline**

Run: `make lint`

Expected: PASS

- [ ] **Step 3: Review the final diff for accidental interface drift**

Run: `git --no-pager diff --stat HEAD~3..HEAD`

Expected: only shared Buffer Layer internals and the new focused tests changed; no algorithm wire-format files, no CLI contract changes, no generated artifacts

- [ ] **Step 4: Commit the validation checkpoint**

```bash
git add algorithms/shared/go/codec/buffer.go algorithms/shared/go/codec/buffer_loop.go algorithms/shared/go/codec/buffer_internal_test.go algorithms/shared/cpp/src/buffer_api.cpp algorithms/shared/cpp/tests/test_lifecycle.cpp algorithms/shared/rust/src/codec/buffer.rs algorithms/shared/rust/tests/buffer.rs
git commit -m "test: validate shared buffer deepening slice"
```

## Self-review

- **Spec coverage:** This plan implements the approved design for slice 1 only and explicitly defers slices 2-5 to follow-on plans.
- **Placeholder scan:** No `TODO`, `TBD`, or “similar to above” shortcuts remain.
- **Type consistency:** The helper names are consistent across tasks: `runBufferStep` in Go, `run_buffer_step` in C++, and `run_buffer_step` in Rust.
