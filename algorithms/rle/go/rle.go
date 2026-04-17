// Package rle provides Run-Length encoding and decoding implementations.
package rle

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

// MaxOutputSize is the maximum allowed output size (1 GiB) to prevent
// decompression bomb attacks.
const MaxOutputSize = 1 * 1024 * 1024 * 1024

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

	first, err := r.ReadByte()
	if err == io.EOF {
		return bw.Flush() // Empty file
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

	buf := make([]byte, 4096)
	var totalWritten uint64

	for {
		var count uint32
		if err := binary.Read(br, binary.LittleEndian, &count); err != nil {
			if err == io.EOF {
				break
			}
			if err == io.ErrUnexpectedEOF {
				return fmt.Errorf("RLE data truncated: cannot read complete count field")
			}
			return fmt.Errorf("failed to read count: %w", err)
		}
		if count == 0 {
			return fmt.Errorf("invalid RLE data: count should not be 0")
		}

		if totalWritten+uint64(count) > MaxOutputSize {
			return fmt.Errorf("output size limit exceeded (max %d bytes)", MaxOutputSize)
		}

		value, err := br.ReadByte()
		if err != nil {
			if err == io.EOF {
				return fmt.Errorf("RLE data truncated: missing value byte")
			}
			return fmt.Errorf("failed to read value: %w", err)
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
				return fmt.Errorf("failed to write decoded data: %w", err)
			}
			totalWritten += uint64(chunk)
			count -= uint32(chunk)
		}
	}

	return bw.Flush()
}

// EncodeFile is a convenience function for file-based encoding.
func EncodeFile(inputPath, outputPath string) error {
	in, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("cannot open input file: %s: %w", inputPath, err)
	}
	defer in.Close()

	out, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("cannot open output file: %s: %w", outputPath, err)
	}
	defer out.Close()

	return Encode(in, out)
}

// DecodeFile is a convenience function for file-based decoding.
func DecodeFile(inputPath, outputPath string) error {
	in, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("cannot open input file: %s: %w", inputPath, err)
	}
	defer in.Close()

	out, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("cannot open output file: %s: %w", outputPath, err)
	}
	defer out.Close()

	return Decode(in, out)
}
