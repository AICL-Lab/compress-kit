import test from 'node:test'
import assert from 'node:assert/strict'
import { readFile } from 'node:fs/promises'
import { getBarHeight, isBestCompressionRatio } from '../.vitepress/theme/components/benchmarkChartMetrics.js'

const loadCompressionRatios = async dataset => {
  const source = await readFile(new URL('../.vitepress/data/benchmarks.json', import.meta.url), 'utf8')
  const { results } = JSON.parse(source)
  return results
    .filter(result => result.dataset === dataset)
    .map(result => result.compressionRatio)
}

test('compression ratio bars reward smaller output/input ratios', async () => {
  const values = await loadCompressionRatios('textlike_10MiB')
  const bestRatio = Math.min(...values)
  const worstRatio = Math.max(...values)

  assert.ok(
    getBarHeight('compressionRatio', bestRatio, worstRatio) >
      getBarHeight('compressionRatio', worstRatio, worstRatio),
    'smaller compression ratios should render taller bars'
  )
})

test('compression ratio bars keep near-identical ratios visually close', () => {
  const better = getBarHeight('compressionRatio', 0.743, 0.746)
  const worse = getBarHeight('compressionRatio', 0.746, 0.746)

  assert.ok(better > worse, 'smaller compression ratios should still render taller bars')
  assert.ok(
    better - worse < 1,
    'near-identical compression ratios should not be stretched across the full chart'
  )
})

test('compression ratio bars bottom out once compression stops helping', () => {
  assert.ok(
    getBarHeight('compressionRatio', 0.95, 1) > getBarHeight('compressionRatio', 1, 1),
    'small-but-positive savings should still render above the no-compression baseline'
  )
  assert.equal(
    getBarHeight('compressionRatio', 1.01, 2),
    2,
    'expanding results should stay at the floor because they save no space'
  )
  assert.equal(
    getBarHeight('compressionRatio', 2, 2),
    2,
    'worse expansion should also stay at the floor'
  )
})

test('best compression ratio stays highlighted even when bars are close', async () => {
  const values = await loadCompressionRatios('textlike_10MiB')
  const bestRatio = Math.min(...values)
  const worstRatio = Math.max(...values)

  assert.equal(
    isBestCompressionRatio('compressionRatio', bestRatio, values),
    true,
    'the lowest compression ratio should be highlighted'
  )
  assert.equal(
    isBestCompressionRatio('compressionRatio', worstRatio, values),
    false,
    'higher compression ratios should not be highlighted'
  )
})

test('best compression ratio still gets highlighted among expanding results', () => {
  assert.equal(
    isBestCompressionRatio('compressionRatio', 1.01, [1.01, 2]),
    true,
    'the smaller expansion ratio should still be highlighted as the better result'
  )
  assert.equal(
    isBestCompressionRatio('compressionRatio', 2, [1.01, 2]),
    false,
    'the larger expansion ratio should not be highlighted'
  )
})

test('encode speed bars still reward faster results', () => {
  assert.ok(
    getBarHeight('encodeSpeed', 277, 277) > getBarHeight('encodeSpeed', 58.8, 277),
    'higher encode speeds should still render taller bars'
  )
})
