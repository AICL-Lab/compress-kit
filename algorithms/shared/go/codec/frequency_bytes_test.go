package codec

import "testing"

func TestAppendFrequenciesReadFrequenciesFromBytes_RoundTrip(t *testing.T) {
	out := []byte{0xAA}
	want := []uint32{7, 11, 13}

	AppendFrequencies(&out, want)

	pos := 1
	got, err := ReadFrequenciesFromBytes(out, &pos)
	if err != nil {
		t.Fatalf("ReadFrequenciesFromBytes failed: %v", err)
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

func TestReadFrequenciesFromBytesRejectsOutOfRangeCount(t *testing.T) {
	in := []byte{0x00, 0x00, 0x00, 0x00}

	pos := 0
	_, err := ReadFrequenciesFromBytes(in, &pos)
	if err == nil {
		t.Fatal("expected error for out-of-range count")
	}
	if err.Error() != "invalid frequency table size: 0" {
		t.Fatalf("err = %q, want %q", err.Error(), "invalid frequency table size: 0")
	}
}

func TestReadFrequenciesFromBytesRejectsTruncatedEntries(t *testing.T) {
	in := []byte{
		0x02, 0x00, 0x00, 0x00,
		0x01, 0x00, 0x00, 0x00,
	}

	pos := 0
	_, err := ReadFrequenciesFromBytes(in, &pos)
	if err == nil {
		t.Fatal("expected error for truncated entries")
	}
	if err.Error() != "failed to read frequency table" {
		t.Fatalf("err = %q, want %q", err.Error(), "failed to read frequency table")
	}
}
