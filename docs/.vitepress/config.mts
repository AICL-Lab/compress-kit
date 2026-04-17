import { defineConfig } from 'vitepress'

// Shared sidebar configuration
const sharedSidebar = {
  '/en/': [
    {
      text: 'Overview',
      items: [{ text: 'Home', link: '/en/' }],
    },
    {
      text: 'Guide',
      items: [
        { text: 'Getting Started', link: '/en/guide/getting-started' },
        { text: 'Algorithms', link: '/en/guide/algorithms' },
      ],
    },
    {
      text: 'Reference',
      items: [
        { text: 'Project Structure', link: '/en/guide/project-structure' },
        { text: 'Specs (SSOT)', link: 'https://github.com/LessUp/encoding/tree/master/specs' },
      ],
    },
  ],
  '/zh/': [
    {
      text: '概览',
      items: [{ text: '首页', link: '/zh/' }],
    },
    {
      text: '指南',
      items: [
        { text: '快速开始', link: '/zh/guide/getting-started' },
        { text: '算法详解', link: '/zh/guide/algorithms' },
      ],
    },
    {
      text: '参考',
      items: [
        { text: '项目结构', link: '/zh/guide/project-structure' },
        { text: '规范文档 (SSOT)', link: 'https://github.com/LessUp/encoding/tree/master/specs' },
      ],
    },
  ],
}

// Shared nav configuration
const sharedNav = (lang: string) => [
  { text: lang === 'zh' ? '概览' : 'Overview', link: lang === 'zh' ? '/zh/' : '/en/' },
  { text: lang === 'zh' ? '快速开始' : 'Get Started', link: lang === 'zh' ? '/zh/guide/getting-started' : '/en/guide/getting-started' },
  { text: lang === 'zh' ? '算法' : 'Algorithms', link: lang === 'zh' ? '/zh/guide/algorithms' : '/en/guide/algorithms' },
  { text: lang === 'zh' ? '参考' : 'Reference', link: lang === 'zh' ? '/zh/guide/project-structure' : '/en/guide/project-structure' },
  { text: 'Changelog', link: 'https://github.com/LessUp/encoding/blob/master/CHANGELOG.md' },
]

export default defineConfig({
  // Default to English
  lang: 'en-US',
  title: 'Encoding',
  description: 'Compression algorithms collection: classic compression algorithms in C++, Go, and Rust for learning, comparison, and verification',
  base: '/encoding/',
  cleanUrls: true,

  sitemap: {
    hostname: 'https://lessup.github.io/encoding/',
  },

  head: [
    ['link', { rel: 'canonical', href: 'https://lessup.github.io/encoding/' }],
    ['meta', { name: 'theme-color', content: '#0f172a' }],
    ['meta', { property: 'og:type', content: 'website' }],
    ['meta', { property: 'og:title', content: 'Encoding | Compression Algorithms Collection' }],
    ['meta', { property: 'og:description', content: 'Classic compression algorithms in C++, Go, and Rust for learning, comparison and cross-language verification' }],
    ['meta', { property: 'og:url', content: 'https://lessup.github.io/encoding/' }],
    ['meta', { name: 'twitter:card', content: 'summary' }],
    ['meta', { name: 'twitter:title', content: 'Encoding | Compression Algorithms Collection' }],
    ['meta', { name: 'twitter:description', content: 'Classic compression algorithms in C++, Go, and Rust for learning, comparison and cross-language verification' }],
    ['link', { rel: 'icon', type: 'image/svg+xml', href: '/encoding/logo.svg' }],
  ],

  markdown: {
    lineNumbers: true,
    languageAlias: {
      cuda: 'cpp',
    },
  },

  lastUpdated: true,

  // Internationalization configuration
  locales: {
    root: {
      label: 'English',
      lang: 'en-US',
      link: '/en/',
      themeConfig: {
        nav: sharedNav('en'),
        sidebar: sharedSidebar['/en/'],
        editLink: {
          pattern: 'https://github.com/LessUp/encoding/edit/master/docs/:path',
          text: 'Edit this page on GitHub',
        },
        footer: {
          message: 'Released under the MIT License',
          copyright: 'Copyright © 2025-2026 LessUp',
        },
        outline: {
          level: [2, 3],
          label: 'On this page',
        },
        lastUpdated: {
          text: 'Last updated',
        },
        docFooter: {
          prev: 'Previous page',
          next: 'Next page',
        },
        returnToTopLabel: 'Return to top',
        sidebarMenuLabel: 'Menu',
        darkModeSwitchLabel: 'Theme',
      },
    },
    zh: {
      label: '简体中文',
      lang: 'zh-CN',
      link: '/zh/',
      themeConfig: {
        nav: sharedNav('zh'),
        sidebar: sharedSidebar['/zh/'],
        editLink: {
          pattern: 'https://github.com/LessUp/encoding/edit/master/docs/:path',
          text: '在 GitHub 上编辑此页',
        },
        footer: {
          message: '基于 MIT 许可证发布',
          copyright: '版权所有 © 2025-2026 LessUp',
        },
        outline: {
          level: [2, 3],
          label: '本页内容',
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
      },
    },
  },

  themeConfig: {
    // Social links (shared)
    socialLinks: [
      { icon: 'github', link: 'https://github.com/LessUp/encoding' },
    ],

    // Search (shared)
    search: {
      provider: 'local',
      options: {
        locales: {
          zh: {
            translations: {
              button: {
                buttonText: '搜索文档',
                buttonAriaLabel: '搜索文档',
              },
              modal: {
                noResultsText: '无法找到相关结果',
                resetButtonTitle: '清除查询条件',
                footer: {
                  selectText: '选择',
                  navigateText: '切换',
                  closeText: '关闭',
                },
              },
            },
          },
        },
      },
    },

    externalLinkIcon: true,
  },
})
