package codec

const encodeBufferInitialSlack = 2048 // Extra room for small outputs before the first retry.

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
	if inputLen < 0 {
		return 0, ErrSizeLimit
	}
	if inputLen > (int(^uint(0)>>1)-encodeBufferInitialSlack)/8 {
		return 0, ErrSizeLimit
	}
	return inputLen*8 + encodeBufferInitialSlack, nil
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

	return encodeBufferWithLimit(encoder, input, len(input)*2+encodeBufferInitialSlack, encodeLimit)
}

func encodeBufferWithLimit(encoder Encoder, input []byte, initialSize int, limit int) ([]byte, error) {
	if initialSize > limit {
		initialSize = limit
	}

	outBuf := make([]byte, initialSize)
	totalWritten := 0

	var err error
	outBuf, totalWritten, err = runBufferStep(outBuf, totalWritten, limit, func(out []byte) (int, error) {
		return encoder.Process(input, out)
	})
	if err != nil {
		return nil, err
	}

	outBuf, totalWritten, err = runBufferStep(outBuf, totalWritten, limit, func(out []byte) (int, error) {
		return encoder.Finish(out)
	})
	if err != nil {
		return nil, err
	}

	if totalWritten > limit {
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

	return decodeBufferWithLimit(decoder, input, len(input)+1024, MaxOutputSize)
}

func decodeBufferWithLimit(decoder Decoder, input []byte, initialSize int, limit int) ([]byte, error) {
	if initialSize > limit {
		initialSize = limit
	}

	outBuf := make([]byte, initialSize)
	totalWritten := 0

	var err error
	outBuf, totalWritten, err = runBufferStep(outBuf, totalWritten, limit, func(out []byte) (int, error) {
		return decoder.Process(input, out)
	})
	if err != nil {
		return nil, err
	}

	outBuf, totalWritten, err = runBufferStep(outBuf, totalWritten, limit, func(out []byte) (int, error) {
		return decoder.Finish(out)
	})
	if err != nil {
		return nil, err
	}

	if totalWritten > limit {
		return nil, ErrSizeLimit
	}

	return outBuf[:totalWritten], nil
}
