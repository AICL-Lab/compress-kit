const localeOrEnglish = locale => (locale === 'zh' ? 'zh' : 'en')

const localize = (value, locale) => value[localeOrEnglish(locale)] ?? value.en

// Root-locale site: links are relative to the docs root, no locale prefix.
const withLocale = (locale, link) => {
  if (/^https?:\/\//.test(link)) {
    return link
  }

  const normalized = link === '/' ? '/' : `/${link.replace(/^\/+/, '')}`
  return normalized
}

export const algorithmCatalog = [
  {
    id: 'huffman',
    slug: 'huffman',
    icon: '🌳',
    name: { en: 'Huffman Coding', zh: '霍夫曼编码' },
    chartLabel: { en: 'Huffman', zh: 'Huffman' },
    description: {
      en: 'Optimal prefix codes based on symbol frequency. The classic approach to lossless compression.',
      zh: '基于符号频率的最优前缀码。经典的无损压缩方法。',
    },
    compression: { en: 'Medium', zh: '中等' },
    speed: { en: 'Fast', zh: '快速' },
    compressionTag: { en: 'Fast Speed', zh: '速度快' },
    speedLevel: 'fast',
    compressionLevel: 'medium',
    bestFor: {
      en: ['Text files', 'General data', 'Natural language'],
      zh: ['文本文件', '通用数据', '自然语言'],
    },
  },
  {
    id: 'arithmetic',
    slug: 'arithmetic',
    icon: '🧮',
    name: { en: 'Arithmetic Coding', zh: '算术编码' },
    chartLabel: { en: 'Arithmetic', zh: 'Arithmetic' },
    description: {
      en: 'Entire message encoded as a single number. Achieves entropy limit for maximum compression.',
      zh: '整个消息编码为单个数字。达到熵极限，实现最大压缩率。',
    },
    compression: { en: 'High', zh: '高' },
    speed: { en: 'Medium', zh: '中速' },
    compressionTag: { en: 'High Compression', zh: '高压缩率' },
    speedLevel: 'medium',
    compressionLevel: 'high',
    bestFor: {
      en: ['Maximum compression', 'Statistical data', 'Archival storage'],
      zh: ['最大压缩率', '统计型数据', '归档存储'],
    },
  },
  {
    id: 'range',
    slug: 'range',
    icon: '🎯',
    name: { en: 'Range Coder', zh: '区间编码' },
    chartLabel: { en: 'Range', zh: 'Range' },
    description: {
      en: 'Integer-based arithmetic coding. Production-ready balance of speed and compression.',
      zh: '基于整数的算术编码。生产级的速度与压缩率平衡。',
    },
    compression: { en: 'High', zh: '高' },
    speed: { en: 'Fast', zh: '快速' },
    compressionTag: { en: 'Fast + High', zh: '快 + 高压缩' },
    speedLevel: 'fast',
    compressionLevel: 'high',
    bestFor: {
      en: ['Production systems', 'Real-time compression', 'Balanced workloads'],
      zh: ['生产系统', '实时压缩', '平衡型负载'],
    },
  },
  {
    id: 'rle',
    slug: 'rle',
    icon: '📏',
    name: { en: 'Run-Length Encoding', zh: '行程编码' },
    chartLabel: { en: 'RLE', zh: 'RLE' },
    description: {
      en: 'Simple and fast compression for repetitive data. Often used as preprocessing.',
      zh: '针对重复数据的简单快速压缩。常作为预处理步骤使用。',
    },
    compression: { en: 'Variable', zh: '可变' },
    speed: { en: 'Very Fast', zh: '极快' },
    compressionTag: { en: 'Very Fast', zh: '极快' },
    speedLevel: 'very-fast',
    compressionLevel: 'variable',
    bestFor: {
      en: ['Bitmap images', 'Log files', 'Preprocessing step'],
      zh: ['位图图像', '日志文件', '预处理步骤'],
    },
  },
]

export const benchmarkCatalog = {
  algorithms: algorithmCatalog.map(entry => ({
    id: entry.id,
    label: entry.chartLabel,
  })),
  languages: [
    { id: 'cpp', color: '#667eea', label: { en: 'C++', zh: 'C++' } },
  ],
  datasets: [
    { id: 'textlike_10MiB', label: { en: 'Text-like (10 MiB)', zh: '类文本 (10 MiB)' } },
    { id: 'repetitive_10MiB', label: { en: 'Repetitive (10 MiB)', zh: '重复数据 (10 MiB)' } },
    { id: 'small_dictionary_like', label: { en: 'Small dictionary-like sample', zh: '小型词典型样本' } },
  ],
  metrics: [
    { id: 'encodeSpeed', label: { en: 'Encode Speed (MiB/s)', zh: '编码速度 (MiB/s)' } },
    { id: 'decodeSpeed', label: { en: 'Decode Speed (MiB/s)', zh: '解码速度 (MiB/s)' } },
    { id: 'compressionRatio', label: { en: 'Size Saved Relative to Input', zh: '相对输入节省的体积' } },
  ],
  metricOptions: [
    { id: 'encodeSpeed', label: { en: 'Encode Speed', zh: '编码速度' } },
    { id: 'decodeSpeed', label: { en: 'Decode Speed', zh: '解码速度' } },
    { id: 'compressionRatio', label: { en: 'Compression Ratio', zh: '压缩比' } },
  ],
}

export const homepageFeatureCatalog = [
  ...algorithmCatalog.map(entry => ({
    id: entry.id,
    kind: 'algorithm',
    algorithmId: entry.id,
    title: entry.name,
    description: entry.description,
    tags: [
      { label: { en: 'Learn More', zh: '了解更多' }, link: `/algorithms/${entry.slug}` },
      { label: entry.compressionTag },
    ],
  })),
  {
    id: 'benchmarks',
    kind: 'guide',
    title: { en: '📊 Benchmarks', zh: '📊 性能基准' },
    description: {
      en: 'Performance benchmarks across all algorithms. Compare speed and compression.',
      zh: '跨所有算法的性能基准测试。比较速度和压缩率。',
    },
    tags: [
      { label: { en: 'View Results', zh: '查看结果' }, link: '/benchmarks/results' },
      { label: { en: 'Run Tests', zh: '运行测试' }, link: '/benchmarks/how-to-run' },
    ],
  },
]

const navCatalog = [
  {
    id: 'home',
    text: { en: 'Home', zh: '首页' },
    link: '/',
    activeMatch: { en: '^/en/$', zh: '^/$' },
  },
  {
    id: 'guide',
    text: { en: 'Get Started', zh: '开始' },
    link: '/guide/getting-started',
    activeMatch: { en: '/en/guide/', zh: '/guide/' },
  },
  {
    id: 'algorithms',
    text: { en: 'Algorithms', zh: '算法' },
    link: '/guide/algorithms',
    activeMatch: { en: '/en/algorithms/', zh: '/algorithms/' },
  },
  {
    id: 'api',
    text: { en: 'API', zh: 'API' },
    link: '/api/cpp',
    activeMatch: { en: '/en/api/', zh: '/api/' },
  },
  {
    id: 'benchmarks',
    text: { en: 'Benchmarks', zh: '基准' },
    link: '/benchmarks/results',
    activeMatch: { en: '/en/benchmarks/', zh: '/benchmarks/' },
  },
]

const sidebarCatalog = [
  {
    title: { en: 'Getting Started', zh: '开始使用' },
    items: [
      { text: { en: 'Introduction', zh: '项目介绍' }, link: '/' },
      { text: { en: 'Quick Start', zh: '快速开始' }, link: '/guide/getting-started' },
      { text: { en: 'Architecture', zh: '架构设计' }, link: '/guide/architecture' },
    ],
  },
  {
    title: { en: 'Academy', zh: '学院' },
    items: [
      { text: { en: 'Algorithm Academy', zh: '算法学院' }, link: '/academy/' },
      { text: { en: 'Huffman Coding', zh: '霍夫曼编码深度解析' }, link: '/academy/huffman' },
    ],
  },
  {
    title: { en: 'Algorithms', zh: '算法详解' },
    items: [
      { text: { en: 'Overview', zh: '算法综述' }, link: '/guide/algorithms' },
      { text: { en: 'Huffman Coding', zh: '霍夫曼编码' }, link: '/algorithms/huffman' },
      { text: { en: 'Arithmetic Coding', zh: '算术编码' }, link: '/algorithms/arithmetic' },
      { text: { en: 'Range Coder', zh: '区间编码' }, link: '/algorithms/range' },
      { text: { en: 'Run-Length Encoding', zh: '行程编码' }, link: '/algorithms/rle' },
    ],
  },
  {
    title: { en: 'API Reference', zh: 'API 参考' },
    items: [
      { text: { en: 'C++ Header', zh: 'C++ 头文件' }, link: '/api/cpp' },
    ],
  },
  {
    title: { en: 'Benchmarks & Testing', zh: '基准测试' },
    items: [
      { text: { en: 'Performance Results', zh: '性能结果' }, link: '/benchmarks/results' },
      { text: { en: 'How to Run', zh: '如何运行' }, link: '/benchmarks/how-to-run' },
    ],
  },
  {
    title: { en: 'Reference', zh: '参考' },
    items: [
      { text: { en: 'Architecture Design', zh: '系统架构设计' }, link: '/architecture/' },
      { text: { en: 'Changelog', zh: '更新日志' }, link: '/release-notes/changelog' },
    ],
  },
]

export function buildNav(locale) {
  return navCatalog.map(item => ({
    text: localize(item.text, locale),
    link: withLocale(locale, item.link),
    activeMatch: localize(item.activeMatch, locale),
  }))
}

export function buildSidebar(locale) {
  return sidebarCatalog.map(section => ({
    text: localize(section.title, locale),
    items: section.items.map(item => ({
      text: localize(item.text, locale),
      link: withLocale(locale, item.link),
    })),
  }))
}

export function getAlgorithmCards(locale) {
  return algorithmCatalog.map(entry => ({
    id: entry.id,
    slug: entry.slug,
    icon: entry.icon,
    name: localize(entry.name, locale),
    description: localize(entry.description, locale),
    compression: localize(entry.compression, locale),
    speed: localize(entry.speed, locale),
    bestFor: localize(entry.bestFor, locale),
    compressionLevel: entry.compressionLevel,
    speedLevel: entry.speedLevel,
    bestForLabel: localize({ en: 'Best for:', zh: '适合场景：' }, locale),
    learnMoreLabel: localize({ en: 'Learn more', zh: '了解更多' }, locale),
    compressionSuffix: localize({ en: 'Compression', zh: '压缩' }, locale),
    speedSuffix: localize({ en: 'Speed', zh: '速度' }, locale),
  }))
}

export function getBenchmarkContent(locale) {
  return {
    algorithms: benchmarkCatalog.algorithms.map(entry => ({
      id: entry.id,
      label: localize(entry.label, locale),
    })),
    languages: benchmarkCatalog.languages.map(entry => ({
      id: entry.id,
      label: localize(entry.label, locale),
      color: entry.color,
    })),
    datasets: benchmarkCatalog.datasets.map(entry => ({
      id: entry.id,
      label: localize(entry.label, locale),
    })),
    metrics: benchmarkCatalog.metrics.map(entry => ({
      id: entry.id,
      label: localize(entry.label, locale),
    })),
    metricOptions: benchmarkCatalog.metricOptions.map(entry => ({
      id: entry.id,
      label: localize(entry.label, locale),
    })),
    title: localize({ en: 'Performance Comparison', zh: '性能对比' }, locale),
    datasetLabel: localize({ en: 'Dataset:', zh: '数据集：' }, locale),
    metricLabel: localize({ en: 'Metric:', zh: '指标：' }, locale),
    compressionNote: localize(
      {
        en: 'Bars show size saved relative to input; labels show the actual output/input ratio, and the best ratio is highlighted.',
        zh: '柱状图展示相对输入节省的体积，标签显示实际输出/输入比值，最佳压缩比会被高亮。',
      },
      locale
    ),
  }
}

export function getHomepageContent(locale) {
  const currentLocale = localeOrEnglish(locale)
  const alternateLocale = currentLocale === 'en' ? 'zh' : 'en'
  const alternateLabel = currentLocale === 'en' ? '中文' : 'English'

  return {
    subtitle: localize(
      { en: 'Lossless Compression Library', zh: '无损压缩算法库' },
      currentLocale
    ),
    intro: localize(
      {
        en: 'CompressKit provides classic lossless compression algorithms in C++17. Read the implementations, run the same inputs through each algorithm, and verify round-trip correctness.',
        zh: 'CompressKit 提供使用 C++17 实现的经典无损压缩算法。阅读实现、对每个算法运行相同输入、验证 round-trip 正确性。',
      },
      currentLocale
    ),
    sections: {
      algorithms: localize({ en: 'Algorithms', zh: '算法' }, currentLocale),
      quickStart: localize({ en: 'Quick Start', zh: '快速开始' }, currentLocale),
    },
    quickStartCommand: 'git clone https://github.com/LessUp/compress-kit.git && cd compress-kit && make build && make test',
    stats: ['C++17', 'CMake'],
    navLinks: [
      { text: localize({ en: 'Get Started', zh: '快速开始' }, currentLocale), link: withLocale(currentLocale, '/guide/getting-started') },
      { text: 'GitHub', link: 'https://github.com/LessUp/compress-kit' },
      { text: alternateLabel, link: withLocale(alternateLocale, '/') },
    ],
    featureCards: homepageFeatureCatalog.map(entry => ({
      id: entry.id,
      title: localize(entry.title, currentLocale),
      description: localize(entry.description, currentLocale),
      tags: entry.tags.map(tag => ({
        label: localize(tag.label, currentLocale),
        link: tag.link ? withLocale(currentLocale, tag.link) : null,
      })),
    })),
  }
}
