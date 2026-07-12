import test from 'node:test'
import assert from 'node:assert/strict'

import {
  LOCALE_STORAGE_KEY,
  getLandingRedirectTarget,
  normalizeBase,
  normalizeLocale,
} from '../.vitepress/theme/utils/language-preference.mjs'

test('normalizeBase produces rooted trailing-slash paths', () => {
  assert.equal(normalizeBase('/compress-kit'), '/compress-kit/')
  assert.equal(normalizeBase('compress-kit'), '/compress-kit/')
  assert.equal(normalizeBase('/'), '/')
})

test('saved locale wins on root landing redirect', () => {
  assert.equal(
    getLandingRedirectTarget({
      pathname: '/compress-kit/',
      base: '/compress-kit',
      savedLocale: 'zh',
      browserLanguage: 'en-US',
    }),
    '/compress-kit/zh/'
  )
})

test('browser language drives redirect when no saved locale exists', () => {
  assert.equal(
    getLandingRedirectTarget({
      pathname: '/compress-kit/',
      base: '/compress-kit/',
      savedLocale: null,
      browserLanguage: 'zh-CN',
    }),
    '/compress-kit/zh/'
  )

  assert.equal(
    getLandingRedirectTarget({
      pathname: '/compress-kit/',
      base: '/compress-kit/',
      savedLocale: null,
      browserLanguage: 'fr-FR',
    }),
    '/compress-kit/en/'
  )
})

test('non-root pages do not trigger landing redirect', () => {
  assert.equal(
    getLandingRedirectTarget({
      pathname: '/compress-kit/en/guide/getting-started',
      base: '/compress-kit/',
      savedLocale: 'zh',
      browserLanguage: 'zh-CN',
    }),
    null
  )
})

test('normalizeLocale accepts only shipped locales', () => {
  assert.equal(LOCALE_STORAGE_KEY, 'compresskit-lang')
  assert.equal(normalizeLocale('zh'), 'zh')
  assert.equal(normalizeLocale('en-US'), 'en')
  assert.equal(normalizeLocale('jp'), null)
  assert.equal(normalizeLocale(''), null)
})
