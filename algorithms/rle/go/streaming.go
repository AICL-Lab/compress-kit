package rle

import (
	"bytes"

	"github.com/LessUp/compress-kit/algorithms/shared/go/codec"
)

// StreamingEncoder implements codec.Encoder for RLE.
type StreamingEncoder struct {
	state      codec.State
	inputBuf   *bytes.Buffer
	totalInput int64
}

// NewStreamingEncoder creates a new streaming RLE encoder.
func NewStreamingEncoder() *StreamingEncoder {
	return &StreamingEncoder{
		state:    codec.StateReady,
		inputBuf: &bytes.Buffer{},
	}
}

// Process buffers input for later encoding.
func (e *StreamingEncoder) Process(in []byte, out []byte) (int, error) {
	if e.state == codec.StateFinished {
		return 0, codec.ErrInvalidState
	}
	if e.state == codec.StateError {
		return 0, codec.ErrInvalidState
	}

	if e.totalInput+int64(len(in)) > codec.MaxInputSize {
		e.state = codec.StateError
		return 0, codec.ErrSizeLimit
	}

	e.inputBuf.Write(in)
	e.totalInput += int64(len(in))
	e.state = codec.StateStreaming
	return 0, nil
}

// Flush is a no-op for RLE.
func (e *StreamingEncoder) Flush(out []byte) (int, error) {
	if e.state == codec.StateFinished {
		return 0, codec.ErrInvalidState
	}
	if e.state == codec.StateError {
		return 0, codec.ErrInvalidState
	}

	if e.state == codec.StateStreaming {
		e.state = codec.StateFlushing
	}
	return 0, nil
}

// Finish encodes all buffered input.
func (e *StreamingEncoder) Finish(out []byte) (int, error) {
	if e.state == codec.StateFinished {
		return 0, codec.ErrInvalidState
	}
	if e.state == codec.StateError {
		return 0, codec.ErrInvalidState
	}

	var outBuf bytes.Buffer
	err := Encode(e.inputBuf, &outBuf)
	if err != nil {
		e.state = codec.StateError
		return 0, err
	}

	encoded := outBuf.Bytes()
	if len(encoded) > len(out) {
		return 0, codec.ErrBufTooSmall
	}

	n := copy(out, encoded)
	e.state = codec.StateFinished
	return n, nil
}

// Reset clears the encoder state.
func (e *StreamingEncoder) Reset() {
	e.state = codec.StateReady
	e.inputBuf.Reset()
	e.totalInput = 0
}

// State returns the current lifecycle state.
func (e *StreamingEncoder) State() codec.State {
	return e.state
}

// StreamingDecoder implements codec.Decoder for RLE.
type StreamingDecoder struct {
	state      codec.State
	inputBuf   *bytes.Buffer
	totalInput int64
}

// NewStreamingDecoder creates a new streaming RLE decoder.
func NewStreamingDecoder() *StreamingDecoder {
	return &StreamingDecoder{
		state:    codec.StateReady,
		inputBuf: &bytes.Buffer{},
	}
}

// Process buffers input for later decoding.
func (d *StreamingDecoder) Process(in []byte, out []byte) (int, error) {
	if d.state == codec.StateFinished {
		return 0, codec.ErrInvalidState
	}
	if d.state == codec.StateError {
		return 0, codec.ErrInvalidState
	}

	if d.totalInput+int64(len(in)) > codec.MaxInputSize {
		d.state = codec.StateError
		return 0, codec.ErrSizeLimit
	}

	d.inputBuf.Write(in)
	d.totalInput += int64(len(in))
	d.state = codec.StateStreaming
	return 0, nil
}

// Flush is a no-op for RLE decoder.
func (d *StreamingDecoder) Flush(out []byte) (int, error) {
	if d.state == codec.StateFinished {
		return 0, codec.ErrInvalidState
	}
	if d.state == codec.StateError {
		return 0, codec.ErrInvalidState
	}

	if d.state == codec.StateStreaming {
		d.state = codec.StateFlushing
	}
	return 0, nil
}

// Finish decodes all buffered input.
func (d *StreamingDecoder) Finish(out []byte) (int, error) {
	if d.state == codec.StateFinished {
		return 0, codec.ErrInvalidState
	}
	if d.state == codec.StateError {
		return 0, codec.ErrInvalidState
	}

	var outBuf bytes.Buffer
	err := Decode(d.inputBuf, &outBuf)
	if err != nil {
		d.state = codec.StateError
		if err.Error() == "RLE data truncated: cannot read complete count field" ||
			err.Error() == "RLE data truncated: missing value byte" {
			return 0, codec.ErrTruncated
		}
		if err.Error() == "invalid RLE data: count should not be 0" {
			return 0, codec.ErrCorrupt
		}
		return 0, err
	}

	decoded := outBuf.Bytes()
	if len(decoded) > codec.MaxOutputSize {
		d.state = codec.StateError
		return 0, codec.ErrSizeLimit
	}

	if len(decoded) > len(out) {
		return 0, codec.ErrBufTooSmall
	}

	n := copy(out, decoded)
	d.state = codec.StateFinished
	return n, nil
}

// Reset clears the decoder state.
func (d *StreamingDecoder) Reset() {
	d.state = codec.StateReady
	d.inputBuf.Reset()
	d.totalInput = 0
}

// State returns the current lifecycle state.
func (d *StreamingDecoder) State() codec.State {
	return d.state
}
