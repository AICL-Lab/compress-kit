// Package codec provides shared streaming encoder/decoder interfaces for CompressKit algorithms.
package codec

import (
	"bytes"
	"errors"
)

// State represents the lifecycle state of an encoder or decoder.
type State int

const (
	// StateReady indicates the codec is ready to accept input.
	StateReady State = iota
	// StateStreaming indicates the codec is actively processing input.
	StateStreaming
	// StateFlushing indicates the codec has been flushed and buffered output has been emitted.
	StateFlushing
	// StateFinished indicates the codec has finished and emitted the end-of-stream marker.
	StateFinished
	// StateError indicates the codec encountered an error and must be reset.
	StateError
)

// EncodeFunc is the signature for a single-shot encoding function.
// It takes the complete input data and returns the complete encoded output.
type EncodeFunc func(input []byte) ([]byte, error)

// DecodeFunc is the signature for a single-shot decoding function.
// It takes the complete encoded input and returns the complete decoded output.
type DecodeFunc func(input []byte) ([]byte, error)

// BufferedEncoder is a deep module that implements the Encoder interface
// using a single-shot encode function. It handles all state machine logic,
// input buffering, and output management, allowing algorithm implementations
// to focus solely on the encoding algorithm.
//
// This is the preferred way to implement streaming encoders for algorithms
// that require complete input (Huffman, Arithmetic, Range, etc.).
type BufferedEncoder struct {
	state      State
	inputBuf   *bytes.Buffer
	totalInput int64
	encode     EncodeFunc
}

// NewBufferedEncoder creates a new streaming encoder backed by a single-shot encode function.
// The encode function will be called once during Finish() with all buffered input.
func NewBufferedEncoder(encode EncodeFunc) *BufferedEncoder {
	return &BufferedEncoder{
		state:    StateReady,
		inputBuf: &bytes.Buffer{},
		encode:   encode,
	}
}

// Process buffers input for later encoding.
func (e *BufferedEncoder) Process(in []byte, out []byte) (int, error) {
	if e.state == StateFinished {
		e.state = StateError
		return 0, ErrInvalidState
	}
	if e.state == StateError {
		return 0, ErrInvalidState
	}

	// Check input size limit
	if e.totalInput+int64(len(in)) > MaxInputSize {
		e.state = StateError
		return 0, ErrSizeLimit
	}

	// Buffer input
	e.inputBuf.Write(in)
	e.totalInput += int64(len(in))
	e.state = StateStreaming

	// No output yet - buffering until finish
	return 0, nil
}

// Flush is a no-op for buffered encoders (need complete input for encoding).
func (e *BufferedEncoder) Flush(out []byte) (int, error) {
	if e.state == StateFinished {
		e.state = StateError
		return 0, ErrInvalidState
	}
	if e.state == StateError {
		return 0, ErrInvalidState
	}

	if e.state == StateStreaming {
		e.state = StateFlushing
	}

	// No output - buffered encoders need complete input
	return 0, nil
}

// Finish encodes all buffered input and writes output.
func (e *BufferedEncoder) Finish(out []byte) (int, error) {
	if e.state == StateFinished {
		e.state = StateError
		return 0, ErrInvalidState
	}
	if e.state == StateError {
		return 0, ErrInvalidState
	}

	// Encode buffered input using the provided function
	encoded, err := e.encode(e.inputBuf.Bytes())
	if err != nil {
		e.state = StateError
		return 0, err
	}

	if len(encoded) > len(out) {
		// Buffer too small - state unchanged (transactional)
		return 0, ErrBufTooSmall
	}

	n := copy(out, encoded)
	e.state = StateFinished
	return n, nil
}

// Reset clears the encoder state.
func (e *BufferedEncoder) Reset() {
	e.state = StateReady
	e.inputBuf.Reset()
	e.totalInput = 0
}

// State returns the current lifecycle state.
func (e *BufferedEncoder) State() State {
	return e.state
}

// BufferedDecoder is a deep module that implements the Decoder interface
// using a single-shot decode function. It handles all state machine logic,
// input buffering, and output management, allowing algorithm implementations
// to focus solely on the decoding algorithm.
//
// This is the preferred way to implement streaming decoders for algorithms
// that require complete input (Huffman, Arithmetic, Range, RLE, etc.).
type BufferedDecoder struct {
	state      State
	inputBuf   *bytes.Buffer
	totalInput int64
	decode     DecodeFunc
}

// NewBufferedDecoder creates a new streaming decoder backed by a single-shot decode function.
// The decode function will be called once during Finish() with all buffered input.
func NewBufferedDecoder(decode DecodeFunc) *BufferedDecoder {
	return &BufferedDecoder{
		state:    StateReady,
		inputBuf: &bytes.Buffer{},
		decode:   decode,
	}
}

// Process buffers input for later decoding.
func (d *BufferedDecoder) Process(in []byte, out []byte) (int, error) {
	if d.state == StateFinished {
		d.state = StateError
		return 0, ErrInvalidState
	}
	if d.state == StateError {
		return 0, ErrInvalidState
	}

	// Check input size limit
	if d.totalInput+int64(len(in)) > MaxInputSize {
		d.state = StateError
		return 0, ErrSizeLimit
	}

	// Buffer input
	d.inputBuf.Write(in)
	d.totalInput += int64(len(in))
	d.state = StateStreaming

	// No output yet - need complete input to decode
	return 0, nil
}

// Flush is a no-op for buffered decoders.
func (d *BufferedDecoder) Flush(out []byte) (int, error) {
	if d.state == StateFinished {
		d.state = StateError
		return 0, ErrInvalidState
	}
	if d.state == StateError {
		return 0, ErrInvalidState
	}

	if d.state == StateStreaming {
		d.state = StateFlushing
	}

	return 0, nil
}

// Finish decodes all buffered input and writes output.
func (d *BufferedDecoder) Finish(out []byte) (int, error) {
	if d.state == StateFinished {
		d.state = StateError
		return 0, ErrInvalidState
	}
	if d.state == StateError {
		return 0, ErrInvalidState
	}

	// Decode buffered input using the provided function
	decoded, err := d.decode(d.inputBuf.Bytes())
	if err != nil {
		d.state = StateError
		// Structured errors (CodecError) are returned directly
		// to preserve their semantic kind
		return 0, err
	}

	// Check output size limit
	if len(decoded) > MaxOutputSize {
		d.state = StateError
		return 0, ErrSizeLimit
	}

	if len(decoded) > len(out) {
		// Buffer too small - state unchanged
		return 0, ErrBufTooSmall
	}

	n := copy(out, decoded)
	d.state = StateFinished
	return n, nil
}

// Reset clears the decoder state.
func (d *BufferedDecoder) Reset() {
	d.state = StateReady
	d.inputBuf.Reset()
	d.totalInput = 0
}

// State returns the current lifecycle state.
func (d *BufferedDecoder) State() State {
	return d.state
}

// Encoder defines the streaming encoder interface.
// All CompressKit algorithms implement this interface.
//
// Lifecycle: READY → STREAMING → FLUSHING → FINISHED
// Error transitions to ERROR state from any state.
// Reset() returns to READY from any state.
//
// Thread safety: Encoders are NOT thread-safe. Callers must manage concurrency.
type Encoder interface {
	// Process encodes input bytes and writes to the output buffer.
	// Returns the number of bytes written to out.
	// If out is too small, returns ErrBufTooSmall and state is unchanged (transactional).
	//
	// State transitions:
	//   READY → STREAMING
	//   STREAMING → STREAMING
	//   FLUSHING → STREAMING
	Process(in []byte, out []byte) (written int, err error)

	// Flush writes all buffered output to the output buffer.
	// Returns the number of bytes written.
	// If out is too small, returns ErrBufTooSmall and state is unchanged.
	//
	// State transitions:
	//   READY → READY (no-op)
	//   STREAMING → FLUSHING
	//   FLUSHING → FLUSHING (idempotent)
	Flush(out []byte) (written int, err error)

	// Finish flushes any remaining buffered data and writes the end-of-stream marker.
	// Returns the number of bytes written.
	// If out is too small, returns ErrBufTooSmall and state is unchanged.
	//
	// State transitions:
	//   READY/STREAMING/FLUSHING → FINISHED
	//   FINISHED → ERROR (invalid)
	Finish(out []byte) (written int, err error)

	// Reset returns the encoder to READY state, clearing all internal buffers.
	// Can be called from any state.
	Reset()

	// State returns the current lifecycle state.
	State() State
}

// Decoder defines the streaming decoder interface.
type Decoder interface {
	// Process decodes input bytes and writes to the output buffer.
	// Returns the number of bytes written to out.
	// If out is too small, returns ErrBufTooSmall and state is unchanged.
	//
	// State transitions: same as Encoder.Process
	Process(in []byte, out []byte) (written int, err error)

	// Flush writes all buffered output to the output buffer.
	// State transitions: same as Encoder.Flush
	Flush(out []byte) (written int, err error)

	// Finish completes decoding and verifies the end-of-stream marker.
	// Returns ErrTruncated if input is incomplete.
	// State transitions: same as Encoder.Finish
	Finish(out []byte) (written int, err error)

	// Reset returns the decoder to READY state.
	Reset()

	// State returns the current lifecycle state.
	State() State
}

// IsCodecError checks if an error is a CodecError with a specific kind.
// This is useful for testing and error handling without importing the errors package.
func IsCodecError(err error, kind ErrorKind) bool {
	var codecErr *CodecError
	if errors.As(err, &codecErr) {
		return codecErr.Kind == kind
	}
	return false
}
