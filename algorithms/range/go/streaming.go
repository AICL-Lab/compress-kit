package rangecoder

import (
	"bytes"

	"github.com/LessUp/compress-kit/algorithms/shared/go/codec"
)

// StreamingEncoder implements codec.Encoder for Range Coding.
type StreamingEncoder struct {
	state      codec.State
	inputBuf   *bytes.Buffer
	totalInput int64
}

// NewStreamingEncoder creates a new streaming Range encoder.
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

// Flush is a no-op for Range Coding.
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

	encoded, err := Encode(e.inputBuf.Bytes())
	if err != nil {
		e.state = codec.StateError
		return 0, err
	}

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

// StreamingDecoder implements codec.Decoder for Range Coding.
type StreamingDecoder struct {
	state      codec.State
	inputBuf   *bytes.Buffer
	totalInput int64
}

// NewStreamingDecoder creates a new streaming Range decoder.
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

// Flush is a no-op for Range decoder.
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

	decoded, err := Decode(d.inputBuf.Bytes())
	if err != nil {
		d.state = codec.StateError
		if err.Error() == "range: bad magic" {
			return 0, codec.ErrCorrupt
		}
		if err.Error() == "range: input too short" || err.Error() == "range: truncated header" {
			return 0, codec.ErrTruncated
		}
		return 0, err
	}

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
