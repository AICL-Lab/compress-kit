# C++ 实现参考

另请参阅: [Streaming API](/zh/api/streaming)

所有 C++ 算法核心仍然保持单文件、零依赖风格，但现在额外共享了 `algorithms/shared/cpp/include/compresskit/` 下的 streaming / buffer 门面层。

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

## Huffman 编码

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

## 算术编码 / 区间编码 / RLE

遵循相同的单文件、零依赖模式。CLI 接口统一为：

```bash
<algorithm>_cpp encode <input> <output>
<algorithm>_cpp decode <input> <output>
```

### 通用模式

| 模式 | 描述 |
|------|------|
| 单文件 | 每个算法在一个 `main.cpp` 中 |
| 零依赖 | 仅使用标准库 |
| 错误处理 | `fprintf(stderr, ...)` + `exit(1)` |
| 内存管理 | `std::unique_ptr` + 自定义删除器 |

## 共享 Streaming 门面

共享头文件位于：

- `compresskit/result.hpp`
- `compresskit/encoder.hpp`
- `compresskit/buffer_api.hpp`
- `compresskit/algorithms.hpp`

这些头文件提供统一的 `Encoder` / `Decoder` 生命周期接口，以及 `compresskit::make_huffman_encoder()` 这类算法工厂。
