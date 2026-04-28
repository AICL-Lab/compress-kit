package codec

func growBuffer(currentLen int, limit int) int {
	if currentLen <= 0 {
		if limit < 1024 {
			return limit
		}
		return 1024
	}
	next := currentLen * 2
	if next < currentLen {
		return limit
	}
	if next > limit {
		return limit
	}
	return next
}

func encodeBufferLimit(inputLen int) (int, error) {
	const overhead = 2048
	if inputLen < 0 {
		return 0, ErrSizeLimit
	}
	if inputLen > (int(^uint(0)>>1)-overhead)/8 {
		return 0, ErrSizeLimit
	}
	return inputLen*8 + overhead, nil
}

// EncodeBuffer is a convenience function that encodes input using the streaming API.
// Equivalent to: new encoder → Process(input) → Finish() → collect output.
//
// Returns the complete encoded output or an error.
func EncodeBuffer(encoder Encoder, input []byte) ([]byte, error) {
	if len(input) > MaxInputSize {
		return nil, ErrSizeLimit
	}

	encodeLimit, err := encodeBufferLimit(len(input))
	if err != nil {
		return nil, err
	}

	// Allocate output buffer using a conservative estimate.
	// Most algorithms need header + encoded data.
	// Use 2x input size + 2KB overhead as a reasonable initial allocation.
	initialSize := len(input)*2 + 2048
	if initialSize > encodeLimit {
		initialSize = encodeLimit
	}
	outBuf := make([]byte, initialSize)
	var totalWritten int

	var n int
	for {
		n, err = encoder.Process(input, outBuf[totalWritten:])
		if err != ErrBufTooSmall {
			break
		}
		totalWritten += n
		if totalWritten > encodeLimit {
			return nil, ErrSizeLimit
		}
		if len(outBuf) >= encodeLimit {
			return nil, ErrSizeLimit
		}
		newSize := growBuffer(len(outBuf), encodeLimit)
		if newSize <= len(outBuf) {
			return nil, ErrSizeLimit
		}
		newBuf := make([]byte, newSize)
		copy(newBuf, outBuf[:totalWritten])
		outBuf = newBuf
	}
	if err != nil {
		return nil, err
	}
	totalWritten += n

	// Finish encoding
	for {
		n, err = encoder.Finish(outBuf[totalWritten:])
		if err != ErrBufTooSmall {
			break
		}
		if len(outBuf) >= encodeLimit {
			return nil, ErrSizeLimit
		}
		newBuf := make([]byte, growBuffer(len(outBuf), encodeLimit))
		copy(newBuf, outBuf[:totalWritten])
		outBuf = newBuf
	}
	if err != nil {
		return nil, err
	}
	totalWritten += n
	if totalWritten > encodeLimit {
		return nil, ErrSizeLimit
	}

	return outBuf[:totalWritten], nil
}

// DecodeBuffer is a convenience function that decodes input using the streaming API.
// Equivalent to: new decoder → Process(input) → Finish() → collect output.
//
// Returns the complete decoded output or an error.
func DecodeBuffer(decoder Decoder, input []byte) ([]byte, error) {
	if len(input) > MaxInputSize {
		return nil, ErrSizeLimit
	}

	// Allocate output buffer.
	// Decode typically expands, so start with input size and grow as needed.
	outBuf := make([]byte, len(input)+1024)
	var totalWritten int

	var n int
	var err error
	for {
		n, err = decoder.Process(input, outBuf[totalWritten:])
		if err != ErrBufTooSmall {
			break
		}
		totalWritten += n
		if totalWritten > MaxOutputSize {
			return nil, ErrSizeLimit
		}
		if len(outBuf) >= MaxOutputSize {
			return nil, ErrSizeLimit
		}
		newSize := growBuffer(len(outBuf), MaxOutputSize)
		if newSize <= len(outBuf) {
			return nil, ErrSizeLimit
		}
		newBuf := make([]byte, newSize)
		copy(newBuf, outBuf[:totalWritten])
		outBuf = newBuf
	}
	if err != nil {
		return nil, err
	}
	totalWritten += n

	// Finish decoding
	for {
		n, err = decoder.Finish(outBuf[totalWritten:])
		if err != ErrBufTooSmall {
			break
		}
		if len(outBuf) >= MaxOutputSize {
			return nil, ErrSizeLimit
		}
		newBuf := make([]byte, growBuffer(len(outBuf), MaxOutputSize))
		copy(newBuf, outBuf[:totalWritten])
		outBuf = newBuf
	}
	if err != nil {
		return nil, err
	}
	totalWritten += n

	if totalWritten > MaxOutputSize {
		return nil, ErrSizeLimit
	}

	return outBuf[:totalWritten], nil
}
