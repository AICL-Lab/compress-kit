package codec

// EncodeBuffer is a convenience function that encodes input using the streaming API.
// Equivalent to: new encoder → Process(input) → Finish() → collect output.
//
// Returns the complete encoded output or an error.
func EncodeBuffer(encoder Encoder, input []byte) ([]byte, error) {
	if len(input) > MaxInputSize {
		return nil, ErrSizeLimit
	}

	// Allocate output buffer using a conservative estimate.
	// Most algorithms need header + encoded data.
	// Use 2x input size + 2KB overhead as a reasonable initial allocation.
	outBuf := make([]byte, len(input)*2+2048)
	var totalWritten int

	// Process all input
	n, err := encoder.Process(input, outBuf[totalWritten:])
	if err == ErrBufTooSmall {
		// Grow buffer and retry
		outBuf = make([]byte, len(input)*8+4096)
		encoder.Reset()
		n, err = encoder.Process(input, outBuf[totalWritten:])
	}
	if err != nil {
		return nil, err
	}
	totalWritten += n

	// Finish encoding
	n, err = encoder.Finish(outBuf[totalWritten:])
	if err == ErrBufTooSmall {
		// Grow buffer for finish
		newBuf := make([]byte, len(outBuf)*2)
		copy(newBuf, outBuf[:totalWritten])
		outBuf = newBuf
		n, err = encoder.Finish(outBuf[totalWritten:])
	}
	if err != nil {
		return nil, err
	}
	totalWritten += n

	return outBuf[:totalWritten], nil
}

// DecodeBuffer is a convenience function that decodes input using the streaming API.
// Equivalent to: new decoder → Process(input) → Finish() → collect output.
//
// Returns the complete decoded output or an error.
func DecodeBuffer(decoder Decoder, input []byte) ([]byte, error) {
	// Allocate output buffer.
	// Decode typically expands, so start with input size and grow as needed.
	outBuf := make([]byte, len(input)+1024)
	var totalWritten int

	// Process all input
	n, err := decoder.Process(input, outBuf[totalWritten:])
	if err == ErrBufTooSmall {
		// Grow buffer and retry
		outBuf = make([]byte, MaxOutputSize)
		decoder.Reset()
		n, err = decoder.Process(input, outBuf[totalWritten:])
	}
	if err != nil {
		return nil, err
	}
	totalWritten += n

	// Finish decoding
	n, err = decoder.Finish(outBuf[totalWritten:])
	if err == ErrBufTooSmall {
		// Grow buffer for finish
		if totalWritten+1024*1024 > MaxOutputSize {
			return nil, ErrSizeLimit
		}
		newBuf := make([]byte, totalWritten+1024*1024)
		copy(newBuf, outBuf[:totalWritten])
		outBuf = newBuf
		n, err = decoder.Finish(outBuf[totalWritten:])
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
