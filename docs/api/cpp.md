# C++ 实现参考

所有 C++ 算法核心保持单文件风格，并共享 `algorithms/shared/cpp/include/compresskit/` 下的薄 buffer 门面层。无状态化 streaming 层；`encode_buffer` / `decode_buffer` 直接接收 `BufferTransform` 函数指针。

## 编译

项目使用 CMake 构建。`make build` 是对 CMake 调用的薄封装：

```bash
cmake -S . -B build && cmake --build build
```

构建产物位于 `build/` 目录下（如 `build/huffman_cpp`）。

### 推荐编译选项

| 选项 | 用途 |
|------|------|
| `-std=c++17` | 启用 C++17 特性 |
| `-O2` | 优化级别 |
| `-Wall -Wextra` | 警告 |
| `-fsanitize=address` | 地址消毒器（调试版本） |
| `-fsanitize=undefined` | 未定义行为消毒器（调试版本） |

## Huffman (`algorithms/huffman/cpp/main.cpp`)

### 用法

```bash
./build/huffman_cpp encode input.bin output.huf
./build/huffman_cpp decode output.huf decoded.bin
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

## Arithmetic (`algorithms/arithmetic/cpp/main.cpp`)

### 用法

```bash
./build/arithmetic_cpp encode input.bin output.aenc
./build/arithmetic_cpp decode output.aenc decoded.bin
```

### 关键类

- `ArithmeticEncoder` — 状态机，包含 `low`、`high`、`pendingBits`
- `ArithmeticDecoder` — 解码器，包含 `code` 初始化

---

## Range Coder (`algorithms/range/cpp/main.cpp`)

### 用法

```bash
./build/rangecoder_cpp encode input.bin output.rcnc
./build/rangecoder_cpp decode output.rcnc decoded.bin
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
./build/rle_cpp encode input.bin output.rle
./build/rle_cpp decode output.rle decoded.bin
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

## 共享 Buffer 门面

共享头文件位于：

- `compresskit/result.hpp` — `StatusCode` 枚举与 `Result<T>` 模板
- `compresskit/buffer_api.hpp` — `BufferTransform`、`encode_buffer`、`decode_buffer`、文件辅助函数
- `compresskit/algorithms.hpp` — 各算法的 `*_encode_buffer` / `*_decode_buffer` 入口

### 状态码

| 码 | 含义 |
|------|---------|
| `OK` | 成功 |
| `ERR_CORRUPT` | 编码数据结构校验失败 |
| `ERR_SIZE_LIMIT` | 输入超过 4 GiB 或解码输出超过 1 GiB |

### Buffer Transform

`BufferTransform` 是函数指针别名，将内存字节缓冲映射为其转换形式。每个算法暴露一个编码 transform 和一个解码 transform。

```cpp
#include <vector>
#include "compresskit/algorithms.hpp"
#include "compresskit/buffer_api.hpp"

std::vector<uint8_t> encode(const std::vector<uint8_t>& input) {
    auto result = compresskit::encode_buffer(huffman_encode_buffer, input);
    return result.ok() ? std::move(result.value) : std::vector<uint8_t>{};
}
```
