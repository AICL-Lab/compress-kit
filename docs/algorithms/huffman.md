# Huffman

Huffman 是一种无损数据压缩算法，使用**变长编码**来表示符号。出现频率较高的符号使用较短的编码，出现频率较低的符号使用较长的编码。

## 工作原理

```cpp
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
- **Huffman 上界**：H ≤ L < H + 1 位每符号
- **最适用于**：频率分布不均匀的数据

## 适用场景

- ✅ **文本文件** — 自然语言的字符频率不均匀
- ✅ **通用二进制数据** — 性能均衡
- ✅ **预处理** — 常用于其他变换之前
- ❌ **随机数据** — 接近 1× 压缩率（仅增加开销）

## 命令行使用

```bash
./build/huffman_cpp encode input.txt output.huf
./build/huffman_cpp decode output.huf restored.txt
```

## 延伸阅读

- [Arithmetic](/algorithms/arithmetic) — 使用分数位的更优压缩
- [Range Coder](/algorithms/range) — 生产级优化实现
- [算法对比](/guide/algorithms)
