import { defineConfig } from 'vitepress'
import { withMermaid } from 'vitepress-plugin-mermaid'
import llmstxt from 'vitepress-plugin-llms'

const rawBase = process.env.VITEPRESS_BASE
const base = rawBase
  ? rawBase.startsWith('/')
    ? rawBase.endsWith('/') ? rawBase : `${rawBase}/`
    : `/${rawBase}/`
  : '/'

// Shared sidebar configuration
const sharedSidebar = {
  '/en/': [
    {
      text: 'Getting Started',
      items: [
        { text: 'Introduction', link: '/en/' },
        { text: 'Quick Start', link: '/en/guide/getting-started' },
        { text: 'Architecture', link: '/en/guide/architecture' },
        { text: 'Project Structure', link: '/en/guide/project-structure' },
      ],
    },
    {
      text: 'Algorithms',
      items: [
        { text: 'Overview', link: '/en/guide/algorithms' },
        { text: 'Huffman Coding', link: '/en/algorithms/huffman' },
        { text: 'Arithmetic Coding', link: '/en/algorithms/arithmetic' },
        { text: 'Range Coder', link: '/en/algorithms/range' },
        { text: 'Run-Length Encoding', link: '/en/algorithms/rle' },
      ],
    },
    {
      text: 'API Reference',
      items: [
        { text: 'Streaming API', link: '/en/api/streaming' },
        { text: 'Go Library', link: '/en/api/go' },
        { text: 'Rust Crate', link: '/en/api/rust' },
        { text: 'C++ Header', link: '/en/api/cpp' },
      ],
    },
    {
      text: 'Benchmarks & Testing',
      items: [
        { text: 'Performance Results', link: '/en/benchmarks/results' },
        { text: 'How to Run', link: '/en/benchmarks/how-to-run' },
        { text: 'Cross-Language Testing', link: '/en/testing/cross-language' },
      ],
    },
    {
      text: 'Reference',
      items: [
        { text: 'OpenSpec Specs', link: 'https://github.com/LessUp/compress-kit/tree/master/openspec/specs' },
        { text: 'Contributing', link: '/en/guide/contributing' },
        { text: 'Changelog', link: '/en/release-notes/changelog' },
      ],
    },
  ],
  '/zh/': [
    {
      text: '开始使用',
      items: [
        { text: '项目介绍', link: '/zh/' },
        { text: '快速开始', link: '/zh/guide/getting-started' },
        { text: '架构设计', link: '/zh/guide/architecture' },
        { text: '项目结构', link: '/zh/guide/project-structure' },
      ],
    },
    {
      text: '算法详解',
      items: [
        { text: '算法综述', link: '/zh/guide/algorithms' },
        { text: '霍夫曼编码', link: '/zh/algorithms/huffman' },
        { text: '算术编码', link: '/zh/algorithms/arithmetic' },
        { text: '区间编码', link: '/zh/algorithms/range' },
        { text: '行程编码', link: '/zh/algorithms/rle' },
      ],
    },
    {
      text: 'API 参考',
      items: [
        { text: 'Streaming API', link: '/zh/api/streaming' },
        { text: 'Go 库', link: '/zh/api/go' },
        { text: 'Rust 包', link: '/zh/api/rust' },
        { text: 'C++ 头文件', link: '/zh/api/cpp' },
      ],
    },
    {
      text: '基准测试',
      items: [
        { text: '性能结果', link: '/zh/benchmarks/results' },
        { text: '如何运行', link: '/zh/benchmarks/how-to-run' },
        { text: '跨语言测试', link: '/zh/testing/cross-language' },
      ],
    },
    {
      text: '参考',
      items: [
        { text: 'OpenSpec 规范', link: 'https://github.com/LessUp/compress-kit/tree/master/openspec/specs' },
        { text: '参与贡献', link: '/zh/guide/contributing' },
        { text: '更新日志', link: '/zh/release-notes/changelog' },
      ],
    },
  ],
}

// Shared nav configuration
const sharedNav = (lang: string) => [
  {
    text: lang === 'zh' ? '首页' : 'Home',
    link: lang === 'zh' ? '/zh/' : '/en/',
    activeMatch: lang === 'zh' ? '^/zh/$' : '^/en/$'
  },
  {
    text: lang === 'zh' ? '开始' : 'Get Started',
    link: lang === 'zh' ? '/zh/guide/getting-started' : '/en/guide/getting-started',
    activeMatch: lang === 'zh' ? '/zh/guide/' : '/en/guide/'
  },
  {
    text: lang === 'zh' ? '算法' : 'Algorithms',
    link: lang === 'zh' ? '/zh/guide/algorithms' : '/en/guide/algorithms',
    activeMatch: lang === 'zh' ? '/zh/algorithms/' : '/en/algorithms/'
  },
  {
    text: 'API',
    link: lang === 'zh' ? '/zh/api/go' : '/en/api/go',
    activeMatch: lang === 'zh' ? '/zh/api/' : '/en/api/'
  },
  {
    text: lang === 'zh' ? '基准' : 'Benchmarks',
    link: lang === 'zh' ? '/zh/benchmarks/results' : '/en/benchmarks/results',
    activeMatch: lang === 'zh' ? '/zh/benchmarks/' : '/en/benchmarks/'
  },
]

export default withMermaid(defineConfig({
  base,
  title: 'CompressKit',
  titleTemplate: ':title | CompressKit',
  description: 'Classic lossless compression algorithms in C++17, Go, and Rust with cross-language binary verification.',
  cleanUrls: true,
  lastUpdated: true,
  appearance: true,

  sitemap: {
    hostname: 'https://lessup.github.io/compress-kit/',
  },

  locales: {
    root: {
      label: 'English',
      lang: 'en-US',
      link: '/en/',
      themeConfig: {
        nav: sharedNav('en'),
        sidebar: sharedSidebar['/en/'],
        editLink: {
          pattern: 'https://github.com/LessUp/compress-kit/edit/master/docs/:path',
          text: 'Edit this page on GitHub',
        },
        footer: false,
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
        search: {
          provider: 'local',
          options: {
            translations: {
              button: {
                buttonText: 'Search',
                buttonAriaLabel: 'Search documentation',
              },
              modal: {
                noResultsText: 'No results found',
                resetButtonTitle: 'Clear search',
                footer: {
                  selectText: 'to select',
                  navigateText: 'to navigate',
                  closeText: 'to close',
                },
              },
            },
          },
        },
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
          pattern: 'https://github.com/LessUp/compress-kit/edit/master/docs/:path',
          text: '在 GitHub 上编辑此页',
        },
        footer: false,
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
        search: {
          provider: 'local',
          options: {
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
  },

  themeConfig: {
    outline: [2, 3],
    search: { provider: 'local' },
    socialLinks: [
      { icon: 'github', link: 'https://github.com/LessUp/compress-kit' },
    ],
    logo: {
      light: '/logo.svg',
      dark: '/logo-dark.svg',
      alt: 'CompressKit Logo'
    },
    siteTitle: 'CompressKit',
    externalLinkIcon: true,
  },

  markdown: {
    lineNumbers: true,
    languageAlias: {
      cuda: 'cpp',
    },
  },

  head: [
    ['link', { rel: 'canonical', href: 'https://lessup.github.io/compress-kit/' }],
    ['meta', { charset: 'UTF-8' }],
    ['meta', { name: 'viewport', content: 'width=device-width, initial-scale=1.0' }],
    ['meta', { name: 'theme-color', content: '#2563eb', media: '(prefers-color-scheme: light)' }],
    ['meta', { name: 'theme-color', content: '#0f172a', media: '(prefers-color-scheme: dark)' }],
    ['meta', { name: 'keywords', content: 'compression algorithms, huffman coding, arithmetic coding, range coder, run-length encoding, C++, Go, Rust, lossless compression, cross-language conformance' }],
    ['meta', { name: 'author', content: 'CompressKit Team' }],
    ['meta', { name: 'robots', content: 'index, follow' }],
    ['meta', { property: 'og:type', content: 'website' }],
    ['meta', { property: 'og:locale', content: 'en_US' }],
    ['meta', { property: 'og:title', content: 'CompressKit | Compression Algorithms Collection' }],
    ['meta', { property: 'og:description', content: 'Classic lossless compression algorithms in C++17, Go, and Rust with cross-language binary verification.' }],
    ['meta', { property: 'og:url', content: 'https://lessup.github.io/compress-kit/' }],
    ['meta', { property: 'og:site_name', content: 'CompressKit' }],
    ['meta', { property: 'og:image', content: '/compress-kit/og-image.svg' }],
    ['link', { rel: 'icon', type: 'image/svg+xml', href: '/compress-kit/logo.svg' }],
  ],

  vite: {
    plugins: [llmstxt()],
    resolve: {
      alias: {
        '@theme': '/.vitepress/theme',
      },
    },
  },
}))
