package system

import (
	"time"

	"github.com/alphadose/haxmap"

	"go.uber.org/zap"

	"github.com/fuermoya/design/server/global"
	"github.com/fuermoya/design/server/model/system"
	"github.com/fuermoya/design/server/utils"
)

type JwtService struct{}
type jwtStruct struct {
	str        string
	expiryTime int64
}

var jwtCache = haxmap.New[string, jwtStruct]()

//@function: JsonInBlacklist
//@description: 拉黑jwt
//@param: jwtList model.JwtBlacklist
//@return: err error

func (jwtService *JwtService) JsonInBlacklist(jwtList system.JwtBlacklist) (err error) {
	err = global.DB.Create(&jwtList).Error
	if err != nil {
		return
	}
	global.BlackCache.SetDefault(jwtList.Jwt, struct{}{})
	return
}

//@function: IsBlacklist
//@description: 判断JWT是否在黑名单内部
//@param: jwt string
//@return: bool

func (jwtService *JwtService) IsBlacklist(jwt string) bool {
	_, ok := global.BlackCache.Get(jwt)
	return ok
	// err := global.DB.Where("jwt = ?", jwt).First(&system.JwtBlacklist{}).Error
	// isNotFound := errors.Is(err, gorm.ErrRecordNotFound)
	// return !isNotFound
}

//@function: GetCacheJWT
//@description: 从缓存取jwt
//@param: userName string
//@return: redisJWT string, ok bool

func (jwtService *JwtService) GetCacheJWT(userName string) (redisJWT string, ok bool) {
	value, ok := jwtCache.Get(userName)
	if !ok {
		return "", false
	}
	unix := time.Now().Unix()
	if value.expiryTime < unix {
		jwtCache.Del(userName)
		return "", false
	}
	return value.str, true
}

//@function: SetCacheJWT
//@description: jwt存入缓存并设置过期时间
//@param: jwt string, userName string
//@return: ok bool

func (jwtService *JwtService) SetCacheJWT(jwt string, userName string) (err error) {
	// 此处过期时间等于jwt过期时间
	dr, err := utils.ParseDuration(global.CONFIG.JWT.ExpiresTime)
	if err != nil {
		return err
	}

	jwtCache.Set(userName, jwtStruct{
		str:        jwt,
		expiryTime: time.Now().Add(dr).Unix(),
	})
	return err
}

func LoadAll() {
	var data []string
	err := global.DB.Model(&system.JwtBlacklist{}).Select("jwt").Find(&data).Error
	if err != nil {
		global.LOG.Error("加载数据库jwt黑名单失败!", zap.Error(err))
		return
	}
	for i := 0; i < len(data); i++ {
		global.BlackCache.SetDefault(data[i], struct{}{})
	} // jwt黑名单 加入 BlackCache 中
}
