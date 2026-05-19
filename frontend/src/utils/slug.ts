import slugifyLib from "@sindresorhus/slugify"

export const SLUG_PATTERN = /^[A-Za-z_][A-Za-z0-9_]*$/

export const SLUG_MAX_LEN = 100

export function isValidSlug(value: string): boolean {
  return value.length > 0 && value.length <= SLUG_MAX_LEN && SLUG_PATTERN.test(value)
}

function toBaseSlug(title: string): string {
  let slug = slugifyLib(title, { separator: "_", decamelize: false })
  if (slug === "") slug = "var"
  if (/^[0-9]/.test(slug)) slug = `_${slug}`
  if (slug.length > SLUG_MAX_LEN) slug = slug.slice(0, SLUG_MAX_LEN)
  return slug
}

export function slugify(title: string, occupied: Iterable<string> = []): string {
  const base = toBaseSlug(title)
  const taken = new Set(occupied)
  if (!taken.has(base)) return base

  for (let i = 2; ; i++) {
    const suffix = `_${i}`
    const trimmedBase = base.slice(0, SLUG_MAX_LEN - suffix.length)
    const candidate = `${trimmedBase}${suffix}`
    if (!taken.has(candidate)) return candidate
  }
}
