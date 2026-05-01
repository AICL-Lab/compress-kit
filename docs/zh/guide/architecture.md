# 架构设计

CompressKit 现有算法核心之上采用两层内存集成模型：

```text
调用方
  -> 缓冲层 (encode_buffer / decode_buffer)
  -> 流式层 (process / flush / finish / reset)
  -> 现有算法核心
```

## 设计目的

- 文件到文件 CLI 现在与内存调用方共享相同的缓冲层路径
- 生命周期语义在 C++17、Go 和 Rust 之间保持一致
- 统一的错误处理使后续帧格式和一致性工作更易构建

## 语言实现

- C++: `algorithms/shared/cpp`
- Go: `algorithms/shared/go`
- Rust: `algorithms/shared/rust`

## 安全边界

- 流式层接受的输入上限为 `4 GiB`
- 流式层产生的解码输出上限为 `1 GiB`

## 验证

`make test` 在运行算法特定测试套件之前，包含所有三种语言的共享流式生命周期测试。
