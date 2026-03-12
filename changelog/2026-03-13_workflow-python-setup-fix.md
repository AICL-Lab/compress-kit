# Workflow Python Setup 修正

日期：2026-03-13

## 变更内容

- 移除 `Test Correctness` job 中 `actions/setup-python` 的 `cache: pip`
- 保留 Python 3.11 运行时与原有跨语言 correctness 验证流程
- 避免在没有 Python 依赖清单的仓库里触发缓存配置错误，恢复正确性测试 job 的可执行性

## 背景

该仓库的 correctness job 只需要 Python 解释器来生成测试数据，并不依赖 `requirements.txt` 或其他 pip 依赖锁文件。此前启用 pip cache 会让 `setup-python` 在 Hosted Runner 上直接失败，导致整条 correctness 验证提前中断。
