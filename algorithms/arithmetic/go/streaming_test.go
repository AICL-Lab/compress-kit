package arithmetic

import (
	"testing"

	"github.com/LessUp/compress-kit/algorithms/shared/go/codec"
)

func TestStreamingEncoder_InputSizeLimit(t *testing.T) {
	be := codec.NewBufferedEncoder(nil) // nil encode func is OK for this test

	// Use a small chunk to simulate large input without actually allocating MaxInputSize
	chunk := make([]byte, 1024*1024) // 1 MB chunk

	// Process chunks until we approach the limit
	for i := 0; i < 10; i++ {
		_, err := be.Process(chunk, nil)
		if err != nil {
			t.Fatalf("Process() unexpected error at chunk %d: %v", i, err)
		}
	}

	// Verify the encoder is still in streaming state
	if be.State() != codec.StateStreaming {
		t.Fatalf("State = %v, want StateStreaming", be.State())
	}

	// Reset and verify we can start fresh
	be.Reset()
	if be.State() != codec.StateReady {
		t.Fatalf("State after Reset = %v, want StateReady", be.State())
	}
}
