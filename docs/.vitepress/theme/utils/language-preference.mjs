export const LOCALE_STORAGE_KEY = 'compresskit-lang'

export function normalizeBase(base = '/') {
  if (!base || base === '/') {
    return '/'
  }

  const rooted = base.startsWith('/') ? base : `/${base}`
  return rooted.endsWith('/') ? rooted : `${rooted}/`
}

export function normalizeLocale(locale) {
  if (!locale) {
    return null
  }

  const normalized = String(locale).trim().toLowerCase()
  if (normalized.startsWith('zh')) {
    return 'zh'
  }
  if (normalized.startsWith('en')) {
    return 'en'
  }
  return null
}

export function buildLocalePath(base, locale) {
  return `${normalizeBase(base)}${locale}/`
}

export function readSavedLocale(storage) {
  if (!storage?.getItem) {
    return null
  }
  return normalizeLocale(storage.getItem(LOCALE_STORAGE_KEY))
}

export function persistLocale(storage, locale) {
  const normalized = normalizeLocale(locale)
  if (!normalized || !storage?.setItem) {
    return null
  }
  storage.setItem(LOCALE_STORAGE_KEY, normalized)
  return normalized
}

export function getPreferredLocale({ savedLocale, browserLanguage }) {
  return normalizeLocale(savedLocale) ?? normalizeLocale(browserLanguage) ?? 'en'
}

function isRootLandingPath(pathname, base) {
  const normalizedBase = normalizeBase(base)
  const trimmedBase = normalizedBase === '/' ? '/' : normalizedBase.slice(0, -1)
  return pathname === normalizedBase || pathname === trimmedBase
}

export function getLandingRedirectTarget({ pathname, base = '/', savedLocale, browserLanguage }) {
  if (!isRootLandingPath(pathname, base)) {
    return null
  }

  const locale = getPreferredLocale({ savedLocale, browserLanguage })
  return buildLocalePath(base, locale)
}
