# C++ 实现参考

另请参阅: [Streaming API](/zh/api/streaming)

所有 C++ 算法核心仍然保持单文件风格，但现在依赖 `algorithms/shared/cpp/include/compresskit/` 下的共享 streaming/buffer 门面层。

## 编译

```bash
g++ -std=c++17 -O2 -Wall -Wextra -o <binary> main.cpp
```

### 推荐编译选项

| 选项 | 用途 |
|------|------|
| `-std=c++17` | 启用 C++17 特性 |
| `-O2` | 优化级别 |
| `-Wall -Wextra` | 警告 |
| `-fsanitize=address` | 地址消毒器（调试版本） |
| `-fsanitize=undefined` | 未定义行为消毒器（调试版本） |

## Huffman 编码 (`algorithms/huffman/cpp/main.cpp`)

```bash
./huffman_cpp encode input.bin output.huf
./huffman_cpp decode output.huf decoded.bin
```

### 内部结构

- `BitWriter` / `BitReader` — 比特级 I/O 类
- `Node` — Huffman 树节点（自定义删除器）
- `compress_file()` / `decompress_file()` — 主编/解码函数

### 文件格式

| 偏移 | 大小 | 字段 |
|------|------|------|
| 0 | 4B | 魔数: `HFMN` |
| 4 | 4B | 频率表大小（始终 257） |
| 8 | 1028B | 频率表（257 × uint32 LE） |
| 1036+ | 可变 | 编码比特流 |

---

## 算术编码 (`algorithms/arithmetic/cpp/main.cpp`)

```bash
./arithmetic_cpp encode input.bin output.aenc
./arithmetic_cpp decode output.aenc decoded.bin
```

### 关键类

- `ArithmeticEncoder` — 状态机，包含 `low`、`high`、`pendingBits`
- `ArithmeticDecoder` — 解码器，包含 `code` 初始化

---

## 区间编码 (`algorithms/range/cpp/main.cpp`)

```bash
./rangecoder_cpp encode input.bin output.rcnc
./rangecoder_cpp decode output.rcnc decoded.bin
```

### 文件格式

| 偏移 | 大小 | 字段 |
|------|------|------|
| 0 | 4B | 魔数: `RCNC` |
| 4 | 4B | 频率表大小 |
| 8 | 可变 | 频率表 |
| ... | 可变 | 字节流（重归一化区间） |

---

## RLE (`algorithms/rle/cpp/main.cpp`)

```bash
./rle_cpp encode input.bin output.rle
./rle_cpp decode output.rle decoded.bin
```

### 文件格式

重复的 `(count: uint32 LE, value: byte)` 对。

---

## 通用模式

| 模式 | 描述 |
|------|------|
| 单文件核心 | 每个算法核心在一个 `main.cpp` 中 |
| 共享依赖 | 使用 `algorithms/shared/cpp/` 中的公共代码 |
| 错误处理 | `fprintf(stderr, ...)` + `exit(1)` |
| 内存管理 | `std::unique_ptr` + 自定义删除器 |

## 共享 Streaming 门面

共享头文件位于：

- `compresskit/result.hpp`
- `compresskit/encoder.hpp`
- `compresskit/buffer_api.hpp`
- `compresskit/algorithms.hpp`

这些头文件提供统一的 `Encoder` / `Decoder` 生命周期接口，以及 `compresskit::make_huffman_encoder()` 这类算法工厂。
