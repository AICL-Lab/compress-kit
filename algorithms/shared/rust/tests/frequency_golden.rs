// Byte-exact golden tests for the shared Frequency Table wire format.
//
// These tests pin the exact bytes emitted by write_frequencies so that
// refactors that silently change endianness, field order, or count encoding
// will fail immediately rather than silently producing incompatible output.
//
// Wire format (from contract-inventory.md § Binary format and Frequency Table facts):
//   count     : uint32 little-endian  -- number of frequency entries
//   freq[0]   : uint32 little-endian
//   ...
//   freq[n-1] : uint32 little-endian
//
// The exact 257-symbol layout used by Huffman, Arithmetic, and Range algorithms
// is not tested here because it depends on input data; instead we test the
// wire-format primitives with small, hand-computed vectors.

use compresskit_codec::codec::frequency::{read_frequencies_exact, write_frequencies};

// ---------------------------------------------------------------------------
// Golden byte vectors (hand-computed, do not edit without updating comments)
// ---------------------------------------------------------------------------

/// write_frequencies([]) must produce count=0 encoded as 4 little-endian zero bytes.
/// Wire: 00 00 00 00
#[test]
fn golden_write_empty_table() {
    let freq: Vec<u32> = vec![];
    let mut out = Vec::new();
    write_frequencies(&mut out, &freq);
    assert_eq!(
        out,
        vec![0x00, 0x00, 0x00, 0x00],
        "empty frequency table must encode as 4-byte LE zero count"
    );
}

/// write_frequencies([1]) must produce count=1 then value=1, each as 4 LE bytes.
/// Wire: 01 00 00 00  01 00 00 00
#[test]
fn golden_write_single_entry() {
    let freq: Vec<u32> = vec![1];
    let mut out = Vec::new();
    write_frequencies(&mut out, &freq);
    assert_eq!(
        out,
        vec![
            0x01, 0x00, 0x00, 0x00, // count = 1
            0x01, 0x00, 0x00, 0x00, // freq[0] = 1
        ],
        "single-entry table wire bytes mismatch"
    );
}

/// write_frequencies([1, 2, 3]) must produce count=3 then each value as 4 LE bytes.
/// Wire: 03 00 00 00  01 00 00 00  02 00 00 00  03 00 00 00
#[test]
fn golden_write_three_entries() {
    let freq: Vec<u32> = vec![1, 2, 3];
    let mut out = Vec::new();
    write_frequencies(&mut out, &freq);
    assert_eq!(
        out,
        vec![
            0x03, 0x00, 0x00, 0x00, // count = 3
            0x01, 0x00, 0x00, 0x00, // freq[0] = 1
            0x02, 0x00, 0x00, 0x00, // freq[1] = 2
            0x03, 0x00, 0x00, 0x00, // freq[2] = 3
        ],
        "three-entry table wire bytes mismatch"
    );
}

/// Max u32 value must be preserved exactly in the wire format.
/// Wire: 01 00 00 00  FF FF FF FF
#[test]
fn golden_write_max_u32_entry() {
    let freq: Vec<u32> = vec![u32::MAX];
    let mut out = Vec::new();
    write_frequencies(&mut out, &freq);
    assert_eq!(
        out,
        vec![
            0x01, 0x00, 0x00, 0x00, // count = 1
            0xFF, 0xFF, 0xFF, 0xFF, // freq[0] = u32::MAX
        ],
        "max u32 entry wire bytes mismatch"
    );
}

/// Multi-byte count value: 256 entries.
/// Wire: 00 01 00 00  (count=256 LE)  followed by 256 * 4 bytes.
#[test]
fn golden_write_256_entries_count_encoding() {
    let freq: Vec<u32> = vec![0u32; 256];
    let mut out = Vec::new();
    write_frequencies(&mut out, &freq);
    // Only check the first 4 bytes (count field).
    assert_eq!(
        &out[..4],
        &[0x00, 0x01, 0x00, 0x00],
        "count=256 must encode as 00 01 00 00 in little-endian"
    );
    assert_eq!(out.len(), 4 + 256 * 4, "total length must be 4 + 256*4");
}

// ---------------------------------------------------------------------------
// Round-trip tests: write then read_frequencies_exact must recover original data
// ---------------------------------------------------------------------------

#[test]
fn roundtrip_three_entries_via_exact_reader() {
    let original: Vec<u32> = vec![7, 11, 13];
    let mut out = Vec::new();
    write_frequencies(&mut out, &original);

    let mut pos = 0;
    let recovered = read_frequencies_exact(
        &out,
        &mut pos,
        3,
        "truncated count",
        "truncated entries",
        "invalid count",
    )
    .expect("read_frequencies_exact must succeed on valid data");

    assert_eq!(recovered, original, "round-trip must recover original frequencies");
    assert_eq!(pos, out.len(), "reader must consume all bytes");
}

#[test]
fn roundtrip_exact_reader_rejects_count_mismatch() {
    let freq: Vec<u32> = vec![1, 2, 3];
    let mut out = Vec::new();
    write_frequencies(&mut out, &freq);

    let mut pos = 0;
    let result = read_frequencies_exact(
        &out,
        &mut pos,
        4, // wrong expected count
        "truncated count",
        "truncated entries",
        "invalid count",
    );

    assert!(
        result.is_err(),
        "read_frequencies_exact must reject count mismatch"
    );
}
