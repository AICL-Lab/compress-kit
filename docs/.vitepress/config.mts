import { defineConfig } from 'vitepress'

export default defineConfig({
  lang: 'zh-CN',
  title: 'Encoding',
  description: '编码算法集合：用 C++、Go、Rust 多语言实现经典压缩算法，用于学习、对比与验证',
  base: '/encoding/',
  cleanUrls: true,

  sitemap: {
    hostname: 'https://lessup.github.io',
  },

  head: [
    ['meta', { property: 'og:type', content: 'website' }],
    ['meta', { property: 'og:title', content: 'Encoding | 编码算法集合' }],
    ['meta', { property: 'og:description', content: '用 C++、Go、Rust 多语言实现经典压缩算法，用于学习、对比与跨语言验证' }],
    ['meta', { property: 'og:url', content: 'https://lessup.github.io/encoding/' }],
    ['meta', { name: 'twitter:card', content: 'summary' }],
    ['meta', { name: 'twitter:title', content: 'Encoding | 编码算法集合' }],
    ['meta', { name: 'twitter:description', content: '用 C++、Go、Rust 多语言实现经典压缩算法，用于学习、对比与跨语言验证' }],
    ['link', { rel: 'icon', type: 'image/svg+xml', href: '/encoding/logo.svg' }],
  ],

  markdown: {
    lineNumbers: true,
    languageAlias: {
      cuda: 'cpp',
    },
  },

  lastUpdated: true,

  themeConfig: {
    nav: [
      { text: '概览', link: '/' },
      { text: '快速开始', link: '/guide/getting-started' },
      { text: '使用指南', link: '/guide/algorithms' },
      { text: '参考', link: '/guide/project-structure' },
      { text: '归档', link: 'https://github.com/LessUp/encoding/blob/master/CHANGELOG.md' },
    ],

    sidebar: [
      {
        text: '概览',
        items: [{ text: '文档首页', link: '/' }],
      },
      {
        text: '快速开始',
        items: [{ text: '快速开始', link: '/guide/getting-started' }],
      },
      {
        text: '使用指南',
        items: [{ text: '算法详解', link: '/guide/algorithms' }],
      },
      {
        text: '参考',
        items: [{ text: '项目结构', link: '/guide/project-structure' }],
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
