package utils

import (
	"github.com/spf13/cast"
	"strings"
)

func ToArrInt32(input string) []int32 {
	strArr := strings.Split(input, ",")
	var intArr []int32
	for _, str := range strArr {
		intArr = append(intArr, cast.ToInt32(strings.Trim(str, " ")))
	}
	return intArr
}
