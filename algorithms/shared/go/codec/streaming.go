package codec

import (
	"bytes"
	"io"
)

type StreamFunc func(io.Reader, io.Writer) error

func NewStreamingEncoderFromIO(encode StreamFunc) Encoder {
	return NewBufferedEncoder(func(input []byte) ([]byte, error) {
		var out bytes.Buffer
		if err := encode(bytes.NewReader(input), &out); err != nil {
			return nil, err
		}
		return out.Bytes(), nil
	})
}

func NewStreamingDecoderFromIO(decode StreamFunc) Decoder {
	return NewBufferedDecoder(func(input []byte) ([]byte, error) {
		var out bytes.Buffer
		if err := decode(bytes.NewReader(input), &out); err != nil {
			return nil, err
		}
		return out.Bytes(), nil
	})
}

func NewStreamingEncoderFromBytes(encode EncodeFunc) Encoder {
	return NewBufferedEncoder(encode)
}

func NewStreamingDecoderFromBytes(decode DecodeFunc) Decoder {
	return NewBufferedDecoder(decode)
}
