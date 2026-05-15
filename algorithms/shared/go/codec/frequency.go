package codec

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

// SymbolLimit is the number of possible symbols (256 bytes + 1 EOF symbol).
const SymbolLimit = 257

// EOFSymbol is the symbol index used to mark end-of-stream.
const EOFSymbol = SymbolLimit - 1

// WriteFrequencies serializes a frequency table to the writer.
// The format is: count (uint32 LE) followed by count frequency values (uint32 LE each).
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
// Expects the format written by WriteFrequencies.
func ReadFrequencies(r io.Reader) ([]uint32, error) {
	var count uint32
	if err := binary.Read(r, binary.LittleEndian, &count); err != nil {
		return nil, WrapError(KindTruncated, "failed to read frequency table", err)
	}
	if count == 0 || count > 1024 {
		return nil, NewError(KindCorrupt, fmt.Sprintf("invalid frequency table size: %d", count))
	}
	freq := make([]uint32, count)
	if err := binary.Read(r, binary.LittleEndian, freq); err != nil {
		return nil, WrapError(KindTruncated, "failed to read frequency table", err)
	}
	return freq, nil
}

// ReadFrequenciesExact deserializes a frequency table from the reader and
// rejects any table whose count does not match expectedCount.
func ReadFrequenciesExact(r io.Reader, expectedCount uint32) ([]uint32, error) {
	var count uint32
	if err := binary.Read(r, binary.LittleEndian, &count); err != nil {
		return nil, WrapError(KindTruncated, "failed to read frequency table", err)
	}
	if count != expectedCount {
		return nil, NewError(KindCorrupt, fmt.Sprintf("invalid frequency table size: %d", count))
	}
	freq := make([]uint32, count)
	if err := binary.Read(r, binary.LittleEndian, freq); err != nil {
		return nil, WrapError(KindTruncated, "failed to read frequency table", err)
	}
	return freq, nil
}

// ScaleFrequencies normalizes frequencies toward maxTotal while preserving every
// nonzero symbol. When the number of observed symbols exceeds maxTotal, the
// scaled total can still exceed maxTotal because each observed symbol keeps at
// least frequency 1. This is needed for Arithmetic and Range coders where
// precision is limited.
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

func accumulateFrequencies(freq []uint32, data []byte) error {
	for _, b := range data {
		if freq[int(b)] == math.MaxUint32 {
			return NewError(KindSizeLimit, fmt.Sprintf("frequency overflow for symbol %d", b))
		}
		freq[int(b)]++
	}
	return nil
}

func buildFrequencies(data []byte) ([]uint32, error) {
	freq := make([]uint32, SymbolLimit)
	if err := accumulateFrequencies(freq, data); err != nil {
		return nil, err
	}
	freq[EOFSymbol] = 1
	return freq, nil
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
// The EOF symbol is always set to 1. It panics if a symbol count would exceed
// math.MaxUint32; use BuildFrequenciesChecked to handle that error explicitly.
func BuildFrequencies(data []byte) []uint32 {
	freq, err := BuildFrequenciesChecked(data)
	if err != nil {
		panic(err)
	}
	return freq
}

// BuildFrequenciesChecked counts byte frequencies in the input data and
// returns an error instead of overflowing if any symbol count would exceed
// math.MaxUint32. The EOF symbol is always set to 1.
func BuildFrequenciesChecked(data []byte) ([]uint32, error) {
	return buildFrequencies(data)
}

// BuildFrequenciesFromReader streams bytes from r and counts symbol frequencies.
func BuildFrequenciesFromReader(r io.Reader) ([]uint32, error) {
	freq := make([]uint32, SymbolLimit)
	buf := make([]byte, 32*1024)
	for {
		n, err := r.Read(buf)
		if countErr := accumulateFrequencies(freq, buf[:n]); countErr != nil {
			return nil, countErr
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
	}
	freq[EOFSymbol] = 1
	return freq, nil
}

// BuildScaledFrequencies counts byte frequencies and scales them with the same
// semantics as BuildFrequencies followed by ScaleFrequencies. It panics if a
// symbol count would exceed math.MaxUint32; use BuildScaledFrequenciesChecked
// to handle that error explicitly.
func BuildScaledFrequencies(data []byte, maxTotal uint32) []uint32 {
	freq, err := BuildScaledFrequenciesChecked(data, maxTotal)
	if err != nil {
		panic(err)
	}
	return freq
}

// BuildScaledFrequenciesChecked counts byte frequencies and scales them with
// the same semantics as BuildFrequenciesChecked followed by ScaleFrequencies.
func BuildScaledFrequenciesChecked(data []byte, maxTotal uint32) ([]uint32, error) {
	freq, err := BuildFrequenciesChecked(data)
	if err != nil {
		return nil, err
	}
	ScaleFrequencies(freq, maxTotal)
	return freq, nil
}

// BuildScaledFrequenciesFromReader streams bytes from r, counts symbol
// frequencies, and scales them with ScaleFrequencies.
func BuildScaledFrequenciesFromReader(r io.Reader, maxTotal uint32) ([]uint32, error) {
	freq, err := BuildFrequenciesFromReader(r)
	if err != nil {
		return nil, err
	}
	ScaleFrequencies(freq, maxTotal)
	return freq, nil
}
