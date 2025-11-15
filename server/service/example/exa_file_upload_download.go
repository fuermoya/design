package example

import (
	"errors"
	"mime/multipart"
	"strings"

	"github.com/fuermoya/design/server/global"
	"github.com/fuermoya/design/server/model/example"
	exampleReq "github.com/fuermoya/design/server/model/example/request"
	"github.com/fuermoya/design/server/utils/upload"
)

type FileUploadAndDownloadService struct{}

//@author: [ccfish86](https://github.com/ccfish86)
//@author: [SliverHorn](https://github.com/SliverHorn)
//@object: *FileUploadAndDownloadService
//@function: UploadFile
//@description: 上传文件
//@param: header *multipart.FileHeader, noSave string
//@return: error, model.ExaFileUploadAndDownload

func (e *FileUploadAndDownloadService) UploadFile(header *multipart.FileHeader, noSave string) (err error, file example.ExaFileUploadAndDownload) {
	oss := upload.NewOss()
	filePath, key, uploadErr := oss.UploadFile(header)
	if uploadErr != nil {
		panic(err)
	}
	if noSave == "0" {
		s := strings.Split(header.Filename, ".")
		f := example.ExaFileUploadAndDownload{
			Url:  filePath,
			Name: header.Filename,
			Tag:  s[len(s)-1],
			Key:  key,
		}
		err = global.DB.Create(&f).Error
		return err, f
	}
	return
}

//@author: [ccfish86](https://github.com/ccfish86)
//@author: [SliverHorn](https://github.com/SliverHorn)
//@object: *FileUploadAndDownloadService
//@function: DeleteFile
//@description: 删除文件
//@param: file model.ExaFileUploadAndDownload
//@return: error

func (e *FileUploadAndDownloadService) DeleteFile(file example.ExaFileUploadAndDownload, id string) (err error) {
	var fileFromDb example.ExaFileUploadAndDownload
	err = global.DB.Where("id = ?", id).First(&fileFromDb).Error
	if err != nil {
		return err
	}
	oss := upload.NewOss()
	if err = oss.DeleteFile(fileFromDb.Key); err != nil {
		return errors.New("文件删除失败")
	}
	err = global.DB.Where("id = ?", id).Unscoped().Delete(&file).Error
	return err
}

//@author: [ccfish86](https://github.com/ccfish86)
//@author: [SliverHorn](https://github.com/SliverHorn)
//@object: *FileUploadAndDownloadService
//@function: GetFileRecordInfoList
//@description: 分页获取数据
//@param: info model.ExaFileUploadAndDownloadSearch
//@return: list interface{}, total int64, err error

func (e *FileUploadAndDownloadService) GetFileRecordInfoList(info exampleReq.ExaFileUploadAndDownloadSearch) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	keyword := info.Name
	var fileLists []example.ExaFileUploadAndDownload
	db := global.DB.Model(&example.ExaFileUploadAndDownload{})
	if keyword != "" {
		db = db.Where("name LIKE ?", "%"+keyword+"%")
	}
	if info.Tag != "" {
		db = db.Where("tag = ?", info.Tag)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Order("updated_at desc").Find(&fileLists).Error
	return fileLists, total, err
}
