package example

import (
	v1 "github.com/fuermoya/design/server/api/v1"
	"github.com/gin-gonic/gin"
)

type ExaFileUploadAndDownloadRouter struct{}

func (e *ExaFileUploadAndDownloadRouter) InitExaFileUploadAndDownloadRouter(Router *gin.RouterGroup) {
	fileUploadAndDownloadRouter := Router.Group("fileUploadAndDownload")
	fileUploadAndDownloadApi := v1.ApiGroupApp.ExampleApiGroup.FileUploadAndDownloadApi
	{
		fileUploadAndDownloadRouter.POST("upload", fileUploadAndDownloadApi.UploadFile)       // 上传文件
		fileUploadAndDownloadRouter.DELETE("deleteFile", fileUploadAndDownloadApi.DeleteFile) // 删除文件
		fileUploadAndDownloadRouter.GET("getFileList", fileUploadAndDownloadApi.GetFileList)  // 获取上传文件列表
		fileUploadAndDownloadRouter.GET("download", fileUploadAndDownloadApi.DownloadFile)    // 下载文件
		fileUploadAndDownloadRouter.GET("preview", fileUploadAndDownloadApi.PreviewFile)      // 预览文件
	}
}
