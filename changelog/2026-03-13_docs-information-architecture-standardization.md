# 文档信息架构规范化（2026-03-13）

## 变更背景

- 继续推进仓库群 GitHub Pages 与文档入口标准化。
- 此前 `README.md` / `README.zh-CN.md` 既承担仓库入口，又内嵌大量算法说明与双语长篇内容；`docs/index.md` 也偏展示页风格。
- 本次目标是保留现有 VitePress 技术栈，只统一 README 与文档首页的职责边界、导航命名和入口层级。

## 导航与目录调整

- `docs/.vitepress/config.mts` 一级导航调整为：`概览`、`快速开始`、`使用指南`、`参考`、`归档`。
- 侧边栏同步收敛为同一层级语义，避免原先 `指南 / 算法详解 / 项目结构 / 相关链接` 混用。
- 保留 `guide/` 目录结构不变，仅通过首页与导航重组阅读顺序。

## 首页调整

- `README.md` / `README.zh-CN.md` 收敛为仓库入口，只保留项目定位、最短命令示例与文档站链接。
- `docs/index.md` 改为文档首页，新增项目定位、适合谁、从哪里开始、推荐阅读路径和核心文档表。
- 站点首页不再直接罗列全部算法卡片细节，而是把算法内容交给 `guide/algorithms.md` 负责。

## Pages / Workflow 调整

- 本次未修改 `pages.yml` 构建流程；继续使用根目录 `package.json` 的 `npm run docs:build` 输出 `docs/.vitepress/dist`。
- 修正了站点配置中的 icon 路径，保持与 `base: /encoding/` 一致。

## 验证结果

- 人工检查 README 与文档首页职责已分离，仓库入口与文档入口不再重复堆叠。
- 人工检查导航和首页链接均指向现有文档页面，快速开始 / 算法详解 / 项目结构两次点击内可达。
- 已在仓库根目录执行 `npm run docs:build`，VitePress 构建成功，产物输出到 `docs/.vitepress/dist/`。

## 后续待办

- 视需要补充一个面向维护者的 `archive/` 或 changelog 索引入口，减少外链到 GitHub 文件的依赖。
- 后续可继续统一 `guide/` 页面中的标题风格与术语中英混排规则。
