# 架构设计

CompressKit 在算法核心之上仅保留一层薄 buffer 层：

```text
调用方
  -> 缓冲层 (encode_buffer / decode_buffer，接收 BufferTransform 函数指针)
  -> 算法核心 (Huffman / Arithmetic / Range / RLE)
```

## 核心设计

- 文件到文件 CLI 与内存调用方共享相同的缓冲层路径
- 单一 `BufferTransform` 函数指针别名将所有算法统一到同一签名
- 共享的体积上限检查与错误码保持接口精简

## 安全边界

- 缓冲层接受的输入上限为 `4 GiB`
- 缓冲层产生的解码输出上限为 `1 GiB`

## 验证

`make test` 先运行共享生命周期测试，再运行算法特定套件。

## 延伸阅读

- [系统架构设计](/architecture/) - 完整分层架构、二进制格式规范与 Deep Module 设计
- [C++ 实现参考](/api/cpp) - 共享 buffer 门面与状态码
- [快速开始](/guide/getting-started) - 构建与测试命令
