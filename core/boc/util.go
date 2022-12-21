package boc

import (
	"math"
	"math/bits"
)

// 计算表示数字num，所需要的位数和字节数
func GetBitsAndBytesNeed(num uint) (int, int) {
	bitsNeed := bits.Len(uint(num))
	bytesNeed := int(math.Max(math.Ceil(float64(bitsNeed)/8), 1))
	return bitsNeed, bytesNeed
}

// 判断slice中是否有某个值
func contains[T comparable](slice []T, item T) bool {
	for i := range slice {
		if slice[i] == item {
			return true
		}
	}
	return false
}
