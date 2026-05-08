# 如何运行基准测试

本项目包含基于 Python 的基准测试框架，测量所有算法和语言的编码/解码速度和压缩比。

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
2. 对每种算法 × 语言 × 数据集运行编码/解码
3. 测量时间和压缩比
4. 保存结果到 `reports/` 目录

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
| `tests/data/repetitive_1MiB.bin` | 重复 256 字节模式 | 1 MiB |
| `tests/data/repetitive_10MiB.bin` | 重复 256 字节模式 | 10 MiB |
| `tests/data/textli_1MiB.bin` | 加权英文字母 | 1 MiB |
| `tests/data/textli_10MiB.bin` | 加权英文字母 | 10 MiB |

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

结果保存在 `reports/` 目录：

```
reports/
├── huffman_cpp_report.txt
├── huffman_go_report.txt
├── huffman_rust_report.txt
├── arithmetic_cpp_report.txt
├── ...
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
4. 编辑相应的 `benchmark/bench.py` 以包含新文件

## 故障排除

### "Binary not found"

```bash
make build  # 重新构建所有实现
```

### "Test data not found"

```bash
make test-data  # 生成测试文件
```

### 区间编码基准测试很慢

::: warning
区间编码解码器对大于 500KB 的文件存在已知性能问题。请使用较小的测试文件。
:::

```bash
# 创建较小的测试文件
dd if=tests/data/random_10MiB.bin of=/tmp/small.bin bs=1024 count=100
```
