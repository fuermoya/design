package portal

import (
	"time"

	"github.com/fuermoya/design/server/global"
	"github.com/fuermoya/design/server/model/system"
)

// SysArticle 文章表
type SysArticle struct {
	global.BASE_MODEL
	Title       string         `json:"title" form:"title" gorm:"column:title;comment:文章标题;size:200;"`
	Content     string         `json:"content" form:"content" gorm:"column:content;comment:文章内容;type:longtext;"`
	Summary     string         `json:"summary" form:"summary" gorm:"column:summary;comment:文章摘要;type:text;"`
	CoverImage  string         `json:"coverImage" form:"coverImage" gorm:"column:cover_image;comment:封面图片;size:255;"`
	CategoryID  uint           `json:"categoryId" form:"categoryId" gorm:"column:category_id;comment:分类ID;"`
	Category    SysCategory    `json:"category" gorm:"foreignKey:CategoryID"`
	Tags        []SysTag       `json:"tags" gorm:"many2many:sys_article_tags;"`
	AuthorID    uint           `json:"authorId" form:"authorId" gorm:"column:author_id;comment:作者ID;"`
	Author      system.SysUser `json:"author" gorm:"foreignKey:AuthorID"`
	Status      int            `json:"status" form:"status" gorm:"column:status;comment:状态 1:发布 2:草稿 3:下架;default:2"`
	PublishTime *time.Time     `json:"publishTime" form:"publishTime" gorm:"column:publish_time;comment:发布时间;"`
	Sort        int            `json:"sort" form:"sort" gorm:"column:sort;comment:排序;default:0"`
	ViewCount   int64          `json:"viewCount" gorm:"-"` // 浏览量，不存储到数据库，通过关联查询计算
	LikeCount   int64          `json:"likeCount" gorm:"-"` // 点赞数，不存储到数据库，通过关联查询计算
}

func (SysArticle) TableName() string {
	return "sys_articles"
}

// SysArticleView 文章浏览记录表
type SysArticleView struct {
	global.BASE_MODEL
	ArticleID uint   `json:"articleId" form:"articleId" gorm:"column:article_id;comment:文章ID;"`
	IP        string `json:"ip" form:"ip" gorm:"column:ip;comment:访问IP;size:50;"`
	UserAgent string `json:"userAgent" form:"userAgent" gorm:"column:user_agent;comment:用户代理;size:500;"`
}

func (SysArticleView) TableName() string {
	return "sys_article_views"
}

// SysArticleLike 文章点赞记录表
type SysArticleLike struct {
	global.BASE_MODEL
	ArticleID uint   `json:"articleId" form:"articleId" gorm:"column:article_id;comment:文章ID;"`
	IP        string `json:"ip" form:"ip" gorm:"column:ip;comment:访问IP;size:50;"`
}

func (SysArticleLike) TableName() string {
	return "sys_article_likes"
}

// SysCategory 分类表
type SysCategory struct {
	global.BASE_MODEL
	Name        string `json:"name" form:"name" gorm:"column:name;comment:分类名称;size:50;"`
	Description string `json:"description" form:"description" gorm:"column:description;comment:分类描述;size:200;"`
	Sort        int    `json:"sort" form:"sort" gorm:"column:sort;comment:排序;default:0"`
	Status      int    `json:"status" form:"status" gorm:"column:status;comment:状态 1:启用 2:禁用;default:1"`
}

func (SysCategory) TableName() string {
	return "sys_categories"
}

// SysTag 标签表
type SysTag struct {
	global.BASE_MODEL
	Name   string `json:"name" form:"name" gorm:"column:name;comment:标签名称;size:50;"`
	Color  string `json:"color" form:"color" gorm:"column:color;comment:标签颜色;size:20;default:#409EFF"`
	Status int    `json:"status" form:"status" gorm:"column:status;comment:状态 1:启用 2:禁用;default:1"`
}

func (SysTag) TableName() string {
	return "sys_tags"
}
