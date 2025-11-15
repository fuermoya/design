package initialize

import (
	"net/http"
	"os"

	"github.com/fuermoya/design/server/docs"
	"github.com/fuermoya/design/server/global"
	"github.com/fuermoya/design/server/middleware"
	"github.com/fuermoya/design/server/router"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type justFilesFilesystem struct {
	fs http.FileSystem
}

func (fs justFilesFilesystem) Open(name string) (http.File, error) {
	f, err := fs.fs.Open(name)
	if err != nil {
		return nil, err
	}

	stat, err := f.Stat()
	if stat.IsDir() {
		return nil, os.ErrPermission
	}

	return f, nil
}

// 初始化总路由
func Routers() *gin.Engine {
	Router := gin.New()
	Router.Use(gin.Recovery())
	if gin.Mode() == gin.DebugMode {
		Router.Use(gin.Logger())
	}

	InstallPlugin(Router) // 安装插件
	systemRouter := router.RouterGroupApp.System
	// 如果想要不使用nginx代理前端网页，可以修改 web/.env.production 下的
	// VUE_APP_BASE_API = /
	// VUE_APP_BASE_PATH = http://localhost
	// 然后执行打包命令 npm run build。在打开下面3行注释
	//Router.Static("/favicon.ico", "./dist/favicon.ico")
	//Router.Static("/assets", "./dist/assets")   // dist里面的静态资源
	//Router.StaticFile("/", "./dist/index.html") // 前端网页入口页面
	// Router.Static("/favicon.ico", "../web/dist/favicon.ico")
	// Router.Static("/assets", "../web/dist/assets")   // dist里面的静态资源
	// Router.StaticFile("/", "../web/dist/index.html") // 前端网页入口页面

	// // 添加 SPA 路由回退支持，所有未匹配的路由都返回 index.html
	// Router.NoRoute(func(c *gin.Context) {
	// 	// 如果是 API 请求，返回 404
	// 	if len(c.Request.URL.Path) >= 4 && c.Request.URL.Path[:4] == "/api" {
	// 		c.JSON(404, gin.H{"error": "API not found"})
	// 		return
	// 	}
	// 	// 其他请求返回 index.html，让前端路由处理
	// 	c.File("../web/dist/index.html")
	// })

	Router.StaticFS(global.CONFIG.Local.Path, justFilesFilesystem{http.Dir(global.CONFIG.Local.StorePath)}) // Router.Use(middleware.LoadTls())  // 如果需要使用https 请打开此中间件 然后前往 core/server.go 将启动模式 更变为 Router.RunTLS("端口","你的cre/pem文件","你的key文件")
	// 跨域，如需跨域可以打开下面的注释
	// Router.Use(middleware.Cors()) // 直接放行全部跨域请求
	// Router.Use(middleware.CorsByRules()) // 按照配置的规则放行跨域请求
	//global.LOG.Info("use middleware cors")
	docs.SwaggerInfo.BasePath = global.CONFIG.System.RouterPrefix
	Router.GET(global.CONFIG.System.RouterPrefix+"/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	global.LOG.Info("register swagger handler")
	// 方便统一添加路由组前缀 多服务器上线使用
	PrivateGroup := Router.Group(global.CONFIG.System.RouterPrefix)
	PrivateGroup.Use(middleware.JWTAuth()).Use(middleware.CasbinHandler())
	PrivateJwtGroup := Router.Group(global.CONFIG.System.RouterPrefix)
	PrivateJwtGroup.Use(middleware.JWTAuth())
	PublicGroup := Router.Group(global.CONFIG.System.RouterPrefix)
	{
		// 健康监测
		PublicGroup.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, "ok")
		})
	}
	{
		systemRouter.InitBaseRouter(PublicGroup) // 注册基础功能路由 不做鉴权
		//systemRouter.InitInitRouter(PublicGroup) // 自动初始化相关
	}

	{
		exampleRouter := router.RouterGroupApp.Example
		exampleRouter.InitExaFileUploadAndDownloadRouter(PrivateGroup) // 文件上传下载功能路由

		systemRouter.InitJwtRouter(PrivateGroup)    // jwt相关路由
		systemRouter.InitUserRouter(PrivateGroup)   // 注册用户路由
		systemRouter.InitCasbinRouter(PrivateGroup) // 权限相关路由
		systemRouter.InitMenuRouter(PrivateGroup)
		systemRouter.InitApiRouter(PrivateGroup, PublicGroup)

		systemRouter.InitAuthorityRouter(PrivateGroup)           // 注册角色路由
		systemRouter.InitSysDictionaryRouter(PrivateGroup)       // 字典管理
		systemRouter.InitSysDictionaryJWTRouter(PrivateJwtGroup) // 字典管理

		systemRouter.InitSysOperationRecordRouter(PrivateGroup)  // 操作记录
		systemRouter.InitSysDictionaryDetailRouter(PrivateGroup) // 字典详情管理
		systemRouter.InitDashboardRouter(PrivateGroup)           // 仪表盘路由

		// 门户网站管理路由
		portalRouter := router.RouterGroupApp.Portal
		portalRouter.InitSysArticleRouter(PrivateGroup)  // 文章管理路由
		portalRouter.InitSysCategoryRouter(PrivateGroup) // 分类管理路由
		portalRouter.InitSysTagRouter(PrivateGroup)      // 标签管理路由
		portalRouter.InitSysThemeRouter(PrivateGroup)    // 主题管理路由
		portalRouter.InitSysMessageRouter(PrivateGroup)  // 留言管理路由
	}

	// 门户网站前台路由（无需认证）
	{
		portalRouter := router.RouterGroupApp.Portal
		portalRouter.InitPortalRouter(PublicGroup)         // 文章前台路由
		portalRouter.InitPortalThemeRouter(PublicGroup)    // 主题前台路由
		portalRouter.InitPortalCategoryRouter(PublicGroup) // 分类前台路由
		portalRouter.InitPortalTagRouter(PublicGroup)      // 标签前台路由
		portalRouter.InitPortalMessageRouter(PublicGroup)  // 留言前台路由
	}

	global.LOG.Info("router register success")
	return Router
}
