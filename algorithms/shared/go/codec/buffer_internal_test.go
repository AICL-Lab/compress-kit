package codec

import (
	"errors"
	"testing"
)

type processRetryEncoder struct {
	processCalls int
}

func (e *processRetryEncoder) Process(in []byte, out []byte) (int, error) {
	e.processCalls++
	if e.processCalls == 1 {
		copy(out, []byte("abc"))
		return 3, ErrBufTooSmall
	}
	copy(out, []byte("def"))
	return 3, nil
}

func (e *processRetryEncoder) Flush(out []byte) (int, error)  { return 0, nil }
func (e *processRetryEncoder) Finish(out []byte) (int, error) { return 0, nil }
func (e *processRetryEncoder) Reset()                         {}
func (e *processRetryEncoder) State() State                   { return StateStreaming }

type processRetryDecoder struct {
	processCalls int
}

func (d *processRetryDecoder) Process(in []byte, out []byte) (int, error) {
	d.processCalls++
	if d.processCalls == 1 {
		copy(out, []byte("ghi"))
		return 3, ErrBufTooSmall
	}
	copy(out, []byte("jkl"))
	return 3, nil
}

func (d *processRetryDecoder) Flush(out []byte) (int, error)  { return 0, nil }
func (d *processRetryDecoder) Finish(out []byte) (int, error) { return 0, nil }
func (d *processRetryDecoder) Reset()                         {}
func (d *processRetryDecoder) State() State                   { return StateStreaming }

func TestEncodeBuffer_PreservesOutputAcrossProcessRetry(t *testing.T) {
	out, err := EncodeBuffer(&processRetryEncoder{}, []byte("ignored"))
	if err != nil {
		t.Fatalf("EncodeBuffer() error = %v", err)
	}
	if string(out) != "abcdef" {
		t.Fatalf("EncodeBuffer() = %q, want %q", out, "abcdef")
	}
}

func TestDecodeBuffer_PreservesOutputAcrossProcessRetry(t *testing.T) {
	out, err := DecodeBuffer(&processRetryDecoder{}, []byte("ignored"))
	if err != nil {
		t.Fatalf("DecodeBuffer() error = %v", err)
	}
	if string(out) != "ghijkl" {
		t.Fatalf("DecodeBuffer() = %q, want %q", out, "ghijkl")
	}
}

type capRetryEncoder struct {
	processCalls int
	finishCalls  int
	stopErr      error
}

func (e *capRetryEncoder) Process(in []byte, out []byte) (int, error) {
	e.processCalls++
	if e.processCalls > 2 {
		if e.stopErr == nil {
			e.stopErr = errors.New("unexpected extra process retry")
		}
		return 0, e.stopErr
	}
	return 0, ErrBufTooSmall
}
func (e *capRetryEncoder) Flush(out []byte) (int, error) { return 0, nil }
func (e *capRetryEncoder) Reset()                        {}
func (e *capRetryEncoder) State() State                  { return StateStreaming }

func (e *capRetryEncoder) Finish(out []byte) (int, error) {
	e.finishCalls++
	if e.finishCalls > 2 {
		if e.stopErr == nil {
			e.stopErr = errors.New("unexpected extra retry")
		}
		return 0, e.stopErr
	}
	return 0, ErrBufTooSmall
}

type scriptedCall struct {
	written int
	err     error
	payload []byte
}

type scriptedEncoder struct {
	process []scriptedCall
	finish  []scriptedCall
}

func (e *scriptedEncoder) Process(in []byte, out []byte) (int, error) {
	call := e.process[0]
	e.process = e.process[1:]
	copy(out, call.payload)
	return call.written, call.err
}

func (e *scriptedEncoder) Flush(out []byte) (int, error) { return 0, nil }
func (e *scriptedEncoder) Reset()                        {}
func (e *scriptedEncoder) State() State                  { return StateStreaming }

func (e *scriptedEncoder) Finish(out []byte) (int, error) {
	call := e.finish[0]
	e.finish = e.finish[1:]
	copy(out, call.payload)
	return call.written, call.err
}

type scriptedDecoder struct {
	process []scriptedCall
	finish  []scriptedCall
}

func (d *scriptedDecoder) Process(in []byte, out []byte) (int, error) {
	call := d.process[0]
	d.process = d.process[1:]
	copy(out, call.payload)
	return call.written, call.err
}

func (d *scriptedDecoder) Flush(out []byte) (int, error) { return 0, nil }
func (d *scriptedDecoder) Reset()                        {}
func (d *scriptedDecoder) State() State                  { return StateStreaming }

func (d *scriptedDecoder) Finish(out []byte) (int, error) {
	call := d.finish[0]
	d.finish = d.finish[1:]
	copy(out, call.payload)
	return call.written, call.err
}

func TestEncodeBuffer_PreservesOutputAcrossFinishRetry(t *testing.T) {
	stub := &scriptedEncoder{
		process: []scriptedCall{{written: 0, err: nil}},
		finish: []scriptedCall{
			{written: 3, err: ErrBufTooSmall, payload: []byte("abc")},
			{written: 3, err: nil, payload: []byte("def")},
		},
	}

	out, err := EncodeBuffer(stub, []byte("ignored"))
	if err != nil {
		t.Fatalf("EncodeBuffer() error = %v", err)
	}
	if string(out) != "abcdef" {
		t.Fatalf("EncodeBuffer() = %q, want %q", out, "abcdef")
	}
}

func TestEncodeBuffer_ReturnsSizeLimitWhenGrowthStops(t *testing.T) {
	// Empty input gives encodeLimit == initialSize == 2048, so growth stops immediately.
	_, err := EncodeBuffer(&capRetryEncoder{}, []byte{})
	if err != ErrSizeLimit {
		t.Fatalf("EncodeBuffer() error = %v, want ErrSizeLimit", err)
	}
}

func TestDecodeBuffer_PreservesOutputAcrossFinishRetry(t *testing.T) {
	stub := &scriptedDecoder{
		process: []scriptedCall{{written: 0, err: nil}},
		finish: []scriptedCall{
			{written: 3, err: ErrBufTooSmall, payload: []byte("ghi")},
			{written: 3, err: nil, payload: []byte("jkl")},
		},
	}

	out, err := DecodeBuffer(stub, []byte("ignored"))
	if err != nil {
		t.Fatalf("DecodeBuffer() error = %v", err)
	}
	if string(out) != "ghijkl" {
		t.Fatalf("DecodeBuffer() = %q, want %q", out, "ghijkl")
	}
}

func TestDecodeBuffer_ReturnsSizeLimitWhenGrowthStops(t *testing.T) {
	stub := &scriptedDecoder{
		process: []scriptedCall{{written: 0, err: ErrBufTooSmall}},
	}

	_, err := decodeBufferWithLimit(stub, []byte("ignored"), 1, 1)
	if err != ErrSizeLimit {
		t.Fatalf("decodeBufferWithLimit() error = %v, want ErrSizeLimit", err)
	}
}
