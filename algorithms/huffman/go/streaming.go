package huffman

import (
	"bytes"

	"github.com/LessUp/compress-kit/algorithms/shared/go/codec"
)

// NewStreamingEncoder creates a new streaming Huffman encoder.
// It uses a buffered encoder that collects all input and encodes in one pass
// during Finish(), since Huffman encoding requires complete input for frequency analysis.
func NewStreamingEncoder() codec.Encoder {
	return codec.NewBufferedEncoder(huffmanEncode)
}

// NewStreamingDecoder creates a new streaming Huffman decoder.
// It uses a buffered decoder that collects all input and decodes in one pass
// during Finish().
func NewStreamingDecoder() codec.Decoder {
	return codec.NewBufferedDecoder(huffmanDecode)
}

// huffmanEncode is the single-shot encode function for the buffered encoder.
func huffmanEncode(input []byte) ([]byte, error) {
	var outBuf bytes.Buffer
	err := Encode(bytes.NewReader(input), &outBuf)
	if err != nil {
		return nil, err
	}
	return outBuf.Bytes(), nil
}

// huffmanDecode is the single-shot decode function for the buffered decoder.
func huffmanDecode(input []byte) ([]byte, error) {
	var outBuf bytes.Buffer
	err := Decode(bytes.NewReader(input), &outBuf)
	if err != nil {
		return nil, err
	}
	return outBuf.Bytes(), nil
}
