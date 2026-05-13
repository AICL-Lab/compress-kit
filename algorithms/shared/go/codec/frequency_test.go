package codec

import (
	"bytes"
	"testing"
)

func TestSymbolLimit(t *testing.T) {
	// SymbolLimit should be 257 (256 bytes + 1 EOF)
	if SymbolLimit != 257 {
		t.Errorf("SymbolLimit = %d, want 257", SymbolLimit)
	}
	if EOFSymbol != 256 {
		t.Errorf("EOFSymbol = %d, want 256", EOFSymbol)
	}
}

func TestBuildFrequencies(t *testing.T) {
	data := []byte("aabbbc")
	freq := BuildFrequencies(data)

	if freq['a'] != 2 {
		t.Errorf("freq['a'] = %d, want 2", freq['a'])
	}
	if freq['b'] != 3 {
		t.Errorf("freq['b'] = %d, want 3", freq['b'])
	}
	if freq['c'] != 1 {
		t.Errorf("freq['c'] = %d, want 1", freq['c'])
	}
	if freq[EOFSymbol] != 1 {
		t.Errorf("freq[EOF] = %d, want 1", freq[EOFSymbol])
	}
}

func TestBuildFrequencies_Empty(t *testing.T) {
	freq := BuildFrequencies([]byte{})

	// All byte frequencies should be 0
	for i := 0; i < 256; i++ {
		if freq[i] != 0 {
			t.Errorf("freq[%d] = %d, want 0", i, freq[i])
		}
	}
	// EOF should always be 1
	if freq[EOFSymbol] != 1 {
		t.Errorf("freq[EOF] = %d, want 1", freq[EOFSymbol])
	}
}

func TestWriteReadFrequencies_RoundTrip(t *testing.T) {
	freq := []uint32{10, 20, 30, 40, 50}

	var buf bytes.Buffer
	if err := WriteFrequencies(&buf, freq); err != nil {
		t.Fatalf("WriteFrequencies failed: %v", err)
	}

	got, err := ReadFrequencies(&buf)
	if err != nil {
		t.Fatalf("ReadFrequencies failed: %v", err)
	}

	if len(got) != len(freq) {
		t.Fatalf("len(got) = %d, want %d", len(got), len(freq))
	}
	for i := range freq {
		if got[i] != freq[i] {
			t.Errorf("got[%d] = %d, want %d", i, got[i], freq[i])
		}
	}
}

func TestReadFrequencies_Truncated(t *testing.T) {
	// Only write count, not the actual frequencies
	buf := bytes.NewBuffer([]byte{3, 0, 0, 0}) // count = 3

	_, err := ReadFrequencies(buf)
	if err == nil {
		t.Error("expected error for truncated frequency table")
	}
}

func TestReadFrequencies_InvalidCount(t *testing.T) {
	// count = 0
	buf := bytes.NewBuffer([]byte{0, 0, 0, 0})

	_, err := ReadFrequencies(buf)
	if err == nil {
		t.Error("expected error for count = 0")
	}
}

func TestScaleFrequencies_NoScalingNeeded(t *testing.T) {
	freq := []uint32{100, 200, 300}
	maxTotal := uint32(10000)

	original := make([]uint32, len(freq))
	copy(original, freq)

	ScaleFrequencies(freq, maxTotal)

	// Should not change since total (600) < maxTotal (10000)
	for i := range freq {
		if freq[i] != original[i] {
			t.Errorf("freq[%d] = %d, want %d (no scaling needed)", i, freq[i], original[i])
		}
	}
}

func TestScaleFrequencies_ScalingNeeded(t *testing.T) {
	freq := []uint32{1000, 2000, 3000}
	maxTotal := uint32(100)

	ScaleFrequencies(freq, maxTotal)

	// Check that total is now <= maxTotal
	var total uint64
	for _, f := range freq {
		total += uint64(f)
	}
	if total > uint64(maxTotal) {
		t.Errorf("total after scaling = %d, want <= %d", total, maxTotal)
	}

	// Check that relative proportions are approximately preserved
	// freq[0] : freq[1] : freq[2] should be approximately 1 : 2 : 3
	if freq[1] < freq[0] || freq[2] < freq[1] {
		t.Errorf("relative proportions not preserved: got [%d, %d, %d]", freq[0], freq[1], freq[2])
	}
}

func TestScaleFrequencies_AllZeros(t *testing.T) {
	freq := []uint32{0, 0, 0, 0}
	maxTotal := uint32(100)

	ScaleFrequencies(freq, maxTotal)

	// All should be set to 1
	for i, f := range freq {
		if f != 1 {
			t.Errorf("freq[%d] = %d, want 1", i, f)
		}
	}
}

func TestBuildCumulative(t *testing.T) {
	freq := []uint32{10, 20, 30}
	cum := BuildCumulative(freq)

	// cum[0] = 0
	// cum[1] = 0 + 10 = 10
	// cum[2] = 10 + 20 = 30
	// cum[3] = 30 + 30 = 60
	expected := []uint32{0, 10, 30, 60}

	if len(cum) != len(expected) {
		t.Fatalf("len(cum) = %d, want %d", len(cum), len(expected))
	}
	for i := range expected {
		if cum[i] != expected[i] {
			t.Errorf("cum[%d] = %d, want %d", i, cum[i], expected[i])
		}
	}
}

func TestBuildCumulative_Empty(t *testing.T) {
	freq := []uint32{0, 0, 0}
	cum := BuildCumulative(freq)

	// When all zeros, should use sequential values
	// cum[0] = 0, cum[1] = 1, cum[2] = 2, cum[3] = 3
	expected := []uint32{0, 1, 2, 3}

	for i := range expected {
		if cum[i] != expected[i] {
			t.Errorf("cum[%d] = %d, want %d", i, cum[i], expected[i])
		}
	}
}

func TestBuildCumulative_SingleElement(t *testing.T) {
	freq := []uint32{42}
	cum := BuildCumulative(freq)

	// cum[0] = 0, cum[1] = 42
	if cum[0] != 0 {
		t.Errorf("cum[0] = %d, want 0", cum[0])
	}
	if cum[1] != 42 {
		t.Errorf("cum[1] = %d, want 42", cum[1])
	}
}
