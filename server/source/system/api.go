package system

import (
	"context"

	"github.com/fuermoya/design/server/global"
	sysModel "github.com/fuermoya/design/server/model/system"
	"github.com/fuermoya/design/server/service/system"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type initApi struct{}

const initOrderApi = system.InitOrderSystem + 1

// auto run
func init() {
	system.RegisterInit(initOrderApi, &initApi{})
}

func (i initApi) InitializerName() string {
	return sysModel.SysApi{}.TableName()
}

func (i *initApi) MigrateTable(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	return ctx, db.AutoMigrate(&sysModel.SysApi{})
}

func (i *initApi) TableCreated(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	return db.Migrator().HasTable(&sysModel.SysApi{})
}

func (i *initApi) InitializeData(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	entities := []sysModel.SysApi{
		{ApiGroup: "jwt", Method: "POST", Path: "/jwt/jsonInBlacklist", Description: "jwt加入黑名单(退出，必选)"},

		{ApiGroup: "系统字典详情", Method: "PUT", Path: "/sysDictionaryDetail/updateSysDictionaryDetail", Description: "更新字典内容"},
		{ApiGroup: "系统字典详情", Method: "POST", Path: "/sysDictionaryDetail/createSysDictionaryDetail", Description: "新增字典内容"},
		{ApiGroup: "系统字典详情", Method: "POST", Path: "/sysDictionaryDetail/createBatchSysDictionaryDetail", Description: "批量新增字典内容"},
		{ApiGroup: "系统字典详情", Method: "DELETE", Path: "/sysDictionaryDetail/deleteSysDictionaryDetail", Description: "删除字典内容"},
		{ApiGroup: "系统字典详情", Method: "GET", Path: "/sysDictionaryDetail/findSysDictionaryDetail", Description: "根据ID获取字典内容"},
		{ApiGroup: "系统字典详情", Method: "GET", Path: "/sysDictionaryDetail/getSysDictionaryDetailList", Description: "获取字典内容列表"},

		{ApiGroup: "系统用户", Method: "DELETE", Path: "/user/deleteUser", Description: "删除用户"},
		{ApiGroup: "系统用户", Method: "POST", Path: "/user/admin_register", Description: "用户注册"},
		{ApiGroup: "系统用户", Method: "POST", Path: "/user/getUserList", Description: "获取用户列表"},
		{ApiGroup: "系统用户", Method: "PUT", Path: "/user/setUserInfo", Description: "设置用户信息"},
		{ApiGroup: "系统用户", Method: "PUT", Path: "/user/setSelfInfo", Description: "设置自身信息(必选)"},
		{ApiGroup: "系统用户", Method: "GET", Path: "/user/getUserInfo", Description: "获取自身信息(必选)"},
		{ApiGroup: "系统用户", Method: "POST", Path: "/user/setUserAuthorities", Description: "设置权限组"},
		{ApiGroup: "系统用户", Method: "POST", Path: "/user/changePassword", Description: "修改密码（建议选择)"},
		{ApiGroup: "系统用户", Method: "POST", Path: "/user/setUserAuthority", Description: "修改用户角色(必选)"},
		{ApiGroup: "系统用户", Method: "POST", Path: "/user/resetPassword", Description: "重置用户密码"},

		{ApiGroup: "角色", Method: "POST", Path: "/authority/copyAuthority", Description: "拷贝角色"},
		{ApiGroup: "角色", Method: "POST", Path: "/authority/createAuthority", Description: "创建角色"},
		{ApiGroup: "角色", Method: "POST", Path: "/authority/deleteAuthority", Description: "删除角色"},
		{ApiGroup: "角色", Method: "PUT", Path: "/authority/updateAuthority", Description: "更新角色信息"},
		{ApiGroup: "角色", Method: "POST", Path: "/authority/getAuthorityList", Description: "获取角色列表"},

		{ApiGroup: "操作记录", Method: "POST", Path: "/sysOperationRecord/createSysOperationRecord", Description: "新增操作记录"},
		{ApiGroup: "操作记录", Method: "GET", Path: "/sysOperationRecord/findSysOperationRecord", Description: "根据ID获取操作记录"},
		{ApiGroup: "操作记录", Method: "GET", Path: "/sysOperationRecord/getSysOperationRecordList", Description: "获取操作记录列表"},
		{ApiGroup: "操作记录", Method: "DELETE", Path: "/sysOperationRecord/deleteSysOperationRecord", Description: "删除操作记录"},
		{ApiGroup: "操作记录", Method: "DELETE", Path: "/sysOperationRecord/deleteSysOperationRecordByIds", Description: "批量删除操作历史"},

		{ApiGroup: "菜单", Method: "POST", Path: "/menu/getMenu", Description: "获取菜单树(必选)"},
		{ApiGroup: "菜单", Method: "POST", Path: "/menu/getBaseMenuTree", Description: "获取用户动态路由"},
		{ApiGroup: "菜单", Method: "POST", Path: "/menu/getMenuAuthority", Description: "获取指定角色menu"},
		{ApiGroup: "菜单", Method: "POST", Path: "/menu/addMenuAuthority", Description: "增加menu和角色关联关系"},

		{ApiGroup: "casbin", Method: "POST", Path: "/casbin/updateCasbin", Description: "更改角色api权限"},
		{ApiGroup: "casbin", Method: "POST", Path: "/casbin/getPolicyPathByAuthorityId", Description: "获取权限列表"},

		{ApiGroup: "api", Method: "POST", Path: "/api/getAllApis", Description: "获取所有api"},

		{ApiGroup: "文件上传下载", Method: "GET", Path: "/fileUploadAndDownload/getFileList", Description: "获取文件列表"},
		{ApiGroup: "文件上传下载", Method: "POST", Path: "/fileUploadAndDownload/upload", Description: "上传文件"},
		{ApiGroup: "文件上传下载", Method: "DELETE", Path: "/fileUploadAndDownload/deleteFile", Description: "删除文件"},
		{ApiGroup: "文件上传下载", Method: "GET", Path: "/fileUploadAndDownload/download", Description: "下载文件"},
		{ApiGroup: "文件上传下载", Method: "GET", Path: "/fileUploadAndDownload/preview", Description: "预览文件"},

		{ApiGroup: "文章管理", Method: "GET", Path: "/sysArticle/getSysArticleList", Description: "获取文章列表"},
		{ApiGroup: "文章管理", Method: "POST", Path: "/sysArticle/createSysArticle", Description: "创建文章"},
		{ApiGroup: "文章管理", Method: "PUT", Path: "/sysArticle/updateSysArticle", Description: "更新文章"},
		{ApiGroup: "文章管理", Method: "DELETE", Path: "/sysArticle/deleteSysArticle", Description: "删除文章"},
		{ApiGroup: "文章管理", Method: "DELETE", Path: "/sysArticle/deleteSysArticleByIds", Description: "批量删除文章"},
		{ApiGroup: "文章管理", Method: "GET", Path: "/sysArticle/findSysArticle", Description: "根据ID获取文章"},

		{ApiGroup: "分类管理", Method: "GET", Path: "/sysCategory/getSysCategoryList", Description: "获取分类列表"},
		{ApiGroup: "分类管理", Method: "POST", Path: "/sysCategory/createSysCategory", Description: "创建分类"},
		{ApiGroup: "分类管理", Method: "PUT", Path: "/sysCategory/updateSysCategory", Description: "更新分类"},
		{ApiGroup: "分类管理", Method: "DELETE", Path: "/sysCategory/deleteSysCategory", Description: "删除分类"},
		{ApiGroup: "分类管理", Method: "DELETE", Path: "/sysCategory/deleteSysCategoryByIds", Description: "批量删除分类"},
		{ApiGroup: "分类管理", Method: "GET", Path: "/sysCategory/findSysCategory", Description: "根据ID获取分类"},

		{ApiGroup: "标签管理", Method: "GET", Path: "/sysTag/getSysTagList", Description: "获取标签列表"},
		{ApiGroup: "标签管理", Method: "POST", Path: "/sysTag/createSysTag", Description: "创建标签"},
		{ApiGroup: "标签管理", Method: "PUT", Path: "/sysTag/updateSysTag", Description: "更新标签"},
		{ApiGroup: "标签管理", Method: "DELETE", Path: "/sysTag/deleteSysTag", Description: "删除标签"},
		{ApiGroup: "标签管理", Method: "DELETE", Path: "/sysTag/deleteSysTagByIds", Description: "批量删除标签"},
		{ApiGroup: "标签管理", Method: "GET", Path: "/sysTag/findSysTag", Description: "根据ID获取标签"},

		{ApiGroup: "主题管理", Method: "GET", Path: "/sysTheme/getSysThemeList", Description: "获取主题列表"},
		{ApiGroup: "主题管理", Method: "POST", Path: "/sysTheme/createSysTheme", Description: "创建主题"},
		{ApiGroup: "主题管理", Method: "PUT", Path: "/sysTheme/updateSysTheme", Description: "更新主题"},
		{ApiGroup: "主题管理", Method: "DELETE", Path: "/sysTheme/deleteSysTheme", Description: "删除主题"},
		{ApiGroup: "主题管理", Method: "DELETE", Path: "/sysTheme/deleteSysThemeByIds", Description: "批量删除主题"},
		{ApiGroup: "主题管理", Method: "GET", Path: "/sysTheme/findSysTheme", Description: "根据ID获取主题"},
		{ApiGroup: "主题管理", Method: "POST", Path: "/sysTheme/activateTheme", Description: "激活主题"},
		{ApiGroup: "主题管理", Method: "GET", Path: "/portal/theme", Description: "获取当前激活的主题"},

		{ApiGroup: "留言管理", Method: "GET", Path: "/sysMessage/getMessageList", Description: "获取留言列表"},
		{ApiGroup: "留言管理", Method: "GET", Path: "/sysMessage/getMessageById", Description: "根据ID获取留言"},
		{ApiGroup: "留言管理", Method: "GET", Path: "/sysMessage/getMessageStats", Description: "获取留言统计"},
		{ApiGroup: "留言管理", Method: "PUT", Path: "/sysMessage/updateMessage", Description: "更新留言"},
		{ApiGroup: "留言管理", Method: "DELETE", Path: "/sysMessage/deleteMessage", Description: "删除留言"},
		{ApiGroup: "留言管理", Method: "POST", Path: "/sysMessage/replyMessage", Description: "回复留言"},
		{ApiGroup: "留言管理", Method: "POST", Path: "/sysMessage/markAsRead", Description: "标记为已读"},

		{ApiGroup: "仪表盘", Method: "GET", Path: "/dashboard/stats", Description: "获取统计数据"},
		{ApiGroup: "仪表盘", Method: "GET", Path: "/dashboard/status", Description: "获取系统状态"},
		{ApiGroup: "仪表盘", Method: "GET", Path: "/dashboard/charts", Description: "获取图表数据"},
		{ApiGroup: "仪表盘", Method: "GET", Path: "/dashboard/activities", Description: "获取最近活动"},
		{ApiGroup: "仪表盘", Method: "GET", Path: "/dashboard/info", Description: "获取系统信息"},
	}
	apis := make([]sysModel.SysApi, 0)
	for _, v := range entities {
		var count int64
		err := db.Raw("select count(1) from sys_apis where path = ? AND method = ?", v.Path, v.Method).Scan(&count).Error
		if err != nil {
			global.LOG.Error("校验API数据错误", zap.Error(err))
		}

		if count <= 0 {
			apis = append(apis, v)
		}
	}

	if len(apis) > 0 {
		if err := db.Create(&apis).Error; err != nil {
			return ctx, errors.Wrap(err, sysModel.SysApi{}.TableName()+"表数据初始化失败!")
		}
	}
	next := context.WithValue(ctx, i.InitializerName(), entities)
	return next, nil
}

func (i *initApi) DataInserted(ctx context.Context) bool {
	return false
}
