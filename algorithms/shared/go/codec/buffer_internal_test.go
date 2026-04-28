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
