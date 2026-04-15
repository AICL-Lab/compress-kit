package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

// Run-Length encoding implementation.
// Format: repeatedly write 4-byte little-endian count + 1-byte value until input ends.
// All three languages (C++/Go/Rust) use the same format for cross-decoding and benchmarking.

// Maximum output size limit (1 GiB) to prevent decompression bomb attacks
const maxOutputSize = 1 * 1024 * 1024 * 1024

// RLEEncodeFile performs Run-Length encoding on the entire file.
func rleEncodeFile(inputPath, outputPath string) error {
	in, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("cannot open input file for reading: %s: %w", inputPath, err)
	}
	defer in.Close()

	out, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("cannot open output file for writing: %s: %w", outputPath, err)
	}
	defer out.Close()

	r := bufio.NewReader(in)
	w := bufio.NewWriter(out)

	first, err := r.ReadByte()
	if err == io.EOF {
		// Empty file: encoding result is also empty.
		if err := w.Flush(); err != nil {
			return err
		}
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}

	current := first
	var count uint32 = 1

	for {
		b, err := r.ReadByte()
		if err == io.EOF {
			// Write last run
			if err := writeRun(w, count, current); err != nil {
				return fmt.Errorf("failed to write RLE data: %w", err)
			}
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read input file: %w", err)
		}

		if b == current && count < ^uint32(0) {
			count++
		} else {
			if err := writeRun(w, count, current); err != nil {
				return fmt.Errorf("failed to write RLE data: %w", err)
			}
			current = b
			count = 1
		}
	}

	if err := w.Flush(); err != nil {
		return err
	}
	return nil
}

// RLEEncodeFile performs Run-Length encoding on the entire file.
func RLEEncodeFile(inputPath, outputPath string) {
	if err := rleEncodeFile(inputPath, outputPath); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

// writeRun writes a single (count, value) pair to the output stream.
func writeRun(w *bufio.Writer, count uint32, value byte) error {
	// Write 4-byte little-endian count
	if err := binary.Write(w, binary.LittleEndian, count); err != nil {
		return err
	}
	// Write 1-byte value
	if err := w.WriteByte(value); err != nil {
		return err
	}
	return nil
}

// RLEDecodeFile decodes an RLE encoded file back to original byte sequence.
func rleDecodeFile(inputPath, outputPath string) error {
	in, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("cannot open input file for reading: %s: %w", inputPath, err)
	}
	defer in.Close()

	out, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("cannot open output file for writing: %s: %w", outputPath, err)
	}
	defer out.Close()

	r := bufio.NewReader(in)
	w := bufio.NewWriter(out)

	buf := make([]byte, 4096)
	var totalWritten uint64

	for {
		var count uint32
		if err := binary.Read(r, binary.LittleEndian, &count); err != nil {
			if err == io.EOF {
				// Normal EOF
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

		// Check output size limit
		if totalWritten+uint64(count) > maxOutputSize {
			return fmt.Errorf("output size limit exceeded (max %d bytes)", maxOutputSize)
		}

		value, err := r.ReadByte()
		if err != nil {
			if err == io.EOF {
				return fmt.Errorf("RLE data truncated: missing value byte")
			}
			return fmt.Errorf("failed to read value: %w", err)
		}

		// Expand (count, value) to output
		for count > 0 {
			chunk := int(count)
			if chunk > len(buf) {
				chunk = len(buf)
			}
			for i := 0; i < chunk; i++ {
				buf[i] = value
			}
			if _, err := w.Write(buf[:chunk]); err != nil {
				return fmt.Errorf("failed to write decoded data: %w", err)
			}
			totalWritten += uint64(chunk)
			count -= uint32(chunk)
		}
	}

	if err := w.Flush(); err != nil {
		return err
	}
	return nil
}

// RLEDecodeFile decodes an RLE encoded file back to original byte sequence.
func RLEDecodeFile(inputPath, outputPath string) {
	if err := rleDecodeFile(inputPath, outputPath); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func main() {
	if len(os.Args) != 4 {
		fmt.Fprintf(os.Stderr, "usage: %s encode|decode input output\n", os.Args[0])
		os.Exit(1)
	}

	mode := os.Args[1]
	inputPath := os.Args[2]
	outputPath := os.Args[3]

	var err error
	switch mode {
	case "encode":
		err = rleEncodeFile(inputPath, outputPath)
	case "decode":
		err = rleDecodeFile(inputPath, outputPath)
	default:
		fmt.Fprintln(os.Stderr, "unknown mode, expected encode or decode")
		os.Exit(1)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
