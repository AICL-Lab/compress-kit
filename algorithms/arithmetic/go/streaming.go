package arithmetic

import (
	"bytes"

	"github.com/LessUp/compress-kit/algorithms/shared/go/codec"
)

// NewStreamingEncoder creates a new streaming Arithmetic encoder.
// It uses a buffered encoder that collects all input and encodes in one pass
// during Finish(), since Arithmetic encoding requires complete input for frequency analysis.
func NewStreamingEncoder() codec.Encoder {
	return codec.NewBufferedEncoder(arithmeticEncode)
}

// NewStreamingDecoder creates a new streaming Arithmetic decoder.
// It uses a buffered decoder that collects all input and decodes in one pass
// during Finish().
func NewStreamingDecoder() codec.Decoder {
	return codec.NewBufferedDecoder(arithmeticDecode)
}

// arithmeticEncode is the single-shot encode function for the buffered encoder.
func arithmeticEncode(input []byte) ([]byte, error) {
	var outBuf bytes.Buffer
	err := Encode(bytes.NewReader(input), &outBuf)
	if err != nil {
		return nil, err
	}
	return outBuf.Bytes(), nil
}

// arithmeticDecode is the single-shot decode function for the buffered decoder.
func arithmeticDecode(input []byte) ([]byte, error) {
	var outBuf bytes.Buffer
	err := Decode(bytes.NewReader(input), &outBuf)
	if err != nil {
		return nil, err
	}
	return outBuf.Bytes(), nil
}
