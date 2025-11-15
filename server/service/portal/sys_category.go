package portal

import (
	"errors"

	"github.com/fuermoya/design/server/global"
	"github.com/fuermoya/design/server/model/portal"
	portalReq "github.com/fuermoya/design/server/model/portal/request"
	"gorm.io/gorm"
)

type SysCategoryService struct{}

// CreateSysCategory 创建分类记录
func (sysCategoryService *SysCategoryService) CreateSysCategory(sysCategory portalReq.SysCategoryCreate) (err error) {
	err = global.DB.Create(&portal.SysCategory{
		Name:        sysCategory.Name,
		Description: sysCategory.Description,
		Sort:        sysCategory.Sort,
		Status:      sysCategory.Status,
	}).Error
	return err
}

// DeleteSysCategory 删除分类记录
func (sysCategoryService *SysCategoryService) DeleteSysCategory(ID string) (err error) {
	err = global.DB.Delete(&portal.SysCategory{}, "id = ?", ID).Error
	return err
}

// DeleteSysCategoryByIds 批量删除分类记录
func (sysCategoryService *SysCategoryService) DeleteSysCategoryByIds(IDs []uint) (err error) {
	err = global.DB.Delete(&[]portal.SysCategory{}, "id in ?", IDs).Error
	return err
}

// UpdateSysCategory 更新分类记录
func (sysCategoryService *SysCategoryService) UpdateSysCategory(sysCategory portalReq.SysCategoryUpdate) (err error) {
	var oldSysCategory portal.SysCategory
	err = global.DB.Where("id = ?", sysCategory.ID).First(&oldSysCategory).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("分类不存在")
	}
	err = global.DB.Model(&oldSysCategory).Updates(&portal.SysCategory{
		Name:        sysCategory.Name,
		Description: sysCategory.Description,
		Sort:        sysCategory.Sort,
		Status:      sysCategory.Status,
	}).Error
	return err
}

// GetSysCategory 根据id获取分类记录
func (sysCategoryService *SysCategoryService) GetSysCategory(ID string) (sysCategory portal.SysCategory, err error) {
	err = global.DB.Where("id = ?", ID).First(&sysCategory).Error
	return
}

// GetSysCategoryInfoList 分页获取分类记录
func (sysCategoryService *SysCategoryService) GetSysCategoryInfoList(info portalReq.SysCategorySearch) (list []portal.SysCategory, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.DB.Model(&portal.SysCategory{})
	var sysCategorys []portal.SysCategory
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.Name != "" {
		db = db.Where("name LIKE ?", "%"+info.Name+"%")
	}
	if info.Keyword != "" {
		db = db.Where("name LIKE ? OR description LIKE ?", "%"+info.Keyword+"%", "%"+info.Keyword+"%")
	}
	if info.Status != 0 {
		db = db.Where("status = ?", info.Status)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Order("sort ASC, id DESC").Find(&sysCategorys).Error
	return sysCategorys, total, err
}

// GetPortalCategories 获取门户网站分类列表
func (sysCategoryService *SysCategoryService) GetPortalCategories() (list []portal.SysCategory, err error) {
	err = global.DB.Where("status = ?", 1).Order("sort ASC, id DESC").Find(&list).Error
	return
}
