package codec

import "io"

// WriterEncoder wraps an Encoder to implement io.Writer.
// Calls to Write() are forwarded to encoder.Process().
// Call Close() to invoke encoder.Finish() and flush final output.
type WriterEncoder struct {
	encoder      Encoder
	w            io.Writer
	outBuf       []byte
	outBufOffset int
}

// NewWriterEncoder creates a WriterEncoder that wraps the given encoder
// and writes encoded output to w.
func NewWriterEncoder(encoder Encoder, w io.Writer) *WriterEncoder {
	return &WriterEncoder{
		encoder: encoder,
		w:       w,
		outBuf:  make([]byte, 64*1024), // 64KB output buffer
	}
}

// Write implements io.Writer by encoding p and writing to the underlying writer.
func (we *WriterEncoder) Write(p []byte) (n int, err error) {
	// Process input through encoder
	var written int
	written, err = we.encoder.Process(p, we.outBuf[we.outBufOffset:])
	for err == ErrBufTooSmall {
		if err := we.flushOutBuf(); err != nil {
			return 0, err
		}
		we.outBuf = make([]byte, growBuffer(len(we.outBuf), MaxOutputSize))
		we.outBufOffset = 0
		written, err = we.encoder.Process(p, we.outBuf)
	}
	if err != nil {
		return 0, err
	}

	we.outBufOffset += written

	// Flush output buffer if getting full
	if we.outBufOffset > len(we.outBuf)/2 {
		if err := we.flushOutBuf(); err != nil {
			return 0, err
		}
	}

	return len(p), nil
}

// Close finishes encoding and flushes all output.
func (we *WriterEncoder) Close() error {
	// Finish encoding
	written, err := we.encoder.Finish(we.outBuf[we.outBufOffset:])
	for err == ErrBufTooSmall {
		if err := we.flushOutBuf(); err != nil {
			return err
		}
		we.outBuf = make([]byte, growBuffer(len(we.outBuf), MaxOutputSize))
		we.outBufOffset = 0
		written, err = we.encoder.Finish(we.outBuf)
	}
	if err != nil {
		return err
	}
	we.outBufOffset += written

	// Flush remaining output
	return we.flushOutBuf()
}

func (we *WriterEncoder) flushOutBuf() error {
	if we.outBufOffset > 0 {
		_, err := we.w.Write(we.outBuf[:we.outBufOffset])
		if err != nil {
			return err
		}
		we.outBufOffset = 0
	}
	return nil
}
