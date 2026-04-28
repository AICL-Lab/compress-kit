package codec_test

import (
	"bytes"
	"testing"

	"github.com/LessUp/compress-kit/algorithms/shared/go/codec"
	"rle"
)

func TestWriterEncoder_RetriesUntilFinished(t *testing.T) {
	input := make([]byte, 40*1024)
	for i := range input {
		input[i] = byte(i % 251)
	}

	var sink bytes.Buffer
	writer := codec.NewWriterEncoder(rle.NewStreamingEncoder(), &sink)

	if _, err := writer.Write(input); err != nil {
		t.Fatalf("Write() error = %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	decoded, err := codec.DecodeBuffer(rle.NewStreamingDecoder(), sink.Bytes())
	if err != nil {
		t.Fatalf("DecodeBuffer() error = %v", err)
	}
	if !bytes.Equal(decoded, input) {
		t.Fatalf("decoded length = %d, want %d", len(decoded), len(input))
	}
}
