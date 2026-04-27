# Huffman Coding

Huffman coding is a lossless data compression algorithm that uses **variable-length codes** to represent symbols. More frequent symbols get shorter codes, while less frequent symbols get longer codes.

## How It Works

::: code-group

```cpp [C++]
void buildHuffmanTree(const vector<uint32_t>& freq) {
    priority_queue<Node*, vector<Node*>, Compare> pq;
    
    // Create leaf nodes
    for (int i = 0; i < 256; i++) {
        if (freq[i] > 0) {
            pq.push(new Node(i, freq[i]));
        }
    }
    
    // Build tree by merging lowest frequency nodes
    while (pq.size() > 1) {
        Node* left = pq.top(); pq.pop();
        Node* right = pq.top(); pq.pop();
        
        Node* parent = new Node(0, left->freq + right->freq);
        parent->left = left;
        parent->right = right;
        
        pq.push(parent);
    }
    
    root = pq.top();
}
```

```go [Go]
func buildHuffmanTree(freq []uint32) *Node {
    pq := make(PriorityQueue, 0)
    heap.Init(&pq)
    
    // Create leaf nodes
    for i, f := range freq {
        if f > 0 {
            heap.Push(&pq, &Node{
                symbol: byte(i),
                freq:   f,
            })
        }
    }
    
    // Build tree
    for pq.Len() > 1 {
        left := heap.Pop(&pq).(*Node)
        right := heap.Pop(&pq).(*Node)
        
        parent := &Node{
            freq:  left.freq + right.freq,
            left:  left,
            right: right,
        }
        heap.Push(&pq, parent)
    }
    
    return heap.Pop(&pq).(*Node)
}
```

```rust [Rust]
fn build_huffman_tree(freq: &[u32; 256]) -> Option<Box<Node>> {
    use std::collections::BinaryHeap;
    
    let mut heap: BinaryHeap<_> = freq.iter()
        .enumerate()
        .filter(|(_, &f)| f > 0)
        .map(|(i, &f)| Node {
            symbol: i as u8,
            freq: f,
            left: None,
            right: None,
        })
        .collect();
    
    while heap.len() > 1 {
        let left = heap.pop().unwrap();
        let right = heap.pop().unwrap();
        
        let parent = Box::new(Node {
            symbol: 0,
            freq: left.freq + right.freq,
            left: Some(left),
            right: Some(right),
        });
        
        heap.push(parent);
    }
    
    heap.pop()
}
```

:::

## Algorithm Steps

1. **Frequency Analysis**: Count the frequency of each byte in the input
2. **Tree Construction**: Build a binary tree where lower frequency symbols have deeper paths
3. **Code Generation**: Generate prefix codes (no code is a prefix of another)
4. **Encoding**: Replace each byte with its corresponding bit code
5. **Decoding**: Traverse the tree according to bits to reconstruct original data

## Complexity

| Aspect | Complexity | Notes |
|--------|-----------|-------|
| Time (build) | O(σ log σ) | σ = alphabet size (256 max) |
| Time (encode) | O(n) | n = input size |
| Time (decode) | O(n) | Single pass |
| Space | O(σ) | Frequency table + tree |

## File Format

| Field | Size | Description |
|-------|------|-------------|
| Magic | 4 bytes | `HFMN` (0x48 0x46 0x4D 0x4E) |
| Frequency Table | 257 × 4 bytes | Little-endian uint32 array |
| Encoded Data | Variable | Bit stream, padded to byte boundary |

## Compression Efficiency

- **Theoretical lower bound**: Average code length ≥ entropy H
- **Huffman upper bound**: H ≤ L < H + 1 bit per symbol
- **Most effective** on data with uneven frequency distribution

## Use Cases

- ✅ **Text files** — Natural language has uneven character frequency
- ✅ **General binary data** — Balanced performance
- ✅ **Preprocessing** — Often used before other transforms
- ❌ **Random data** — Near-1× compression (only overhead)

## CLI Usage

::: code-group

```bash [C++]
./huffman_cpp encode input.txt output.huf
./huffman_cpp decode output.huf restored.txt
```

```bash [Go]
./huffman_go encode input.txt output.huf
./huffman_go decode output.huf restored.txt
```

```bash [Rust]
./huffman_rust encode input.txt output.huf
./huffman_rust decode output.huf restored.txt
```

:::

## Library Usage

### Go

```go
package main

import (
    "github.com/LessUp/compress-kit/algorithms/huffman/go/huffman"
)

func main() {
    // Encode
    encoded, err := huffman.Encode(inputData)
    if err != nil {
        log.Fatal(err)
    }
    
    // Decode
    decoded, err := huffman.Decode(encoded)
    if err != nil {
        log.Fatal(err)
    }
}
```

### Rust

```rust
use compress_kit::huffman::{encode, decode};

fn main() -> Result<(), Box<dyn Error>> {
    let encoded = encode(&input)?;
    let decoded = decode(&encoded)?;
    Ok(())
}
```

## Further Reading

- [Arithmetic Coding](/en/algorithms/arithmetic) — Better compression with fractional bits
- [Range Coder](/en/algorithms/range) — Production-optimized arithmetic coding
- [Algorithm Comparison](/en/guide/algorithms)
