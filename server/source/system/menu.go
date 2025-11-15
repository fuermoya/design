package system

import (
	"context"

	"github.com/fuermoya/design/server/global"
	. "github.com/fuermoya/design/server/model/system"
	"github.com/fuermoya/design/server/service/system"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const initOrderMenu = initOrderDictDetail + 1

type initMenu struct{}

// auto run
func init() {
	system.RegisterInit(initOrderMenu, &initMenu{})
}

func (i initMenu) InitializerName() string {
	return SysBaseMenu{}.TableName()
}

func (i *initMenu) MigrateTable(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	return ctx, db.AutoMigrate(
		&SysBaseMenu{},
	)
}

func (i *initMenu) TableCreated(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	m := db.Migrator()
	return m.HasTable(&SysBaseMenu{})
}

func (i *initMenu) InitializeData(ctx context.Context) (next context.Context, err error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	entities := []SysBaseMenu{
		{BASE_MODEL: global.BASE_MODEL{ID: 1}, MenuLevel: 0, Hidden: false, ParentId: "0", Path: "admin", Name: "admin", Sort: 2, Meta: Meta{Title: "系统", Icon: "user"}},
		{BASE_MODEL: global.BASE_MODEL{ID: 2}, MenuLevel: 0, Hidden: false, ParentId: "1", Path: "admin/authority/authority", Name: "authority", Sort: 1, Meta: Meta{Title: "角色管理", Icon: "avatar"}},
		{BASE_MODEL: global.BASE_MODEL{ID: 3}, MenuLevel: 0, Hidden: false, ParentId: "1", Path: "admin/user/user", Name: "user", Sort: 2, Meta: Meta{Title: "用户管理", Icon: "coordinate"}},
		{BASE_MODEL: global.BASE_MODEL{ID: 4}, MenuLevel: 0, Hidden: false, ParentId: "1", Path: "admin/operation/sysOperationRecord", Name: "operation", Sort: 3, Meta: Meta{Title: "操作历史", Icon: "pie-chart"}},
		{BASE_MODEL: global.BASE_MODEL{ID: 5}, MenuLevel: 0, Hidden: true, ParentId: "1", Path: "person/person", Name: "person", Sort: 4, Meta: Meta{Title: "个人信息", Icon: "message"}},
		{BASE_MODEL: global.BASE_MODEL{ID: 6}, MenuLevel: 0, Hidden: false, ParentId: "0", Path: "substance", Name: "substance", Sort: 3, Meta: Meta{Title: "内容", Icon: "tools"}},
		{BASE_MODEL: global.BASE_MODEL{ID: 7}, MenuLevel: 0, Hidden: false, ParentId: "6", Path: "substance/upload/upload", Name: "upload", Sort: 1, Meta: Meta{Title: "媒体库（上传下载）", Icon: "upload"}},
		{BASE_MODEL: global.BASE_MODEL{ID: 8}, MenuLevel: 0, Hidden: false, ParentId: "6", Path: "substance/article/article", Name: "article", Sort: 2, Meta: Meta{Title: "文章管理", Icon: "document"}},
		{BASE_MODEL: global.BASE_MODEL{ID: 9}, MenuLevel: 0, Hidden: false, ParentId: "6", Path: "substance/category/category", Name: "category", Sort: 3, Meta: Meta{Title: "分类管理", Icon: "folder"}},
		{BASE_MODEL: global.BASE_MODEL{ID: 10}, MenuLevel: 0, Hidden: false, ParentId: "6", Path: "substance/tag/tag", Name: "tag", Sort: 4, Meta: Meta{Title: "标签管理", Icon: "collection"}},
		{BASE_MODEL: global.BASE_MODEL{ID: 11}, MenuLevel: 0, Hidden: false, ParentId: "6", Path: "substance/theme/theme", Name: "theme", Sort: 5, Meta: Meta{Title: "主题管理", Icon: "brush"}},
		{BASE_MODEL: global.BASE_MODEL{ID: 12}, MenuLevel: 0, Hidden: false, ParentId: "6", Path: "substance/message/message", Name: "message", Sort: 6, Meta: Meta{Title: "留言管理", Icon: "message"}},
		//TODO max id 12

	}
	menus := make([]SysBaseMenu, 0)
	for _, v := range entities {
		var count int64
		err = db.Raw("select count(1) from sys_base_menus where id = ?", v.ID).Scan(&count).Error
		if err != nil {
			global.LOG.Error("校验菜单数据错误", zap.Error(err))
		}
		if count <= 0 { // 判断是否存在数据
			menus = append(menus, v)
		}
	}

	if len(menus) > 0 {
		if err = db.Create(&menus).Error; err != nil {
			return ctx, errors.Wrap(err, SysBaseMenu{}.TableName()+"表数据初始化失败!")
		}
	}

	next = context.WithValue(ctx, i.InitializerName(), entities)
	return next, nil
}

func (i *initMenu) DataInserted(ctx context.Context) bool {
	//db, ok := ctx.Value("db").(*gorm.DB)
	//if !ok {
	//	return false
	//}
	//if errors.Is(db.Where("path = ?", "autoPkg").First(&SysBaseMenu{}).Error, gorm.ErrRecordNotFound) { // 判断是否存在数据
	//	return false
	//}
	return false
}
