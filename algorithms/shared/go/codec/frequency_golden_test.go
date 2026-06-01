// Byte-exact golden tests for the shared Frequency Table wire format.
//
// These tests pin the exact bytes emitted by WriteFrequenciesToBytes / AppendFrequencies
// so that refactors that silently change endianness, field order, or count encoding
// will fail immediately rather than silently producing incompatible output.
//
// Wire format (from contract-inventory.md § Binary format and Frequency Table facts):
//   count     : uint32 little-endian  -- number of frequency entries
//   freq[0]   : uint32 little-endian
//   ...
//   freq[n-1] : uint32 little-endian
package codec

import (
	"bytes"
	"testing"
)

// TestGolden_WriteFrequenciesToBytes_Empty verifies that an empty table encodes
// as exactly 4 zero bytes (count=0 in little-endian).
// Wire: 00 00 00 00
func TestGolden_WriteFrequenciesToBytes_Empty(t *testing.T) {
	var out []byte
	WriteFrequenciesToBytes(&out, nil)
	want := []byte{0x00, 0x00, 0x00, 0x00}
	if !bytes.Equal(out, want) {
		t.Errorf("empty table: got %x, want %x", out, want)
	}
}

// TestGolden_WriteFrequenciesToBytes_SingleEntry verifies count=1 then value=1.
// Wire: 01 00 00 00  01 00 00 00
func TestGolden_WriteFrequenciesToBytes_SingleEntry(t *testing.T) {
	var out []byte
	WriteFrequenciesToBytes(&out, []uint32{1})
	want := []byte{
		0x01, 0x00, 0x00, 0x00, // count = 1
		0x01, 0x00, 0x00, 0x00, // freq[0] = 1
	}
	if !bytes.Equal(out, want) {
		t.Errorf("single-entry table: got %x, want %x", out, want)
	}
}

// TestGolden_WriteFrequenciesToBytes_ThreeEntries verifies count=3 then each value.
// Wire: 03 00 00 00  01 00 00 00  02 00 00 00  03 00 00 00
func TestGolden_WriteFrequenciesToBytes_ThreeEntries(t *testing.T) {
	var out []byte
	WriteFrequenciesToBytes(&out, []uint32{1, 2, 3})
	want := []byte{
		0x03, 0x00, 0x00, 0x00, // count = 3
		0x01, 0x00, 0x00, 0x00, // freq[0] = 1
		0x02, 0x00, 0x00, 0x00, // freq[1] = 2
		0x03, 0x00, 0x00, 0x00, // freq[2] = 3
	}
	if !bytes.Equal(out, want) {
		t.Errorf("three-entry table: got %x, want %x", out, want)
	}
}

// TestGolden_WriteFrequenciesToBytes_MaxU32 verifies that u32::MAX is preserved.
// Wire: 01 00 00 00  FF FF FF FF
func TestGolden_WriteFrequenciesToBytes_MaxU32(t *testing.T) {
	var out []byte
	WriteFrequenciesToBytes(&out, []uint32{^uint32(0)})
	want := []byte{
		0x01, 0x00, 0x00, 0x00, // count = 1
		0xFF, 0xFF, 0xFF, 0xFF, // freq[0] = 0xFFFFFFFF
	}
	if !bytes.Equal(out, want) {
		t.Errorf("max u32 entry: got %x, want %x", out, want)
	}
}

// TestGolden_WriteFrequenciesToBytes_Count256 verifies multi-byte count encoding.
// count=256 must encode as 00 01 00 00 in little-endian.
func TestGolden_WriteFrequenciesToBytes_Count256(t *testing.T) {
	freq := make([]uint32, 256)
	var out []byte
	WriteFrequenciesToBytes(&out, freq)
	wantCountBytes := []byte{0x00, 0x01, 0x00, 0x00}
	if !bytes.Equal(out[:4], wantCountBytes) {
		t.Errorf("count=256 field: got %x, want %x", out[:4], wantCountBytes)
	}
	wantLen := 4 + 256*4
	if len(out) != wantLen {
		t.Errorf("total length: got %d, want %d", len(out), wantLen)
	}
}

// TestGolden_AppendFrequencies_PreservesLeadingBytes verifies that AppendFrequencies
// appends to (not replaces) an existing byte slice.
func TestGolden_AppendFrequencies_PreservesLeadingBytes(t *testing.T) {
	out := []byte{0xDE, 0xAD}
	AppendFrequencies(&out, []uint32{5})
	want := []byte{
		0xDE, 0xAD, // original prefix
		0x01, 0x00, 0x00, 0x00, // count = 1
		0x05, 0x00, 0x00, 0x00, // freq[0] = 5
	}
	if !bytes.Equal(out, want) {
		t.Errorf("AppendFrequencies: got %x, want %x", out, want)
	}
}
