package portal

import (
	"github.com/fuermoya/design/server/global"
	"github.com/fuermoya/design/server/model/portal"
	"github.com/fuermoya/design/server/model/portal/request"
	"gorm.io/gorm"
)

type SysThemeService struct{}

// CreateSysTheme 创建主题
func (sysThemeService *SysThemeService) CreateSysTheme(sysTheme request.SysThemeCreate) (err error) {
	var theme portal.SysTheme
	theme.Name = sysTheme.Name
	theme.Description = sysTheme.Description
	theme.Config = sysTheme.Config
	theme.Preview = sysTheme.Preview
	theme.Sort = sysTheme.Sort
	theme.IsActive = false // 新创建的主题默认不激活

	err = global.DB.Create(&theme).Error
	return err
}

// DeleteSysTheme 删除主题
func (sysThemeService *SysThemeService) DeleteSysTheme(id string) (err error) {
	// 检查是否为激活主题
	var theme portal.SysTheme
	if err := global.DB.Where("id = ?", id).First(&theme).Error; err != nil {
		return err
	}
	if theme.IsActive {
		return gorm.ErrInvalidData
	}

	err = global.DB.Delete(&portal.SysTheme{}, "id = ?", id).Error
	return err
}

// DeleteSysThemeByIds 批量删除主题
func (sysThemeService *SysThemeService) DeleteSysThemeByIds(ids []string) (err error) {
	// 检查是否包含激活主题
	var count int64
	global.DB.Model(&portal.SysTheme{}).Where("id IN ? AND is_active = ?", ids, true).Count(&count)
	if count > 0 {
		return gorm.ErrInvalidData
	}

	err = global.DB.Delete(&[]portal.SysTheme{}, "id in ?", ids).Error
	return err
}

// UpdateSysTheme 更新主题
func (sysThemeService *SysThemeService) UpdateSysTheme(sysTheme request.SysThemeUpdate) (err error) {
	var theme portal.SysTheme
	if err := global.DB.Where("id = ?", sysTheme.ID).First(&theme).Error; err != nil {
		return err
	}

	theme.Name = sysTheme.Name
	theme.Description = sysTheme.Description
	theme.Config = sysTheme.Config
	theme.Preview = sysTheme.Preview
	theme.Sort = sysTheme.Sort

	err = global.DB.Save(&theme).Error
	return err
}

// GetSysTheme 根据id获取主题
func (sysThemeService *SysThemeService) GetSysTheme(id string) (sysTheme portal.SysTheme, err error) {
	err = global.DB.Where("id = ?", id).First(&sysTheme).Error
	return
}

// GetSysThemeInfoList 分页获取主题列表
func (sysThemeService *SysThemeService) GetSysThemeInfoList(info request.SysThemeSearch) (list []portal.SysTheme, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.DB.Model(&portal.SysTheme{})
	var sysThemes []portal.SysTheme

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	err = db.Limit(limit).Offset(offset).Order("sort ASC, created_at DESC").Find(&sysThemes).Error
	return sysThemes, total, err
}

// ActivateTheme 激活主题
func (sysThemeService *SysThemeService) ActivateTheme(id string) (err error) {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		// 先取消所有主题的激活状态
		if err := tx.Model(&portal.SysTheme{}).Where("is_active = ?", true).Update("is_active", false).Error; err != nil {
			return err
		}

		// 激活指定主题
		if err := tx.Model(&portal.SysTheme{}).Where("id = ?", id).Update("is_active", true).Error; err != nil {
			return err
		}

		return nil
	})
}

// GetActiveTheme 获取当前激活的主题
func (sysThemeService *SysThemeService) GetActiveTheme() (sysTheme portal.SysTheme, err error) {
	err = global.DB.Where("is_active = ?", true).First(&sysTheme).Error
	return
}

// GetPortalThemes 获取门户网站主题列表
func (sysThemeService *SysThemeService) GetPortalThemes() (list []portal.SysTheme, err error) {
	err = global.DB.Order("sort ASC, created_at DESC").Find(&list).Error
	return
}
