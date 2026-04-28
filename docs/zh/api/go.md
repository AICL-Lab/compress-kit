# Go 库 API 参考

另请参阅: [Streaming API](/zh/api/streaming)

所有 Go 实现都提供库 API 和命令行工具。现在每种算法都遵循统一的分层模式：

- 传统的 `Encode`、`Decode`、`EncodeFile`、`DecodeFile`
- 共享 streaming 实现 `NewStreamingEncoder()` / `NewStreamingDecoder()`
- 共享 buffer helper：`github.com/LessUp/compress-kit/algorithms/shared/go/codec`

## Huffman 编码 (`algorithms/huffman/go`)

### 导入

```go
import "huffman"
```

### 函数

#### `Encode(input io.Reader, w io.Writer) error`

从 `input` 读取数据，将 Huffman 编码输出写入 `w`。

```go
in, _ := os.Open("input.bin")
defer in.Close()
out, _ := os.Create("output.huf")
defer out.Close()

if err := huffman.Encode(in, out); err != nil {
    log.Fatal(err)
}
```

#### `Decode(r io.Reader, w io.Writer) error`

从 `r` 读取编码数据，将解码输出写入 `w`。

#### `EncodeFile(inputPath, outputPath string) error`

文件编码便捷函数。

#### `DecodeFile(inputPath, outputPath string) error`

文件解码便捷函数。

### 常量

| 常量 | 值 | 描述 |
|------|-----|------|
| `SymbolLimit` | 257 | 256 字节 + 1 EOF 符号 |
| `MaxInputSize` | 4 GiB | 最大允许输入大小 |

---

## 算术编码 (`algorithms/arithmetic/go`)

### 导入

```go
import "arithmetic"
```

### 函数

与 Huffman 相同的 API：

- `Encode(input io.Reader, w io.Writer) error`
- `Decode(r io.Reader, w io.Writer) error`
- `EncodeFile(inputPath, outputPath string) error`
- `DecodeFile(inputPath, outputPath string) error`

### 额外函数

#### `ScaleFrequencies(freq []uint32)`

将频率归一化到 `MaxTotal`（16M）以内。

#### `BuildCumulative(freq []uint32) []uint32`

从原始频率构建累积频率表。

---

## 区间编码 (`algorithms/range/go`)

### 导入

```go
import "rangecoder"
```

### 函数

- `Encode(data []byte) ([]byte, error)` — 返回编码后的字节切片
- `Decode(data []byte) ([]byte, error)` — 返回解码后的字节切片

::: tip 提示
区间编码使用字节级 I/O（而非比特级），比算术编码更高效，同时保持几乎相同的压缩率。
:::

---

## RLE (`algorithms/rle/go`)

### 导入

```go
import "rle"
```

### 函数

- `Encode(input io.Reader, w io.Writer) error`
- `Decode(r io.Reader, w io.Writer) error`
- `EncodeFile(inputPath, outputPath string) error`
- `DecodeFile(inputPath, outputPath string) error`

### 错误处理

RLE 解码验证 `count > 0`，拒绝无效数据：

- `"invalid RLE data: count should not be 0"` — 数据损坏
- `"RLE data truncated"` — 不完整的 (计数, 值) 对
- `"output size limit exceeded"` — 检测到解压缩炸弹

---

## CLI 使用

每个 Go 模块都有 `cmd/main.go` 入口：

```bash
# 构建
cd algorithms/huffman/go && go build -o huffman_go ./cmd

# 使用
./huffman_go encode input.bin output.huf
./huffman_go decode output.huf decoded.bin
```

退出码：
- `0` — 成功
- `1` — 错误（消息输出到 stderr）
