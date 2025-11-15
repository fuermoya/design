package system

import (
	"context"
	"fmt"
	"github.com/fuermoya/design/server/global"
	sysModel "github.com/fuermoya/design/server/model/system"
	"github.com/fuermoya/design/server/service/system"

	"go.uber.org/zap"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const initOrderDictDetail = initOrderDict + 1

type initDictDetail struct{}

// auto run
func init() {
	system.RegisterInit(initOrderDictDetail, &initDictDetail{})
}

func (i *initDictDetail) MigrateTable(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	return ctx, db.AutoMigrate(&sysModel.SysDictionaryDetail{})
}

func (i *initDictDetail) TableCreated(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	return db.Migrator().HasTable(&sysModel.SysDictionaryDetail{})
}

func (i initDictDetail) InitializerName() string {
	return sysModel.SysDictionaryDetail{}.TableName()
}

func (i *initDictDetail) InitializeData(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	dicts, ok := ctx.Value(initDict{}.InitializerName()).([]sysModel.SysDictionary)
	if !ok {
		return ctx, errors.Wrap(system.ErrMissingDependentContext,
			fmt.Sprintf("未找到 %s 表初始化数据", sysModel.SysDictionary{}.TableName()))
	}
	True := true
	dicts[0].SysDictionaryDetails = []sysModel.SysDictionaryDetail{
		{Label: "男", Value: 1, Status: &True, Sort: 1},
		{Label: "女", Value: 2, Status: &True, Sort: 2},
	}

	newDict := make([]sysModel.SysDictionary, 0)
	for _, v := range dicts {
		var count int64
		err := db.Raw("select count(1) from sys_dictionaries where id = ?", v.ID).Scan(&count).Error
		if err != nil {
			global.LOG.Error("校验菜单数据错误", zap.Error(err))
		}
		if count <= 0 { // 判断是否存在数据
			newDict = append(newDict, v)
		}
	}

	if len(newDict) > 0 {
		err := db.Create(&newDict).Error
		if err != nil {
			return ctx, errors.Wrap(err, sysModel.SysDictionary{}.TableName()+"表数据初始化失败!")
		}
		for _, dict := range newDict {
			if err := db.Model(&dict).Association("SysDictionaryDetails").
				Replace(dict.SysDictionaryDetails); err != nil {
				return ctx, errors.Wrap(err, sysModel.SysDictionaryDetail{}.TableName()+"表数据初始化失败!")
			}
		}
	}

	return ctx, nil
}

func (i *initDictDetail) DataInserted(ctx context.Context) bool {
	//db, ok := ctx.Value("db").(*gorm.DB)
	//if !ok {
	//	return false
	//}
	//var dict sysModel.SysDictionary
	//if err := db.Preload("SysDictionaryDetails").
	//	First(&dict, &sysModel.SysDictionary{Name: "数据库bool类型"}).Error; err != nil {
	//	return false
	//}
	//return len(dict.SysDictionaryDetails) > 0 && dict.SysDictionaryDetails[0].Label == "tinyint"
	return false
}
