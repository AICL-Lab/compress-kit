package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

// Arithmetic coding Go implementation.
// File format is fully compatible with C++ implementation, supports cross-language encode/decode verification.
// Magic: AENC (4 bytes)
// Frequency table: count(4 bytes LE) + count × freq(4 bytes LE)
// Arithmetic coding bitstream

const (
	SymbolLimit  = 257
	EOFSymbol    = SymbolLimit - 1
	MaxTotal     = uint32(1) << 24
	MaxInputSize = 4 * 1024 * 1024 * 1024 // 4 GiB max to prevent frequency overflow

	stateBits    = 32
	fullRange    = uint64(1) << stateBits
	halfRange    = fullRange >> 1
	firstQuarter = halfRange >> 1
	thirdQuarter = firstQuarter * 3
)

// ---------------------------------------------------------------------------
// BitWriter / BitReader
// ---------------------------------------------------------------------------

type BitWriter struct {
	w            *bufio.Writer
	buffer       byte
	bitsInBuffer uint8
}

func NewBitWriter(w io.Writer) *BitWriter {
	return &BitWriter{w: bufio.NewWriter(w)}
}

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

type BitReader struct {
	r             *bufio.Reader
	currentByte   byte
	bitsRemaining uint8
	reachedEOF    bool
}

func NewBitReader(r *bufio.Reader) *BitReader {
	return &BitReader{r: r}
}

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

// ---------------------------------------------------------------------------
// ArithmeticEncoder
// ---------------------------------------------------------------------------

type ArithmeticEncoder struct {
	writer      *BitWriter
	low         uint64
	high        uint64
	pendingBits uint64
}

func NewArithmeticEncoder(w *BitWriter) *ArithmeticEncoder {
	return &ArithmeticEncoder{
		writer: w,
		low:    0,
		high:   fullRange - 1,
	}
}

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

// ---------------------------------------------------------------------------
// ArithmeticDecoder
// ---------------------------------------------------------------------------

type ArithmeticDecoder struct {
	reader *BitReader
	low    uint64
	high   uint64
	code   uint64
}

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

// ---------------------------------------------------------------------------
// Frequency table processing
// ---------------------------------------------------------------------------

func scaleFrequencies(freq []uint32) {
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

func buildFrequenciesFromFile(path string) ([]uint32, error) {
	freq := make([]uint32, SymbolLimit)
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("cannot open input file for reading: %s: %w", path, err)
	}
	defer f.Close()

	// Check file size to prevent frequency overflow
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
	scaleFrequencies(freq)
	return freq, nil
}

func buildCumulative(freq []uint32) []uint32 {
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

func writeFrequencies(w io.Writer, freq []uint32) error {
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

func readFrequencies(r io.Reader) ([]uint32, error) {
	var count uint32
	if err := binary.Read(r, binary.LittleEndian, &count); err != nil {
		return nil, fmt.Errorf("failed to read frequency table: %w", err)
	}
	if count != uint32(SymbolLimit) {
		return nil, fmt.Errorf("invalid frequency table size: %d", count)
	}
	freq := make([]uint32, count)
	if err := binary.Read(r, binary.LittleEndian, freq); err != nil {
		return nil, fmt.Errorf("failed to read frequency table: %w", err)
	}
	return freq, nil
}

// ---------------------------------------------------------------------------
// Compression / Decompression
// ---------------------------------------------------------------------------

func compressFile(inputPath, outputPath string) error {
	freq, err := buildFrequenciesFromFile(inputPath)
	if err != nil {
		return err
	}
	cumulative := buildCumulative(freq)

	inFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("cannot open input file for reading: %s: %w", inputPath, err)
	}
	defer inFile.Close()

	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("cannot open output file for writing: %s: %w", outputPath, err)
	}
	defer outFile.Close()

	if _, err := outFile.Write([]byte{'A', 'E', 'N', 'C'}); err != nil {
		return err
	}
	if err := writeFrequencies(outFile, freq); err != nil {
		return err
	}

	bw := NewBitWriter(outFile)
	encoder := NewArithmeticEncoder(bw)

	r := bufio.NewReader(inFile)
	for {
		b, err := r.ReadByte()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read input file: %w", err)
		}
		if err := encoder.EncodeSymbol(uint32(b), cumulative); err != nil {
			return err
		}
	}
	if err := encoder.EncodeSymbol(uint32(EOFSymbol), cumulative); err != nil {
		return err
	}
	return encoder.Finish()
}

func decompressFile(inputPath, outputPath string) error {
	inFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("cannot open input file for reading: %s: %w", inputPath, err)
	}
	defer inFile.Close()
	r := bufio.NewReader(inFile)

	magic := make([]byte, 4)
	if _, err := io.ReadFull(r, magic); err != nil || magic[0] != 'A' || magic[1] != 'E' || magic[2] != 'N' || magic[3] != 'C' {
		return fmt.Errorf("invalid input file format")
	}

	freq, err := readFrequencies(r)
	if err != nil {
		return err
	}
	cumulative := buildCumulative(freq)

	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("cannot open output file for writing: %s: %w", outputPath, err)
	}
	defer outFile.Close()
	w := bufio.NewWriter(outFile)

	br := NewBitReader(r)
	decoder := NewArithmeticDecoder(br)

	for {
		sym := decoder.DecodeSymbol(cumulative)
		if sym == uint32(EOFSymbol) {
			break
		}
		if err := w.WriteByte(byte(sym)); err != nil {
			return err
		}
	}

	return w.Flush()
}

func main() {
	if len(os.Args) != 4 {
		fmt.Fprintf(os.Stderr, "Usage: %s encode|decode input output\n", os.Args[0])
		os.Exit(1)
	}
	mode := os.Args[1]
	inputPath := os.Args[2]
	outputPath := os.Args[3]

	var err error

	if mode == "encode" {
		err = compressFile(inputPath, outputPath)
	} else if mode == "decode" {
		err = decompressFile(inputPath, outputPath)
	} else {
		fmt.Fprintln(os.Stderr, "unknown mode, expected encode or decode")
		os.Exit(1)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
