package codec

// AppendFrequencies appends a frequency table using the shared LE wire format.
func AppendFrequencies(out *[]byte, freq []uint32) {
	WriteFrequenciesToBytes(out, freq)
}

// ReadFrequenciesFromBytesExact reads a frequency table from a byte slice and
// rejects any table whose count does not match expectedCount.
func ReadFrequenciesFromBytesExact(in []byte, pos *int, expectedCount int) ([]uint32, error) {
	return ReadFrequenciesFromBytes(in, pos, uint32(expectedCount))
}
