package response

import "github.com/fuermoya/design/server/config"

type SysConfigResponse struct {
	Config config.Server `json:"config"`
}
