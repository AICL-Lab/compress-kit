package rangecoder

import (
	"github.com/LessUp/compress-kit/algorithms/shared/go/codec"
)

// NewStreamingEncoder creates a new streaming Range encoder.
// It uses a buffered encoder that collects all input and encodes in one pass
// during Finish().
func NewStreamingEncoder() codec.Encoder {
	return codec.NewBufferedEncoder(Encode)
}

// NewStreamingDecoder creates a new streaming Range decoder.
// It uses a buffered decoder that collects all input and decodes in one pass
// during Finish().
func NewStreamingDecoder() codec.Decoder {
	return codec.NewBufferedDecoder(Decode)
}
