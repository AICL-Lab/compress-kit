# 项目结构

CompressKit 按“算法优先”的方式组织。每个算法保持单文件 C++17 核心，
并共享 streaming/buffer 门面层，便于直接对比和维护，无需语言分叉。

## 源码布局

```text
algorithms/
├── shared/        # streaming 与 buffer API 基础层
├── huffman/       # cpp/
├── arithmetic/    # cpp/
├── range/         # cpp/
└── rle/           # cpp/

tests/
├── gen_testdata.py
└── streaming_api_contract/

docs/              # VitePress 站点：根门户 + en/ + zh/
```

## 职责边界

| 区域 | 负责 | 不负责 |
|------|------|--------|
| `algorithms/<algo>/cpp/` | 算法实现、CLI 入口、算法测试 | 全局文档编排 |
| `algorithms/shared/` | Streaming 生命周期、buffer API、共享契约测试 | 算法私有文件格式 |
| `tests/data/` | `make test-data` 生成的本地语料 | 需要提交的源文件 fixture |
| `docs/` | 用户指南、API 说明、已知限制 | 营销文案 |

## 二进制格式

当前终局基线保持各算法既有格式稳定：

| 算法 | 魔数/头部 | 扩展名 | 载荷 |
|------|-----------|--------|------|
| Huffman | `HFMN` + 频率表 | `.huf` | 比特流 |
| Arithmetic | `AENC` + 频率表 | `.aenc` | 比特流 |
| Range Coder | `RCNC` + 频率表 | `.rcnc` | 字节流 |
| RLE | `RLE\x00` | `.rle` | 魔数 + `(count: uint32 LE, value: byte)` 对 |

## 生成产物

以下内容是生成产物，默认被忽略：

- `build/` 下的算法二进制（如 `build/huffman_cpp`）
- `tests/data/*.bin`
- benchmark 报告和临时目录
- `docs/.vitepress/dist/`

打包或审查仓库形态前可运行 `make clean`。

## 相关页面

- [快速开始](/zh/guide/getting-started)
- [Streaming API](/zh/api/streaming)
