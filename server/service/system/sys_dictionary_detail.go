package system

import (
	"github.com/fuermoya/design/server/global"
	"github.com/fuermoya/design/server/model/system"
	"github.com/fuermoya/design/server/model/system/request"
)

//@function: CreateSysDictionaryDetail
//@description: 创建字典详情数据
//@param: sysDictionaryDetail model.SysDictionaryDetail
//@return: err error

type DictionaryDetailService struct{}

var DictionaryDetailServiceApi = new(DictionaryDetailService)

func (dictionaryDetailService *DictionaryDetailService) CreateSysDictionaryDetail(sysDictionaryDetail system.SysDictionaryDetail) (err error) {
	err = global.DB.Create(&sysDictionaryDetail).Error
	return err
}

//@function: DeleteSysDictionaryDetail
//@description: 删除字典详情数据
//@param: sysDictionaryDetail model.SysDictionaryDetail
//@return: err error

func (dictionaryDetailService *DictionaryDetailService) DeleteSysDictionaryDetail(sysDictionaryDetail system.SysDictionaryDetail) (err error) {
	err = global.DB.Delete(&sysDictionaryDetail).Error
	return err
}

//@function: UpdateSysDictionaryDetail
//@description: 更新字典详情数据
//@param: sysDictionaryDetail *model.SysDictionaryDetail
//@return: err error

func (dictionaryDetailService *DictionaryDetailService) UpdateSysDictionaryDetail(sysDictionaryDetail *system.SysDictionaryDetail) (err error) {
	err = global.DB.Save(sysDictionaryDetail).Error
	return err
}

//@function: GetSysDictionaryDetail
//@description: 根据id获取字典详情单条数据
//@param: id uint
//@return: sysDictionaryDetail system.SysDictionaryDetail, err error

func (dictionaryDetailService *DictionaryDetailService) GetSysDictionaryDetail(id uint) (sysDictionaryDetail system.SysDictionaryDetail, err error) {
	err = global.DB.Where("id = ?", id).First(&sysDictionaryDetail).Error
	return
}

//@function: GetSysDictionaryDetailInfoList
//@description: 分页获取字典详情列表
//@param: info request.SysDictionaryDetailSearch
//@return: list interface{}, total int64, err error

func (dictionaryDetailService *DictionaryDetailService) GetSysDictionaryDetailInfoList(info request.SysDictionaryDetailSearch) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.DB.Model(&system.SysDictionaryDetail{})
	var sysDictionaryDetails []system.SysDictionaryDetail
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.Label != "" {
		db = db.Where("label LIKE ?", "%"+info.Label+"%")
	}
	if info.Value != 0 {
		db = db.Where("value = ?", info.Value)
	}
	if info.Extend != "" {
		db = db.Where("extend = ?", info.Extend)
	}
	if info.Status != nil {
		db = db.Where("status = ?", info.Status)
	}
	if info.SysDictionaryID != 0 {
		db = db.Where("sys_dictionary_id = ?", info.SysDictionaryID)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Order("sort").Find(&sysDictionaryDetails).Error
	return sysDictionaryDetails, total, err
}

// 按照字典id获取字典全部内容的方法
func (dictionaryDetailService *DictionaryDetailService) GetDictionaryList(dictionaryID uint) (list []system.SysDictionaryDetail, err error) {
	var sysDictionaryDetails []system.SysDictionaryDetail
	err = global.DB.Find(&sysDictionaryDetails, "sys_dictionary_id = ?", dictionaryID).Error
	return sysDictionaryDetails, err
}

// 按照字典type获取字典全部内容的方法
func (dictionaryDetailService *DictionaryDetailService) GetDictionaryListByType(t string) (list []system.SysDictionaryDetail, err error) {
	var sysDictionaryDetails []system.SysDictionaryDetail
	db := global.DB.Model(&system.SysDictionaryDetail{}).Joins("JOIN sys_dictionaries ON sys_dictionaries.id = sys_dictionary_details.sys_dictionary_id")
	err = db.Find(&sysDictionaryDetails, "type = ?", t).Error
	return sysDictionaryDetails, err
}

func (dictionaryDetailService *DictionaryDetailService) GetDictionaryMapByType(t string) map[int]string {
	var sysDictionaryDetails []system.SysDictionaryDetail
	db := global.DB.Model(&system.SysDictionaryDetail{}).Joins("JOIN sys_dictionaries ON sys_dictionaries.id = sys_dictionary_details.sys_dictionary_id")
	db.Find(&sysDictionaryDetails, "type = ?", t)
	m := make(map[int]string, len(sysDictionaryDetails))
	for i := range sysDictionaryDetails {
		v := sysDictionaryDetails[i]
		m[v.Value] = v.Label
	}
	return m
}

// 按照字典id+字典内容value获取单条字典内容
func (dictionaryDetailService *DictionaryDetailService) GetDictionaryInfoByValue(dictionaryID uint, value uint) (detail system.SysDictionaryDetail, err error) {
	var sysDictionaryDetail system.SysDictionaryDetail
	err = global.DB.First(&sysDictionaryDetail, "sys_dictionary_id = ? and value = ?", dictionaryID, value).Error
	return sysDictionaryDetail, err
}

// 按照字典type+字典内容value获取单条字典内容
func (dictionaryDetailService *DictionaryDetailService) GetDictionaryInfoByTypeValue(t string, value uint) (detail system.SysDictionaryDetail, err error) {
	var sysDictionaryDetails system.SysDictionaryDetail
	db := global.DB.Model(&system.SysDictionaryDetail{}).Joins("JOIN sys_dictionaries ON sys_dictionaries.id = sys_dictionary_details.sys_dictionary_id")
	err = db.First(&sysDictionaryDetails, "sys_dictionaries.type = ? and sys_dictionary_details.value = ?", t, value).Error
	return sysDictionaryDetails, err
}
