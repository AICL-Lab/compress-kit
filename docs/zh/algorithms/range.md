# 区间编码

区间编码是算术编码的**整数实现等价物**。它使用整数区间运算而非浮点运算，更适合生产环境，同时实现相同的压缩率。

## 工作原理

区间编码使用固定宽度整数维护区间 [low, low + range)。与算术编码的位输出不同，区间编码输出**字节**，显著提高 I/O 效率。

::: code-group

```cpp [C++]
void encode(const vector<uint8_t>& data,
            const uint32_t cumFreq[257]) {
    uint64_t low = 0;
    uint64_t range = MAX_RANGE;
    uint64_t total = cumFreq[256];
    
    for (uint8_t symbol : data) {
        uint32_t symLow = cumFreq[symbol];
        uint32_t symHigh = cumFreq[symbol + 1];
        
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

```go [Go]
func (rc *RangeCoder) Encode(data []byte, cumFreq [257]uint32) error {
    const maxRange uint64 = 1 << 32
    const minRange uint64 = 1 << 24
    total := cumFreq[256]
    
    for _, symbol := range data {
        symLow := cumFreq[symbol]
        symHigh := cumFreq[symbol+1]
        
        rc.range /= uint64(total)
        rc.low += uint64(symLow) * rc.range
        rc.range *= uint64(symHigh - symLow)
        
        for rc.range < minRange {
            rc.outputByte(byte(rc.low >> 56))
            rc.low <<= 8
            rc.range <<= 8
        }
    }
    rc.flush()
    return nil
}
```

```rust [Rust]
pub fn encode(&mut self, data: &[u8], cum_freq: &[u32; 257]) {
    const MAX_RANGE: u64 = 1u64 << 32;
    const MIN_RANGE: u64 = 1u64 << 24;
    let total = cum_freq[256];
    
    for &symbol in data {
        let sym_low = cum_freq[symbol as usize];
        let sym_high = cum_freq[symbol as usize + 1];
        
        self.range /= total as u64;
        self.low += (sym_low as u64) * self.range;
        self.range *= (sym_high - sym_low) as u64;
        
        while self.range < MIN_RANGE {
            self.output((self.low >> 56) as u8);
            self.low <<= 8;
            self.range <<= 8;
        }
    }
    self.flush();
}
```

:::

## 算术编码 vs 区间编码

| 方面 | 算术编码 | 区间编码 |
|------|----------|----------|
| 运算 | 浮点数 | 固定宽度整数 |
| 输出单位 | 位 | 字节 |
| I/O 效率 | 较低 | 较高 |
| 压缩率 | 几乎相同 | 几乎相同 |
| 专利状态 | 历史上有专利 | 无限制 |
| 生产使用 | 学术 | 工业标准 |

## 复杂度

| 方面 | 复杂度 | 说明 |
|------|--------|------|
| 时间（编码） | O(n) | 与算术编码类似 |
| 时间（解码） | O(n) | 字节级 I/O 更快 |
| 空间 | O(σ) | 累积频率表 |
| 精度 | 固定 | 64 位整数 |

## 性能特征

| 语言 | 文本压缩率 | 速度 | 内存 |
|------|------------|------|------|
| C++ | 1.90× | 58 MiB/s | 低 |
| Go | 1.90× | 50 MiB/s | 低 |
| Rust | 1.90× | 64 MiB/s | 低 |

## 适用场景

- ✅ **生产系统** — 最广泛部署的熵编码器
- ✅ **均衡工作负载** — 良好的速度和压缩率
- ✅ **视频编解码器** — H.264、HEVC 使用区间编码
- ✅ **压缩工具** — 用于现代归档工具

## 库 API

### Go

```go
import "github.com/LessUp/compress-kit/algorithms/range/go/rangecoder"

// 使用频率表创建编码器
freq := rangecoder.BuildFrequencyTable(data)
cumFreq := rangecoder.BuildCumulative(freq)

// 编码
encoded, err := rangecoder.Encode(data, cumFreq)

// 解码
decoded, err := rangecoder.Decode(encoded, cumFreq, len(data))
```

### Rust

```rust
use rangecoder;

let freq = rangecoder::build_frequency_table(&data);
let cum_freq = rangecoder::cumulative(&freq);

let encoded = rangecoder::encode(&data, &cum_freq)?;
let decoded = rangecoder::decode(&encoded, &cum_freq, data.len())?;
```

## 延伸阅读

- [算术编码](/zh/algorithms/arithmetic) — 浮点数等价实现
- [基准测试](/zh/benchmarks/results) — 性能对比
- [OpenSpec 架构规范](https://github.com/LessUp/compress-kit/tree/master/openspec/specs/core-architecture)

## 已知限制

::: warning 大文件性能问题

当前区间编码实现存在一个**已知的解码性能问题**：当文件大于 **500 KB** 时，解码操作可能会变得非常缓慢或出现卡顿。

**临时解决方案**：测试时请使用小于 100 KB 的文件。CI 管道中已使用 100 KB 测试文件进行区间编码验证。

**状态**：这是一个已知问题，已记录以便未来改进。编码操作对所有文件大小均正常工作。

:::
