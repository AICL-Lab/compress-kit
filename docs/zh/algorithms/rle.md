# RLE

RLE（行程编码）是**最简单的压缩算法**，将连续重复的字节替换为 `(计数, 值)` 对。它速度极快，对具有长连续相同值的数据效果非常好。

## 工作原理

::: code-group

```cpp [C++]
vector<uint8_t> encode(const vector<uint8_t>& data) {
    vector<uint8_t> result;
    
    for (size_t i = 0; i < data.size();) {
        uint8_t current = data[i];
        uint32_t count = 1;
        
        // 统计连续相同字节
        while (i + count < data.size() && 
               data[i + count] == current && 
               count < UINT32_MAX) {
            count++;
        }
        
        // 写入计数（4 字节，小端序）+ 值
        result.push_back(count & 0xFF);
        result.push_back((count >> 8) & 0xFF);
        result.push_back((count >> 16) & 0xFF);
        result.push_back((count >> 24) & 0xFF);
        result.push_back(current);
        
        i += count;
    }
    
    return result;
}
```

```go [Go]
func Encode(data []byte) []byte {
    var result bytes.Buffer
    
    for i := 0; i < len(data); {
        current := data[i]
        count := uint32(1)
        
        for i+int(count) < len(data) && 
            data[i+int(count)] == current && 
            count < math.MaxUint32 {
            count++
        }
        
        // 将计数写为 uint32（小端序）+ 值
        binary.Write(&result, binary.LittleEndian, count)
        result.WriteByte(current)
        
        i += int(count)
    }
    
    return result.Bytes()
}
```

```rust [Rust]
pub fn encode(data: &[u8]) -> Vec<u8> {
    let mut result = Vec::new();
    let mut i = 0;
    
    while i < data.len() {
        let current = data[i];
        let mut count: u32 = 1;
        
        while i + count as usize < data.len() && 
              data[i + count as usize] == current && 
              count < u32::MAX {
            count += 1;
        }
        
        // 写入计数（小端序）+ 值
        result.extend_from_slice(&count.to_le_bytes());
        result.push(current);
        
        i += count as usize;
    }
    
    result
}
```

:::

## 文件格式

| 字段 | 大小 | 描述 |
|------|------|------|
| Magic | 4 字节 | `RLE\x00` (0x52 0x4C 0x45 0x00) |
| 计数 | 4 字节 | 小端序无符号整数（行程长度） |
| 值 | 1 字节 | 重复的字节值 |

## 复杂度

| 方面 | 复杂度 | 说明 |
|------|--------|------|
| 时间（编码） | O(n) | 单次遍历 |
| 时间（解码） | O(n) | 单次遍历，非常快 |
| 空间 | O(1) | 无辅助结构 |

## 性能

| 数据类型 | 压缩率 | 速度 |
|----------|--------|------|
| 重复数据 (10 MB) | 25× | 300+ MiB/s |
| 文本 | 1.1× | 快 |
| 随机 | 0.2×（膨胀） | 快 |

::: warning 最坏情况
对于随机数据，RLE 将每个字节扩展为 5 字节（4 字节计数 + 1 字节值），导致 **5 倍大小增加**。
:::

## 适用场景

- ✅ **位图图像** — 长连续相同颜色
- ✅ **日志文件** — 重复模式
- ✅ **预处理** — 在 BWT 或其他变换之前
- ✅ **传真传输** — 标准压缩方法
- ❌ **随机数据** — 严重膨胀
- ❌ **已压缩数据** — 无收益

## 常见用途

### 作为预处理

RLE 常作为更复杂压缩流程的预处理步骤：

```
原始数据 → BWT → MTF → RLE → Arithmetic → 压缩数据
```

这种组合（Burrows-Wheeler + Move-to-Front + RLE + Arithmetic）是 **bzip2** 的基础。

### 图像格式中

- **BMP**：简单 RLE 变体
- **PCX**：RLE 压缩
- **TIFF**：可选 RLE packbits

## 与其他算法对比

| 算法 | 初始化开销 | 压缩率 | 速度 |
|------|------------|--------|------|
| RLE | 无 | 可变 | 最快 |
| Huffman | O(σ log σ) | 中等 | 快 |
| Range Coder | O(n) | 高 | 快 |

## 延伸阅读

- [算法对比](/zh/guide/algorithms) — 完整对比矩阵
- [基准测试](/zh/benchmarks/results) — 性能数据
- [Burrows-Wheeler 变换](https://en.wikipedia.org/wiki/Burrows%E2%80%93Wheeler_transform) — 常与 RLE 配合使用
