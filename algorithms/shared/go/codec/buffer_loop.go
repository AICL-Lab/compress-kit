package codec

// bufferStep appends output into the provided slice window.
// If it returns ErrBufTooSmall, the reported byte count must already be copied
// into out so runBufferStep can preserve that partial output before retrying.
type bufferStep func(out []byte) (int, error)

// runBufferStep owns the buffer-layer retry policy for BUF_TOO_SMALL handling.
// It grows the caller-owned output buffer, preserves transactional partial
// writes reported with ErrBufTooSmall, and stops at ErrSizeLimit.
func runBufferStep(outBuf []byte, totalWritten int, limit int, step bufferStep) ([]byte, int, error) {
	for {
		n, err := step(outBuf[totalWritten:])
		if err != ErrBufTooSmall {
			if err != nil {
				return nil, totalWritten, err
			}
			return outBuf, totalWritten + n, nil
		}

		totalWritten += n
		if totalWritten > limit || len(outBuf) >= limit {
			return nil, totalWritten, ErrSizeLimit
		}

		newSize := growBuffer(len(outBuf), limit)
		if newSize <= len(outBuf) {
			return nil, totalWritten, ErrSizeLimit
		}

		newBuf := make([]byte, newSize)
		copy(newBuf, outBuf[:totalWritten])
		outBuf = newBuf
	}
}
