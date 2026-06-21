# 如何运行基准测试

本项目包含基于 Python 的基准测试框架，测量所有算法的编码/解码速度和压缩比。

## 前置条件

- Python 3.8+
- 所有实现已构建：`make build`
- 测试数据已生成：`make test-data`

## 运行基准测试

### 全部基准测试

```bash
make bench
```

这将运行 `scripts/run_all_bench.py`，会：
1. 生成测试数据（如果 `tests/data/` 为空）
2. 运行已验证的算法 × 数据集矩阵
3. 测量时间和压缩比
4. 把报告保存到 `reports/`，并刷新 `docs/.vitepress/data/benchmarks.json`

### 单个算法基准测试

```bash
cd algorithms/huffman/benchmark
python3 bench.py

cd algorithms/arithmetic/benchmark
python3 bench.py

cd algorithms/range/benchmark
python3 bench.py

cd algorithms/rle/benchmark
python3 bench.py
```

## 基准测试配置

### 测试数据

| 文件 | 生成方式 | 大小 |
|------|----------|------|
| `tests/data/random_1MiB.bin` | `os.urandom(1024*1024)` | 1 MiB |
| `tests/data/random_10MiB.bin` | `os.urandom(10*1024*1024)` | 10 MiB |
| `tests/data/repetitive_10MiB.bin` | 重复 256 字节模式 | 10 MiB |
| `tests/data/textlike_10MiB.bin` | 加权英文字母 | 10 MiB |
| `tests/data/small_dictionary_like.bin` | 小型重复词典风格样本 | 约 8 KiB |

重新生成：

```bash
make test-data
# 或
python3 tests/gen_testdata.py
```

### 测量指标

| 指标 | 描述 |
|------|------|
| 编码时间 | 压缩输入的挂钟时间 |
| 解码时间 | 恢复原始的挂钟时间 |
| 编码速度 | MiB/s = 输入大小 / 编码时间 |
| 解码速度 | MiB/s = 输入大小 / 解码时间 |
| 压缩比 | 输出大小 / 输入大小（越小越好） |

### 输出格式

结果保存在 `reports/`，文档快照写入
`docs/.vitepress/data/benchmarks.json`：

```text
reports/
├── huffman_report_<timestamp>.txt
├── arithmetic_report_<timestamp>.txt
├── range_report_<timestamp>.txt
└── rle_report_<timestamp>.txt

docs/.vitepress/data/
└── benchmarks.json
```

每个报告包含：

```
Algorithm: Huffman
Language: C++
Input: 10 MiB random data
Encode: 245 ms (40.8 MiB/s)
Decode: 198 ms (50.5 MiB/s)
Compression ratio: 1.23
```

## 添加新基准测试

添加新测试数据集：

1. 编辑 `tests/gen_testdata.py`
2. 在 `generate_random_file()` 中添加生成代码或创建新生成器
3. 运行 `make test-data`
4. 同步更新 `scripts/run_all_bench.py`，确保生成的文档快照保持一致

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
