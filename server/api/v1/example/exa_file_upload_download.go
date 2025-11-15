package example

import (
	"os"

	"github.com/fuermoya/design/server/global"
	"github.com/fuermoya/design/server/model/common/response"
	"github.com/fuermoya/design/server/model/example"
	exampleReq "github.com/fuermoya/design/server/model/example/request"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type FileUploadAndDownloadApi struct{}

// UploadFile 上传文件
// @Tags ExaFileUploadAndDownload
// @Summary 上传文件
// @Security ApiKeyAuth
// @accept multipart/form-data
// @Produce application/json
// @Param file formData file true "上传文件"
// @Success 200 {object} response.Response{data=example.ExaFileUploadAndDownload,msg=string} "上传成功"
// @Router /fileUploadAndDownload/upload [post]
func (u *FileUploadAndDownloadApi) UploadFile(c *gin.Context) {
	var file example.ExaFileUploadAndDownload
	noSave := c.DefaultQuery("noSave", "0")
	_, header, err := c.Request.FormFile("file")
	if err != nil {
		global.LOG.Error("接收文件失败!", zap.Error(err))
		response.FailWithMessage("接收文件失败", c)
		return
	}
	err, file = fileUploadAndDownloadService.UploadFile(header, noSave) // 文件上传后拿到文件路径
	if err != nil {
		global.LOG.Error("修改数据库链接失败!", zap.Error(err))
		response.FailWithMessage("修改数据库链接失败", c)
		return
	}
	response.OkWithDetailed(file, "上传成功", c)
}

// DeleteFile 删除文件
// @Tags ExaFileUploadAndDownload
// @Summary 删除文件
// @Security ApiKeyAuth
// @Produce application/json
// @Param id query uint true "文件ID"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /fileUploadAndDownload/deleteFile [delete]
func (u *FileUploadAndDownloadApi) DeleteFile(c *gin.Context) {
	var file example.ExaFileUploadAndDownload
	id := c.Query("id")
	err := fileUploadAndDownloadService.DeleteFile(file, id)
	if err != nil {
		global.LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// GetFileList 分页文件列表
// @Tags ExaFileUploadAndDownload
// @Summary 分页文件列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body example.ExaFileUploadAndDownload true "查询参数"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /fileUploadAndDownload/getFileList [post]
func (u *FileUploadAndDownloadApi) GetFileList(c *gin.Context) {
	var pageInfo exampleReq.ExaFileUploadAndDownloadSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := fileUploadAndDownloadService.GetFileRecordInfoList(pageInfo)
	if err != nil {
		global.LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     list,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, "获取成功", c)
	}
}

// DownloadFile 下载文件
// @Tags ExaFileUploadAndDownload
// @Summary 下载文件
// @Security ApiKeyAuth
// @Produce application/octet-stream
// @Param id query uint true "文件ID"
// @Success 200 {file} file "文件下载"
// @Router /fileUploadAndDownload/download [get]
func (u *FileUploadAndDownloadApi) DownloadFile(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		response.FailWithMessage("文件ID不能为空", c)
		return
	}

	var fileFromDb example.ExaFileUploadAndDownload
	err := global.DB.Where("id = ?", id).First(&fileFromDb).Error
	if err != nil {
		global.LOG.Error("文件不存在!", zap.Error(err))
		response.FailWithMessage("文件不存在", c)
		return
	}

	// 构建文件完整路径
	filePath := global.CONFIG.Local.StorePath + "/" + fileFromDb.Key

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		global.LOG.Error("文件不存在!", zap.Error(err))
		response.FailWithMessage("文件不存在", c)
		return
	}

	// 设置响应头
	c.Header("Content-Disposition", "attachment; filename="+fileFromDb.Name)
	c.Header("Content-Type", "application/octet-stream")

	// 发送文件
	c.File(filePath)
}

// PreviewFile 预览文件
// @Tags ExaFileUploadAndDownload
// @Summary 预览文件
// @Security ApiKeyAuth
// @Produce application/octet-stream
// @Param id query uint true "文件ID"
// @Success 200 {file} file "文件预览"
// @Router /fileUploadAndDownload/preview [get]
func (u *FileUploadAndDownloadApi) PreviewFile(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		response.FailWithMessage("文件ID不能为空", c)
		return
	}

	var fileFromDb example.ExaFileUploadAndDownload
	err := global.DB.Where("id = ?", id).First(&fileFromDb).Error
	if err != nil {
		global.LOG.Error("文件不存在!", zap.Error(err))
		response.FailWithMessage("文件不存在", c)
		return
	}

	// 构建文件完整路径
	filePath := global.CONFIG.Local.StorePath + "/" + fileFromDb.Key

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		global.LOG.Error("文件不存在!", zap.Error(err))
		response.FailWithMessage("文件不存在", c)
		return
	}

	// 根据文件扩展名设置Content-Type
	contentType := getContentType(fileFromDb.Name)
	c.Header("Content-Type", contentType)

	// 设置缓存头，允许浏览器缓存预览文件
	c.Header("Cache-Control", "public, max-age=3600")

	// 设置文件名（可选）
	c.Header("Content-Disposition", "inline; filename="+fileFromDb.Name)

	// 发送文件
	c.File(filePath)
}

// getContentType 根据文件扩展名获取Content-Type
func getContentType(filename string) string {
	ext := ""
	if idx := len(filename) - 1; idx >= 0 && filename[idx] != '.' {
		// 找到最后一个点的位置
		for i := idx; i >= 0; i-- {
			if filename[i] == '.' {
				ext = filename[i:]
				break
			}
		}
	}

	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".bmp":
		return "image/bmp"
	case ".webp":
		return "image/webp"
	case ".svg":
		return "image/svg+xml"
	case ".mp4":
		return "video/mp4"
	case ".avi":
		return "video/x-msvideo"
	case ".mov":
		return "video/quicktime"
	case ".wmv":
		return "video/x-ms-wmv"
	case ".flv":
		return "video/x-flv"
	case ".webm":
		return "video/webm"
	case ".mkv":
		return "video/x-matroska"
	case ".m4v":
		return "video/x-m4v"
	case ".pdf":
		return "application/pdf"
	case ".txt":
		return "text/plain"
	case ".html", ".htm":
		return "text/html"
	case ".css":
		return "text/css"
	case ".js":
		return "application/javascript"
	case ".json":
		return "application/json"
	case ".xml":
		return "application/xml"
	default:
		return "application/octet-stream"
	}
}
