// Package rle provides Run-Length encoding and decoding implementations.
package rle

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/LessUp/compress-kit/algorithms/shared/go/codec"
)

// MaxOutputSize is the maximum allowed output size (1 GiB) to prevent
// decompression bomb attacks.
const MaxOutputSize = 1 * 1024 * 1024 * 1024

// rleMagic is the 4-byte magic number for RLE format identification
const rleMagic = "RLE\x00"

// writeRun writes a single (count, value) pair to the output stream.
func writeRun(w *bufio.Writer, count uint32, value byte) error {
	if err := binary.Write(w, binary.LittleEndian, count); err != nil {
		return err
	}
	if err := w.WriteByte(value); err != nil {
		return err
	}
	return nil
}

// Encode reads from input and writes the RLE encoded output to w.
func Encode(input io.Reader, w io.Writer) error {
	r := bufio.NewReader(input)
	bw := bufio.NewWriter(w)

	// Write magic number
	if _, err := bw.WriteString(rleMagic); err != nil {
		return fmt.Errorf("failed to write RLE magic: %w", err)
	}

	first, err := r.ReadByte()
	if err == io.EOF {
		return bw.Flush() // Empty file (magic only)
	}
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}

	current := first
	var count uint32 = 1

	for {
		b, err := r.ReadByte()
		if err == io.EOF {
			if err := writeRun(bw, count, current); err != nil {
				return fmt.Errorf("failed to write RLE data: %w", err)
			}
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}

		if b == current && count < ^uint32(0) {
			count++
		} else {
			if err := writeRun(bw, count, current); err != nil {
				return fmt.Errorf("failed to write RLE data: %w", err)
			}
			current = b
			count = 1
		}
	}

	return bw.Flush()
}

// Decode reads from r and writes the decoded output to w.
func Decode(r io.Reader, w io.Writer) error {
	br := bufio.NewReader(r)
	bw := bufio.NewWriter(w)

	// Verify magic number
	magic := make([]byte, 4)
	if _, err := io.ReadFull(br, magic); err != nil {
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return codec.NewError(codec.KindTruncated, "cannot read magic number")
		}
		return codec.WrapError(codec.KindTruncated, "failed to read magic", err)
	}
	if string(magic) != rleMagic {
		return codec.NewError(codec.KindCorrupt, "invalid RLE file: bad magic number")
	}

	buf := make([]byte, 4096)
	var totalWritten uint64

	for {
		var count uint32
		if err := binary.Read(br, binary.LittleEndian, &count); err != nil {
			if err == io.EOF {
				break
			}
			if err == io.ErrUnexpectedEOF {
				return codec.NewError(codec.KindTruncated, "RLE data truncated: cannot read complete count field")
			}
			return codec.WrapError(codec.KindTruncated, "failed to read count", err)
		}
		if count == 0 {
			return codec.NewError(codec.KindCorrupt, "invalid RLE data: count should not be 0")
		}

		if totalWritten+uint64(count) > MaxOutputSize {
			return codec.NewError(codec.KindSizeLimit, fmt.Sprintf("output size limit exceeded (max %d bytes)", MaxOutputSize))
		}

		value, err := br.ReadByte()
		if err != nil {
			if err == io.EOF {
				return codec.NewError(codec.KindTruncated, "RLE data truncated: missing value byte")
			}
			return codec.WrapError(codec.KindTruncated, "failed to read value", err)
		}

		for count > 0 {
			chunk := int(count)
			if chunk > len(buf) {
				chunk = len(buf)
			}
			for i := 0; i < chunk; i++ {
				buf[i] = value
			}
			if _, err := bw.Write(buf[:chunk]); err != nil {
				return codec.WrapError(codec.KindCorrupt, "failed to write decoded data", err)
			}
			totalWritten += uint64(chunk)
			count -= uint32(chunk)
		}
	}

	return bw.Flush()
}

// EncodeFile is a convenience function for file-based encoding.
func EncodeFile(inputPath, outputPath string) error {
	input, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("cannot open input file: %s: %w", inputPath, err)
	}

	encoded, err := codec.EncodeBuffer(NewStreamingEncoder(), input)
	if err != nil {
		return err
	}

	if err := os.WriteFile(outputPath, encoded, 0o644); err != nil {
		return fmt.Errorf("cannot open output file: %s: %w", outputPath, err)
	}

	return nil
}

// DecodeFile is a convenience function for file-based decoding.
func DecodeFile(inputPath, outputPath string) error {
	input, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("cannot open input file: %s: %w", inputPath, err)
	}

	decoded, err := codec.DecodeBuffer(NewStreamingDecoder(), input)
	if err != nil {
		return err
	}

	if err := os.WriteFile(outputPath, decoded, 0o644); err != nil {
		return fmt.Errorf("cannot open output file: %s: %w", outputPath, err)
	}

	return nil
}
