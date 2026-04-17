package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestCompressDecompressRoundTrip(t *testing.T) {
	data := bytes.Repeat([]byte("arithmetic-test-data-"), 256)
	data = append(data, []byte{0, 1, 2, 3, 254, 255}...)

	tmpDir := t.TempDir()
	inputPath := filepath.Join(tmpDir, "input.bin")
	encodedPath := filepath.Join(tmpDir, "encoded.aenc")
	outputPath := filepath.Join(tmpDir, "output.bin")

	if err := os.WriteFile(inputPath, data, 0o644); err != nil {
		t.Fatalf("write input: %v", err)
	}
	if err := compressFile(inputPath, encodedPath); err != nil {
		t.Fatalf("compress: %v", err)
	}
	if err := decompressFile(encodedPath, outputPath); err != nil {
		t.Fatalf("decompress: %v", err)
	}

	decoded, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	if !bytes.Equal(decoded, data) {
		t.Fatalf("round-trip mismatch: got %d bytes, want %d bytes", len(decoded), len(data))
	}
}

func TestCompressDecompressEmptyFile(t *testing.T) {
	tmpDir := t.TempDir()
	inputPath := filepath.Join(tmpDir, "empty.bin")
	encodedPath := filepath.Join(tmpDir, "empty.aenc")
	outputPath := filepath.Join(tmpDir, "empty.out")

	if err := os.WriteFile(inputPath, nil, 0o644); err != nil {
		t.Fatalf("write input: %v", err)
	}
	if err := compressFile(inputPath, encodedPath); err != nil {
		t.Fatalf("compress empty: %v", err)
	}
	if err := decompressFile(encodedPath, outputPath); err != nil {
		t.Fatalf("decompress empty: %v", err)
	}

	decoded, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	if len(decoded) != 0 {
		t.Fatalf("expected empty output, got %d bytes", len(decoded))
	}
}

func TestCompressDecompressSingleByte(t *testing.T) {
	tmpDir := t.TempDir()
	inputPath := filepath.Join(tmpDir, "single.bin")
	encodedPath := filepath.Join(tmpDir, "single.aenc")
	outputPath := filepath.Join(tmpDir, "single.out")

	data := []byte{0x42}
	if err := os.WriteFile(inputPath, data, 0o644); err != nil {
		t.Fatalf("write input: %v", err)
	}
	if err := compressFile(inputPath, encodedPath); err != nil {
		t.Fatalf("compress: %v", err)
	}
	if err := decompressFile(encodedPath, outputPath); err != nil {
		t.Fatalf("decompress: %v", err)
	}

	decoded, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	if !bytes.Equal(decoded, data) {
		t.Fatalf("round-trip mismatch")
	}
}

func TestCompressDecompressAllBytes(t *testing.T) {
	tmpDir := t.TempDir()
	inputPath := filepath.Join(tmpDir, "allbytes.bin")
	encodedPath := filepath.Join(tmpDir, "allbytes.aenc")
	outputPath := filepath.Join(tmpDir, "allbytes.out")

	data := make([]byte, 256)
	for i := range data {
		data[i] = byte(i)
	}
	if err := os.WriteFile(inputPath, data, 0o644); err != nil {
		t.Fatalf("write input: %v", err)
	}
	if err := compressFile(inputPath, encodedPath); err != nil {
		t.Fatalf("compress: %v", err)
	}
	if err := decompressFile(encodedPath, outputPath); err != nil {
		t.Fatalf("decompress: %v", err)
	}

	decoded, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	if !bytes.Equal(decoded, data) {
		t.Fatalf("round-trip mismatch")
	}
}
