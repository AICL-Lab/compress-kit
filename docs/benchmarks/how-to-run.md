# 如何运行基准测试

本仓库的基准数据通过仓库内的 Python 脚本与 CLI 二进制生成，文档站读取
`docs/.vitepress/data/benchmarks.json` 这一快照展示交互式图表。

## 前置条件

- Python 3.8+
- 所有实现已构建：`make build`
- 测试数据已生成：`make test-data`

## 仓库现状

当前仓库**未内置一键基准脚本**（无 `make bench` 目标，亦无
`scripts/run_all_bench.py` 与 `algorithms/<algo>/benchmark/bench.py`）。
基准快照 `docs/.vitepress/data/benchmarks.json` 由维护者手工生成并提交。

## 手动复现基准

如需本地复现，可参照以下流程：

1. 构建所有算法二进制：

   ```bash
   make build
   ```

2. 生成测试语料：

   ```bash
   make test-data
   ```

3. 对每个算法 × 数据集组合手工运行编码与解码，记录挂钟时间、输出大小：

   ```bash
   ./build/huffman_cpp encode tests/data/textlike_10MiB.bin out.huf
   ./build/huffman_cpp decode out.huf restored.bin
   # 同理运行 arithmetic_cpp / rangecoder_cpp / rle_cpp
   ```

4. 汇总结果写入 `docs/.vitepress/data/benchmarks.json`（字段含义见
   [基准测试结果](/benchmarks/results)），随后运行 `cd docs && npm run build`
   验证文档图表可正确渲染。

## 测试数据

`make test-data` 调用 `tests/gen_testdata.py`，在 `tests/data/` 下生成：

| 文件 | 生成方式 | 大小 |
|------|----------|------|
| `random_1MiB.bin` | `os.urandom(1024*1024)` | 1 MiB |
| `random_10MiB.bin` | `os.urandom(10*1024*1024)` | 10 MiB |
| `repetitive_10MiB.bin` | 重复 256 字节模式 | 10 MiB |
| `textlike_10MiB.bin` | 加权英文字母 | 10 MiB |
| `small_dictionary_like.bin` | 小型重复词典风格样本 | 约 8 KiB |

## 测量指标

| 指标 | 描述 |
|------|------|
| 编码时间 | 压缩输入的挂钟时间 |
| 解码时间 | 恢复原始的挂钟时间 |
| 编码速度 | MiB/s = 输入大小 / 编码时间 |
| 解码速度 | MiB/s = 输入大小 / 解码时间 |
| 压缩比 | 输出大小 / 输入大小（越小越好） |

## Range Coder 注意事项

Range Coder 解码器对大于 **500 KiB** 的输入存在已知性能问题。基准快照
中 Range 结果使用 `small_dictionary_like` 数据集，手工复现时请同样保持
在小样本范围内（建议 < 100 KiB）。

## 故障排除

### "Binary not found"

```bash
make build  # 重新构建所有实现
```

### "Test data not found"

```bash
make test-data  # 生成测试文件
```


### Range Coder 基准测试很慢

::: warning
Range Coder 解码器对大于 500 KiB 的文件存在已知性能问题。因此仓库内置的
基准流程会为 Range 结果使用 `small_dictionary_like.bin`。
:::

```bash
# 创建较小的测试文件
dd if=tests/data/random_10MiB.bin of=small.bin bs=1024 count=100
```
