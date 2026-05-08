// Package arithmetic provides arithmetic encoding and decoding implementations.
package arithmetic

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/LessUp/compress-kit/algorithms/shared/go/codec"
)

const (
	// SymbolLimit is the number of possible symbols (256 bytes + 1 EOF symbol).
	SymbolLimit = 257
	// EOFSymbol is the symbol index used to mark end-of-stream.
	EOFSymbol = SymbolLimit - 1
	// MaxTotal is the maximum total frequency value for numerical stability.
	MaxTotal = uint32(1) << 24
	// MaxInputSize is the maximum allowed input file size (4 GiB).
	MaxInputSize = 4 * 1024 * 1024 * 1024

	stateBits    = 32
	fullRange    = uint64(1) << stateBits
	halfRange    = fullRange >> 1
	firstQuarter = halfRange >> 1
	thirdQuarter = firstQuarter * 3
)

// BitWriter writes individual bits to an underlying writer.
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

// ArithmeticEncoder encodes symbols using arithmetic coding.
type ArithmeticEncoder struct {
	writer      *BitWriter
	low         uint64
	high        uint64
	pendingBits uint64
}

// NewArithmeticEncoder creates a new encoder wrapping a BitWriter.
func NewArithmeticEncoder(w *BitWriter) *ArithmeticEncoder {
	return &ArithmeticEncoder{
		writer: w,
		low:    0,
		high:   fullRange - 1,
	}
}

// EncodeSymbol encodes a single symbol using the cumulative frequency table.
func (e *ArithmeticEncoder) EncodeSymbol(symbol uint32, cumulative []uint32) error {
	rangeVal := e.high - e.low + 1
	total := uint64(cumulative[len(cumulative)-1])
	symLow := uint64(cumulative[symbol])
	symHigh := uint64(cumulative[symbol+1])

	e.high = e.low + (rangeVal*symHigh)/total - 1
	e.low = e.low + (rangeVal*symLow)/total

	for {
		if e.high < halfRange {
			if err := e.outputBit(0); err != nil {
				return err
			}
		} else if e.low >= halfRange {
			if err := e.outputBit(1); err != nil {
				return err
			}
			e.low -= halfRange
			e.high -= halfRange
		} else if e.low >= firstQuarter && e.high < thirdQuarter {
			e.pendingBits++
			e.low -= firstQuarter
			e.high -= firstQuarter
		} else {
			break
		}
		e.low <<= 1
		e.high = (e.high << 1) | 1
	}
	return nil
}

// Finish flushes the encoder and writes any remaining bits.
func (e *ArithmeticEncoder) Finish() error {
	e.pendingBits++
	if e.low < firstQuarter {
		if err := e.outputBit(0); err != nil {
			return err
		}
	} else {
		if err := e.outputBit(1); err != nil {
			return err
		}
	}
	return e.writer.Flush()
}

func (e *ArithmeticEncoder) outputBit(bit int) error {
	if err := e.writer.WriteBit(bit); err != nil {
		return err
	}
	complement := bit ^ 1
	for e.pendingBits > 0 {
		if err := e.writer.WriteBit(complement); err != nil {
			return err
		}
		e.pendingBits--
	}
	return nil
}

// ArithmeticDecoder decodes symbols using arithmetic coding.
type ArithmeticDecoder struct {
	reader *BitReader
	low    uint64
	high   uint64
	code   uint64
}

// NewArithmeticDecoder creates a new decoder wrapping a BitReader.
func NewArithmeticDecoder(r *BitReader) *ArithmeticDecoder {
	d := &ArithmeticDecoder{
		reader: r,
		low:    0,
		high:   fullRange - 1,
	}
	for i := uint64(0); i < stateBits; i++ {
		d.code = (d.code << 1) | uint64(r.ReadBit())
	}
	return d
}

// DecodeSymbol decodes the next symbol using the cumulative frequency table.
func (d *ArithmeticDecoder) DecodeSymbol(cumulative []uint32) uint32 {
	rangeVal := d.high - d.low + 1
	total := uint64(cumulative[len(cumulative)-1])
	offset := d.code - d.low
	value := ((offset+1)*total - 1) / rangeVal

	lo := uint32(0)
	hi := uint32(len(cumulative) - 1)
	for lo+1 < hi {
		mid := lo + (hi-lo)/2
		if uint64(cumulative[mid]) > value {
			hi = mid
		} else {
			lo = mid
		}
	}
	symbol := lo

	symLow := uint64(cumulative[symbol])
	symHigh := uint64(cumulative[symbol+1])

	d.high = d.low + (rangeVal*symHigh)/total - 1
	d.low = d.low + (rangeVal*symLow)/total

	for {
		if d.high < halfRange {
			// nothing
		} else if d.low >= halfRange {
			d.low -= halfRange
			d.high -= halfRange
			d.code -= halfRange
		} else if d.low >= firstQuarter && d.high < thirdQuarter {
			d.low -= firstQuarter
			d.high -= firstQuarter
			d.code -= firstQuarter
		} else {
			break
		}
		d.low <<= 1
		d.high = (d.high << 1) | 1
		d.code = (d.code << 1) | uint64(d.reader.ReadBit())
	}

	return symbol
}

// ScaleFrequencies normalizes frequencies to fit within MaxTotal.
func ScaleFrequencies(freq []uint32) {
	var total uint64
	for _, f := range freq {
		total += uint64(f)
	}
	if total == 0 {
		for i := range freq {
			freq[i] = 1
		}
		return
	}
	if total <= uint64(MaxTotal) {
		return
	}
	var newTotal uint64
	for i, f := range freq {
		if f == 0 {
			continue
		}
		scaled := uint64(f) * uint64(MaxTotal) / total
		if scaled == 0 {
			scaled = 1
		}
		freq[i] = uint32(scaled)
		newTotal += scaled
	}
	if newTotal == 0 {
		base := MaxTotal / uint32(len(freq))
		if base == 0 {
			base = 1
		}
		for i := range freq {
			freq[i] = base
		}
	}
}

// BuildFrequenciesFromFile reads a file and counts byte frequencies.
func BuildFrequenciesFromFile(path string) ([]uint32, error) {
	freq := make([]uint32, SymbolLimit)
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("cannot open input file for reading: %s: %w", path, err)
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return nil, fmt.Errorf("cannot stat input file: %w", err)
	}
	if stat.Size() > MaxInputSize {
		return nil, fmt.Errorf("input file too large (max %d bytes)", MaxInputSize)
	}

	r := bufio.NewReader(f)
	for {
		b, err := r.ReadByte()
		if err != nil {
			break
		}
		freq[int(b)]++
	}
	freq[EOFSymbol] = 1
	ScaleFrequencies(freq)
	return freq, nil
}

// BuildCumulative builds a cumulative frequency table from frequencies.
func BuildCumulative(freq []uint32) []uint32 {
	cum := make([]uint32, len(freq)+1)
	for i, f := range freq {
		cum[i+1] = cum[i] + f
	}
	if cum[len(cum)-1] == 0 {
		for i := range freq {
			cum[i+1] = uint32(i + 1)
		}
	}
	return cum
}

// WriteFrequencies serializes a frequency table to the writer.
func WriteFrequencies(w io.Writer, freq []uint32) error {
	count := uint32(len(freq))
	if err := binary.Write(w, binary.LittleEndian, count); err != nil {
		return err
	}
	for _, v := range freq {
		if err := binary.Write(w, binary.LittleEndian, v); err != nil {
			return err
		}
	}
	return nil
}

// ReadFrequencies deserializes a frequency table from the reader.
func ReadFrequencies(r io.Reader) ([]uint32, error) {
	var count uint32
	if err := binary.Read(r, binary.LittleEndian, &count); err != nil {
		return nil, codec.WrapError(codec.KindTruncated, "failed to read frequency table", err)
	}
	if count != uint32(SymbolLimit) {
		return nil, codec.NewError(codec.KindCorrupt, fmt.Sprintf("invalid frequency table size: %d", count))
	}
	freq := make([]uint32, count)
	if err := binary.Read(r, binary.LittleEndian, freq); err != nil {
		return nil, codec.WrapError(codec.KindTruncated, "failed to read frequency table", err)
	}
	return freq, nil
}

// Encode reads from input and writes the arithmetic encoded output to w.
func Encode(input io.Reader, w io.Writer) error {
	data, err := io.ReadAll(input)
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}
	if int64(len(data)) > MaxInputSize {
		return fmt.Errorf("input too large (max %d bytes)", MaxInputSize)
	}

	freq := make([]uint32, SymbolLimit)
	for _, b := range data {
		freq[int(b)]++
	}
	freq[EOFSymbol] = 1
	ScaleFrequencies(freq)
	cumulative := BuildCumulative(freq)

	if _, err := w.Write([]byte{'A', 'E', 'N', 'C'}); err != nil {
		return err
	}
	if err := WriteFrequencies(w, freq); err != nil {
		return err
	}

	bw := NewBitWriter(w)
	encoder := NewArithmeticEncoder(bw)

	for _, b := range data {
		if err := encoder.EncodeSymbol(uint32(b), cumulative); err != nil {
			return err
		}
	}
	if err := encoder.EncodeSymbol(uint32(EOFSymbol), cumulative); err != nil {
		return err
	}
	return encoder.Finish()
}

// Decode reads from r and writes the decoded output to w.
func Decode(r io.Reader, w io.Writer) error {
	br := bufio.NewReader(r)

	magic := make([]byte, 4)
	if _, err := io.ReadFull(br, magic); err != nil || string(magic) != "AENC" {
		return codec.NewError(codec.KindCorrupt, "invalid input file format")
	}

	freq, err := ReadFrequencies(br)
	if err != nil {
		return err
	}
	cumulative := BuildCumulative(freq)

	bw := bufio.NewWriter(w)
	bitReader := NewBitReader(br)
	decoder := NewArithmeticDecoder(bitReader)
	var totalWritten uint64

	for {
		sym := decoder.DecodeSymbol(cumulative)
		if sym == uint32(EOFSymbol) {
			break
		}
		totalWritten++
		if totalWritten > codec.MaxOutputSize {
			return codec.NewError(codec.KindSizeLimit, "output size limit exceeded")
		}
		if err := bw.WriteByte(byte(sym)); err != nil {
			return err
		}
	}

	return bw.Flush()
}

// EncodeFile is a convenience function for file-based encoding.
func EncodeFile(inputPath, outputPath string) error {
	input, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("cannot open input file: %s: %w", inputPath, err)
	}

	encoded, err := codec.EncodeBuffer(NewStreamingEncoder(), input)
	if err != nil {
		return err
	}

	if err := os.WriteFile(outputPath, encoded, 0o644); err != nil {
		return fmt.Errorf("cannot open output file: %s: %w", outputPath, err)
	}

	return nil
}

// DecodeFile is a convenience function for file-based decoding.
func DecodeFile(inputPath, outputPath string) error {
	input, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("cannot open input file: %s: %w", inputPath, err)
	}

	decoded, err := codec.DecodeBuffer(NewStreamingDecoder(), input)
	if err != nil {
		return err
	}

	if err := os.WriteFile(outputPath, decoded, 0o644); err != nil {
		return fmt.Errorf("cannot open output file: %s: %w", outputPath, err)
	}

	return nil
}
