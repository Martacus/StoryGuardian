/**
 * Converts a timestamp to a human-readable relative time string.
 * e.g. "Just now", "5m ago", "3h ago", "2d ago"
 *
 * Accepts an ISO string (as returned by Wails Go time.Time bindings),
 * a Date object, or null/undefined.
 */
export function timeAgo(date: string | Date | null | undefined): string {
  if (!date) return 'Never'
  const d = typeof date === 'string' ? new Date(date) : date
  if (isNaN(d.getTime())) return 'Unknown'

  const seconds = Math.floor((Date.now() - d.getTime()) / 1000)

  if (seconds < 5)  return 'Just now'
  if (seconds < 60) return `${seconds}s ago`

  const minutes = Math.floor(seconds / 60)
  if (minutes < 60) return `${minutes}m ago`

  const hours = Math.floor(minutes / 60)
  if (hours < 24) return `${hours}h ago`

  const days = Math.floor(hours / 24)
  if (days < 30) return `${days}d ago`

  const months = Math.floor(days / 30)
  if (months < 12) return `${months}mo ago`

  return `${Math.floor(months / 12)}y ago`
}
