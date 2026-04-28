# Rust 包 API 参考

另请参阅: [Streaming API](/zh/api/streaming)

每种 Rust 实现都提供可复用的库包，现在四种算法也统一提供 shared streaming adapter 与 `compresskit-codec` buffer helper。

## Huffman 编码

### 添加依赖

```toml
[dependencies]
huffman = { path = "../algorithms/huffman/rust" }
```

### API

```rust
use huffman;

fn main() -> std::io::Result<()> {
    huffman::huffman_encode_file("input.bin", "output.huf")?;
    huffman::huffman_decode_file("output.huf", "decoded.bin")?;
    Ok(())
}
```

---

## 算术编码

```rust
use arithmetic;

fn main() -> std::io::Result<()> {
    arithmetic::arithmetic_encode_file("input.bin", "output.aenc")?;
    arithmetic::arithmetic_decode_file("output.aenc", "decoded.bin")?;
    Ok(())
}
```

---

## 区间编码

### 添加依赖

```toml
[dependencies]
rangecoder = { path = "../algorithms/range/rust" }
```

### API

```rust
use rangecoder;

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let encoded = rangecoder::encode(&input_data)?;
    let decoded = rangecoder::decode(&encoded)?;
    Ok(())
}
```

### CLI 二进制文件

```bash
cargo build --bin rangecoder --release
./rangecoder encode input.bin output.rcnc
./rangecoder decode output.rcnc decoded.bin
```

---

## RLE

```rust
use rle;

fn main() -> std::io::Result<()> {
    rle::rle_encode_file("input.bin", "output.rle")?;
    rle::rle_decode_file("output.rle", "decoded.bin")?;
    Ok(())
}
```

---

## 通用模式

所有 Rust 实现都遵循以下约定：

| 模式 | 描述 |
|------|------|
| `*_encode_file` | 文件编码，返回 `io::Result<()>` |
| `*_decode_file` | 文件解码，返回 `io::Result<()>` |
| 错误传播 | 全程使用 `?` 运算符 |
| 缓冲 I/O | 所有文件操作使用 `BufReader`/`BufWriter` |
