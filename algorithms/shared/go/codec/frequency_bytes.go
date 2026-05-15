package codec

import "fmt"

func appendU32LE(out *[]byte, v uint32) {
	*out = append(*out,
		byte(v&0xFF),
		byte((v>>8)&0xFF),
		byte((v>>16)&0xFF),
		byte((v>>24)&0xFF),
	)
}

func readU32LE(in []byte, pos *int) (uint32, bool) {
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

// AppendFrequencies appends a frequency table using the shared LE wire format.
func AppendFrequencies(out *[]byte, freq []uint32) {
	appendU32LE(out, uint32(len(freq)))
	for _, value := range freq {
		appendU32LE(out, value)
	}
}

// ReadFrequenciesFromBytes reads a bounded frequency table from a byte slice.
func ReadFrequenciesFromBytes(in []byte, pos *int) ([]uint32, error) {
	count, ok := readU32LE(in, pos)
	if !ok {
		return nil, NewError(KindTruncated, "failed to read frequency table")
	}
	if count == 0 || count > 1024 {
		return nil, NewError(KindCorrupt, fmt.Sprintf("invalid frequency table size: %d", count))
	}

	freq := make([]uint32, count)
	for i := range freq {
		value, ok := readU32LE(in, pos)
		if !ok {
			return nil, NewError(KindTruncated, "failed to read frequency table")
		}
		freq[i] = value
	}
	return freq, nil
}

// ReadFrequenciesFromBytesExact reads a frequency table from a byte slice and
// rejects any table whose count does not match expectedCount.
func ReadFrequenciesFromBytesExact(in []byte, pos *int, expectedCount int) ([]uint32, error) {
	count, ok := readU32LE(in, pos)
	if !ok {
		return nil, NewError(KindTruncated, "failed to read frequency table")
	}
	if count != uint32(expectedCount) {
		return nil, NewError(KindCorrupt, fmt.Sprintf("invalid frequency table size: %d", count))
	}

	freq := make([]uint32, count)
	for i := range freq {
		value, ok := readU32LE(in, pos)
		if !ok {
			return nil, NewError(KindTruncated, "failed to read frequency table")
		}
		freq[i] = value
	}
	return freq, nil
}
