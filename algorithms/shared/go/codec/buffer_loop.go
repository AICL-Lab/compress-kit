package codec

type bufferStep func(out []byte) (int, error)

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
