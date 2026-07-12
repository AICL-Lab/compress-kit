const MIN_BAR_HEIGHT = 5
const MIN_COMPRESSION_BAR_HEIGHT = 2
const MIN_SPEED_SCALE = 300
const COMPRESSION_RATIO_TOLERANCE = 1e-9

export function getBarHeight(metric, value, maxValue) {
  if (metric === 'compressionRatio') {
    const safeRatio = Math.max(value, 0)
    const sizeSaved = 1 - Math.min(safeRatio, 1)
    if (sizeSaved <= 0) {
      return MIN_COMPRESSION_BAR_HEIGHT
    }

    return MIN_COMPRESSION_BAR_HEIGHT + (sizeSaved * (100 - MIN_COMPRESSION_BAR_HEIGHT))
  }

  const max = Math.max(maxValue, MIN_SPEED_SCALE)
  return Math.max((value / max) * 100, MIN_BAR_HEIGHT)
}

export function isBestCompressionRatio(metric, value, values) {
  if (metric !== 'compressionRatio' || values.length === 0) {
    return false
  }

  const bestRatio = Math.min(...values)
  return Math.abs(value - bestRatio) <= COMPRESSION_RATIO_TOLERANCE
}
