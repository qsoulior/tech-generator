const rtf = new Intl.RelativeTimeFormat("ru", { numeric: "auto" })

const thresholds: ReadonlyArray<{ unit: Intl.RelativeTimeFormatUnit; seconds: number }> = [
  { unit: "year", seconds: 60 * 60 * 24 * 365 },
  { unit: "month", seconds: 60 * 60 * 24 * 30 },
  { unit: "day", seconds: 60 * 60 * 24 },
  { unit: "hour", seconds: 60 * 60 },
  { unit: "minute", seconds: 60 },
  { unit: "second", seconds: 1 },
]

export function formatRelativeTime(date: Date, now: Date = new Date()): string {
  const diffSec = Math.round((date.getTime() - now.getTime()) / 1000)
  const abs = Math.abs(diffSec)
  if (abs < 5) return "только что"

  for (const { unit, seconds } of thresholds) {
    if (abs >= seconds) {
      return rtf.format(Math.round(diffSec / seconds), unit)
    }
  }
  return rtf.format(diffSec, "second")
}
