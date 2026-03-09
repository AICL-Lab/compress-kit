import { defineConfig } from 'vitepress'

export default defineConfig({
  lang: 'zh-CN',
  title: 'Encoding',
  description: '编码算法集合 — 用 C++、Go、Rust 多语言实现经典压缩编码算法',

  // GitHub Pages 部署：base 需要与仓库名一致
  base: '/encoding/',

  markdown: {
    lineNumbers: true,
    languageAlias: {
      cuda: 'cpp',
    },
  },

  lastUpdated: true,

  themeConfig: {
    nav: [
      { text: '指南', link: '/guide/getting-started' },
      { text: '算法详解', link: '/guide/algorithms' },
      { text: '项目结构', link: '/guide/project-structure' },
    ],

    sidebar: [
      {
        text: '指南',
        items: [
          { text: '快速开始', link: '/guide/getting-started' },
          { text: '算法详解', link: '/guide/algorithms' },
          { text: '项目结构', link: '/guide/project-structure' },
        ],
      },
    ],

    editLink: {
      pattern: 'https://github.com/LessUp/encoding/edit/master/docs/:path',
      text: '在 GitHub 上编辑此页',
    },

    socialLinks: [
      { icon: 'github', link: 'https://github.com/LessUp/encoding' },
    ],

    footer: {
      message: '基于 MIT 许可发布',
      copyright: 'Copyright © 2025-2026 LessUp',
    },

    search: {
      provider: 'local',
    },

    outline: {
      level: [2, 3],
      label: '目录',
    },

    lastUpdated: {
      text: '最后更新',
    },

    docFooter: {
      prev: '上一页',
      next: '下一页',
    },

    returnToTopLabel: '返回顶部',
    sidebarMenuLabel: '菜单',
    darkModeSwitchLabel: '主题',
    externalLinkIcon: true,
  },
})
