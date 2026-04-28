// Package codec provides shared streaming encoder/decoder interfaces for CompressKit algorithms.
package codec

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
