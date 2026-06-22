# CompressKit 领域词汇表

本文档定义 CompressKit 项目的核心领域概念，为架构讨论提供统一语言。

## 核心实体

### Magic Number（魔数）

二进制格式的标识符，出现在编码文件的开头，用于识别编码算法类型。

| 算法 | Magic Number | 说明 |
|------|-------------|------|
| Huffman | `HFMN` | 4 字节 ASCII |
| Arithmetic | `AENC` | 4 字节 ASCII |
| Range Coder | `RCNC` | 4 字节 ASCII |
| RLE | `RLE\x00` | 4 字节，含空终止符 |

**设计决策：** Magic Number 不可更改，是二进制格式兼容性的基础。

### Frequency Table（频率表）

静态模型算法（Huffman、Arithmetic、Range）使用的符号频率数据结构。

- **格式：** 257 个 uint32 LE 值（256 字节值 + 1 EOF 标记）
- **总大小：** 4 字节（符号计数）+ 257 × 4 字节 = 1032 字节
- **编码顺序：** 符号 0-255（字节值），符号 256（EOF）

**约束：** 频率表格式不可更改。

### Symbol Limit（符号限制）

编码器处理的符号数量上限。

- **值：** 257（256 字节值 + 1 EOF 符号）
- **EOF 符号索引：** 256

### Streaming Layer（流式层）

提供增量处理能力的编码/解码接口，通过生命周期状态机管理。

**生命周期：**
```
READY → STREAMING → FLUSHING → FINISHED
         ↓
       ERROR（任意状态可转入）
```

**核心接口：**
- `process(in, out)` - 增量处理输入
- `flush(out)` - 刷新缓冲输出
- `finish(out)` - 完成处理并写入结束标记
- `reset()` - 重置到 READY 状态
- `state()` - 查询当前状态

**语义保证：**
- `process()` 在编码器中缓冲输入，在解码器中缓冲输入
- `flush()` 对 Huffman/Arithmetic/Range 是无操作（需要完整输入）
- `finish()` 触发最终编码/解码并写入输出

### Buffer Layer（缓冲层）

提供一次性处理完整数据的便捷 API。

**接口：**
- `encode_buffer(encoder, input)` - 编码完整输入
- `decode_buffer(decoder, input)` - 解码完整输入

**语义：** 等价于 `encoder.process(input) → encoder.finish()`

### State Machine（状态机）

Streaming Layer 的核心控制逻辑，定义有效的状态转换。

**状态：**
- `State::READY` - 初始状态，准备接收输入
- `State::STREAMING` - 正在处理输入
- `State::FLUSHING` - 已刷新，等待完成
- `State::FINISHED` - 处理完成，不可再接收输入
- `State::ERROR` - 错误状态，需要 reset

**转换规则：**
| 操作 | 有效前置状态 | 结果状态 |
|------|-------------|---------|
| process | READY, STREAMING, FLUSHING | STREAMING |
| flush | READY, STREAMING, FLUSHING | FLUSHING |
| finish | READY, STREAMING, FLUSHING | FINISHED |
| reset | 任意 | READY |

## 安全边界

### Input Size Limit（输入大小限制）

- **上限：** 4 GiB
- **目的：** 防止频率溢出和解压缩炸弹攻击

### Output Size Limit（输出大小限制）

- **上限：** 1 GiB（仅解码）
- **目的：** 防止解压缩炸弹攻击

## 错误类型

### 标准错误码

| 错误 | 语义 | 恢复方式 |
|------|------|---------|
| `BUF_TOO_SMALL` | 输出缓冲区不足 | 使用更大缓冲区重试 |
| `ERR_CORRUPT` | 数据损坏或校验失败 | 输入数据无效 |
| `ERR_INVALID_STATE` | 当前状态不支持此操作 | reset 后重试 |
| `ERR_SIZE_LIMIT` | 超过安全限制 | 输入/输出过大 |

### 事务性保证

当 `BUF_TOO_SMALL` 返回时，内部状态保持不变，调用者可以使用更大的缓冲区重试。

## 算法特性

### Huffman 编码

- **类型：** 基于前缀码的熵编码
- **模型：** 静态（需要完整输入构建频率表）
- **输出：** 位流（需要位对齐处理）
- **限制：** 需要完整输入才能编码

### Arithmetic 编码

- **类型：** 区间编码
- **模型：** 静态
- **精度：** 区间缩放到 2^24
- **限制：** 需要完整输入才能编码

### Range Coder

- **类型：** 区间编码
- **模型：** 静态
- **精度：** 区间缩放到 2^24

### RLE（Run-Length Encoding）

- **类型：** 游程编码
- **格式：** `(count: uint32 LE, value: byte)` 对序列
- **特点：** 不需要频率表，可增量处理

## 架构分层

```
┌─────────────────────────────────────┐
│           CLI Entry Point           │  命令行接口
├─────────────────────────────────────┤
│           Buffer Layer              │  便捷 API（内存 transform）
├─────────────────────────────────────┤
│          Streaming Layer            │  生命周期管理
├─────────────────────────────────────┤
│          Algorithm Core             │  编码/解码逻辑
├─────────────────────────────────────┤
│          Shared Utilities           │  共享工具（C++）
└─────────────────────────────────────┘
```

**调用流向：** CLI → Buffer Layer → Streaming Layer → Algorithm Core

### Shared Utilities（共享工具层）

C++ 实现的共享工具模块，位于 `algorithms/shared/cpp/`：

| 模块 | 功能 |
|------|------|
| `encoder.hpp` | Encoder/Decoder 接口、状态机定义 |
| `buffer_api.hpp/cpp` | BufferEncoder/BufferDecoder、encode_buffer/decode_buffer |
| `frequency_table.hpp/cpp` | 频率表读写、count_frequencies、scale_frequencies、build_cumulative |
| `cli_launcher.hpp/cpp` | CLI 入口统一处理 |
| `constants.hpp` | Magic Number、SYMBOL_LIMIT、EOF_SYMBOL、安全限制常量 |
| `result.hpp` | Result 模板、StatusCode 枚举 |

## 文档层次

| 文档 | 用途 | 受众 |
|------|------|------|
| README.md | 项目入口 | 新用户 |
| VitePress 文档 | 产品门户 | 最终用户 |
| CHANGELOG.md | 用户可见变更 | 用户 |
| CONTEXT.md | 领域词汇 | 贡献者/AI |

## 参考资料

- `docs/adr/0001-validation-metadata-module-shape.md` - 验证元数据模块形状决策
- `docs/adr/0003-range-coder-corpus-cap-policy.md` - Range Coder 语料库上限策略决策
- `docs/adr/0004-rle-buffered-streaming-stance.md` - RLE 缓冲流式语义决策
- `docs/adr/0005-range-cpp-bench-mode-migration-path.md` - Range C++ bench 模式迁移路径决策

## 文档术语表

本文节定义 Git Pages 文档（VitePress）中的术语翻译规范。

### 算法名称

| 英文 | 中文文档用法 | 说明 |
|------|-------------|------|
| Huffman | Huffman | 保留英文 |
| Arithmetic | Arithmetic | 保留英文 |
| Range Coder | Range Coder | 保留英文 |
| RLE | RLE | 保留英文，首次出现可标注"行程编码" |

### API 术语

| 英文 | 中文翻译 | 说明 |
|------|----------|------|
| Streaming API | Streaming API | 保留英文，代码术语不翻译 |
| Buffer Layer | 缓冲层 | 翻译 |
| Streaming Layer | 流式层 | 翻译 |
| CLI | CLI | 保留英文，通用术语 |
| Import | 导入 | 翻译 |
| Functions | 函数 | 翻译 |
| Constants | 常量 | 翻译 |
| Error Handling | 错误处理 | 翻译 |
| Error Type | 错误类型 | 翻译 |
| File Format | 文件格式 | 翻译 |
| Usage | 用法 | 翻译 |
| Internal Structure | 内部结构 | 翻译 |
| Key Classes | 关键类 | 翻译 |
| Common Patterns | 通用模式 | 翻译 |

### 文档结构

| 英文 | 中文翻译 | 说明 |
|------|----------|------|
| Getting Started | 快速开始 | 统一使用"快速开始" |
| API Reference | API 参考 | 保留 API 英文 |
| Benchmarks | 基准测试 | 翻译 |
| Prerequisites | 前置条件 | 统一使用"前置条件" |
| Troubleshooting | 故障排除 | 翻译 |

### 性能术语

| 英文 | 中文翻译 | 说明 |
|------|----------|------|
| encode/decode | 编码/解码 | 翻译 |
| Encoder/Decoder | 编码器/解码器 | 翻译 |
| encode time | 编码时间 | 翻译 |
| decode time | 解码时间 | 翻译 |
| compression ratio | 压缩比 | 翻译 |
| throughput | 吞吐量 | 翻译 |
| MiB/s | MiB/s | 保留单位 |
