package codec

import "errors"

// StatusCode represents codec operation result codes.
type StatusCode int

const (
	// StatusOK indicates success.
	StatusOK StatusCode = iota
	// StatusBufTooSmall indicates the output buffer is too small.
	StatusBufTooSmall
	// StatusTruncated indicates input stream ends prematurely.
	StatusTruncated
	// StatusCorrupt indicates data corruption or integrity check failure.
	StatusCorrupt
	// StatusInvalidState indicates the operation is not valid in the current state.
	StatusInvalidState
	// StatusSizeLimit indicates input or output exceeds security limits.
	StatusSizeLimit
	// StatusVersionUnsupported indicates frame version is not supported.
	StatusVersionUnsupported
	// StatusUnknownAlgo indicates unknown algorithm ID in frame header.
	StatusUnknownAlgo
)

// Error sentinel values for common codec errors.
var (
	// ErrBufTooSmall indicates the caller-supplied output buffer is too small.
	// The operation is transactional: internal state is unchanged, caller may retry with larger buffer.
	ErrBufTooSmall = errors.New("output buffer too small")

	// ErrTruncated indicates the input stream ends prematurely during decode.
	ErrTruncated = errors.New("input stream truncated")

	// ErrCorrupt indicates checksum or structural integrity check failed.
	ErrCorrupt = errors.New("data corrupted")

	// ErrInvalidState indicates the call is not valid in the current lifecycle state.
	ErrInvalidState = errors.New("invalid state for operation")

	// ErrSizeLimit indicates input or output exceeds security limits.
	// Input limit: 4 GiB. Decode output limit: 1 GiB.
	ErrSizeLimit = errors.New("size limit exceeded")

	// ErrVersionUnsupported indicates frame version byte is not supported.
	ErrVersionUnsupported = errors.New("unsupported version")

	// ErrUnknownAlgo indicates unknown algorithm ID in frame header.
	ErrUnknownAlgo = errors.New("unknown algorithm")
)

// Security limits
const (
	// MaxInputSize is the maximum allowed input size (4 GiB).
	MaxInputSize = 4 * 1024 * 1024 * 1024

	// MaxOutputSize is the maximum allowed decode output size (1 GiB).
	MaxOutputSize = 1 * 1024 * 1024 * 1024
)
