package boc

type BitString struct {
	buf []byte // 容纳这些bit位
	cap int    // 总容量
	len int    // 实际使用的bit位长
}

func NewBitString(bitLen int) BitString {
	return BitString{
		buf: make([]byte, int(math.Ceil(float64(bitLen)/float64(8)))),
		cap: bitLen,
		len: 0,
	}
}
