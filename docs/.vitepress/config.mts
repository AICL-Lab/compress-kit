import { defineConfig } from 'vitepress'

// Shared sidebar configuration
const sharedSidebar = {
  '/en/': [
    {
      text: 'Getting Started',
      items: [
        { text: 'Introduction', link: '/en/' },
        { text: 'Quick Start', link: '/en/guide/getting-started' },
        { text: 'Project Structure', link: '/en/guide/project-structure' },
        { text: 'Architecture', link: '/en/guide/architecture' },
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
        { text: 'Specifications', link: 'https://github.com/LessUp/compress-kit/tree/master/specs' },
        { text: 'Contributing', link: '/en/guide/contributing' },
        { text: 'Changelog', link: 'https://github.com/LessUp/compress-kit/blob/master/CHANGELOG.md' },
      ],
    },
  ],
  '/zh/': [
    {
      text: '开始使用',
      items: [
        { text: '项目介绍', link: '/zh/' },
        { text: '快速开始', link: '/zh/guide/getting-started' },
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
        { text: '规范文档', link: 'https://github.com/LessUp/compress-kit/tree/master/specs' },
        { text: '参与贡献', link: '/zh/guide/contributing' },
        { text: '更新日志', link: 'https://github.com/LessUp/compress-kit/blob/master/CHANGELOG.md' },
      ],
    },
  ],
}

// Shared nav configuration
const sharedNav = (lang: string) => [
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

export default defineConfig({
  // Brand Configuration
  lang: 'en-US',
  title: 'CompressKit',
  titleTemplate: ':title | CompressKit',
  description: 'Production-ready compression algorithms in C++17, Go, and Rust. Learn, compare, and verify across languages with identical binary formats.',
  
  // Base URL
  base: '/compress-kit/',
  cleanUrls: true,
  
  // Appearance
  appearance: true,
  
  // Sitemap
  sitemap: {
    hostname: 'https://lessup.github.io/compress-kit/',
  },
  
  // Head Meta Tags
  head: [
    ['link', { rel: 'canonical', href: 'https://lessup.github.io/compress-kit/' }],
    ['meta', { charset: 'UTF-8' }],
    ['meta', { name: 'viewport', content: 'width=device-width, initial-scale=1.0' }],
    ['meta', { name: 'theme-color', content: '#2563eb', media: '(prefers-color-scheme: light)' }],
    ['meta', { name: 'theme-color', content: '#0f172a', media: '(prefers-color-scheme: dark)' }],
    
    // SEO
    ['meta', { name: 'keywords', content: 'compression algorithms, huffman coding, arithmetic coding, range coder, run-length encoding, C++, Go, Rust, lossless compression, cross-language, benchmark' }],
    ['meta', { name: 'author', content: 'CompressKit Team' }],
    ['meta', { name: 'robots', content: 'index, follow' }],
    
    // Open Graph
    ['meta', { property: 'og:type', content: 'website' }],
    ['meta', { property: 'og:locale', content: 'en_US' }],
    ['meta', { property: 'og:title', content: 'CompressKit | Compression Algorithms Collection' }],
    ['meta', { property: 'og:description', content: 'Production-ready compression algorithms in C++17, Go, and Rust. Learn, compare, and verify across languages.' }],
    ['meta', { property: 'og:url', content: 'https://lessup.github.io/compress-kit/' }],
    ['meta', { property: 'og:site_name', content: 'CompressKit' }],
    ['meta', { property: 'og:image', content: '/compress-kit/og-image.png' }],
    ['meta', { property: 'og:image:width', content: '1200' }],
    ['meta', { property: 'og:image:height', content: '630' }],
    
    // Twitter
    ['meta', { name: 'twitter:card', content: 'summary_large_image' }],
    ['meta', { name: 'twitter:site', content: '@compresskit' }],
    ['meta', { name: 'twitter:title', content: 'CompressKit | Compression Algorithms Collection' }],
    ['meta', { name: 'twitter:description', content: 'Production-ready compression algorithms in C++17, Go, and Rust. Learn, compare, and verify across languages.' }],
    ['meta', { name: 'twitter:image', content: '/compress-kit/og-image.png' }],
    
    // Favicon
    ['link', { rel: 'icon', type: 'image/svg+xml', href: '/compress-kit/logo.svg' }],
    ['link', { rel: 'mask-icon', href: '/compress-kit/logo.svg', color: '#2563eb' }],
    ['link', { rel: 'apple-touch-icon', href: '/compress-kit/logo.svg' }],
    
    // Preconnect for fonts
    ['link', { rel: 'preconnect', href: 'https://fonts.googleapis.com' }],
    ['link', { rel: 'preconnect', href: 'https://fonts.gstatic.com', crossorigin: '' }],
    ['link', { rel: 'stylesheet', href: 'https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500;600&display=swap' }],
    
    // Structured Data (JSON-LD)
    ['script', { type: 'application/ld+json' }, JSON.stringify({
      '@context': 'https://schema.org',
      '@type': 'SoftwareApplication',
      'name': 'CompressKit',
      'description': 'Production-ready compression algorithms in C++17, Go, and Rust',
      'url': 'https://lessup.github.io/compress-kit/',
      'applicationCategory': 'DeveloperApplication',
      'operatingSystem': 'Linux, macOS, Windows',
      'softwareVersion': '1.0.0',
      'license': 'https://opensource.org/licenses/MIT',
      'programmingLanguage': ['C++', 'Go', 'Rust'],
      'author': {
        '@type': 'Organization',
        'name': 'LessUp'
      },
      'codeRepository': 'https://github.com/LessUp/compress-kit',
      'featureList': [
        'Huffman Coding implementation',
        'Arithmetic Coding implementation', 
        'Range Coder implementation',
        'Run-Length Encoding implementation',
        'Cross-language binary compatibility',
        'Comprehensive benchmark suite'
      ]
    })],
  ],
  
  // Markdown Configuration
  markdown: {
    lineNumbers: true,
    languageAlias: {
      cuda: 'cpp',
    },
    config: (md) => {
      // Custom markdown enhancements can be added here
    }
  },
  
  // Last Updated
  lastUpdated: true,
  
  // Internationalization Configuration
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
        footer: {
          message: 'Released under the MIT License',
          copyright: 'Copyright © 2025-2026 LessUp. Built with VitePress.',
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
        footer: {
          message: '基于 MIT 许可证发布',
          copyright: '版权所有 © 2025-2026 LessUp. 使用 VitePress 构建。',
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
  
  // Theme Configuration
  themeConfig: {
    // Logo
    logo: {
      light: '/logo.svg',
      dark: '/logo-dark.svg',
      alt: 'CompressKit Logo'
    },
    
    // Site Title
    siteTitle: 'CompressKit',
    
    // Social links
    socialLinks: [
      { icon: 'github', link: 'https://github.com/LessUp/compress-kit' },
    ],
    
    // External link indicator
    externalLinkIcon: true,
  },
  
  // Vite Configuration
  vite: {
    resolve: {
      alias: {
        '@theme': '/.vitepress/theme',
      },
    },
    css: {
      preprocessorOptions: {
        scss: {
          additionalData: '',
        },
      },
    },
  },
})
