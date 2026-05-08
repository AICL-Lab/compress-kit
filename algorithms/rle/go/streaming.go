package rle

import (
	"bytes"

	"github.com/LessUp/compress-kit/algorithms/shared/go/codec"
)

// NewStreamingEncoder creates a new streaming RLE encoder.
// It uses a buffered encoder that collects all input and encodes in one pass
// during Finish().
func NewStreamingEncoder() codec.Encoder {
	return codec.NewBufferedEncoder(rleEncode)
}

// NewStreamingDecoder creates a new streaming RLE decoder.
// It uses a buffered decoder that collects all input and decodes in one pass
// during Finish().
func NewStreamingDecoder() codec.Decoder {
	return codec.NewBufferedDecoder(rleDecode)
}

// rleEncode is the single-shot encode function for the buffered encoder.
func rleEncode(input []byte) ([]byte, error) {
	var outBuf bytes.Buffer
	err := Encode(bytes.NewReader(input), &outBuf)
	if err != nil {
		return nil, err
	}
	return outBuf.Bytes(), nil
}

// rleDecode is the single-shot decode function for the buffered decoder.
func rleDecode(input []byte) ([]byte, error) {
	var outBuf bytes.Buffer
	err := Decode(bytes.NewReader(input), &outBuf)
	if err != nil {
		return nil, err
	}
	return outBuf.Bytes(), nil
}
