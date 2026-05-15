package codec

import (
	"encoding/binary"
	"fmt"
	"io"
	"sort"
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

// BuildScaledFrequencies counts byte frequencies, preserves an EOF symbol
// frequency of 1, and scales the remaining frequencies to fit within maxTotal.
func BuildScaledFrequencies(data []byte, maxTotal uint32) []uint32 {
	freq := BuildFrequencies(data)
	if maxTotal == 0 {
		return make([]uint32, SymbolLimit)
	}
	if maxTotal == 1 {
		for i := 0; i < EOFSymbol; i++ {
			freq[i] = 0
		}
		freq[EOFSymbol] = 1
		return freq
	}

	var total uint64
	var dataTotal uint64
	nonZero := make([]int, 0, EOFSymbol)
	for i, f := range freq {
		total += uint64(f)
		if i < EOFSymbol && f > 0 {
			dataTotal += uint64(f)
			nonZero = append(nonZero, i)
		}
	}
	if total <= uint64(maxTotal) || dataTotal == 0 {
		return freq
	}

	remaining := int(maxTotal - 1)
	if len(nonZero) > remaining {
		sort.Slice(nonZero, func(i, j int) bool {
			left := nonZero[i]
			right := nonZero[j]
			if freq[left] == freq[right] {
				return left < right
			}
			return freq[left] > freq[right]
		})
		scaled := make([]uint32, SymbolLimit)
		for _, idx := range nonZero[:remaining] {
			scaled[idx] = 1
		}
		scaled[EOFSymbol] = 1
		return scaled
	}

	scaled := make([]uint32, SymbolLimit)
	scaled[EOFSymbol] = 1

	var scaledTotal uint64
	for _, idx := range nonZero {
		value := uint64(freq[idx]) * uint64(maxTotal-1) / dataTotal
		if value == 0 {
			value = 1
		}
		scaled[idx] = uint32(value)
		scaledTotal += value
	}
	for scaledTotal > uint64(maxTotal-1) {
		candidate := -1
		for _, idx := range nonZero {
			if scaled[idx] <= 1 {
				continue
			}
			if candidate == -1 || scaled[idx] > scaled[candidate] || (scaled[idx] == scaled[candidate] && idx < candidate) {
				candidate = idx
			}
		}
		if candidate == -1 {
			break
		}
		scaled[candidate]--
		scaledTotal--
	}

	return scaled
}
