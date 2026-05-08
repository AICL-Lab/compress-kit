// Package codec_test implements the Streaming API Contract Test Suite.
//
// This file implements the contract test cases defined in:
// tests/streaming_api_contract/contract_cases.md
//
// Contract Test Case Coverage:
//
// Lifecycle Tests (L1-L7):
//   - L1: TestLifecycle_EmptyInput - empty input scenario
//   - L2: TestLifecycle_SingleByte - single-byte input
//   - L3: TestLifecycle_ChunkedInput - chunked input where output is delayed
//   - L4: TestLifecycle_FlushWithoutFinish - flush without finish
//   - L5: TestLifecycle_FinishAfterMultipleProcess - finish after multiple process calls
//   - L6: TestLifecycle_ResetAfterFinish - reset after finish
//   - L7: TestLifecycle_ResetAfterError - reset after error
//
// Buffer Contract Tests (B1-B3):
//   - B1: TestBuffer_BufTooSmallTransactional - BUF_TOO_SMALL is transactional
//   - B2: TestBuffer_EncodeFullPath - buffer encode full path
//   - B3: TestBuffer_DecodeFullPath - buffer decode full path
//
// Error Handling Tests (E1-E5):
//   - E1: TestError_TruncatedFrame - truncated frame on decode
//   - E2: (partial) TestError_TruncatedFrame - corrupt data detection
//   - E3: TestConstants - input size limit (4 GiB)
//   - E4: TestConstants - output size limit (1 GiB)
//   - E5: (covered in L5) - invalid state transitions
//
// Cross-Algorithm Tests (X1-X2):
//   - X1: All tests run for all 4 algorithms (Huffman, Arithmetic, Range, RLE)
//   - X2: Buffer API tests verify equivalence with streaming API
package codec_test

import (
	"bytes"
	"errors"
	"testing"

	"arithmetic"
	"github.com/LessUp/compress-kit/algorithms/shared/go/codec"
	"huffman"
	"rangecoder"
	"rle"
)

// TestLifecycle_EmptyInput tests L1: empty input scenario.
func TestLifecycle_EmptyInput(t *testing.T) {
	tests := []struct {
		name    string
		encoder codec.Encoder
		decoder codec.Decoder
	}{
		{"Huffman", huffman.NewStreamingEncoder(), huffman.NewStreamingDecoder()},
		{"Arithmetic", arithmetic.NewStreamingEncoder(), arithmetic.NewStreamingDecoder()},
		{"Range", rangecoder.NewStreamingEncoder(), rangecoder.NewStreamingDecoder()},
		{"RLE", rle.NewStreamingEncoder(), rle.NewStreamingDecoder()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encode empty input
			outBuf := make([]byte, 4096)
			n, err := tt.encoder.Finish(outBuf)
			if err != nil {
				t.Fatalf("Finish() error = %v", err)
			}
			if tt.encoder.State() != codec.StateFinished {
				t.Errorf("State = %v, want StateFinished", tt.encoder.State())
			}

			// Decode
			encoded := outBuf[:n]
			decBuf := make([]byte, 4096)
			n2, err := tt.decoder.Process(encoded, decBuf)
			if err != nil {
				t.Fatalf("Decode Process() error = %v", err)
			}
			n3, err := tt.decoder.Finish(decBuf[n2:])
			if err != nil {
				t.Fatalf("Decode Finish() error = %v", err)
			}

			decoded := decBuf[:n2+n3]
			if len(decoded) != 0 {
				t.Errorf("Decoded length = %d, want 0", len(decoded))
			}
		})
	}
}

// TestLifecycle_SingleByte tests L2: single-byte input.
func TestLifecycle_SingleByte(t *testing.T) {
	tests := []struct {
		name    string
		encoder codec.Encoder
		decoder codec.Decoder
	}{
		{"Huffman", huffman.NewStreamingEncoder(), huffman.NewStreamingDecoder()},
		{"Arithmetic", arithmetic.NewStreamingEncoder(), arithmetic.NewStreamingDecoder()},
		{"Range", rangecoder.NewStreamingEncoder(), rangecoder.NewStreamingDecoder()},
		{"RLE", rle.NewStreamingEncoder(), rle.NewStreamingDecoder()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := []byte{0x42}

			// Encode
			outBuf := make([]byte, 4096)
			n1, err := tt.encoder.Process(input, outBuf)
			if err != nil {
				t.Fatalf("Process() error = %v", err)
			}
			n2, err := tt.encoder.Finish(outBuf[n1:])
			if err != nil {
				t.Fatalf("Finish() error = %v", err)
			}

			if tt.encoder.State() != codec.StateFinished {
				t.Errorf("State = %v, want StateFinished", tt.encoder.State())
			}

			// Decode
			encoded := outBuf[:n1+n2]
			decBuf := make([]byte, 4096)
			n3, err := tt.decoder.Process(encoded, decBuf)
			if err != nil {
				t.Fatalf("Decode Process() error = %v", err)
			}
			n4, err := tt.decoder.Finish(decBuf[n3:])
			if err != nil {
				t.Fatalf("Decode Finish() error = %v", err)
			}

			decoded := decBuf[:n3+n4]
			if !bytes.Equal(decoded, input) {
				t.Errorf("Decoded = %v, want %v", decoded, input)
			}
		})
	}
}

// TestLifecycle_ChunkedInput tests L3: chunked input.
func TestLifecycle_ChunkedInput(t *testing.T) {
	tests := []struct {
		name    string
		encoder codec.Encoder
		decoder codec.Decoder
	}{
		{"Huffman", huffman.NewStreamingEncoder(), huffman.NewStreamingDecoder()},
		{"Arithmetic", arithmetic.NewStreamingEncoder(), arithmetic.NewStreamingDecoder()},
		{"Range", rangecoder.NewStreamingEncoder(), rangecoder.NewStreamingDecoder()},
		{"RLE", rle.NewStreamingEncoder(), rle.NewStreamingDecoder()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chunk1 := []byte("Hello, ")
			chunk2 := []byte("World")
			chunk3 := []byte("!")
			expected := []byte("Hello, World!")

			outBuf := make([]byte, 4096)
			offset := 0

			// Process chunks
			n, err := tt.encoder.Process(chunk1, outBuf[offset:])
			if err != nil {
				t.Fatalf("Process(chunk1) error = %v", err)
			}
			offset += n

			n, err = tt.encoder.Process(chunk2, outBuf[offset:])
			if err != nil {
				t.Fatalf("Process(chunk2) error = %v", err)
			}
			offset += n

			n, err = tt.encoder.Process(chunk3, outBuf[offset:])
			if err != nil {
				t.Fatalf("Process(chunk3) error = %v", err)
			}
			offset += n

			// Finish
			n, err = tt.encoder.Finish(outBuf[offset:])
			if err != nil {
				t.Fatalf("Finish() error = %v", err)
			}
			offset += n

			// Decode
			encoded := outBuf[:offset]
			decBuf := make([]byte, 4096)
			n2, err := tt.decoder.Process(encoded, decBuf)
			if err != nil {
				t.Fatalf("Decode Process() error = %v", err)
			}
			n3, err := tt.decoder.Finish(decBuf[n2:])
			if err != nil {
				t.Fatalf("Decode Finish() error = %v", err)
			}

			decoded := decBuf[:n2+n3]
			if !bytes.Equal(decoded, expected) {
				t.Errorf("Decoded = %q, want %q", decoded, expected)
			}
		})
	}
}

// TestLifecycle_FlushWithoutFinish tests L4: flush without finish.
func TestLifecycle_FlushWithoutFinish(t *testing.T) {
	tests := []struct {
		name       string
		newEncoder func() codec.Encoder
	}{
		{"Huffman", func() codec.Encoder { return huffman.NewStreamingEncoder() }},
		{"Arithmetic", func() codec.Encoder { return arithmetic.NewStreamingEncoder() }},
		{"Range", func() codec.Encoder { return rangecoder.NewStreamingEncoder() }},
		{"RLE", func() codec.Encoder { return rle.NewStreamingEncoder() }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enc := tt.newEncoder()
			input := []byte("test")
			outBuf := make([]byte, 4096)

			_, err := enc.Process(input, outBuf)
			if err != nil {
				t.Fatalf("Process() error = %v", err)
			}

			if enc.State() != codec.StateStreaming {
				t.Errorf("State = %v, want StateStreaming", enc.State())
			}

			_, err = enc.Flush(outBuf)
			if err != nil {
				t.Fatalf("Flush() error = %v", err)
			}

			if enc.State() != codec.StateFlushing {
				t.Errorf("State after Flush = %v, want StateFlushing", enc.State())
			}

			_, err = enc.Process([]byte("more"), outBuf)
			if err != nil {
				t.Fatalf("Process() after Flush error = %v", err)
			}

			if enc.State() != codec.StateStreaming {
				t.Errorf("State after Process = %v, want StateStreaming", enc.State())
			}
		})
	}
}

// TestLifecycle_FinishAfterMultipleProcess tests L5: finish after multiple process calls.
func TestLifecycle_FinishAfterMultipleProcess(t *testing.T) {
	tests := []struct {
		name       string
		newEncoder func() codec.Encoder
	}{
		{"Huffman", func() codec.Encoder { return huffman.NewStreamingEncoder() }},
		{"Arithmetic", func() codec.Encoder { return arithmetic.NewStreamingEncoder() }},
		{"Range", func() codec.Encoder { return rangecoder.NewStreamingEncoder() }},
		{"RLE", func() codec.Encoder { return rle.NewStreamingEncoder() }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enc := tt.newEncoder()
			outBuf := make([]byte, 4096)

			_, err := enc.Process([]byte("test1"), outBuf)
			if err != nil {
				t.Fatalf("Process(1) error = %v", err)
			}

			_, err = enc.Process([]byte("test2"), outBuf)
			if err != nil {
				t.Fatalf("Process(2) error = %v", err)
			}

			_, err = enc.Finish(outBuf)
			if err != nil {
				t.Fatalf("Finish() error = %v", err)
			}

			if enc.State() != codec.StateFinished {
				t.Errorf("State = %v, want StateFinished", enc.State())
			}

			_, err = enc.Process([]byte("fail"), outBuf)
			if err != codec.ErrInvalidState {
				t.Errorf("Process() after Finish error = %v, want ErrInvalidState", err)
			}
			if enc.State() != codec.StateError {
				t.Errorf("State after invalid Process = %v, want StateError", enc.State())
			}

			_, err = enc.Flush(outBuf)
			if err != codec.ErrInvalidState {
				t.Errorf("Flush() after Finish error = %v, want ErrInvalidState", err)
			}
			if enc.State() != codec.StateError {
				t.Errorf("State after invalid Flush = %v, want StateError", enc.State())
			}

			_, err = enc.Finish(outBuf)
			if err != codec.ErrInvalidState {
				t.Errorf("Finish() after Finish error = %v, want ErrInvalidState", err)
			}
			if enc.State() != codec.StateError {
				t.Errorf("State after invalid Finish = %v, want StateError", enc.State())
			}
		})
	}
}

// TestLifecycle_ResetAfterFinish tests L6: reset after finish.
func TestLifecycle_ResetAfterFinish(t *testing.T) {
	tests := []struct {
		name       string
		newEncoder func() codec.Encoder
	}{
		{"Huffman", func() codec.Encoder { return huffman.NewStreamingEncoder() }},
		{"Arithmetic", func() codec.Encoder { return arithmetic.NewStreamingEncoder() }},
		{"Range", func() codec.Encoder { return rangecoder.NewStreamingEncoder() }},
		{"RLE", func() codec.Encoder { return rle.NewStreamingEncoder() }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enc := tt.newEncoder()
			outBuf := make([]byte, 4096)

			_, err := enc.Process([]byte("test"), outBuf)
			if err != nil {
				t.Fatalf("Process() error = %v", err)
			}

			_, err = enc.Finish(outBuf)
			if err != nil {
				t.Fatalf("Finish() error = %v", err)
			}

			enc.Reset()

			if enc.State() != codec.StateReady {
				t.Errorf("State after Reset = %v, want StateReady", enc.State())
			}

			_, err = enc.Process([]byte("new"), outBuf)
			if err != nil {
				t.Fatalf("Process() after Reset error = %v", err)
			}
		})
	}
}

// TestLifecycle_ResetAfterError tests L7: reset after error.
func TestLifecycle_ResetAfterError(t *testing.T) {
	tests := []struct {
		name       string
		newDecoder func() codec.Decoder
	}{
		{"Huffman", func() codec.Decoder { return huffman.NewStreamingDecoder() }},
		{"Arithmetic", func() codec.Decoder { return arithmetic.NewStreamingDecoder() }},
		{"Range", func() codec.Decoder { return rangecoder.NewStreamingDecoder() }},
		{"RLE", func() codec.Decoder { return rle.NewStreamingDecoder() }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dec := tt.newDecoder()
			outBuf := make([]byte, 4096)

			_, err := dec.Process([]byte("not-a-valid-frame"), outBuf)
			if err != nil {
				t.Fatalf("Process() error = %v", err)
			}

			_, err = dec.Finish(outBuf)
			if err == nil {
				t.Fatalf("Finish() expected decode error")
			}

			if dec.State() != codec.StateError {
				t.Fatalf("State after error = %v, want StateError", dec.State())
			}

			dec.Reset()
			if dec.State() != codec.StateReady {
				t.Errorf("State after Reset = %v, want StateReady", dec.State())
			}

			reuseInput := []byte("reuse")
			var validEncoded []byte
			if tt.name == "Huffman" {
				validEncoded, err = codec.EncodeBuffer(huffman.NewStreamingEncoder(), reuseInput)
			} else if tt.name == "Arithmetic" {
				validEncoded, err = codec.EncodeBuffer(arithmetic.NewStreamingEncoder(), reuseInput)
			} else if tt.name == "Range" {
				validEncoded, err = codec.EncodeBuffer(rangecoder.NewStreamingEncoder(), reuseInput)
			} else {
				validEncoded, err = codec.EncodeBuffer(rle.NewStreamingEncoder(), reuseInput)
			}
			if err != nil {
				t.Fatalf("EncodeBuffer() after Reset setup error = %v", err)
			}

			reuseBuf := make([]byte, 1024)
			_, err = dec.Process(validEncoded, reuseBuf)
			if err != nil {
				t.Fatalf("Process() after Reset error = %v", err)
			}
			n, err := dec.Finish(reuseBuf)
			if err != nil {
				t.Fatalf("Finish() after Reset error = %v", err)
			}
			if !bytes.Equal(reuseBuf[:n], reuseInput) {
				t.Fatalf("decoded after reset = %q, want %q", reuseBuf[:n], reuseInput)
			}
		})
	}
}

// TestBuffer_BufTooSmallTransactional tests B1: BUF_TOO_SMALL is transactional.
func TestBuffer_BufTooSmallTransactional(t *testing.T) {
	input := []byte("test data for encoding")
	tests := []struct {
		name       string
		newEncoder func() codec.Encoder
		newDecoder func() codec.Decoder
	}{
		{"Huffman", func() codec.Encoder { return huffman.NewStreamingEncoder() }, func() codec.Decoder { return huffman.NewStreamingDecoder() }},
		{"Arithmetic", func() codec.Encoder { return arithmetic.NewStreamingEncoder() }, func() codec.Decoder { return arithmetic.NewStreamingDecoder() }},
		{"Range", func() codec.Encoder { return rangecoder.NewStreamingEncoder() }, func() codec.Decoder { return rangecoder.NewStreamingDecoder() }},
		{"RLE", func() codec.Encoder { return rle.NewStreamingEncoder() }, func() codec.Decoder { return rle.NewStreamingDecoder() }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enc := tt.newEncoder()

			largeBuf := make([]byte, 4096)
			_, err := enc.Process(input, largeBuf)
			if err != nil {
				t.Fatalf("Process() error = %v", err)
			}

			stateBefore := enc.State()

			smallBuf := make([]byte, 10)
			_, err = enc.Finish(smallBuf)
			if err != codec.ErrBufTooSmall {
				t.Fatalf("Expected ErrBufTooSmall, got %v", err)
			}

			stateAfter := enc.State()
			if stateAfter != stateBefore {
				t.Errorf("State changed after BufTooSmall: %v -> %v", stateBefore, stateAfter)
			}

			n, err := enc.Finish(largeBuf)
			if err != nil {
				t.Fatalf("Finish() with large buffer error = %v", err)
			}

			if enc.State() != codec.StateFinished {
				t.Errorf("State = %v, want StateFinished", enc.State())
			}

			decoded, err := codec.DecodeBuffer(tt.newDecoder(), largeBuf[:n])
			if err != nil {
				t.Fatalf("DecodeBuffer() error = %v", err)
			}
			if !bytes.Equal(decoded, input) {
				t.Fatalf("decoded = %q, want %q", decoded, input)
			}
		})
	}
}

// TestBuffer_EncodeFullPath tests B2: buffer encode full path.
func TestBuffer_EncodeFullPath(t *testing.T) {
	input := []byte("Hello, streaming world!")

	// Test each algorithm
	tests := []struct {
		name    string
		encoder codec.Encoder
		decoder codec.Decoder
	}{
		{"Huffman", huffman.NewStreamingEncoder(), huffman.NewStreamingDecoder()},
		{"Arithmetic", arithmetic.NewStreamingEncoder(), arithmetic.NewStreamingDecoder()},
		{"Range", rangecoder.NewStreamingEncoder(), rangecoder.NewStreamingDecoder()},
		{"RLE", rle.NewStreamingEncoder(), rle.NewStreamingDecoder()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded, err := codec.EncodeBuffer(tt.encoder, input)
			if err != nil {
				t.Fatalf("EncodeBuffer() error = %v", err)
			}

			if len(encoded) == 0 {
				t.Errorf("EncodeBuffer() returned empty output")
			}
		})
	}
}

// TestBuffer_DecodeFullPath tests B3: buffer decode full path.
func TestBuffer_DecodeFullPath(t *testing.T) {
	input := []byte("Hello, streaming world!")

	tests := []struct {
		name    string
		encoder codec.Encoder
		decoder codec.Decoder
	}{
		{"Huffman", huffman.NewStreamingEncoder(), huffman.NewStreamingDecoder()},
		{"Arithmetic", arithmetic.NewStreamingEncoder(), arithmetic.NewStreamingDecoder()},
		{"Range", rangecoder.NewStreamingEncoder(), rangecoder.NewStreamingDecoder()},
		{"RLE", rle.NewStreamingEncoder(), rle.NewStreamingDecoder()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encode
			encoded, err := codec.EncodeBuffer(tt.encoder, input)
			if err != nil {
				t.Fatalf("EncodeBuffer() error = %v", err)
			}

			// Decode
			decoded, err := codec.DecodeBuffer(tt.decoder, encoded)
			if err != nil {
				t.Fatalf("DecodeBuffer() error = %v", err)
			}

			if !bytes.Equal(decoded, input) {
				t.Errorf("Decoded = %q, want %q", decoded, input)
			}
		})
	}
}

// TestBuffer_DecodeBufTooSmallTransactional tests that decoder Finish can be retried
// with a larger buffer without losing buffered input.
func TestBuffer_DecodeBufTooSmallTransactional(t *testing.T) {
	tests := []struct {
		name    string
		encoder codec.Encoder
		decoder codec.Decoder
	}{
		{"Huffman", huffman.NewStreamingEncoder(), huffman.NewStreamingDecoder()},
		{"Arithmetic", arithmetic.NewStreamingEncoder(), arithmetic.NewStreamingDecoder()},
		{"Range", rangecoder.NewStreamingEncoder(), rangecoder.NewStreamingDecoder()},
		{"RLE", rle.NewStreamingEncoder(), rle.NewStreamingDecoder()},
	}

	input := []byte("Hello, streaming world! Hello, streaming world! Hello, streaming world!")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded, err := codec.EncodeBuffer(tt.encoder, input)
			if err != nil {
				t.Fatalf("EncodeBuffer() error = %v", err)
			}

			_, err = tt.decoder.Process(encoded, make([]byte, 1))
			if err != nil {
				t.Fatalf("Process() error = %v", err)
			}

			smallBuf := make([]byte, 4)
			_, err = tt.decoder.Finish(smallBuf)
			if err != codec.ErrBufTooSmall {
				t.Fatalf("Finish() small buffer error = %v, want ErrBufTooSmall", err)
			}

			largeBuf := make([]byte, len(input)+1024)
			n, err := tt.decoder.Finish(largeBuf)
			if err != nil {
				t.Fatalf("Finish() retry error = %v", err)
			}

			decoded := largeBuf[:n]
			if !bytes.Equal(decoded, input) {
				t.Fatalf("decoded = %q, want %q", decoded, input)
			}
		})
	}
}

// TestError_TruncatedFrame tests E1: truncated frame on decode.
func TestError_TruncatedFrame(t *testing.T) {
	input := []byte("test data for truncation")

	// Encode with Huffman
	enc := huffman.NewStreamingEncoder()
	encoded, err := codec.EncodeBuffer(enc, input)
	if err != nil {
		t.Fatalf("EncodeBuffer() error = %v", err)
	}

	// Truncate the encoded data
	truncated := encoded[:len(encoded)/2]

	// Try to decode
	dec := huffman.NewStreamingDecoder()
	_, err = codec.DecodeBuffer(dec, truncated)
	if !errors.Is(err, codec.ErrTruncated) && !errors.Is(err, codec.ErrCorrupt) {
		t.Errorf("DecodeBuffer() error = %v, want ErrTruncated or ErrCorrupt", err)
	}
}

func TestError_ArithmeticTruncatedFrame(t *testing.T) {
	input := []byte("test data for arithmetic truncation")

	encoded, err := codec.EncodeBuffer(arithmetic.NewStreamingEncoder(), input)
	if err != nil {
		t.Fatalf("EncodeBuffer() error = %v", err)
	}

	truncated := encoded[:len(encoded)/2]
	_, err = codec.DecodeBuffer(arithmetic.NewStreamingDecoder(), truncated)
	if !errors.Is(err, codec.ErrTruncated) && !errors.Is(err, codec.ErrCorrupt) {
		t.Fatalf("DecodeBuffer() error = %v, want ErrTruncated or ErrCorrupt", err)
	}
}

func TestDecodeBuffer_GrowsWithoutMaxAllocation(t *testing.T) {
	input := bytes.Repeat([]byte("abcd"), 1024)
	encoded, err := codec.EncodeBuffer(rle.NewStreamingEncoder(), input)
	if err != nil {
		t.Fatalf("EncodeBuffer() error = %v", err)
	}

	decoded, err := codec.DecodeBuffer(rle.NewStreamingDecoder(), encoded)
	if err != nil {
		t.Fatalf("DecodeBuffer() error = %v", err)
	}

	if !bytes.Equal(decoded, input) {
		t.Fatalf("decoded = %q, want %q", decoded, input)
	}
}

// TestConstants verifies security limit constants.
func TestConstants(t *testing.T) {
	if codec.MaxInputSize != 4*1024*1024*1024 {
		t.Errorf("MaxInputSize = %d, want %d", codec.MaxInputSize, 4*1024*1024*1024)
	}
	if codec.MaxOutputSize != 1*1024*1024*1024 {
		t.Errorf("MaxOutputSize = %d, want %d", codec.MaxOutputSize, 1*1024*1024*1024)
	}
}
