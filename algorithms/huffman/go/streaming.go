package huffman

import (
	"github.com/LessUp/compress-kit/algorithms/shared/go/codec"
)

// NewStreamingEncoder creates a new streaming Huffman encoder.
// It uses a buffered encoder that collects all input and encodes in one pass
// during Finish(), since Huffman encoding requires complete input for frequency analysis.
func NewStreamingEncoder() codec.Encoder {
	return codec.NewStreamingEncoderFromIO(Encode)
}

// NewStreamingDecoder creates a new streaming Huffman decoder.
// It uses a buffered decoder that collects all input and decodes in one pass
// during Finish().
func NewStreamingDecoder() codec.Decoder {
	return codec.NewStreamingDecoderFromIO(Decode)
}
