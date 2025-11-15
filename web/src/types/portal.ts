// 文章相关类型
export interface SysArticle {
  ID: number
  title: string
  content: string
  summary: string
  coverImage?: string
  categoryId: number
  category?: SysCategory
  tags?: SysTag[]
  authorID: number
  author?: SysUser
  status: number
  viewCount: number
  likeCount: number
  publishTime?: string
  sort: number
  CreatedAt: string
  UpdatedAt: string
  selected?: boolean
}

// 分类类型
export interface SysCategory {
  ID: number
  name: string
  description: string
  sort: number
  status: number
  CreatedAt: string
  UpdatedAt: string
}

// 标签类型
export interface SysTag {
  ID: number
  name: string
  color: string
  status: number
  CreatedAt: string
  UpdatedAt: string
}

// 用户类型
export interface SysUser {
  ID: number
  username: string
  nickname: string
  avatar?: string
  CreatedAt: string
  UpdatedAt: string
}

// 主题类型
export interface SysTheme {
  ID: number
  name: string
  description: string
  isActive: boolean
  config: string
  preview?: string
  sort: number
  CreatedAt: string
  UpdatedAt: string
  selected?: boolean
}

// 主题配置类型
export interface ThemeConfig {
  primaryColor: string
  secondaryColor: string
  backgroundColor: string
  textColor: string
  fontFamily: string
  fontSize: string
  borderRadius: string
  shadow: string
  logo?: string
  favicon?: string
  siteName: string
  siteDescription: string
  footerText: string
  contactPhone?: string
  contactEmail?: string
  contactAddress?: string
}
