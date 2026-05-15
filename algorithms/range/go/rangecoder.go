package rangecoder

import (
	"github.com/LessUp/compress-kit/algorithms/shared/go/codec"
)

const (
	maxTotal        uint32 = 1 << 24
	renormThreshold        = uint32(1) << 24
	maxOutputSize          = 1 << 30 // 1 GiB maximum decoded output
)

func scaleFrequencies(freq []uint32) {
	codec.ScaleFrequencies(freq, maxTotal)
}

func buildFrequencies(data []byte) ([]uint32, error) {
	return codec.BuildScaledFrequenciesChecked(data, maxTotal)
}

func buildCumulative(freq []uint32) []uint32 {
	return codec.BuildCumulative(freq)
}

func writeHeader(out *[]byte, freq []uint32) {
	*out = append(*out, 'R', 'C', 'N', 'C')
	codec.AppendFrequencies(out, freq)
}

func readHeader(in []byte, pos *int) ([]uint32, error) {
	if len(in) < 8 {
		return nil, codec.NewError(codec.KindTruncated, "range: input too short")
	}
	if in[0] != 'R' || in[1] != 'C' || in[2] != 'N' || in[3] != 'C' {
		return nil, codec.NewError(codec.KindCorrupt, "range: bad magic")
	}
	*pos = 4
	freq, err := codec.ReadFrequenciesFromBytesExact(in, pos, codec.SymbolLimit)
	if err != nil {
		if codecErr, ok := err.(*codec.CodecError); ok && codecErr.Kind == codec.KindTruncated {
			return nil, codec.NewError(codec.KindTruncated, "range: truncated header")
		}
		return nil, codec.NewError(codec.KindCorrupt, "range: bad header")
	}
	return freq, nil
}

type encoder struct {
	low  uint32
	high uint32
	out  *[]byte
}

func newEncoder(out *[]byte) *encoder {
	return &encoder{low: 0, high: 0xFFFFFFFF, out: out}
}

func (e *encoder) encodeSymbol(symbol uint32, cumulative []uint32) {
	rangeVal := uint64(e.high) - uint64(e.low) + 1
	total := uint64(cumulative[len(cumulative)-1])
	symLow := uint64(cumulative[symbol])
	symHigh := uint64(cumulative[symbol+1])

	e.high = e.low + uint32((rangeVal*symHigh)/total-1)
	e.low = e.low + uint32((rangeVal*symLow)/total)

	for (e.low ^ e.high) < renormThreshold {
		b := byte(e.low >> 24)
		*e.out = append(*e.out, b)
		e.low <<= 8
		e.high = (e.high << 8) | 0xFF
	}
}

func (e *encoder) finish() {
	for i := 0; i < 4; i++ {
		b := byte(e.low >> 24)
		*e.out = append(*e.out, b)
		e.low <<= 8
	}
}

type decoder struct {
	low  uint32
	high uint32
	code uint32
	in   []byte
	pos  int
}

func newDecoder(in []byte) *decoder {
	d := &decoder{low: 0, high: 0xFFFFFFFF, in: in}
	for i := 0; i < 4; i++ {
		b := d.readByte()
		d.code = (d.code << 8) | uint32(b)
	}
	return d
}

func (d *decoder) readByte() byte {
	if d.pos < len(d.in) {
		b := d.in[d.pos]
		d.pos++
		return b
	}
	return 0
}

// Note: This binary search has O(log n) complexity per symbol.
// For better performance with large files, consider using a lookup table
// or prefix-sum index to achieve O(1) symbol lookup.
func (d *decoder) decodeSymbol(cumulative []uint32) uint32 {
	rangeVal := uint64(d.high) - uint64(d.low) + 1
	total := uint64(cumulative[len(cumulative)-1])
	offset := uint64(d.code - d.low)
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

	d.high = d.low + uint32((rangeVal*symHigh)/total-1)
	d.low = d.low + uint32((rangeVal*symLow)/total)

	for (d.low ^ d.high) < renormThreshold {
		d.low <<= 8
		d.high = (d.high << 8) | 0xFF
		d.code = (d.code << 8) | uint32(d.readByte())
	}

	return symbol
}

func Encode(input []byte) ([]byte, error) {
	freq, err := buildFrequencies(input)
	if err != nil {
		return nil, err
	}
	cum := buildCumulative(freq)

	out := make([]byte, 0, len(input))
	writeHeader(&out, freq)

	enc := newEncoder(&out)
	for _, b := range input {
		enc.encodeSymbol(uint32(b), cum)
	}
	enc.encodeSymbol(codec.EOFSymbol, cum)
	enc.finish()

	return out, nil
}

func Decode(encoded []byte) ([]byte, error) {
	pos := 0
	freq, err := readHeader(encoded, &pos)
	if err != nil {
		return nil, err
	}
	if len(freq) != codec.SymbolLimit {
		return nil, codec.NewError(codec.KindCorrupt, "range: unexpected symbol count")
	}
	cum := buildCumulative(freq)
	if pos >= len(encoded) {
		return []byte{}, nil
	}

	dec := newDecoder(encoded[pos:])
	out := make([]byte, 0, len(encoded))
	for {
		sym := dec.decodeSymbol(cum)
		if sym == uint32(codec.EOFSymbol) {
			break
		}
		out = append(out, byte(sym))
		if len(out) > maxOutputSize {
			return nil, codec.NewError(codec.KindSizeLimit, "range: output exceeds maximum size (1 GiB)")
		}
	}
	return out, nil
}

// EncodeFile is a convenience function for file-based encoding.
func EncodeFile(inputPath, outputPath string) error {
	return codec.EncodeFile(NewStreamingEncoder(), inputPath, outputPath)
}

// DecodeFile is a convenience function for file-based decoding.
func DecodeFile(inputPath, outputPath string) error {
	return codec.DecodeFile(NewStreamingDecoder(), inputPath, outputPath)
}
