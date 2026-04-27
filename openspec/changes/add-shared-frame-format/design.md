# Design: add-shared-frame-format

## Frame Layout

All CompressKit compressed files SHALL conform to the following layout:

```
┌──────────────────────────────────────────────────────────┐
│  Frame Header  (fixed 24 bytes)                          │
├──────────────────────────────────────────────────────────┤
│  Extension Blocks  (0 or more, variable)                 │
├──────────────────────────────────────────────────────────┤
│  Compressed Payload  (variable)                          │
├──────────────────────────────────────────────────────────┤
│  Frame Trailer  (8 bytes)                                │
└──────────────────────────────────────────────────────────┘
```

## Frame Header (24 bytes, all fields little-endian)

| Offset | Size | Field | Description |
|--------|------|-------|-------------|
| 0 | 4 | `magic` | `0x434B5A4D` (ASCII "CKZM") |
| 4 | 1 | `version` | Format version; current = `0x01` |
| 5 | 1 | `algo_id` | Algorithm identifier (see registry) |
| 6 | 2 | `flags` | Feature flags bitmap |
| 8 | 8 | `content_size` | Uncompressed size in bytes; `0` = unknown |
| 16 | 8 | `checksum` | xxHash-64 of uncompressed content; `0` if `FLAG_NO_CHECKSUM` set |

Total: 24 bytes.

### Magic: `0x434B5A4D`

Bytes: `43 4B 5A 4D` = "CKZM" (CompressKit Z-family Magic).

Chosen to be unique and human-readable when hexdumped.

### Algorithm ID Registry

| `algo_id` | Algorithm |
|-----------|-----------|
| `0x01` | Huffman |
| `0x02` | Arithmetic Coding |
| `0x03` | Range Coder |
| `0x04` | RLE |
| `0x05–0xEF` | Reserved for future CompressKit algorithms |
| `0xF0–0xFF` | Reserved for private/experimental use |

### Flags Bitmap (uint16 LE)

| Bit | Name | Meaning |
|-----|------|---------|
| 0 | `FLAG_HAS_DICT` | Extension block with dictionary_id is present |
| 1 | `FLAG_HAS_META` | One or more skippable metadata blocks follow extensions |
| 2 | `FLAG_NO_CHECKSUM` | Checksum field is zero and MUST be ignored |
| 3–15 | Reserved | MUST be 0; decoders MUST reject frames with unknown flag bits set |

## Extension Blocks

Extension blocks appear between the frame header and the compressed payload. Each block has a 4-byte type tag and 4-byte length prefix (little-endian uint32).

### Dictionary ID Block (`type = 0x44494354`, "DICT")

```
┌──────────┬──────────┬──────────────────┐
│ type (4) │ len  (4) │ dict_id  (8)     │
└──────────┴──────────┴──────────────────┘
```

- `len` = 8 (the dict_id field only)
- `dict_id`: uint64 LE identifier matching a pre-shared dictionary

### Skippable Metadata Block (`type = 0x536B6970`, "Skip")

```
┌──────────┬──────────┬─────────────────────┐
│ type (4) │ len  (4) │ opaque payload (len) │
└──────────┴──────────┴─────────────────────┘
```

- Decoders that do not understand this block MUST skip `len` bytes and continue.
- `len` MUST be ≤ 65535 bytes.

## Frame Trailer (8 bytes)

| Offset | Size | Field | Description |
|--------|------|-------|-------------|
| 0 | 4 | `end_magic` | `0x454E4400` ("END\0") |
| 4 | 4 | `payload_crc32` | CRC-32/ISO-HDLC of compressed payload bytes only |

The trailer allows quick detection of payload corruption independent of the content checksum.

## Integrity Check Semantics

The frame contains two independent integrity checks that serve different purposes:

| Check | Location | Covers | Disabled by |
|-------|----------|--------|-------------|
| `payload_crc32` (CRC-32/ISO-HDLC) | Trailer | Compressed payload bytes | — (always present) |
| `checksum` (xxHash-64) | Header | Uncompressed content bytes | `FLAG_NO_CHECKSUM` |

### Rules

1. **Trailer CRC is always mandatory.** `FLAG_NO_CHECKSUM` affects only the header xxHash-64 field; the trailer `payload_crc32` MUST always be computed and checked.
2. **Check order:** Decoders SHALL verify `payload_crc32` first (before or immediately after decompression), then verify `checksum` against the decompressed bytes (unless `FLAG_NO_CHECKSUM` is set).
3. **Error handling:** Any checksum mismatch SHALL return `ERR_CORRUPT` and produce no output. The first failing check ends processing; the second check is skipped.

## Backward Compatibility

- Version `0x00` is invalid (reserved to detect zero-initialised buffers).
- Decoders MUST check `version == 0x01` before parsing; reject with `ERR_VERSION_UNSUPPORTED` otherwise.
- New flag bits MUST only be assigned in future spec versions. Current decoders MUST reject frames with unknown flags to prevent silent misinterpretation.

## Checksum Selection: xxHash-64

xxHash-64 is chosen for:
- Speed: ~10 GB/s on modern hardware (faster than SHA, CRC-32)
- Zero external dependencies in Go and Rust (crate/module available)
- Well-defined portable output across platforms

Reference implementation: https://github.com/Cyan4973/xxHash

## Interaction with Existing Frequency Table Format

`REQ-ARCH-003` (frequency table) is unchanged. The frequency table is part of the compressed payload, immediately following the frame header + extension blocks. The frame format wraps the payload; it does not replace per-algorithm internal structures.

## Migration from Pre-Frame Files

Existing files encoded without this header are identified by the absence of the `CKZM` magic. A migration utility task (tasks.md § Phase D) will re-wrap legacy files.
