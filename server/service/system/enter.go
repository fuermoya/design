package system

type ServiceGroup struct {
	JwtService
	MenuService
	UserService
	CasbinService
	InitDBService
	AuthorityService
	DictionaryService
	OperationRecordService
	DictionaryDetailService
	ApiService
}
