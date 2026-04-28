# Streaming API

CompressKit 现在在 C++17、Go 和 Rust 三种语言里都提供一致的两层内存 API：

- 有状态的 streaming 层：`process`、`flush`、`finish`、`reset`
- 无状态的 buffer 层：封装完整生命周期，适合一次性处理字节切片

## 生命周期

streaming 层遵循统一状态机：

| 状态 | 允许调用 | 说明 |
|------|----------|------|
| `READY` | `process`、`flush`、`finish`、`reset` | 此时 `flush` 为 no-op |
| `STREAMING` | `process`、`flush`、`finish`、`reset` | 算法可以内部缓冲输入 |
| `FLUSHING` | `process`、`flush`、`finish`、`reset` | 再次 `process` 会回到 `STREAMING` |
| `FINISHED` | `reset` | 其他调用都返回 `ERR_INVALID_STATE` |
| `ERROR` | `reset` | 错误态在 `reset` 前不可恢复 |

`finish()` 会自动执行最终 flush。

## 错误码

| 代码 | 含义 |
|------|------|
| `OK` | 成功 |
| `BUF_TOO_SMALL` | 调用方输出缓冲区过小，内部状态保持不变 |
| `ERR_TRUNCATED` | 输入在完整解码前提前结束 |
| `ERR_CORRUPT` | 编码数据结构校验失败 |
| `ERR_INVALID_STATE` | 调用顺序违反生命周期约定 |
| `ERR_SIZE_LIMIT` | 输入超过 4 GiB 或解码输出超过 1 GiB |
| `ERR_VERSION_UNSUPPORTED` | 预留给后续 frame 层版本校验 |
| `ERR_UNKNOWN_ALGO` | 预留给后续 frame 层算法校验 |

## Buffer 层

buffer 层等价于：

```text
new encoder -> process(input) -> finish()
new decoder -> process(input) -> finish()
```

这样可以让 file-to-file 路径与内存调用路径共用同一套实现。

## Go 示例

```go
import (
    "github.com/LessUp/compress-kit/algorithms/shared/go/codec"
    "huffman"
)

func encode(data []byte) ([]byte, error) {
    return codec.EncodeBuffer(huffman.NewStreamingEncoder(), data)
}

func decode(encoded []byte) ([]byte, error) {
    return codec.DecodeBuffer(huffman.NewStreamingDecoder(), encoded)
}
```

## Rust 示例

```rust
use compresskit_codec::codec::{decode_buffer, encode_buffer};
use huffman::{StreamingDecoder, StreamingEncoder};

fn roundtrip(input: &[u8]) -> Result<Vec<u8>, compresskit_codec::codec::CodecError> {
    let mut encoder = StreamingEncoder::new();
    let encoded = encode_buffer(&mut encoder, input)?;

    let mut decoder = StreamingDecoder::new();
    decode_buffer(&mut decoder, &encoded)
}
```

## C++ 示例

```cpp
#include <vector>

#include "compresskit/algorithms.hpp"

std::vector<uint8_t> encode(const std::vector<uint8_t>& input) {
    auto encoder = compresskit::make_huffman_encoder();
    auto result = compresskit::encode_buffer(encoder, input);
    return result.value;
}
```

## 验证方式

`make test` 现在会先运行 C++、Go、Rust 三种语言的共享 streaming 层测试，再运行各算法测试套件。
