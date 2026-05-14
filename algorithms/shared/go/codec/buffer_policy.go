package codec

import "errors"

type bufferStep func(out []byte) (int, error)

type resizingBuffer struct {
	buf     []byte
	written int
	limit   int
}

func newResizingBuffer(initialSize int, limit int) *resizingBuffer {
	if initialSize > limit {
		initialSize = limit
	}
	return &resizingBuffer{
		buf:   make([]byte, initialSize),
		limit: limit,
	}
}

func (r *resizingBuffer) run(step bufferStep) error {
	for {
		n, err := step(r.buf[r.written:])
		if !errors.Is(err, ErrBufTooSmall) {
			if err != nil {
				return err
			}
			r.written += n
			if r.written > r.limit {
				return ErrSizeLimit
			}
			return nil
		}

		r.written += n
		if r.written > r.limit || len(r.buf) >= r.limit {
			return ErrSizeLimit
		}

		newSize := growBuffer(len(r.buf), r.limit)
		if newSize <= len(r.buf) {
			return ErrSizeLimit
		}

		next := make([]byte, newSize)
		copy(next, r.buf[:r.written])
		r.buf = next
	}
}

func (r *resizingBuffer) bytes() []byte {
	return r.buf[:r.written]
}

func runBufferStep(outBuf []byte, totalWritten int, limit int, step bufferStep) ([]byte, int, error) {
	runner := &resizingBuffer{
		buf:     outBuf,
		written: totalWritten,
		limit:   limit,
	}

	if err := runner.run(step); err != nil {
		return nil, runner.written, err
	}

	return runner.buf, runner.written, nil
}
