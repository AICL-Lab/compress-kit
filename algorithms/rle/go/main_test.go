package main

import (
	"bytes"
	"encoding/binary"
	"os"
	"path/filepath"
	"testing"
)

func TestRLERoundTrip(t *testing.T) {
	data := bytes.Repeat([]byte{0xAA}, 1024)
	data = append(data, bytes.Repeat([]byte("run-length-"), 128)...)
	data = append(data, bytes.Repeat([]byte{0x00, 0xFF}, 64)...)

	tmpDir := t.TempDir()
	inputPath := filepath.Join(tmpDir, "input.bin")
	encodedPath := filepath.Join(tmpDir, "encoded.rle")
	outputPath := filepath.Join(tmpDir, "output.bin")

	if err := os.WriteFile(inputPath, data, 0o644); err != nil {
		t.Fatalf("write input: %v", err)
	}
	if err := rleEncodeFile(inputPath, encodedPath); err != nil {
		t.Fatalf("encode: %v", err)
	}
	if err := rleDecodeFile(encodedPath, outputPath); err != nil {
		t.Fatalf("decode: %v", err)
	}

	decoded, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	if !bytes.Equal(decoded, data) {
		t.Fatalf("round-trip mismatch")
	}
}

func TestRLEDecodeRejectsZeroCount(t *testing.T) {
	tmpDir := t.TempDir()
	inputPath := filepath.Join(tmpDir, "invalid.rle")
	outputPath := filepath.Join(tmpDir, "output.bin")

	buf := make([]byte, 5)
	binary.LittleEndian.PutUint32(buf[:4], 0)
	buf[4] = 0x42
	if err := os.WriteFile(inputPath, buf, 0o644); err != nil {
		t.Fatalf("write invalid input: %v", err)
	}

	if err := rleDecodeFile(inputPath, outputPath); err == nil {
		t.Fatalf("expected decode error for zero count")
	}
}
