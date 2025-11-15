// 格式化日期
export function formatDate(date: string | Date): string {
  if (!date)
    return '-'

  const d = new Date(date)
  if (Number.isNaN(d.getTime()))
    return '-'

  const year = d.getFullYear()
  const month = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  const hours = String(d.getHours()).padStart(2, '0')
  const minutes = String(d.getMinutes()).padStart(2, '0')
  const seconds = String(d.getSeconds()).padStart(2, '0')

  return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`
}

// 格式化相对时间
export function formatRelativeTime(date: string | Date): string {
  if (!date)
    return '-'

  const d = new Date(date)
  if (Number.isNaN(d.getTime()))
    return '-'

  const now = new Date()
  const diff = now.getTime() - d.getTime()
  const diffMinutes = Math.floor(diff / (1000 * 60))
  const diffHours = Math.floor(diff / (1000 * 60 * 60))
  const diffDays = Math.floor(diff / (1000 * 60 * 60 * 24))

  if (diffMinutes < 1) {
    return '刚刚'
  }
  else if (diffMinutes < 60) {
    return `${diffMinutes}分钟前`
  }
  else if (diffHours < 24) {
    return `${diffHours}小时前`
  }
  else if (diffDays < 7) {
    return `${diffDays}天前`
  }
  else {
    return formatDate(date)
  }
}

// 格式化文件大小
export function formatFileSize(bytes: number): string {
  if (bytes === 0)
    return '0 B'

  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))

  return `${Number.parseFloat((bytes / k ** i).toFixed(2))} ${sizes[i]}`
}

// 格式化数字
export function formatNumber(num: number): string {
  if (num >= 1000000) {
    return `${(num / 1000000).toFixed(1)}M`
  }
  else if (num >= 1000) {
    return `${(num / 1000).toFixed(1)}K`
  }
  return num.toString()
}
