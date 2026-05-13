package codec

import (
	"bufio"
	"io"
)

// BitWriter writes individual bits to an underlying writer, buffering until
// a full byte is available.
type BitWriter struct {
	w            *bufio.Writer
	buffer       byte
	bitsInBuffer uint8
}

// NewBitWriter creates a BitWriter wrapping the given io.Writer.
func NewBitWriter(w io.Writer) *BitWriter {
	return &BitWriter{w: bufio.NewWriter(w)}
}

// WriteBit queues a single bit for writing.
func (b *BitWriter) WriteBit(bit int) error {
	b.buffer = (b.buffer << 1) | byte(bit&1)
	b.bitsInBuffer++
	if b.bitsInBuffer == 8 {
		if err := b.w.WriteByte(b.buffer); err != nil {
			return err
		}
		b.bitsInBuffer = 0
		b.buffer = 0
	}
	return nil
}

// Flush writes any pending bits and flushes the underlying writer.
func (b *BitWriter) Flush() error {
	if b.bitsInBuffer > 0 {
		b.buffer <<= (8 - b.bitsInBuffer)
		if err := b.w.WriteByte(b.buffer); err != nil {
			return err
		}
		b.bitsInBuffer = 0
		b.buffer = 0
	}
	return b.w.Flush()
}

// BitReader reads individual bits from an underlying buffered reader.
type BitReader struct {
	r             *bufio.Reader
	currentByte   byte
	bitsRemaining uint8
	reachedEOF    bool
}

// NewBitReader creates a BitReader wrapping the given bufio.Reader.
func NewBitReader(r *bufio.Reader) *BitReader {
	return &BitReader{r: r}
}

// ReadBit returns the next bit (0 or 1).
func (b *BitReader) ReadBit() int {
	if b.bitsRemaining == 0 {
		c, err := b.r.ReadByte()
		if err != nil {
			b.reachedEOF = true
			return 0
		}
		b.currentByte = c
		b.bitsRemaining = 8
	}
	b.bitsRemaining--
	return int((b.currentByte >> b.bitsRemaining) & 1)
}

// EOF returns true if the underlying reader has been exhausted.
func (b *BitReader) EOF() bool {
	return b.reachedEOF
}
