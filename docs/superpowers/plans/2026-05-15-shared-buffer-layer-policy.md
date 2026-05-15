# Shared Buffer Layer Policy Module Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Deepen the Go and Rust Shared Buffer Layer internals so buffer resizing, retry, and limit policy live in one module per language while CompressKit's Buffer Layer interface stays unchanged.

**Architecture:** Extract one internal resizing-buffer runner per language, then migrate `EncodeBuffer` / `DecodeBuffer` onto it without changing the Buffer Layer surface. Keep Phase 1 scoped to full-buffer orchestration only; do not fold `WriterEncoder` into the new seam yet because that is Phase 5.

**Tech Stack:** Go, Rust, repository Make targets, `go test`, `cargo test`, `make test`, `make lint`

---

## Scope note

This plan covers **Phase 1 only** from `docs/superpowers/specs/2026-05-15-ordered-architecture-deepening-design.md`.

Phases 2-5 should each get their own plan after Phase 1 lands cleanly:

1. Streaming Layer wrapper module
2. Frequency Table module
3. Magic Number and header parsing module
4. Writer adapter module

## File structure

**Create**

- `algorithms/shared/go/codec/buffer_policy.go` - internal Go resizing-buffer runner for Buffer Layer full-buffer orchestration
- `algorithms/shared/rust/src/codec/buffer_policy.rs` - internal Rust resizing-buffer runner for Buffer Layer full-buffer orchestration

**Modify**

- `algorithms/shared/go/codec/buffer.go` - delegate Go Buffer Layer encode/decode paths to the new runner
- `algorithms/shared/go/codec/buffer_internal_test.go` - cover the deeper Go policy module directly and keep contract-focused Buffer Layer tests
- `algorithms/shared/rust/src/lib.rs` - register the new internal Rust codec module
- `algorithms/shared/rust/src/codec/buffer.rs` - delegate Rust Buffer Layer encode/decode paths to the new runner

**Delete**

- `algorithms/shared/go/codec/buffer_loop.go` - remove the old shallow retry helper once `buffer_policy.go` owns the policy

**Test**

- `algorithms/shared/go/codec/buffer_internal_test.go`
- `algorithms/shared/rust/src/codec/buffer_policy.rs`
- `algorithms/shared/rust/src/codec/buffer.rs`

### Task 1: Deepen the Go Buffer Layer orchestration

**Files:**
- Create: `algorithms/shared/go/codec/buffer_policy.go`
- Modify: `algorithms/shared/go/codec/buffer.go`
- Modify: `algorithms/shared/go/codec/buffer_internal_test.go`
- Delete: `algorithms/shared/go/codec/buffer_loop.go`
- Test: `algorithms/shared/go/codec/buffer_internal_test.go`

- [ ] **Step 1: Write the failing Go policy tests**

```go
func TestResizingBuffer_PreservesWrittenBytesAcrossRetry(t *testing.T) {
	runner := newResizingBuffer(3, 12)
	calls := 0

	err := runner.run(func(out []byte) (int, error) {
		calls++
		if calls == 1 {
			copy(out, []byte("abc"))
			return 3, ErrBufTooSmall
		}
		copy(out, []byte("def"))
		return 3, nil
	})
	if err != nil {
		t.Fatalf("run() error = %v", err)
	}
	if got := string(runner.bytes()); got != "abcdef" {
		t.Fatalf("bytes() = %q, want %q", got, "abcdef")
	}
}

func TestNewResizingBuffer_ClampsInitialSizeToLimit(t *testing.T) {
	runner := newResizingBuffer(8, 3)
	if got := cap(runner.buf); got != 3 {
		t.Fatalf("cap(buf) = %d, want 3", got)
	}
}

func TestEncodeBufferWithLimit_UsesResizingBufferOutput(t *testing.T) {
	stub := &scriptedEncoder{
		process: []scriptedCall{{written: 0, err: nil}},
		finish: []scriptedCall{
			{written: 3, err: ErrBufTooSmall, payload: []byte("abc")},
			{written: 3, err: nil, payload: []byte("def")},
		},
	}

	out, err := encodeBufferWithLimit(stub, []byte("ignored"), 3, 12)
	if err != nil {
		t.Fatalf("encodeBufferWithLimit() error = %v", err)
	}
	if string(out) != "abcdef" {
		t.Fatalf("encodeBufferWithLimit() = %q, want %q", out, "abcdef")
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./algorithms/shared/go/... -run 'TestResizingBuffer_|TestNewResizingBuffer_|TestEncodeBufferWithLimit_UsesResizingBufferOutput'`
Expected: FAIL with `undefined: newResizingBuffer`

- [ ] **Step 3: Write the minimal Go policy implementation**

```go
type resizingBuffer struct {
	buf     []byte
	written int
	limit   int
}

func newResizingBuffer(initialSize int, limit int) *resizingBuffer {
	if initialSize > limit {
		initialSize = limit
	}
	return &resizingBuffer{
		buf:   make([]byte, initialSize),
		limit: limit,
	}
}

func (r *resizingBuffer) run(step bufferStep) error {
	for {
		n, err := step(r.buf[r.written:])
		if !errors.Is(err, ErrBufTooSmall) {
			if err != nil {
				return err
			}
			r.written += n
			if r.written > r.limit {
				return ErrSizeLimit
			}
			return nil
		}

		r.written += n
		if r.written > r.limit || len(r.buf) >= r.limit {
			return ErrSizeLimit
		}

		newSize := growBuffer(len(r.buf), r.limit)
		if newSize <= len(r.buf) {
			return ErrSizeLimit
		}

		next := make([]byte, newSize)
		copy(next, r.buf[:r.written])
		r.buf = next
	}
}

func (r *resizingBuffer) bytes() []byte {
	return r.buf[:r.written]
}

func encodeBufferWithLimit(encoder Encoder, input []byte, initialSize int, limit int) ([]byte, error) {
	runner := newResizingBuffer(initialSize, limit)

	if err := runner.run(func(out []byte) (int, error) {
		return encoder.Process(input, out)
	}); err != nil {
		return nil, err
	}
	if err := runner.run(func(out []byte) (int, error) {
		return encoder.Finish(out)
	}); err != nil {
		return nil, err
	}
	return runner.bytes(), nil
}

func decodeBufferWithLimit(decoder Decoder, input []byte, initialSize int, limit int) ([]byte, error) {
	runner := newResizingBuffer(initialSize, limit)

	if err := runner.run(func(out []byte) (int, error) {
		return decoder.Process(input, out)
	}); err != nil {
		return nil, err
	}
	if err := runner.run(func(out []byte) (int, error) {
		return decoder.Finish(out)
	}); err != nil {
		return nil, err
	}
	return runner.bytes(), nil
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./algorithms/shared/go/...`
Expected: PASS with `ok` for `algorithms/shared/go/codec`

- [ ] **Step 5: Commit**

```bash
git add algorithms/shared/go/codec/buffer.go algorithms/shared/go/codec/buffer_policy.go algorithms/shared/go/codec/buffer_internal_test.go
git rm algorithms/shared/go/codec/buffer_loop.go
git commit -m "refactor(shared-go): deepen buffer layer orchestration" -m "Co-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>"
```

### Task 2: Deepen the Rust Buffer Layer orchestration

**Files:**
- Create: `algorithms/shared/rust/src/codec/buffer_policy.rs`
- Modify: `algorithms/shared/rust/src/lib.rs`
- Modify: `algorithms/shared/rust/src/codec/buffer.rs`
- Test: `algorithms/shared/rust/src/codec/buffer_policy.rs`
- Test: `algorithms/shared/rust/src/codec/buffer.rs`

- [ ] **Step 1: Write the failing Rust policy tests**

```rust
#[test]
fn resizing_buffer_retries_after_buf_too_small() {
    let mut runner = ResizingBuffer::new(1, 8);
    let mut calls = 0;

    runner
        .run(&mut |out| {
            calls += 1;
            if calls == 1 {
                return Err(CodecError::BufTooSmall);
            }
            out[..3].copy_from_slice(b"def");
            Ok(3)
        })
        .unwrap();

    assert_eq!(runner.into_vec(), b"def");
}

#[test]
fn resizing_buffer_clamps_initial_size_to_limit() {
    let runner = ResizingBuffer::new(8, 3);
    assert_eq!(runner.capacity(), 3);
}

#[test]
fn encode_buffer_retries_finish_after_buffer_growth() {
    struct RetryEncoder {
        calls: usize,
    }

    impl Encoder for RetryEncoder {
        fn process(&mut self, _: &[u8], _: &mut [u8]) -> Result<usize, CodecError> { Ok(0) }
        fn flush(&mut self, _: &mut [u8]) -> Result<usize, CodecError> { Ok(0) }
        fn finish(&mut self, out: &mut [u8]) -> Result<usize, CodecError> {
            self.calls += 1;
            if self.calls == 1 {
                return Err(CodecError::BufTooSmall);
            }
            out[..6].copy_from_slice(b"abcdef");
            Ok(6)
        }
        fn reset(&mut self) {}
        fn state(&self) -> State { State::Streaming }
    }

    let mut encoder = RetryEncoder { calls: 0 };
    let out = encode_buffer(&mut encoder, b"ignored").unwrap();
    assert_eq!(out, b"abcdef");
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cargo test --manifest-path algorithms/shared/rust/Cargo.toml 'resizing_buffer|encode_buffer_retries_finish_after_buffer_growth'`
Expected: FAIL with `cannot find type 'ResizingBuffer'`

- [ ] **Step 3: Write the minimal Rust policy implementation**

```rust
pub mod codec {
    pub mod bits;
    pub mod buffer;
    pub mod buffer_policy;
    pub mod buffered;
    pub mod encoder;
    pub mod error;
    pub mod write;
}

use crate::codec::buffer_policy::ResizingBuffer;

pub(crate) struct ResizingBuffer {
    buf: Vec<u8>,
    written: usize,
    limit: usize,
}

impl ResizingBuffer {
    pub(crate) fn new(initial_size: usize, limit: usize) -> Self {
        let size = initial_size.min(limit);
        Self { buf: vec![0; size], written: 0, limit }
    }

    pub(crate) fn capacity(&self) -> usize {
        self.buf.len()
    }

    pub(crate) fn run<F>(&mut self, step: &mut F) -> Result<(), CodecError>
    where
        F: FnMut(&mut [u8]) -> Result<usize, CodecError>,
    {
        loop {
            match step(&mut self.buf[self.written..]) {
                Ok(n) => {
                    self.written = self.written.checked_add(n).ok_or(CodecError::SizeLimit)?;
                    if self.written > self.limit {
                        return Err(CodecError::SizeLimit);
                    }
                    return Ok(());
                }
                Err(CodecError::BufTooSmall) => {
                    if self.written > self.limit || self.buf.len() >= self.limit {
                        return Err(CodecError::SizeLimit);
                    }
                    let new_size = grow_buffer(self.buf.len(), self.limit);
                    if new_size <= self.buf.len() {
                        return Err(CodecError::SizeLimit);
                    }
                    self.buf.resize(new_size, 0);
                }
                Err(err) => return Err(err),
            }
        }
    }

    pub(crate) fn into_vec(mut self) -> Vec<u8> {
        self.buf.truncate(self.written);
        self.buf
    }
}

pub fn encode_buffer(encoder: &mut dyn Encoder, input: &[u8]) -> Result<Vec<u8>, CodecError> {
    if input.len() > MAX_INPUT_SIZE {
        return Err(CodecError::SizeLimit);
    }

    let encode_limit = encode_buffer_limit(input.len())?;
    let mut runner = ResizingBuffer::new(
        input.len().saturating_mul(2).saturating_add(2048),
        encode_limit,
    );

    runner.run(&mut |output| encoder.process(input, output))?;
    runner.run(&mut |output| encoder.finish(output))?;
    Ok(runner.into_vec())
}

pub(crate) fn decode_buffer_with_limit(
    decoder: &mut dyn Decoder,
    input: &[u8],
    limit: usize,
) -> Result<Vec<u8>, CodecError> {
    if input.len() > MAX_INPUT_SIZE {
        return Err(CodecError::SizeLimit);
    }

    let mut runner = ResizingBuffer::new(input.len().saturating_add(1024), limit);
    runner.run(&mut |output| decoder.process(input, output))?;
    runner.run(&mut |output| decoder.finish(output))?;
    Ok(runner.into_vec())
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `cargo test --manifest-path algorithms/shared/rust/Cargo.toml`
Expected: PASS with `test result: ok`

- [ ] **Step 5: Commit**

```bash
git add algorithms/shared/rust/src/lib.rs algorithms/shared/rust/src/codec/buffer.rs algorithms/shared/rust/src/codec/buffer_policy.rs
git commit -m "refactor(shared-rust): deepen buffer layer orchestration" -m "Co-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>"
```

### Task 3: Verify the phase and prepare the next plan

**Files:**
- Test: repository root validation commands

- [ ] **Step 1: Run the shared-language validation commands**

Run:

```bash
go test ./algorithms/shared/go/... && cargo test --manifest-path algorithms/shared/rust/Cargo.toml
```

Expected: PASS with `ok` / `test result: ok`

- [ ] **Step 2: Run repository validation**

Run:

```bash
make test && make lint
```

Expected: PASS across the existing repository baseline

- [ ] **Step 3: Inspect the diff for scope control**

Run:

```bash
git --no-pager diff --stat HEAD~2..HEAD
```

Expected: only Shared Buffer Layer files and their direct tests are touched
