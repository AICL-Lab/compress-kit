package codec

import (
	"encoding/binary"
	"fmt"
	"io"
)

// SymbolLimit is the number of possible symbols (256 bytes + 1 EOF symbol).
const SymbolLimit = 257

// EOFSymbol is the symbol index used to mark end-of-stream.
const EOFSymbol = SymbolLimit - 1

// WriteFrequencies serializes a frequency table to the writer.
// The format is: count (uint32 LE) followed by count frequency values (uint32 LE each).
// This is the standard format shared by all entropy-based algorithms (Huffman, Arithmetic, Range).
func WriteFrequencies(w io.Writer, freq []uint32) error {
	count := uint32(len(freq))
	if err := binary.Write(w, binary.LittleEndian, count); err != nil {
		return WrapError(KindIO, "failed to write frequency count", err)
	}
	for i, v := range freq {
		if err := binary.Write(w, binary.LittleEndian, v); err != nil {
			return WrapError(KindIO, fmt.Sprintf("failed to write frequency[%d]", i), err)
		}
	}
	return nil
}

// ReadFrequencies deserializes a frequency table from the reader.
// Expects the format written by WriteFrequencies.
// Returns an error if the count is 0, exceeds 1024, or doesn't match expectedCount.
// Pass 0 for expectedCount to skip the count validation.
func ReadFrequencies(r io.Reader, expectedCount uint32) ([]uint32, error) {
	var count uint32
	if err := binary.Read(r, binary.LittleEndian, &count); err != nil {
		return nil, WrapError(KindTruncated, "failed to read frequency table", err)
	}
	if count == 0 || count > 1024 {
		return nil, NewError(KindCorrupt, fmt.Sprintf("invalid frequency table size: %d", count))
	}
	if expectedCount > 0 && count != expectedCount {
		return nil, NewError(KindCorrupt, fmt.Sprintf("unexpected frequency table size: got %d, want %d", count, expectedCount))
	}
	freq := make([]uint32, count)
	if err := binary.Read(r, binary.LittleEndian, freq); err != nil {
		return nil, WrapError(KindTruncated, "failed to read frequency table values", err)
	}
	return freq, nil
}

// ScaleFrequencies normalizes frequencies to fit within maxTotal.
// This is needed for Arithmetic and Range coders where precision is limited.
func ScaleFrequencies(freq []uint32, maxTotal uint32) {
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
	if total <= uint64(maxTotal) {
		return
	}
	var newTotal uint64
	for i, f := range freq {
		if f == 0 {
			continue
		}
		scaled := uint64(f) * uint64(maxTotal) / total
		if scaled == 0 {
			scaled = 1
		}
		freq[i] = uint32(scaled)
		newTotal += scaled
	}
	if newTotal == 0 {
		base := maxTotal / uint32(len(freq))
		if base == 0 {
			base = 1
		}
		for i := range freq {
			freq[i] = base
		}
	}
}

// BuildCumulative builds a cumulative frequency table from frequencies.
// The result has len(freq)+1 elements, where cum[i+1] = cum[i] + freq[i].
func BuildCumulative(freq []uint32) []uint32 {
	cum := make([]uint32, len(freq)+1)
	for i, f := range freq {
		cum[i+1] = cum[i] + f
	}
	// Handle empty case
	if cum[len(cum)-1] == 0 {
		for i := range freq {
			cum[i+1] = uint32(i + 1)
		}
	}
	return cum
}

// BuildFrequencies counts byte frequencies in the input data.
// The EOF symbol is always set to 1.
func BuildFrequencies(data []byte) []uint32 {
	freq := make([]uint32, SymbolLimit)
	for _, b := range data {
		freq[int(b)]++
	}
	freq[EOFSymbol] = 1
	return freq
}

// WriteU32LE appends a uint32 in little-endian format to a byte slice.
func WriteU32LE(out *[]byte, v uint32) {
	*out = append(*out,
		byte(v&0xFF),
		byte((v>>8)&0xFF),
		byte((v>>16)&0xFF),
		byte((v>>24)&0xFF),
	)
}

// ReadU32LE reads a uint32 in little-endian format from a byte slice at the given position.
// Returns the value and true if successful, or 0 and false if out of bounds.
func ReadU32LE(in []byte, pos *int) (uint32, bool) {
	if *pos+4 > len(in) {
		return 0, false
	}
	v := uint32(in[*pos]) |
		uint32(in[*pos+1])<<8 |
		uint32(in[*pos+2])<<16 |
		uint32(in[*pos+3])<<24
	*pos += 4
	return v, true
}

// WriteFrequenciesToBytes appends a frequency table to a byte slice.
// This is useful for algorithms that work with byte slices directly.
func WriteFrequenciesToBytes(out *[]byte, freq []uint32) {
	WriteU32LE(out, uint32(len(freq)))
	for _, v := range freq {
		WriteU32LE(out, v)
	}
}

// ReadFrequenciesFromBytes reads a frequency table from a byte slice at the given position.
// Returns the frequencies and the number of bytes read.
func ReadFrequenciesFromBytes(in []byte, pos *int, expectedCount uint32) ([]uint32, error) {
	count, ok := ReadU32LE(in, pos)
	if !ok {
		return nil, NewError(KindTruncated, "failed to read frequency count")
	}
	if count == 0 || count > 1024 {
		return nil, NewError(KindCorrupt, fmt.Sprintf("invalid frequency table size: %d", count))
	}
	if expectedCount > 0 && count != expectedCount {
		return nil, NewError(KindCorrupt, fmt.Sprintf("unexpected frequency table size: got %d, want %d", count, expectedCount))
	}
	freq := make([]uint32, count)
	for i := uint32(0); i < count; i++ {
		v, ok := ReadU32LE(in, pos)
		if !ok {
			return nil, NewError(KindTruncated, "failed to read frequency values")
		}
		freq[i] = v
	}
	return freq, nil
}
