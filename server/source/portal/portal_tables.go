package portal

import (
	"context"

	"github.com/fuermoya/design/server/global"
	"github.com/fuermoya/design/server/model/portal"
	"github.com/fuermoya/design/server/service/system"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const initOrderPortal = system.InitOrderExternal + 1

type initPortal struct{}

// auto run
func init() {
	system.RegisterInit(initOrderPortal, &initPortal{})
}

func (i *initPortal) MigrateTable(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	return ctx, db.AutoMigrate(
		&portal.SysArticle{},
		&portal.SysCategory{},
		&portal.SysTag{},
		&portal.SysTheme{},
	)
}

func (i *initPortal) TableCreated(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	return db.Migrator().HasTable(&portal.SysArticle{})
}

func (i initPortal) InitializerName() string {
	return "portal_tables"
}

func (i *initPortal) InitializeData(ctx context.Context) (next context.Context, err error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	next = ctx

	// 初始化默认分类
	if err = i.initDefaultCategories(db); err != nil {
		return ctx, errors.Wrap(err, "初始化默认分类失败")
	}

	// 初始化默认标签
	if err = i.initDefaultTags(db); err != nil {
		return ctx, errors.Wrap(err, "初始化默认标签失败")
	}

	// 初始化默认主题
	if err = i.initDefaultTheme(db); err != nil {
		return ctx, errors.Wrap(err, "初始化默认主题失败")
	}

	return next, nil
}

func (i *initPortal) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	var count int64
	if errors.Is(db.Model(&portal.SysCategory{}).Count(&count).Error, gorm.ErrRecordNotFound) {
		return false
	}
	return count > 0
}

// 初始化默认分类
func (i *initPortal) initDefaultCategories(db *gorm.DB) error {
	var count int64
	db.Model(&portal.SysCategory{}).Count(&count)
	if count == 0 {
		categories := []portal.SysCategory{
			{
				Name:        "历程",
				Description: "公司历程",
				Sort:        1,
				Status:      1,
			},
			{
				Name:        "生活",
				Description: "生活动态",
				Sort:        2,
				Status:      1,
			},
			{
				Name:        "新闻",
				Description: "新闻资讯",
				Sort:        3,
				Status:      1,
			},
		}

		if err := db.Create(&categories).Error; err != nil {
			return err
		}
		global.LOG.Info("init default categories success")
	}
	return nil
}

// 初始化默认标签
func (i *initPortal) initDefaultTags(db *gorm.DB) error {
	var count int64
	db.Model(&portal.SysTag{}).Count(&count)
	if count == 0 {
		tags := []portal.SysTag{
			{
				Name:   "可靠",
				Color:  "#00ADD8",
				Status: 1,
			},
			{
				Name:   "性价比高",
				Color:  "#4FC08D",
				Status: 1,
			},
			{
				Name:   "环保",
				Color:  "#F7DF1E",
				Status: 1,
			},
			{
				Name:   "可持续",
				Color:  "#3776AB",
				Status: 1,
			},
			{
				Name:   "热门",
				Color:  "#FF6B6B",
				Status: 1,
			},
		}

		if err := db.Create(&tags).Error; err != nil {
			return err
		}
		global.LOG.Info("init default tags success")
	}
	return nil
}

// 初始化默认主题
func (i *initPortal) initDefaultTheme(db *gorm.DB) error {
	var count int64
	db.Model(&portal.SysTheme{}).Count(&count)
	if count == 0 {
		defaultConfig := `{
			"primaryColor": "#409EFF",
			"secondaryColor": "#67C23A",
			"backgroundColor": "#ffffff",
			"textColor": "#333333",
			"fontFamily": "Arial, sans-serif",
			"fontSize": "14px",
			"borderRadius": "4px",
			"shadow": "0 2px 4px rgba(0,0,0,0.1)",
			"siteName": "为你提供被广泛验证的生产力工具",
			"siteDescription": "强大易用，值得信赖",
			"footerText": "© 2025 门户网站. All rights reserved."
		}`

		theme := portal.SysTheme{
			Name:        "默认主题",
			Description: "系统默认主题",
			IsActive:    true,
			Config:      defaultConfig,
			Sort:        1,
		}

		if err := db.Create(&theme).Error; err != nil {
			return err
		}
		global.LOG.Info("init default theme success")
	}
	return nil
}
