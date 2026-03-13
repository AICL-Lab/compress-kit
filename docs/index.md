---
layout: home

hero:
  name: Encoding
  text: 编码算法集合
  tagline: 用 C++17、Go、Rust 三种语言实现经典压缩算法，用于学习原理、对比实现与跨语言验证
  actions:
    - theme: brand
      text: 快速开始
      link: /guide/getting-started
    - theme: alt
      text: 算法详解
      link: /guide/algorithms
    - theme: alt
      text: 项目结构
      link: /guide/project-structure

features:
  - title: 多语言对照
    details: 同一算法同时提供 C++17、Go、Rust 三套实现，便于比较代码风格、工程组织和性能取舍。
  - title: 统一文件格式
    details: 同一算法的不同语言实现共享二进制格式，可直接交叉编码、解码与验收。
  - title: 学习导向
    details: 文档优先解释算法适用场景、原理差异和上手路径，而不是只堆砌命令列表。
  - title: 可验证
    details: 保留构建、测试和基准流程，方便你在学习时同时验证正确性与性能。
---

## 项目定位

Encoding 是一个围绕经典压缩算法构建的学习型仓库：一方面保留可直接运行的实现，另一方面把算法背景、适用场景和目录结构整理成可以顺着读下去的文档入口。

## 适合谁

- 想通过多语言对照理解压缩编码算法的学习者
- 想比较 C++、Go、Rust 在同一算法上的实现差异的工程师
- 需要做跨语言格式兼容验证或构建基准测试的维护者

## 从哪里开始

1. 先看 [快速开始](/guide/getting-started)，确认依赖、构建方式和测试入口。
2. 再看 [算法详解](/guide/algorithms)，理解四种算法的适用场景与核心差异。
3. 需要定位代码和目录时，再查看 [项目结构](/guide/project-structure)。

## 推荐阅读路径

### 我只想尽快跑起来

- [快速开始](/guide/getting-started)
- [项目结构](/guide/project-structure)

### 我想先理解算法差异

- [算法详解](/guide/algorithms)
- [快速开始](/guide/getting-started)

### 我准备做实现对比或扩展

- [项目结构](/guide/project-structure)
- [算法详解](/guide/algorithms)
- [仓库 CHANGELOG](https://github.com/LessUp/encoding/blob/master/CHANGELOG.md)

## 核心文档

| 类别 | 页面 | 说明 |
|------|------|------|
| 概览 | [文档首页](/) | 项目定位、阅读路径与核心入口 |
| 快速开始 | [快速开始](/guide/getting-started) | 环境要求、构建、测试与基准命令 |
| 使用指南 | [算法详解](/guide/algorithms) | 四种算法的场景、原理与实现差异 |
| 参考 | [项目结构](/guide/project-structure) | 目录组织、CLI 约定与文件格式兼容性 |
| 归档 | [仓库 CHANGELOG](https://github.com/LessUp/encoding/blob/master/CHANGELOG.md) | 历史变更记录 |
