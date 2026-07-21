# 基准测试结果

本页面用于展示随文档一起提交的生成型基准快照。

::: tip 以本地结果为准
交互式图表直接读取 `docs/.vitepress/data/benchmarks.json`，该文件为仓库
维护的基准快照。当前仓库快照中的基准事实，应以这个生成 JSON 文件和
下方图表为准。
:::

## 交互式性能图表

<BenchmarkChart />

## 生成数据模型

| JSON 位置 | 字段 | 含义 |
|-----------|------|------|
| 顶层 | `generated` | 文档当前展示的基准快照日期 |
| 每条 `results[]` 记录 | `algorithm`、`dataset` | 图表中的基准坐标 |
| 每条 `results[]` 记录 | `encodeTime`、`decodeTime` | 挂钟时间，单位毫秒 |
| 每条 `results[]` 记录 | `encodeSpeed`、`decodeSpeed` | 吞吐量，单位 MiB/s |
| 每条 `results[]` 记录 | `compressionRatio`、`throughput` | 压缩比，以及 UI 使用的粗粒度吞吐标签 |

## 当前数据集覆盖范围

| 数据集键 | 图表标签 | 出现原因 |
|----------|----------|----------|
| `textlike_10MiB` | Text-like (10 MiB) | Huffman 与 Arithmetic 的主对比输入 |
| `repetitive_10MiB` | Repetitive (10 MiB) | RLE 有意义的高重复输入 |
| `small_dictionary_like` | Small dictionary-like sample | 为避开已知的大文件 Range 解码限制而使用的缩小样本 |

由于图表会直接读取生成 JSON，本页面不再维护手写静态数字表；本地重跑后，
快照会自动替换。

## 正确解读快照

- 得出结论前，请先确保比较的是**同一数据集**。
- Range 结果故意使用更小的输入，不能直接与 10 MiB 数据集做绝对吞吐对比。
- RLE 主要适用于高重复输入，对其他数据可能产生膨胀。

## 刷新事实来源

基准快照由维护者手工生成并提交。如需本地复现并刷新快照，请参考
[如何运行基准测试](/benchmarks/how-to-run) 中的手动流程，更新
`docs/.vitepress/data/benchmarks.json` 后执行：

```bash
cd docs && npm run build
```

- `npm run build` 会先同步 changelog，再构建文档站以验证刷新后的数据集
- 如果需要提交基准更新，请提交生成的 JSON 文件，而不是手工改写本页叙述数字

## 另见

- [如何运行基准测试](/benchmarks/how-to-run) — 详细说明
- [算法指南](/guide/algorithms) — 对比和选择指南
