package codec

import (
	"errors"
	"fmt"
)

// ErrorKind represents the category of a codec error.
// Used for structured error reporting across algorithms.
type ErrorKind int

const (
	// KindBufTooSmall indicates the caller-supplied output buffer is too small.
	KindBufTooSmall ErrorKind = iota
	// KindTruncated indicates the input stream ends prematurely.
	KindTruncated
	// KindCorrupt indicates data corruption or integrity check failure.
	KindCorrupt
	// KindInvalidState indicates the call is not valid in the current lifecycle state.
	KindInvalidState
	// KindSizeLimit indicates input or output exceeds security limits.
	KindSizeLimit
	// KindVersionUnsupported indicates frame version is not supported.
	KindVersionUnsupported
	// KindUnknownAlgo indicates unknown algorithm ID in frame header.
	KindUnknownAlgo
	// KindIO indicates an I/O error (read/write failure).
	KindIO
)

// CodecError is a structured error with a semantic kind and optional cause.
// Algorithms should return CodecError to enable precise error mapping
// without relying on string comparison.
type CodecError struct {
	// Kind is the semantic category of the error.
	Kind ErrorKind
	// Message is a human-readable description.
	Message string
	// Cause is the underlying error, if any.
	Cause error
}

// Error implements the error interface.
func (e *CodecError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

// Unwrap returns the underlying cause for errors.Is/errors.As compatibility.
func (e *CodecError) Unwrap() error {
	return e.Cause
}

// Is enables errors.Is comparison with sentinel errors.
func (e *CodecError) Is(target error) bool {
	switch e.Kind {
	case KindBufTooSmall:
		return target == ErrBufTooSmall
	case KindTruncated:
		return target == ErrTruncated
	case KindCorrupt:
		return target == ErrCorrupt
	case KindInvalidState:
		return target == ErrInvalidState
	case KindSizeLimit:
		return target == ErrSizeLimit
	case KindVersionUnsupported:
		return target == ErrVersionUnsupported
	case KindUnknownAlgo:
		return target == ErrUnknownAlgo
	default:
		return false
	}
}

// NewError creates a new CodecError with the given kind and message.
func NewError(kind ErrorKind, message string) *CodecError {
	return &CodecError{Kind: kind, Message: message}
}

// WrapError creates a new CodecError with the given kind, message, and cause.
func WrapError(kind ErrorKind, message string, cause error) *CodecError {
	return &CodecError{Kind: kind, Message: message, Cause: cause}
}

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
