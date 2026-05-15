// Package huffman provides Huffman encoding and decoding implementations.
package huffman

import (
	"bufio"
	"container/heap"
	"fmt"
	"io"
	"os"

	"github.com/LessUp/compress-kit/algorithms/shared/go/codec"
)

const (
	// SymbolLimit is the number of possible symbols (256 bytes + 1 EOF symbol).
	// This is an alias for codec.SymbolLimit for backward compatibility.
	SymbolLimit = codec.SymbolLimit
	// EOFSymbol is the symbol index used to mark end-of-stream.
	// This is an alias for codec.EOFSymbol for backward compatibility.
	EOFSymbol = codec.EOFSymbol
	// MaxInputSize is the maximum allowed input file size (4 GiB) to prevent
	// frequency overflow and decompression bomb attacks.
	MaxInputSize = 4 * 1024 * 1024 * 1024
)

// Node represents a node in the Huffman tree.
type Node struct {
	Symbol uint32
	Freq   uint64
	Left   *Node
	Right  *Node
}

// IsLeaf returns true if the node has no children.
func (n *Node) IsLeaf() bool {
	return n.Left == nil && n.Right == nil
}

type nodeHeap []*Node

func (h nodeHeap) Len() int { return len(h) }
func (h nodeHeap) Less(i, j int) bool {
	if h[i].Freq != h[j].Freq {
		return h[i].Freq < h[j].Freq
	}
	return h[i].Symbol < h[j].Symbol
}
func (h nodeHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *nodeHeap) Push(x any) {
	*h = append(*h, x.(*Node))
}

func (h *nodeHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

// BuildTree constructs a Huffman tree from the given frequency table.
func BuildTree(freq []uint32) *Node {
	h := &nodeHeap{}
	heap.Init(h)
	for s := 0; s < SymbolLimit; s++ {
		if freq[s] == 0 {
			continue
		}
		n := &Node{
			Symbol: uint32(s),
			Freq:   uint64(freq[s]),
		}
		heap.Push(h, n)
	}
	if h.Len() == 0 {
		return &Node{Symbol: uint32(EOFSymbol), Freq: 1}
	}
	if h.Len() == 1 {
		only := heap.Pop(h).(*Node)
		parent := &Node{Symbol: only.Symbol, Freq: only.Freq, Left: only, Right: nil}
		heap.Push(h, parent)
	}
	for h.Len() > 1 {
		a := heap.Pop(h).(*Node)
		b := heap.Pop(h).(*Node)
		minSymbol := a.Symbol
		if b.Symbol < minSymbol {
			minSymbol = b.Symbol
		}
		parent := &Node{
			Symbol: minSymbol,
			Freq:   a.Freq + b.Freq,
			Left:   a,
			Right:  b,
		}
		heap.Push(h, parent)
	}
	return heap.Pop(h).(*Node)
}

// BuildFrequenciesFromFile reads the file and counts byte frequencies.
func BuildFrequenciesFromFile(path string) ([]uint32, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("cannot open input file: %s: %w", path, err)
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return nil, fmt.Errorf("cannot stat input file: %s: %w", path, err)
	}
	if stat.Size() > MaxInputSize {
		return nil, fmt.Errorf("input file too large (max %d bytes)", MaxInputSize)
	}

	data, err := io.ReadAll(bufio.NewReader(f))
	if err != nil {
		return nil, fmt.Errorf("cannot read input file: %s: %w", path, err)
	}
	return codec.BuildFrequencies(data), nil
}

// WriteFrequencies serializes a frequency table to the writer.
func WriteFrequencies(w io.Writer, freq []uint32) error {
	return codec.WriteFrequencies(w, freq)
}

// ReadFrequencies deserializes a frequency table from the reader.
func ReadFrequencies(r io.Reader) ([]uint32, error) {
	return codec.ReadFrequenciesExact(r, SymbolLimit)
}

// BuildCodes generates Huffman codes for each symbol by traversing the tree.
func BuildCodes(node *Node, codes []string, prefix []byte) {
	if node == nil {
		return
	}
	if node.IsLeaf() {
		if len(prefix) == 0 {
			codes[int(node.Symbol)] = "0"
		} else {
			codes[int(node.Symbol)] = string(append([]byte(nil), prefix...))
		}
		return
	}
	prefix = append(prefix, '0')
	BuildCodes(node.Left, codes, prefix)
	prefix[len(prefix)-1] = '1'
	BuildCodes(node.Right, codes, prefix)
}

// Encode reads from input, writes the encoded output to w.
func Encode(input io.Reader, w io.Writer) error {
	data, err := io.ReadAll(input)
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}
	if int64(len(data)) > MaxInputSize {
		return fmt.Errorf("input too large (max %d bytes)", MaxInputSize)
	}

	freq := codec.BuildFrequencies(data)

	root := BuildTree(freq)
	codes := make([]string, SymbolLimit)
	BuildCodes(root, codes, nil)

	if _, err := w.Write([]byte{'H', 'F', 'M', 'N'}); err != nil {
		return err
	}
	if err := WriteFrequencies(w, freq); err != nil {
		return err
	}
	bw := codec.NewBitWriter(w)
	for _, b := range data {
		code := codes[int(b)]
		for i := 0; i < len(code); i++ {
			bit := 0
			if code[i] == '1' {
				bit = 1
			}
			if err := bw.WriteBit(bit); err != nil {
				return err
			}
		}
	}
	eofCode := codes[EOFSymbol]
	for i := 0; i < len(eofCode); i++ {
		bit := 0
		if eofCode[i] == '1' {
			bit = 1
		}
		if err := bw.WriteBit(bit); err != nil {
			return err
		}
	}
	return bw.Flush()
}

// Decode reads from r and writes the decoded output to w.
func Decode(r io.Reader, w io.Writer) error {
	br := bufio.NewReader(r)

	magic := make([]byte, 4)
	if _, err := io.ReadFull(br, magic); err != nil || magic[0] != 'H' || magic[1] != 'F' || magic[2] != 'M' || magic[3] != 'N' {
		return codec.NewError(codec.KindCorrupt, "invalid input file format")
	}

	freq, err := ReadFrequencies(br)
	if err != nil {
		return err
	}
	root := BuildTree(freq)
	if root == nil {
		return codec.NewError(codec.KindCorrupt, "decode failed")
	}

	bw := bufio.NewWriter(w)
	bitReader := codec.NewBitReader(br)
	node := root
	sawEOF := false
	var totalWritten uint64

	for {
		bit := bitReader.ReadBit()
		if bit == 0 {
			if node.Left != nil {
				node = node.Left
			} else {
				return codec.NewError(codec.KindCorrupt, "input data corrupted or truncated")
			}
		} else {
			if node.Right != nil {
				node = node.Right
			} else {
				return codec.NewError(codec.KindCorrupt, "input data corrupted or truncated")
			}
		}
		if node.IsLeaf() {
			if node.Symbol == uint32(EOFSymbol) {
				sawEOF = true
				break
			}
			totalWritten++
			if totalWritten > codec.MaxOutputSize {
				return codec.NewError(codec.KindSizeLimit, "output size limit exceeded")
			}
			if err := bw.WriteByte(byte(node.Symbol)); err != nil {
				return err
			}
			node = root
		}
		if bitReader.EOF() && node == root {
			break
		}
	}

	if !sawEOF {
		return codec.NewError(codec.KindTruncated, "input data corrupted or truncated")
	}
	return bw.Flush()
}

// EncodeFile is a convenience function for file-based encoding.
func EncodeFile(inputPath, outputPath string) error {
	return codec.EncodeFile(NewStreamingEncoder(), inputPath, outputPath)
}

// DecodeFile is a convenience function for file-based decoding.
func DecodeFile(inputPath, outputPath string) error {
	return codec.DecodeFile(NewStreamingDecoder(), inputPath, outputPath)
}
