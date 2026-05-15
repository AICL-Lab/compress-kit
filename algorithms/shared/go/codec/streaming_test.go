package codec

import (
	"io"
	"testing"
)

func TestNewStreamingEncoderFromIO_AdaptsReaderWriterEncode(t *testing.T) {
	enc := NewStreamingEncoderFromIO(func(r io.Reader, w io.Writer) error {
		data, err := io.ReadAll(r)
		if err != nil {
			return err
		}
		_, err = w.Write(append([]byte("io:"), data...))
		return err
	})

	out, err := EncodeBuffer(enc, []byte("abc"))
	if err != nil {
		t.Fatalf("EncodeBuffer() error = %v", err)
	}
	if string(out) != "io:abc" {
		t.Fatalf("EncodeBuffer() = %q, want %q", out, "io:abc")
	}
}

func TestNewStreamingDecoderFromBytes_UsesDecodeFunc(t *testing.T) {
	dec := NewStreamingDecoderFromBytes(func(input []byte) ([]byte, error) {
		return append([]byte("bytes:"), input...), nil
	})

	out, err := DecodeBuffer(dec, []byte("abc"))
	if err != nil {
		t.Fatalf("DecodeBuffer() error = %v", err)
	}
	if string(out) != "bytes:abc" {
		t.Fatalf("DecodeBuffer() = %q, want %q", out, "bytes:abc")
	}
}
