package main

import (
	"rangecoder"

	"github.com/LessUp/compress-kit/algorithms/shared/go/cli"
)

type RangeProcessor struct{}

func (p *RangeProcessor) EncodeFile(inputPath, outputPath string) error {
	return rangecoder.EncodeFile(inputPath, outputPath)
}

func (p *RangeProcessor) DecodeFile(inputPath, outputPath string) error {
	return rangecoder.DecodeFile(inputPath, outputPath)
}

func main() {
	cli.Run("rangecoder", &RangeProcessor{})
}
