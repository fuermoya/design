package utils

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/mozillazg/go-pinyin"
)

// CreateSignMd5 创建签名
func CreateSignMd5(params map[string]interface{}, str ...string) string {
	// 对参数进行ASCII排序
	keys := make([]string, 0, len(params))
	for key := range params {
		keys = append(keys, key)
	}

	sort.Strings(keys)
	// 构建签名字符串
	var signStrings []string
	for _, key := range keys {
		signStrings = append(signStrings, fmt.Sprintf("%s=%s", key, Strval(params[key])))
	}

	signStrings = append(signStrings, str...)
	sign := strings.Join(signStrings, "&")
	return MD5V([]byte(sign))
}

// CheckGameSignMd5 校验签名
func CheckGameSignMd5(params map[string]interface{}, chunkMd5 string) (check bool) {
	//signMd5 := CreateSignMd5(params, global.CONFIG.XYSGJWT.Issuer)
	//if signMd5 == chunkMd5 {
	//	return true
	//} else {
	//	return false
	//}
	return false
}

func FirstSpell(str string) string {
	// 创建拼音实例
	p := pinyin.NewArgs()

	// 将字符串转换为字符数组
	runes := []rune(str)

	// 存储结果的切片
	result := make([]string, len(runes))

	for i, r := range runes {
		// 判断字符是否为中文
		if isChinese(r) {
			// 将中文字符转换为拼音首字母
			pinyinSlice := pinyin.Pinyin(string(r), p)
			initial := pinyinSlice[0][0][0:1]
			result[i] = initial
		} else {
			// 非中文字符直接保留原字符
			result[i] = string(r)
		}
	}

	// 返回结果字符串
	return strings.Join(result, "")
}

func isChinese(r rune) bool {
	return r >= '\u4e00' && r <= '\u9fff'
}

func GenerateGiftCodes(length, count int, prefixStr string) ([]string, error) {
	giftCodes := make([]string, count)
	for i := 0; i < count; i++ {
		randomBytes := make([]byte, length)
		_, err := rand.Read(randomBytes)
		if err != nil {
			return nil, err
		}
		giftCode := base64.URLEncoding.EncodeToString(randomBytes)
		giftCode = giftCode[:length]

		giftCodes[i] = prefixStr + giftCode
	}
	return giftCodes, nil
}

// Strval 获取变量的字符串值
// 浮点型 3.0将会转换成字符串3, "3"
// 非数值或字符类型的变量将会被转换成JSON格式字符串
func Strval(value interface{}) string {
	var key string
	if value == nil {
		return key
	}

	switch value.(type) {
	case float64:
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		key = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		key = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		key = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		key = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		key = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
	case string:
		key = value.(string)
	case []byte:
		key = string(value.([]byte))
	default:
		newValue, _ := json.Marshal(value)
		key = string(newValue)
	}

	return key
}
