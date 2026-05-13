package codec

import (
	"bufio"
	"bytes"
	"testing"
)

func TestBitWriter_SingleBit(t *testing.T) {
	var buf bytes.Buffer
	bw := NewBitWriter(&buf)

	if err := bw.WriteBit(1); err != nil {
		t.Fatalf("WriteBit failed: %v", err)
	}
	if err := bw.Flush(); err != nil {
		t.Fatalf("Flush failed: %v", err)
	}

	// Single bit 1 should be padded to 0x80 (10000000)
	if buf.Len() != 1 {
		t.Errorf("expected 1 byte, got %d", buf.Len())
	}
	if buf.Bytes()[0] != 0x80 {
		t.Errorf("expected 0x80, got 0x%02x", buf.Bytes()[0])
	}
}

func TestBitWriter_FullByte(t *testing.T) {
	var buf bytes.Buffer
	bw := NewBitWriter(&buf)

	// Write 8 bits: 10110010
	bits := []int{1, 0, 1, 1, 0, 0, 1, 0}
	for _, bit := range bits {
		if err := bw.WriteBit(bit); err != nil {
			t.Fatalf("WriteBit failed: %v", err)
		}
	}

	// Need to flush to see the data
	if err := bw.Flush(); err != nil {
		t.Fatalf("Flush failed: %v", err)
	}

	if buf.Len() != 1 {
		t.Errorf("expected 1 byte after 8 bits, got %d", buf.Len())
	}
	// 10110010 = 0xB2
	if buf.Bytes()[0] != 0xB2 {
		t.Errorf("expected 0xB2, got 0x%02x", buf.Bytes()[0])
	}
}

func TestBitWriter_MultipleBytes(t *testing.T) {
	var buf bytes.Buffer
	bw := NewBitWriter(&buf)

	// Write 16 bits: all 1s
	for i := 0; i < 16; i++ {
		if err := bw.WriteBit(1); err != nil {
			t.Fatalf("WriteBit failed: %v", err)
		}
	}

	if err := bw.Flush(); err != nil {
		t.Fatalf("Flush failed: %v", err)
	}

	if buf.Len() != 2 {
		t.Errorf("expected 2 bytes after 16 bits, got %d", buf.Len())
	}
	if buf.Bytes()[0] != 0xFF || buf.Bytes()[1] != 0xFF {
		t.Errorf("expected 0xFFFF, got 0x%04x", buf.Bytes())
	}
}

func TestBitWriter_FlushPadding(t *testing.T) {
	var buf bytes.Buffer
	bw := NewBitWriter(&buf)

	// Write 3 bits: 101
	for _, bit := range []int{1, 0, 1} {
		if err := bw.WriteBit(bit); err != nil {
			t.Fatalf("WriteBit failed: %v", err)
		}
	}
	if err := bw.Flush(); err != nil {
		t.Fatalf("Flush failed: %v", err)
	}

	// 101 padded to 8 bits: 10100000 = 0xA0
	if buf.Bytes()[0] != 0xA0 {
		t.Errorf("expected 0xA0, got 0x%02x", buf.Bytes()[0])
	}
}

func TestBitReader_SingleByte(t *testing.T) {
	buf := bytes.NewBuffer([]byte{0xB2}) // 10110010
	br := NewBitReader(bufio.NewReader(buf))

	expected := []int{1, 0, 1, 1, 0, 0, 1, 0}
	for i, exp := range expected {
		if got := br.ReadBit(); got != exp {
			t.Errorf("bit %d: expected %d, got %d", i, exp, got)
		}
	}
}

func TestBitReader_MultipleBytes(t *testing.T) {
	buf := bytes.NewBuffer([]byte{0xFF, 0x00})
	br := NewBitReader(bufio.NewReader(buf))

	// First 8 bits should be all 1s
	for i := 0; i < 8; i++ {
		if got := br.ReadBit(); got != 1 {
			t.Errorf("bit %d: expected 1, got %d", i, got)
		}
	}
	// Next 8 bits should be all 0s
	for i := 0; i < 8; i++ {
		if got := br.ReadBit(); got != 0 {
			t.Errorf("bit %d: expected 0, got %d", i+8, got)
		}
	}
}

func TestBitReader_EOF(t *testing.T) {
	buf := bytes.NewBuffer([]byte{0x80})
	br := NewBitReader(bufio.NewReader(buf))

	// Read first bit (1)
	if got := br.ReadBit(); got != 1 {
		t.Errorf("expected 1, got %d", got)
	}
	if br.EOF() {
		t.Error("EOF should be false after reading from valid data")
	}

	// Read remaining 7 bits
	for i := 0; i < 7; i++ {
		br.ReadBit()
	}

	// Now try to read past end
	br.ReadBit()
	if !br.EOF() {
		t.Error("EOF should be true after exhausting buffer")
	}
}

func TestBitWriterReader_RoundTrip(t *testing.T) {
	var buf bytes.Buffer
	bw := NewBitWriter(&buf)

	// Write a sequence of bits
	bits := []int{1, 0, 1, 1, 0, 0, 1, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 0, 1}
	for _, bit := range bits {
		if err := bw.WriteBit(bit); err != nil {
			t.Fatalf("WriteBit failed: %v", err)
		}
	}
	if err := bw.Flush(); err != nil {
		t.Fatalf("Flush failed: %v", err)
	}

	// Read back
	br := NewBitReader(bufio.NewReader(bytes.NewReader(buf.Bytes())))
	for i, exp := range bits {
		if got := br.ReadBit(); got != exp {
			t.Errorf("bit %d: expected %d, got %d", i, exp, got)
		}
	}
}

func TestBitWriter_ZeroBit(t *testing.T) {
	var buf bytes.Buffer
	bw := NewBitWriter(&buf)

	// Write 0 bit
	if err := bw.WriteBit(0); err != nil {
		t.Fatalf("WriteBit failed: %v", err)
	}
	if err := bw.Flush(); err != nil {
		t.Fatalf("Flush failed: %v", err)
	}

	// 0 padded: 00000000
	if buf.Bytes()[0] != 0x00 {
		t.Errorf("expected 0x00, got 0x%02x", buf.Bytes()[0])
	}
}

func TestBitReader_EmptyBuffer(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	br := NewBitReader(bufio.NewReader(buf))

	// Reading from empty buffer should set EOF
	br.ReadBit()
	if !br.EOF() {
		t.Error("expected EOF on empty buffer")
	}
}
