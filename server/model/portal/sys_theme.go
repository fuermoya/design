package portal

import (
	"github.com/fuermoya/design/server/global"
)

// SysTheme 主题表
type SysTheme struct {
	global.BASE_MODEL
	Name        string `json:"name" form:"name" gorm:"column:name;comment:主题名称;size:50;"`
	Description string `json:"description" form:"description" gorm:"column:description;comment:主题描述;size:200;"`
	IsActive    bool   `json:"isActive" form:"isActive" gorm:"column:is_active;comment:是否激活;default:false"`
	Config      string `json:"config" form:"config" gorm:"column:config;comment:主题配置;type:text;"`
	Preview     string `json:"preview" form:"preview" gorm:"column:preview;comment:预览图片;size:255;"`
	Sort        int    `json:"sort" form:"sort" gorm:"column:sort;comment:排序;default:0"`
}

func (SysTheme) TableName() string {
	return "sys_themes"
}

// ThemeConfig 主题配置结构
type ThemeConfig struct {
	PrimaryColor    string `json:"primaryColor"`    // 主色调
	SecondaryColor  string `json:"secondaryColor"`  // 次要色调
	BackgroundColor string `json:"backgroundColor"` // 背景色
	TextColor       string `json:"textColor"`       // 文字颜色
	FontFamily      string `json:"fontFamily"`      // 字体
	FontSize        string `json:"fontSize"`        // 字体大小
	BorderRadius    string `json:"borderRadius"`    // 圆角
	Shadow          string `json:"shadow"`          // 阴影
	Logo            string `json:"logo"`            // Logo
	Favicon         string `json:"favicon"`         // 网站图标
	SiteName        string `json:"siteName"`        // 网站名称
	SiteDescription string `json:"siteDescription"` // 网站描述
	FooterText      string `json:"footerText"`      // 页脚文字
}
