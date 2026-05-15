// Package arithmetic provides arithmetic encoding and decoding implementations.
package arithmetic

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/LessUp/compress-kit/algorithms/shared/go/codec"
)

const (
	// SymbolLimit is the number of possible symbols (256 bytes + 1 EOF symbol).
	// This is an alias for codec.SymbolLimit for backward compatibility.
	SymbolLimit = codec.SymbolLimit
	// EOFSymbol is the symbol index used to mark end-of-stream.
	// This is an alias for codec.EOFSymbol for backward compatibility.
	EOFSymbol = codec.EOFSymbol
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

// ArithmeticEncoder encodes symbols using arithmetic coding.
type ArithmeticEncoder struct {
	writer      *codec.BitWriter
	low         uint64
	high        uint64
	pendingBits uint64
}

// NewArithmeticEncoder creates a new encoder wrapping a BitWriter.
func NewArithmeticEncoder(w *codec.BitWriter) *ArithmeticEncoder {
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
	reader *codec.BitReader
	low    uint64
	high   uint64
	code   uint64
}

// NewArithmeticDecoder creates a new decoder wrapping a BitReader.
func NewArithmeticDecoder(r *codec.BitReader) *ArithmeticDecoder {
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
// This is an alias for codec.ScaleFrequencies for backward compatibility.
func ScaleFrequencies(freq []uint32) {
	codec.ScaleFrequencies(freq, MaxTotal)
}

// BuildCumulative builds a cumulative frequency table from frequencies.
// This is an alias for codec.BuildCumulative for backward compatibility.
func BuildCumulative(freq []uint32) []uint32 {
	return codec.BuildCumulative(freq)
}

// WriteFrequencies serializes a frequency table to the writer.
// This is an alias for codec.WriteFrequencies for backward compatibility.
func WriteFrequencies(w io.Writer, freq []uint32) error {
	return codec.WriteFrequencies(w, freq)
}

// ReadFrequencies deserializes a frequency table from the reader.
// This is an alias for codec.ReadFrequenciesExact for backward compatibility.
func ReadFrequencies(r io.Reader) ([]uint32, error) {
	return codec.ReadFrequenciesExact(r, SymbolLimit)
}

// BuildFrequenciesFromFile reads a file and counts byte frequencies.
// The frequencies are scaled to fit within MaxTotal.
func BuildFrequenciesFromFile(path string) ([]uint32, error) {
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

	freq, err := codec.BuildScaledFrequenciesFromReader(bufio.NewReader(f), MaxTotal)
	if err != nil {
		return nil, fmt.Errorf("cannot read input file: %s: %w", path, err)
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

	freq, err := codec.BuildScaledFrequenciesChecked(data, MaxTotal)
	if err != nil {
		return fmt.Errorf("failed to count input frequencies: %w", err)
	}
	cumulative := codec.BuildCumulative(freq)

	if _, err := w.Write([]byte{'A', 'E', 'N', 'C'}); err != nil {
		return err
	}
	if err := WriteFrequencies(w, freq); err != nil {
		return err
	}

	bw := codec.NewBitWriter(w)
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
	cumulative, err := codec.BuildCumulativeStrict(freq, "invalid frequency table")
	if err != nil {
		return err
	}

	bw := bufio.NewWriter(w)
	bitReader := codec.NewBitReader(br)
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

// NewStreamingEncoder creates a new streaming Arithmetic encoder.
// It uses a buffered encoder that collects all input and encodes in one pass
// during Finish(), since Arithmetic encoding requires complete input for frequency analysis.
func NewStreamingEncoder() codec.Encoder {
	return codec.NewBufferedEncoder(func(input []byte) ([]byte, error) {
		var outBuf bytes.Buffer
		if err := Encode(bytes.NewReader(input), &outBuf); err != nil {
			return nil, err
		}
		return outBuf.Bytes(), nil
	})
}

// NewStreamingDecoder creates a new streaming Arithmetic decoder.
// It uses a buffered decoder that collects all input and decodes in one pass
// during Finish().
func NewStreamingDecoder() codec.Decoder {
	return codec.NewBufferedDecoder(func(input []byte) ([]byte, error) {
		var outBuf bytes.Buffer
		if err := Decode(bytes.NewReader(input), &outBuf); err != nil {
			return nil, err
		}
		return outBuf.Bytes(), nil
	})
}

// EncodeFile is a convenience function for file-based encoding.
func EncodeFile(inputPath, outputPath string) error {
	return codec.EncodeFile(NewStreamingEncoder(), inputPath, outputPath)
}

// DecodeFile is a convenience function for file-based decoding.
func DecodeFile(inputPath, outputPath string) error {
	return codec.DecodeFile(NewStreamingDecoder(), inputPath, outputPath)
}
