# Range Coder

Range Coder 是 Arithmetic 的**整数实现等价物**。它使用整数区间运算而非浮点运算，更适合生产环境，同时实现相同的压缩率。

## 工作原理

Range Coder 使用固定宽度整数维护区间 [low, low + range)。与 Arithmetic 的位输出不同，Range Coder 输出**字节**，显著提高 I/O 效率。

```cpp
void encode(const vector<uint8_t>& data,
            const uint32_t cumFreq[257]) {
    uint64_t low = 0;
    uint64_t range = MAX_RANGE;
    uint64_t total = cumFreq[256];
    
    for (uint8_t symbol : data) {
        uint32_t symLow = cumFreq[symbol];
        uint32_t symHigh = cumFreq[symbol + 1];
        
        // total 必须非零：空频率表是非法输入，需在编码前校验
        range /= total;
        low += symLow * range;
        range *= (symHigh - symLow);
        
        // 重归一化
        while (range < MIN_RANGE) {
            output_byte(low >> 56);
            low <<= 8;
            range <<= 8;
        }
    }
    // 输出最终字节
    flush(low);
}
```

## Arithmetic vs Range Coder

| 方面 | Arithmetic | Range Coder |
|------|----------|----------|
| 运算 | 浮点数 | 固定宽度整数 |
| 输出单位 | 位 | 字节 |
| I/O 效率 | 较低 | 较高 |
| 压缩率 | 几乎相同 | 几乎相同 |
| 专利状态 | 历史上有专利 | 无限制 |
| 生产使用 | 学术 | 工业标准 |

## 文件格式

| 字段 | 大小 | 描述 |
|------|------|------|
| 魔数 | 4 字节 | `RCNC` (0x52 0x43 0x4E 0x43) |
| 频率表大小 | 4 字节 | 小端 uint32（始终 257） |
| 频率表 | 257 × 4 字节 | 小端 uint32 数组（符号 0-255 + EOF） |
| 编码数据 | 可变 | 字节流（重归一化区间输出） |

通用结构参见 [架构设计 - 二进制格式规范](/architecture/#各算法格式)。

## 复杂度

| 方面 | 复杂度 | 说明 |
|------|--------|------|
| 时间（编码） | O(n) | 与 Arithmetic 类似 |
| 时间（解码） | O(n) | 字节级 I/O 更快 |
| 空间 | O(σ) | 累积频率表 |
| 精度 | 固定 | 64 位整数 |

## 性能特征

| 输入类型 | 压缩率 | 速度 | 内存 |
|----------|--------|------|------|
| 文本 | 1.90× | 58 MiB/s | 低 |

## 适用场景

- ✅ **生产系统** — 最广泛部署的熵编码器
- ✅ **均衡工作负载** — 良好的速度和压缩率
- ✅ **视频编解码器** — H.264、HEVC 使用 Range Coder
- ✅ **压缩工具** — 用于现代归档工具

## 命令行使用

```bash
./build/rangecoder_cpp encode input.bin output.rcnc
./build/rangecoder_cpp decode output.rcnc restored.bin
```

## 延伸阅读

- [Arithmetic](/algorithms/arithmetic) — 浮点数等价实现
- [基准测试](/benchmarks/results) — 性能对比

## 已知限制

::: warning 大文件性能问题

当前 Range Coder 实现存在一个**已知的解码性能问题**：当文件大于 **500 KiB** 时，解码操作可能会变得非常缓慢或出现卡顿。

**临时解决方案**：测试时请使用小于 100 KiB 的文件。CI 管道中已使用 100 KiB 测试文件进行 Range Coder 验证。

**状态**：这是一个已知问题，已记录以便未来改进。编码操作对所有文件大小均正常工作。

:::
