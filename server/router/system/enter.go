package system

type RouterGroup struct {
	JwtRouter
	BaseRouter
	UserRouter
	CasbinRouter
	AuthorityRouter
	DictionaryRouter
	OperationRecordRouter
	DictionaryDetailRouter
	MenuRouter
	ApiRouter
	DashboardRouter
}
