package codec

import (
	"bytes"
	"testing"
)

func TestWriterEncoder_ReturnsSizeLimitAtCap(t *testing.T) {
	stub := &capRetryEncoder{}
	var sink bytes.Buffer
	writer := NewWriterEncoder(stub, &sink)
	writer.outBuf = make([]byte, MaxOutputSize)

	err := writer.Close()
	if err != ErrSizeLimit {
		t.Fatalf("Close() error = %v, want ErrSizeLimit", err)
	}
}

func TestWriterEncoder_WriteReturnsSizeLimitAtCap(t *testing.T) {
	stub := &capRetryEncoder{}
	var sink bytes.Buffer
	writer := NewWriterEncoder(stub, &sink)
	writer.outBuf = make([]byte, MaxOutputSize)

	_, err := writer.Write([]byte("payload"))
	if err != ErrSizeLimit {
		t.Fatalf("Write() error = %v, want ErrSizeLimit", err)
	}
}
