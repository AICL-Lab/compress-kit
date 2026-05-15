package codec

import "testing"

func TestAppendFrequenciesReadFrequenciesFromBytesExact_RoundTrip(t *testing.T) {
	out := []byte{0xAA}
	want := []uint32{7, 11, 13}

	AppendFrequencies(&out, want)

	pos := 1
	got, err := ReadFrequenciesFromBytesExact(out, &pos, len(want))
	if err != nil {
		t.Fatalf("ReadFrequenciesFromBytesExact failed: %v", err)
	}
	if pos != len(out) {
		t.Fatalf("pos = %d, want %d", pos, len(out))
	}
	if len(got) != len(want) {
		t.Fatalf("len(got) = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got[%d] = %d, want %d", i, got[i], want[i])
		}
	}
}
