package portal

import (
	"errors"

	"github.com/fuermoya/design/server/global"
	"github.com/fuermoya/design/server/model/portal"
	portalReq "github.com/fuermoya/design/server/model/portal/request"
	"gorm.io/gorm"
)

type SysTagService struct{}

// CreateSysTag 创建标签记录
func (sysTagService *SysTagService) CreateSysTag(sysTag portalReq.SysTagCreate) (err error) {
	err = global.DB.Create(&portal.SysTag{
		Name:   sysTag.Name,
		Color:  sysTag.Color,
		Status: sysTag.Status,
	}).Error
	return err
}

// DeleteSysTag 删除标签记录
func (sysTagService *SysTagService) DeleteSysTag(ID string) (err error) {
	err = global.DB.Delete(&portal.SysTag{}, "id = ?", ID).Error
	return err
}

// DeleteSysTagByIds 批量删除标签记录
func (sysTagService *SysTagService) DeleteSysTagByIds(IDs []uint) (err error) {
	err = global.DB.Delete(&[]portal.SysTag{}, "id in ?", IDs).Error
	return err
}

// UpdateSysTag 更新标签记录
func (sysTagService *SysTagService) UpdateSysTag(sysTag portalReq.SysTagUpdate) (err error) {
	var oldSysTag portal.SysTag
	err = global.DB.Where("id = ?", sysTag.ID).First(&oldSysTag).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("标签不存在")
	}
	err = global.DB.Model(&oldSysTag).Updates(&portal.SysTag{
		Name:   sysTag.Name,
		Color:  sysTag.Color,
		Status: sysTag.Status,
	}).Error
	return err
}

// GetSysTag 根据id获取标签记录
func (sysTagService *SysTagService) GetSysTag(ID string) (sysTag portal.SysTag, err error) {
	err = global.DB.Where("id = ?", ID).First(&sysTag).Error
	return
}

// GetSysTagInfoList 分页获取标签记录
func (sysTagService *SysTagService) GetSysTagInfoList(info portalReq.SysTagSearch) (list []portal.SysTag, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.DB.Model(&portal.SysTag{})
	var sysTags []portal.SysTag
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.Name != "" {
		db = db.Where("name LIKE ?", "%"+info.Name+"%")
	}
	if info.Keyword != "" {
		db = db.Where("name LIKE ?", "%"+info.Keyword+"%")
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

	err = db.Order("id DESC").Find(&sysTags).Error
	return sysTags, total, err
}

// GetPortalTags 获取门户网站标签列表
func (sysTagService *SysTagService) GetPortalTags() (list []portal.SysTag, err error) {
	err = global.DB.Where("status = ?", 1).Order("id DESC").Find(&list).Error
	return
}
