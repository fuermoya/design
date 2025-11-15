package system

import (
	"context"
	"fmt"
	"github.com/fuermoya/design/server/utils"

	sysModel "github.com/fuermoya/design/server/model/system"
	"github.com/fuermoya/design/server/service/system"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const initOrderUserCasbin = initOrderUser + 1

var casbinServiceApp = new(system.CasbinService)

type initUserCasbin struct{}

// auto run
func init() {
	system.RegisterInit(initOrderUserCasbin, &initUserCasbin{})
}

func (i *initUserCasbin) MigrateTable(ctx context.Context) (context.Context, error) {
	_, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	return ctx, nil
}

func (i *initUserCasbin) TableCreated(ctx context.Context) bool {
	_, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	return true
}

func (i initUserCasbin) InitializerName() string {
	return "initAuthority888"
}

func (i *initUserCasbin) InitializeData(ctx context.Context) (next context.Context, err error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	next = ctx
	// 888 角色满权限
	entities := sysModel.SysAuthority{AuthorityId: 888, AuthorityName: "普通用户", ParentId: utils.Pointer[uint](0), DefaultRouter: "dashboard"}

	menus, ok := ctx.Value(initMenu{}.InitializerName()).([]sysModel.SysBaseMenu)
	if !ok {
		return next, errors.Wrap(errors.New(""), "创建 [888菜单-权限] 关联失败, 未找到菜单表初始化数据")
	}

	api, ok := ctx.Value(initApi{}.InitializerName()).([]sysModel.SysApi)
	if !ok {
		return next, errors.Wrap(errors.New(""), "创建 [888api-权限] 关联失败, 未找到api表初始化数据")
	}

	err = db.Model(&entities).Association("SysBaseMenus").Replace(menus)
	if err != nil {
		fmt.Println(err)
	}

	//做权限去重处理
	deduplicateMap := make(map[string]bool)
	rules := [][]string{}
	for _, v := range api {
		key := "888" + v.Path + v.Method
		if _, ok := deduplicateMap[key]; !ok {
			deduplicateMap[key] = true
			rules = append(rules, []string{"888", v.Path, v.Method})
		}
	}
	casbinServiceApp.ClearCasbin(0, "888")
	e := casbinServiceApp.Casbin()
	success, _ := e.AddPolicies(rules)
	if !success {
		return next, fmt.Errorf("888 权限存在相同api,添加失败")
	}

	return next, err
}

func (i *initUserCasbin) DataInserted(ctx context.Context) bool {
	return false
}
