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

**设计决策：** Magic Number 不可更改，是跨语言二进制兼容性的基础。

### Frequency Table（频率表）

静态模型算法（Huffman、Arithmetic、Range）使用的符号频率数据结构。

- **格式：** 257 个 uint32 LE 值（256 字节值 + 1 EOF 标记）
- **总大小：** 4 字节（符号计数）+ 257 × 4 字节 = 1032 字节
- **编码顺序：** 符号 0-255（字节值），符号 256（EOF）

**约束：** 频率表格式跨语言统一，不可更改。

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
- `Process(in []byte, out []byte)` - 增量处理输入
- `Flush(out []byte)` - 刷新缓冲输出
- `Finish(out []byte)` - 完成处理并写入结束标记
- `Reset()` - 重置到 READY 状态
- `State()` - 查询当前状态

**语义保证：**
- `Process()` 在编码器中缓冲输入，在解码器中增量输出
- `Flush()` 对 Huffman/Arithmetic/Range 是无操作（需要完整输入）
- `Finish()` 触发最终编码/解码并写入输出

### Buffer Layer（缓冲层）

提供一次性处理完整数据的便捷 API。

**接口：**
- `EncodeBuffer(input []byte) ([]byte, error)`
- `DecodeBuffer(input []byte) ([]byte, error)`

**语义：** 等价于 `new encoder → process(input) → finish()`

### State Machine（状态机）

Streaming Layer 的核心控制逻辑，定义有效的状态转换。

**状态：**
- `StateReady` - 初始状态，准备接收输入
- `StateStreaming` - 正在处理输入
- `StateFlushing` - 已刷新，等待完成
- `StateFinished` - 处理完成，不可再接收输入
- `StateError` - 错误状态，需要 Reset

**转换规则：**
| 操作 | 有效前置状态 | 结果状态 |
|------|-------------|---------|
| Process | READY, STREAMING, FLUSHING | STREAMING |
| Flush | READY, STREAMING, FLUSHING | FLUSHING |
| Finish | READY, STREAMING, FLUSHING | FINISHED |
| Reset | 任意 | READY |

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
| `ErrBufTooSmall` | 输出缓冲区不足 | 使用更大缓冲区重试 |
| `ErrTruncated` | 输入流过早结束 | 检查输入完整性 |
| `ErrCorrupt` | 数据损坏或校验失败 | 输入数据无效 |
| `ErrInvalidState` | 当前状态不支持此操作 | Reset 后重试 |
| `ErrSizeLimit` | 超过安全限制 | 输入/输出过大 |
| `ErrVersionUnsupported` | 不支持的版本 | 检查版本兼容性 |
| `ErrUnknownAlgo` | 未知的算法标识 | 检查 Magic Number |

### 事务性保证

当 `ErrBufTooSmall` 返回时，内部状态保持不变，调用者可以使用更大的缓冲区重试。

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
- **已知问题：** 大文件（>500KB）性能下降，推荐限制在 100KB 以内

### RLE（Run-Length Encoding）

- **类型：** 游程编码
- **格式：** `(count: uint32 LE, value: byte)` 对序列
- **特点：** 不需要频率表，可增量处理

## 架构分层

```
┌─────────────────────────────────────┐
│           CLI Entry Point           │  命令行接口
├─────────────────────────────────────┤
│           Buffer Layer              │  便捷 API
├─────────────────────────────────────┤
│          Streaming Layer            │  生命周期管理
├─────────────────────────────────────┤
│          Algorithm Core             │  编码/解码逻辑
└─────────────────────────────────────┘
```

**调用流向：** CLI → Buffer Layer → Streaming Layer → Algorithm Core

## 跨语言一致性

### 必须一致的部分

1. **Magic Number** - 二进制格式标识
2. **Frequency Table 格式** - 字节序、顺序、大小
3. **状态机语义** - 状态转换规则
4. **错误码语义** - 错误类型和恢复方式
5. **安全限制** - 输入/输出大小上限

### 可不一致的部分

1. **内部实现** - 数据结构、算法优化
2. **API 命名风格** - 遵循各语言惯例
3. **错误消息文本** - 不影响错误码语义

## 文档层次

| 文档 | 用途 | 受众 |
|------|------|------|
| README.md | 项目入口 | 新用户 |
| VitePress 文档 | 产品门户 | 最终用户 |
| OpenSpec | 需求来源 | 开发者 |
| CHANGELOG.md | 用户可见变更 | 用户 |
| CONTEXT.md | 领域词汇 | 贡献者/AI |

## 参考资料

- `openspec/specs/core-architecture/spec.md` - 核心架构规范
- `openspec/specs/encoding-project/spec.md` - 项目需求规范
- `openspec/specs/cross-language-testing/spec.md` - 测试规范
