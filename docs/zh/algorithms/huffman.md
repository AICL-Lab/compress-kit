# 霍夫曼编码

霍夫曼编码是一种无损数据压缩算法，使用**变长编码**来表示符号。出现频率较高的符号使用较短的编码，出现频率较低的符号使用较长的编码。

## 工作原理

::: code-group

```cpp [C++]
void buildHuffmanTree(const vector<uint32_t>& freq) {
    priority_queue<Node*, vector<Node*>, Compare> pq;
    
    // 创建叶子节点
    for (int i = 0; i < 256; i++) {
        if (freq[i] > 0) {
            pq.push(new Node(i, freq[i]));
        }
    }
    
    // 通过合并最低频率节点构建树
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
    
    // 创建叶子节点
    for i, f := range freq {
        if f > 0 {
            heap.Push(&pq, &Node{
                symbol: byte(i),
                freq:   f,
            })
        }
    }
    
    // 构建树
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

## 算法步骤

1. **频率统计**：统计输入中每个字节的出现频率
2. **构建树**：构建二叉树，频率较低的符号路径更深
3. **生成编码**：生成前缀编码（没有任何编码是另一个编码的前缀）
4. **编码**：将每个字节替换为对应的位编码
5. **解码**：根据位序列遍历树来重建原始数据

## 复杂度

| 方面 | 复杂度 | 说明 |
|------|--------|------|
| 时间（构建） | O(σ log σ) | σ = 字母表大小（最大 256） |
| 时间（编码） | O(n) | n = 输入大小 |
| 时间（解码） | O(n) | 单次遍历 |
| 空间 | O(σ) | 频率表 + 树 |

## 文件格式

| 字段 | 大小 | 描述 |
|------|------|------|
| Magic | 4 字节 | `HFMN` (0x48 0x46 0x4D 0x4E) |
| 频率表 | 257 × 4 字节 | 小端序 uint32 数组 |
| 编码数据 | 可变 | 位流，填充到字节边界 |

## 压缩效率

- **理论上界**：平均编码长度 ≥ 熵 H
- **霍夫曼上界**：H ≤ L < H + 1 位每符号
- **最适用于**：频率分布不均匀的数据

## 适用场景

- ✅ **文本文件** — 自然语言的字符频率不均匀
- ✅ **通用二进制数据** — 性能均衡
- ✅ **预处理** — 常用于其他变换之前
- ❌ **随机数据** — 接近 1× 压缩率（仅增加开销）

## 命令行使用

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

## 库使用

### Go

```go
package main

import (
    "github.com/LessUp/compress-kit/algorithms/huffman/go/huffman"
)

func main() {
    // 编码
    encoded, err := huffman.Encode(inputData)
    if err != nil {
        log.Fatal(err)
    }
    
    // 解码
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

## 延伸阅读

- [算术编码](/zh/algorithms/arithmetic) — 使用分数位的更优压缩
- [区间编码](/zh/algorithms/range) — 生产级优化的算术编码
- [算法对比](/zh/guide/algorithms)
