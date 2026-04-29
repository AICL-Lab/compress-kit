# 项目结构

CompressKit 按“算法优先、语言其次”的方式组织。这样可以直接对比同一算法
在 C++17、Go、Rust 中的实现，而不会掩盖各语言自己的工程习惯。

## 源码布局

```text
algorithms/
├── shared/        # streaming 与 buffer API 基础层
├── huffman/       # cpp/, go/, rust/
├── arithmetic/    # cpp/, go/, rust/
├── range/         # cpp/, go/, rust/
└── rle/           # cpp/, go/, rust/

tests/
├── gen_testdata.py
├── streaming_api_contract/
└── conformance/

docs/              # VitePress 站点：根门户 + en/ + zh/
openspec/          # 稳定规范与已归档设计变更
```

## 职责边界

| 区域 | 负责 | 不负责 |
|------|------|--------|
| `algorithms/<algo>/<lang>/` | 算法实现、CLI 入口、语言内测试 | 全局文档或跨语言编排 |
| `algorithms/shared/` | Streaming 生命周期、buffer API、共享契约测试 | 算法私有文件格式 |
| `tests/conformance/` | 稳定格式的跨语言编码/解码矩阵 | 未来 shared-frame 校验 |
| `tests/data/` | `make test-data` 生成的本地语料 | 需要提交的源文件 fixture |
| `docs/` | 用户指南、API 说明、已知限制 | OpenSpec 变更追踪 |
| `openspec/` | 规范要求与提案归档历史 | 营销文案 |

## 二进制格式

当前终局基线保持各算法既有格式稳定：

| 算法 | 魔数/头部 | 扩展名 | 载荷 |
|------|-----------|--------|------|
| Huffman | `HFMN` + 频率表 | `.huf` | 比特流 |
| 算术编码 | `AENC` + 频率表 | `.aenc` | 比特流 |
| 区间编码 | `RCNC` + 频率表 | `.rcnc` | 字节流 |
| RLE | 无魔数头 | `.rle` | `(count: uint32 LE, value: byte)` 对 |

未来 shared-frame 提案已归档到 `openspec/changes/archive/`，不属于当前活跃文件格式契约。

## 生成产物

以下内容是生成产物，默认被忽略：

- `huffman_cpp`、`huffman_go`、`huffman_rust` 等算法二进制
- Rust `target/` 目录
- `tests/data/*.bin`
- benchmark 报告和 conformance 临时目录
- `docs/.vitepress/dist/`

打包或审查仓库形态前可运行 `make clean`。

## 相关页面

- [快速开始](/zh/guide/getting-started)
- [Streaming API](/zh/api/streaming)
- [跨语言测试](/zh/testing/cross-language)
