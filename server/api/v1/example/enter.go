package example

import "github.com/fuermoya/design/server/service"

type ApiGroup struct {
	FileUploadAndDownloadApi
}

var (
	fileUploadAndDownloadService = service.ServiceGroupApp.ExampleServiceGroup.FileUploadAndDownloadService
)
