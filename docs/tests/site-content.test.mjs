import test from 'node:test'
import assert from 'node:assert/strict'

import {
  algorithmCatalog,
  benchmarkCatalog,
  buildNav,
  buildSidebar,
  homepageFeatureCatalog,
} from '../.vitepress/data/site-content.mjs'

const stripLocalePrefix = link => link.replace(/^\/(en|zh)\//, '/')

test('sidebar keeps one canonical link graph across locales', () => {
  const english = buildSidebar('en')
  const chinese = buildSidebar('zh')

  assert.equal(english.length, chinese.length, 'both locales should expose same sidebar sections')

  for (let index = 0; index < english.length; index += 1) {
    const left = english[index]
    const right = chinese[index]
    assert.equal(left.items.length, right.items.length, `section ${index} should keep same item count`)
    assert.deepEqual(
      left.items.map(item => stripLocalePrefix(item.link)),
      right.items.map(item => stripLocalePrefix(item.link)),
      `section ${index} should map to same canonical destinations`
    )
  }
})

test('top nav keeps one canonical destination set across locales', () => {
  const english = buildNav('en')
  const chinese = buildNav('zh')

  assert.deepEqual(
    english.map(item => stripLocalePrefix(item.link)),
    chinese.map(item => stripLocalePrefix(item.link)),
    'both locales should keep same nav destination order'
  )
})

test('homepage and benchmark metadata share one canonical algorithm catalog', () => {
  const canonicalAlgorithms = algorithmCatalog.map(entry => entry.id).sort()
  const homepageAlgorithms = homepageFeatureCatalog
    .filter(entry => entry.kind === 'algorithm')
    .map(entry => entry.algorithmId)
    .sort()
  const benchmarkAlgorithms = benchmarkCatalog.algorithms.map(entry => entry.id).sort()

  assert.deepEqual(homepageAlgorithms, canonicalAlgorithms, 'homepage algorithm cards should derive from canonical catalog')
  assert.deepEqual(benchmarkAlgorithms, canonicalAlgorithms, 'benchmark labels should derive from canonical catalog')
})

test('benchmark catalog preserves shipped language order', () => {
  assert.deepEqual(
    benchmarkCatalog.languages.map(entry => entry.id),
    ['cpp'],
    'benchmark legend should stay aligned with shipped languages'
  )
})
