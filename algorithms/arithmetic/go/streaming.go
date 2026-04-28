package arithmetic

import (
	"bytes"

	"github.com/LessUp/compress-kit/algorithms/shared/go/codec"
)

// StreamingEncoder implements codec.Encoder for Arithmetic Coding.
type StreamingEncoder struct {
	state      codec.State
	inputBuf   *bytes.Buffer
	totalInput int64
}

// NewStreamingEncoder creates a new streaming Arithmetic encoder.
func NewStreamingEncoder() *StreamingEncoder {
	return &StreamingEncoder{
		state:    codec.StateReady,
		inputBuf: &bytes.Buffer{},
	}
}

// Process buffers input for later encoding.
func (e *StreamingEncoder) Process(in []byte, out []byte) (int, error) {
	if e.state == codec.StateFinished {
		e.state = codec.StateError
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

// Flush is a no-op for Arithmetic Coding.
func (e *StreamingEncoder) Flush(out []byte) (int, error) {
	if e.state == codec.StateFinished {
		e.state = codec.StateError
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
		e.state = codec.StateError
		return 0, codec.ErrInvalidState
	}
	if e.state == codec.StateError {
		return 0, codec.ErrInvalidState
	}

	var outBuf bytes.Buffer
	err := Encode(bytes.NewReader(e.inputBuf.Bytes()), &outBuf)
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

// StreamingDecoder implements codec.Decoder for Arithmetic Coding.
type StreamingDecoder struct {
	state      codec.State
	inputBuf   *bytes.Buffer
	totalInput int64
}

// NewStreamingDecoder creates a new streaming Arithmetic decoder.
func NewStreamingDecoder() *StreamingDecoder {
	return &StreamingDecoder{
		state:    codec.StateReady,
		inputBuf: &bytes.Buffer{},
	}
}

// Process buffers input for later decoding.
func (d *StreamingDecoder) Process(in []byte, out []byte) (int, error) {
	if d.state == codec.StateFinished {
		d.state = codec.StateError
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

// Flush is a no-op for Arithmetic decoder.
func (d *StreamingDecoder) Flush(out []byte) (int, error) {
	if d.state == codec.StateFinished {
		d.state = codec.StateError
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
		d.state = codec.StateError
		return 0, codec.ErrInvalidState
	}
	if d.state == codec.StateError {
		return 0, codec.ErrInvalidState
	}

	var outBuf bytes.Buffer
	err := Decode(bytes.NewReader(d.inputBuf.Bytes()), &outBuf)
	if err != nil {
		d.state = codec.StateError
		errStr := err.Error()
		if err.Error() == "invalid input file format" {
			return 0, codec.ErrCorrupt
		}
		if errStr == "failed to read frequency table: unexpected EOF" ||
			errStr == "failed to read frequency table: EOF" {
			return 0, codec.ErrTruncated
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
