package codec

import (
	"fmt"
	"os"
)

// EncodeFile encodes a file using the provided encoder and writes to the output file.
// This is a convenience function for file-based encoding workflows.
func EncodeFile(enc Encoder, inputPath, outputPath string) error {
	input, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("cannot open input file: %s: %w", inputPath, err)
	}

	encoded, err := EncodeBuffer(enc, input)
	if err != nil {
		return err
	}

	if err := os.WriteFile(outputPath, encoded, 0o644); err != nil {
		return fmt.Errorf("cannot open output file: %s: %w", outputPath, err)
	}

	return nil
}

// DecodeFile decodes a file using the provided decoder and writes to the output file.
// This is a convenience function for file-based decoding workflows.
func DecodeFile(dec Decoder, inputPath, outputPath string) error {
	input, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("cannot open input file: %s: %w", inputPath, err)
	}

	decoded, err := DecodeBuffer(dec, input)
	if err != nil {
		return err
	}

	if err := os.WriteFile(outputPath, decoded, 0o644); err != nil {
		return fmt.Errorf("cannot open output file: %s: %w", outputPath, err)
	}

	return nil
}
