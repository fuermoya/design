package utils

import "math"

// RoundToDecimal 小数四舍五入
func RoundToDecimal(num float64, decimalPlaces int) float64 {
	// 将浮点数乘以10的指定位数次方，然后进行四舍五入
	rounded := math.Round(num*math.Pow(10, float64(decimalPlaces))) / math.Pow(10, float64(decimalPlaces))
	return rounded
}
