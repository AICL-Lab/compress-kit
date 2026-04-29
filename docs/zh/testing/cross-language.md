# 跨语言测试

CompressKit 的核心特性之一是**所有语言实现使用相同的二进制格式**。这使得跨语言编码和解码成为可能。

## 测试方法

我们的 CI 流水线验证所有 12 种可能的编码/解码组合的正确性：

```
4 种算法 × 3 种语言 = 12 个实现
12 个编码器 × 3 个解码器 = 36 个跨语言配对
```

## 测试矩阵

| 编码 ↓ / 解码 → | C++ | Go | Rust |
|----------------|-----|-----|------|
| **C++** | ✅ | ✅ | ✅ |
| **Go** | ✅ | ✅ | ✅ |
| **Rust** | ✅ | ✅ | ✅ |

## 工作原理

每种算法使用**严格定义的二进制格式**：

### 霍夫曼编码示例

```
+----------+-------------------+---------------+
|  魔数    | 频率表            | 编码数据      |
| (4 字节) | (257 × 4 字节)    | (可变长度)    |
+----------+-------------------+---------------+
| "HFMN"   | uint32[257]       | 比特流        |
+----------+-------------------+---------------+
```

三种实现都写入完全相同的结构：
- 相同的魔数字节 (`HFMN`)
- 相同的字节序（小端）
- 相同的表格布局
- 相同的比特流编码

## 运行跨语言测试

### 手动验证

```bash
# 生成测试数据
dd if=/dev/urandom of=test.bin bs=1M count=1

# C++ 编码 → Go 解码
./algorithms/huffman/cpp/huffman_cpp encode test.bin encoded.huf
./algorithms/huffman/go/huffman_go decode encoded.huf restored.bin
diff test.bin restored.bin  # 应该没有输出

# Go 编码 → Rust 解码
./algorithms/huffman/go/huffman_go encode test.bin encoded.huf
./algorithms/huffman/rust/huffman_rust decode encoded.huf restored.bin
diff test.bin restored.bin

# 尝试任意组合：
# C++ ↔ Go, C++ ↔ Rust, Go ↔ Rust
# 适用于所有 4 种算法
```

### 自动测试

```bash
make test
```

这会在所有算法上运行完整的跨语言测试套件。

## 文件格式规范

详细的格式规范可在 specs 目录找到：

- [核心架构规范](https://github.com/LessUp/compress-kit/tree/master/openspec/specs/core-architecture)
- [跨语言测试规范](https://github.com/LessUp/compress-kit/tree/master/openspec/specs/cross-language-testing)

## 跨语言的重要性

1. **数据可移植性**：在一个环境中编码，在另一个环境中解码
2. **渐进迁移**：逐步从一种语言迁移到另一种
3. **测试验证**：多个独立实现可以发现错误
4. **学习对比**：比较同一算法的不同实现方式

## 已知限制

| 算法 | 限制 | 变通方法 |
|------|------|----------|
| 区间编码 | 解码性能 > 500KB 时下降 | 使用较小的数据块 |

## 报告问题

如果你发现跨语言不兼容：

1. 使用 `diff` 验证二进制不匹配
2. 使用 `xxd -l 20 encoded.huf` 检查文件头
3. 在 [GitHub Issues](https://github.com/LessUp/compress-kit/issues) 报告并附上：
   - 算法名称
   - 编码语言
   - 解码语言
   - 输入文件类型
