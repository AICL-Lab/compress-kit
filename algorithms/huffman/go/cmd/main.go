package main

import (
	"fmt"
	"os"

	"huffman"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Fprintf(os.Stderr, "Usage: %s encode|decode input output\n", os.Args[0])
		os.Exit(1)
	}
	mode := os.Args[1]
	inputPath := os.Args[2]
	outputPath := os.Args[3]

	var err error

	switch mode {
	case "encode":
		err = huffman.EncodeFile(inputPath, outputPath)
	case "decode":
		err = huffman.DecodeFile(inputPath, outputPath)
	default:
		fmt.Fprintln(os.Stderr, "Unknown mode")
		os.Exit(1)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
