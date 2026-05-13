package codec

import (
	"os"
	"path/filepath"
	"testing"
)

// mockEncoder is a test double that records calls and returns preset output.
type mockEncoder struct {
	input  []byte
	output []byte
	err    error
}

func (m *mockEncoder) Process(in []byte, out []byte) (int, error) {
	m.input = append(m.input, in...)
	return 0, nil
}

func (m *mockEncoder) Flush(out []byte) (int, error) {
	return 0, nil
}

func (m *mockEncoder) Finish(out []byte) (int, error) {
	if m.err != nil {
		return 0, m.err
	}
	if len(m.output) > len(out) {
		return 0, ErrBufTooSmall
	}
	return copy(out, m.output), nil
}

func (m *mockEncoder) Reset() {
	m.input = nil
}

func (m *mockEncoder) State() State {
	return StateFinished
}

// mockDecoder is a test double that records calls and returns preset output.
type mockDecoder struct {
	input  []byte
	output []byte
	err    error
}

func (m *mockDecoder) Process(in []byte, out []byte) (int, error) {
	m.input = append(m.input, in...)
	return 0, nil
}

func (m *mockDecoder) Flush(out []byte) (int, error) {
	return 0, nil
}

func (m *mockDecoder) Finish(out []byte) (int, error) {
	if m.err != nil {
		return 0, m.err
	}
	if len(m.output) > len(out) {
		return 0, ErrBufTooSmall
	}
	return copy(out, m.output), nil
}

func (m *mockDecoder) Reset() {
	m.input = nil
}

func (m *mockDecoder) State() State {
	return StateFinished
}

func TestEncodeFile_Success(t *testing.T) {
	// Create temp input file
	tmpDir := t.TempDir()
	inputPath := filepath.Join(tmpDir, "input.bin")
	outputPath := filepath.Join(tmpDir, "output.bin")

	testData := []byte("hello world")
	if err := os.WriteFile(inputPath, testData, 0o644); err != nil {
		t.Fatalf("failed to create test input: %v", err)
	}

	// Use mock encoder that returns fixed output
	enc := &mockEncoder{output: []byte("encoded")}

	err := EncodeFile(enc, inputPath, outputPath)
	if err != nil {
		t.Fatalf("EncodeFile failed: %v", err)
	}

	// Verify encoder received the input
	if string(enc.input) != string(testData) {
		t.Errorf("encoder received %q, want %q", enc.input, testData)
	}

	// Verify output file
	output, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}
	if string(output) != "encoded" {
		t.Errorf("output file = %q, want %q", output, "encoded")
	}
}

func TestEncodeFile_InputNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	inputPath := filepath.Join(tmpDir, "nonexistent.bin")
	outputPath := filepath.Join(tmpDir, "output.bin")

	enc := &mockEncoder{output: []byte("encoded")}

	err := EncodeFile(enc, inputPath, outputPath)
	if err == nil {
		t.Error("expected error for nonexistent input file")
	}
}

func TestEncodeFile_OutputNotWritable(t *testing.T) {
	tmpDir := t.TempDir()
	inputPath := filepath.Join(tmpDir, "input.bin")

	if err := os.WriteFile(inputPath, []byte("test"), 0o644); err != nil {
		t.Fatalf("failed to create test input: %v", err)
	}

	// Use a directory path as output (should fail)
	outputPath := tmpDir // This is a directory, not a file

	enc := &mockEncoder{output: []byte("encoded")}

	err := EncodeFile(enc, inputPath, outputPath)
	if err == nil {
		t.Error("expected error for unwritable output path")
	}
}

func TestDecodeFile_Success(t *testing.T) {
	tmpDir := t.TempDir()
	inputPath := filepath.Join(tmpDir, "input.bin")
	outputPath := filepath.Join(tmpDir, "output.bin")

	encodedData := []byte("encoded data")
	if err := os.WriteFile(inputPath, encodedData, 0o644); err != nil {
		t.Fatalf("failed to create test input: %v", err)
	}

	dec := &mockDecoder{output: []byte("decoded")}

	err := DecodeFile(dec, inputPath, outputPath)
	if err != nil {
		t.Fatalf("DecodeFile failed: %v", err)
	}

	// Verify decoder received the input
	if string(dec.input) != string(encodedData) {
		t.Errorf("decoder received %q, want %q", dec.input, encodedData)
	}

	// Verify output file
	output, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}
	if string(output) != "decoded" {
		t.Errorf("output file = %q, want %q", output, "decoded")
	}
}

func TestDecodeFile_InputNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	inputPath := filepath.Join(tmpDir, "nonexistent.bin")
	outputPath := filepath.Join(tmpDir, "output.bin")

	dec := &mockDecoder{output: []byte("decoded")}

	err := DecodeFile(dec, inputPath, outputPath)
	if err == nil {
		t.Error("expected error for nonexistent input file")
	}
}

func TestEncodeFile_EncoderError(t *testing.T) {
	tmpDir := t.TempDir()
	inputPath := filepath.Join(tmpDir, "input.bin")
	outputPath := filepath.Join(tmpDir, "output.bin")

	if err := os.WriteFile(inputPath, []byte("test"), 0o644); err != nil {
		t.Fatalf("failed to create test input: %v", err)
	}

	enc := &mockEncoder{err: ErrCorrupt}

	err := EncodeFile(enc, inputPath, outputPath)
	if err != ErrCorrupt {
		t.Errorf("expected ErrCorrupt, got %v", err)
	}
}

func TestDecodeFile_DecoderError(t *testing.T) {
	tmpDir := t.TempDir()
	inputPath := filepath.Join(tmpDir, "input.bin")
	outputPath := filepath.Join(tmpDir, "output.bin")

	if err := os.WriteFile(inputPath, []byte("test"), 0o644); err != nil {
		t.Fatalf("failed to create test input: %v", err)
	}

	dec := &mockDecoder{err: ErrTruncated}

	err := DecodeFile(dec, inputPath, outputPath)
	if err != ErrTruncated {
		t.Errorf("expected ErrTruncated, got %v", err)
	}
}
