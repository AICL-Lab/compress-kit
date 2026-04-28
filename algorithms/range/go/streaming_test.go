package rangecoder

import (
	"testing"

	"github.com/LessUp/compress-kit/algorithms/shared/go/codec"
)

func TestStreamingEncoder_InputSizeLimit(t *testing.T) {
	enc := NewStreamingEncoder()
	enc.totalInput = codec.MaxInputSize

	_, err := enc.Process([]byte{0x01}, nil)
	if err != codec.ErrSizeLimit {
		t.Fatalf("Process() error = %v, want ErrSizeLimit", err)
	}
	if enc.State() != codec.StateError {
		t.Fatalf("State = %v, want StateError", enc.State())
	}
}
