package arithmetic

import (
	"bytes"
	"encoding/binary"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/LessUp/compress-kit/algorithms/shared/go/codec"
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
	if err := EncodeFile(inputPath, encodedPath); err != nil {
		t.Fatalf("compress: %v", err)
	}
	if err := DecodeFile(encodedPath, outputPath); err != nil {
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
	if err := EncodeFile(inputPath, encodedPath); err != nil {
		t.Fatalf("compress empty: %v", err)
	}
	if err := DecodeFile(encodedPath, outputPath); err != nil {
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
	if err := EncodeFile(inputPath, encodedPath); err != nil {
		t.Fatalf("compress: %v", err)
	}
	if err := DecodeFile(encodedPath, outputPath); err != nil {
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
	if err := EncodeFile(inputPath, encodedPath); err != nil {
		t.Fatalf("compress: %v", err)
	}
	if err := DecodeFile(encodedPath, outputPath); err != nil {
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

func TestDecodeRejectsAllZeroFrequencyTable(t *testing.T) {
	var encoded bytes.Buffer
	if _, err := encoded.Write([]byte("AENC")); err != nil {
		t.Fatalf("write magic: %v", err)
	}
	if err := binary.Write(&encoded, binary.LittleEndian, uint32(codec.SymbolLimit)); err != nil {
		t.Fatalf("write count: %v", err)
	}
	for i := 0; i < codec.SymbolLimit; i++ {
		if err := binary.Write(&encoded, binary.LittleEndian, uint32(0)); err != nil {
			t.Fatalf("write freq[%d]: %v", i, err)
		}
	}
	if _, err := encoded.Write([]byte{0xFF, 0xFF, 0xFF, 0xFF}); err != nil {
		t.Fatalf("write trailer: %v", err)
	}

	var decoded bytes.Buffer
	err := Decode(bytes.NewReader(encoded.Bytes()), &decoded)
	if err == nil {
		t.Fatal("expected decode to reject all-zero frequency table")
	}
	if !errors.Is(err, codec.ErrCorrupt) {
		t.Fatalf("expected corrupt error, got %v", err)
	}
}
