package boc

// 提供一种map结构
// key是一定长度的bit位（可以不是8的整数倍）
// value是
// keys和values是
type HashMap[T any] struct {
	keySize int
	keys    []BitString
	values  []T
}
