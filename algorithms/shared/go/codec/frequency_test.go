package codec

import (
	"bytes"
	"errors"
	"io"
	"math"
	"testing"
)

type chunkedReader struct {
	data      []byte
	chunkSize int
	offset    int
}

func (r *chunkedReader) Read(p []byte) (int, error) {
	if r.offset >= len(r.data) {
		return 0, io.EOF
	}
	n := r.chunkSize
	if n > len(p) {
		n = len(p)
	}
	remaining := len(r.data) - r.offset
	if n > remaining {
		n = remaining
	}
	copy(p, r.data[r.offset:r.offset+n])
	r.offset += n
	return n, nil
}

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
	freq := MustBuildFrequencies(data)

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
	freq := MustBuildFrequencies([]byte{})

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

func TestAccumulateFrequencies_RejectsSymbolOverflow(t *testing.T) {
	freq := make([]uint32, SymbolLimit)
	freq['a'] = math.MaxUint32

	err := accumulateFrequencies(freq, []byte{'a'})
	if err == nil {
		t.Fatal("expected overflow error")
	}
	if !errors.Is(err, ErrSizeLimit) {
		t.Fatalf("expected size limit error, got %v", err)
	}
	if freq['a'] != math.MaxUint32 {
		t.Fatalf("freq['a'] = %d, want %d", freq['a'], uint32(math.MaxUint32))
	}
}

func TestBuildScaledFrequencies_ClampsTotalAndPreservesEOF(t *testing.T) {
	data := append(bytes.Repeat([]byte{'a'}, 12), bytes.Repeat([]byte{'b'}, 6)...)

	freq := MustBuildScaledFrequencies(data, 8)

	var total uint32
	for _, f := range freq {
		total += f
	}
	if total > 8 {
		t.Fatalf("total = %d, want <= 8", total)
	}
	if freq[EOFSymbol] != 1 {
		t.Fatalf("freq[EOF] = %d, want 1", freq[EOFSymbol])
	}
	if freq['a'] <= freq['b'] {
		t.Fatalf("expected scaled frequencies to preserve ordering, got a=%d b=%d", freq['a'], freq['b'])
	}
}

func TestBuildScaledFrequencies_MatchesScaleFrequenciesSemantics(t *testing.T) {
	data := []byte{'a', 'b', 'c', 'd'}

	want := MustBuildFrequencies(data)
	ScaleFrequencies(want, 4)

	got := MustBuildScaledFrequencies(data, 4)

	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got[%d] = %d, want %d", i, got[i], want[i])
		}
	}
}

func TestScaleFrequencies_PreservesObservedSymbolsWhenBudgetIsTooSmall(t *testing.T) {
	freq := []uint32{4, 3, 2}

	ScaleFrequencies(freq, 2)

	var total uint32
	for i, f := range freq {
		if f == 0 {
			t.Fatalf("freq[%d] = 0, want observed symbol preserved", i)
		}
		total += f
	}
	if total <= 2 {
		t.Fatalf("total = %d, want > 2 when preserving all observed symbols", total)
	}
}

func TestScaleFrequencies_ClampsOvershootWhenReductionIsFeasible(t *testing.T) {
	freq := []uint32{100, 1, 1, 1}

	ScaleFrequencies(freq, 5)

	var total uint32
	for i, f := range freq {
		if f == 0 {
			t.Fatalf("freq[%d] = 0, want observed symbol preserved", i)
		}
		total += f
	}
	if total > 5 {
		t.Fatalf("total = %d, want <= 5 after feasible reduction", total)
	}
	if freq[0] <= 1 {
		t.Fatalf("freq[0] = %d, want dominant symbol to remain > 1", freq[0])
	}
}

func TestBuildFrequenciesFromReader_MatchesSliceHelper(t *testing.T) {
	data := append(bytes.Repeat([]byte{'a'}, 7), bytes.Repeat([]byte{'b'}, 3)...)
	reader := &chunkedReader{data: data, chunkSize: 2}

	got, err := BuildFrequenciesFromReader(reader)
	if err != nil {
		t.Fatalf("BuildFrequenciesFromReader failed: %v", err)
	}

	want := MustBuildFrequencies(data)
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got[%d] = %d, want %d", i, got[i], want[i])
		}
	}
}

func TestBuildScaledFrequenciesFromReader_MatchesSliceHelper(t *testing.T) {
	data := []byte{'a', 'b', 'c', 'd'}
	reader := &chunkedReader{data: data, chunkSize: 1}

	got, err := BuildScaledFrequenciesFromReader(reader, 4)
	if err != nil {
		t.Fatalf("BuildScaledFrequenciesFromReader failed: %v", err)
	}

	want := MustBuildScaledFrequencies(data, 4)
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got[%d] = %d, want %d", i, got[i], want[i])
		}
	}
}

func TestWriteReadFrequencies_RoundTrip(t *testing.T) {
	freq := []uint32{10, 20, 30, 40, 50}

	var buf bytes.Buffer
	if err := WriteFrequencies(&buf, freq); err != nil {
		t.Fatalf("WriteFrequencies failed: %v", err)
	}

	got, err := ReadFrequencies(&buf, 0)
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

	_, err := ReadFrequencies(buf, 0)
	if err == nil {
		t.Error("expected error for truncated frequency table")
	}
}

func TestReadFrequencies_InvalidCount(t *testing.T) {
	// count = 0
	buf := bytes.NewBuffer([]byte{0, 0, 0, 0})

	_, err := ReadFrequencies(buf, 0)
	if err == nil {
		t.Error("expected error for count = 0")
	}
}

func TestReadFrequencies_ExpectedCountMismatch(t *testing.T) {
	// Write count = 5
	buf := bytes.NewBuffer([]byte{5, 0, 0, 0, 1, 0, 0, 0, 2, 0, 0, 0, 3, 0, 0, 0, 4, 0, 0, 0, 5, 0, 0, 0})

	_, err := ReadFrequencies(buf, 3) // expect 3, but got 5
	if err == nil {
		t.Error("expected error for count mismatch")
	}
}

func TestReadFrequenciesExact_RejectsWrongCounts(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteFrequencies(&buf, []uint32{10, 20, 30}); err != nil {
		t.Fatalf("WriteFrequencies failed: %v", err)
	}

	_, err := ReadFrequenciesExact(&buf, 4)
	if err == nil {
		t.Fatal("expected error for wrong frequency count")
	}
	if !errors.Is(err, ErrCorrupt) {
		t.Fatalf("expected corrupt error, got %v", err)
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

func TestBuildCumulativeStrict_RejectsAllZeroTable(t *testing.T) {
	_, err := BuildCumulativeStrict([]uint32{0, 0, 0}, "invalid table")
	if err == nil {
		t.Fatal("expected error for all-zero table")
	}
	if !errors.Is(err, ErrCorrupt) {
		t.Fatalf("expected corrupt error, got %v", err)
	}
	if err.Error() != "invalid table" {
		t.Fatalf("err = %q, want %q", err.Error(), "invalid table")
	}
}

func TestBuildCumulativeStrict_PreservesFallbackBehaviorForNonZeroTable(t *testing.T) {
	got, err := BuildCumulativeStrict([]uint32{0, 2, 0}, "invalid table")
	if err != nil {
		t.Fatalf("BuildCumulativeStrict failed: %v", err)
	}
	want := []uint32{0, 0, 2, 2}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got[%d] = %d, want %d", i, got[i], want[i])
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
