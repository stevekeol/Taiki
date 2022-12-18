package boc

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"math/big"
	"strings"
)

// BitString提供位操作
type BitString struct {
	buf    []byte // 容纳这些bit位
	cap    int    // 总容量(一定是小于等于字节数*8的，不一定占满，但字节的分配单位为8位)
	len    int    // 实际使用的bit位长(同时也是指明下一个将要写入的Bit的位置)
	cursor int    // 位串中逐个bit位读取的游标位置（即下一个可以读取的位置）
}

// 创建一个新的待操作位串
func NewBitString(bitLen int) BitString {
	return BitString{
		// buf: make([]byte, int(math.Ceil(float64(bitLen)/float64(8)))),
		buf:    make([]byte, (bitLen+7)/8),
		cap:    bitLen,
		len:    0,
		cursor: 0,
	}
}

//////////////////////////////////////////////////
///											   ///
///                   Write                    ///
///                                            ///
//////////////////////////////////////////////////

// 在当前位串中写入另一个位串
// 在当前位串的剩余可写空间不足时，会自动扩容
func (bs *BitString) WriteAppend(s BitString) error {
	bitsNeed := s.len - bs.bitsRemainingForWrite()
	if bitsNeed > 0 {
		bs.grow(bitsNeed)
	}
	err := bs.WriteBitString(s)
	if err != nil {
		return err
	}
	return nil
}

// 在当前位串的剩余可写空间内写入另一个位串
// 在当前位串剩余空间不足的情况下，则报错
func (bs *BitString) WriteBitString(s BitString) error {
	if bs.bitsRemainingForWrite() < s.len {
		return errors.New("not enough bits to write in bitstring")
	}

	s.cursor = 0
	for i := 0; i < s.len; i++ {
		bit, err := s.ReadBit()
		if err != nil {
			return err
		}
		err = bs.WriteBit(bit)
		if err != nil {
			return err
		}
	}
	return nil
}

// 一次性向BitString写入多个字节（内部是借助WriteByte逐个字节写入的）
func (bs *BitString) WriteBytes(data []byte) error {
	for _, item := range data {
		err := bs.WriteByte(item)
		if err != nil {
			return err
		}
	}
	return nil
}

// 向BitString中写入单个字节（内部是借助WriteUint逐个写入8个bit位）
func (bs *BitString) WriteByte(value byte) error {
	err := bs.WriteUint(uint64(value), 8)
	if err != nil {
		return err
	}
	return nil
}

// 将uint64类型的数值value，以bitLen长度的位串形式写入BitString
// 1.对于bitlen长度足够表示该value，很明显前面会用0补足
// 2.对于bitlen长度不够表示该value，则会舍弃高位的bits
// TODO: 当bitLen设定为uint会发生什么？发现uint, int32, int64等在测试时都会报错！只有int不报错
func (bs *BitString) WriteUint(value uint64, bitLen int) error {
	for i := bitLen - 1; i >= 0; i-- {
		err := bs.WriteBit(((value >> i) & 1) > 0)
		if err != nil {
			return err
		}
	}
	return nil
}

// func (bs *BitString) WriteBigUint(value *big.Int, bitLen int) error {}

// 将int64类型的数值value，以bitLen长度的位串形式写入BitString
// 感悟：bitLen决定了用多少位来表示该数值
// 1. 对于bitLen足够时，前面会用0(正数)或1(负数)补足
// 2. 对于bitLen不够时，在正确形式的负数(取反加1)的基础上会直接截断高位
func (bs *BitString) WriteInt(value int64, bitLen int) error {
	if bitLen == 1 {
		if value == -1 {
			err := bs.WriteBit(true)
			if err != nil {
				return err
			}
		}
		if value == 0 {
			err := bs.WriteBit(false)
			if err != nil {
				return err
			}
		}
	} else {
		if value < 0 {
			err := bs.WriteBit(true)
			if err != nil {
				return err
			}
			// 对于负数而言，bitLen决定了用多少个bit来表示
			// 如Write(-17, 8)写入的位串是：17是00010001，-17是(取反加1)11101111
			// 如Wirte(-17, 10)写入的位串是：17是00 00010001，-17是(取反加1)11 11101111
			err = bs.WriteUint(uint64(1<<(bitLen-1)+value), bitLen-1)
			if err != nil {
				return err
			}
		} else {
			err := bs.WriteBit(false)
			if err != nil {
				return err
			}
			err = bs.WriteUint(uint64(value), bitLen-1)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// func (bs *BitString) WriteBigInt() {}

// 一次性写入多个bool类型的bit值
func (bs *BitString) WriteBits(vals []bool) error {
	for _, val := range vals {
		err := bs.WriteBit(val)
		if err != nil {
			return err
		}
	}
	return nil
}

// 在下一个可写位置（即索引len指明的位置），写入0/1状态（以布尔类型表示）的一个bit值
// 因此常用于向BitString中逐个数据位的写入
func (bs *BitString) WriteBit(on bool) error {
	if on {
		err := bs.On(bs.len)
		if err != nil {
			return err
		}
	} else {
		err := bs.Off(bs.len)
		if err != nil {
			return err
		}
	}
	bs.len++
	return nil
}

// 将BitString中从左往右位置为index上的Bit置1（类似开关的On状态）
// Notice: 是否应归为helper一类
func (bs *BitString) On(index int) error {
	err := bs.checkValid(index)
	if err != nil {
		return err
	}
	bs.buf[index/8] |= 1 << (7 - (index % 8))
	return nil
}

// 将BitString中从左往右位置为index上的Bit清0（类似开关的off状态）
// Notice: 是否应归为helper一类
func (bs *BitString) Off(index int) error {
	err := bs.checkValid(index)
	if err != nil {
		return err
	}
	bs.buf[index/8] &= ^(1 << (7 - (index % 8)))
	return nil
}

//////////////////////////////////////////////////
///											   ///
///                   Read                     ///
///                                            ///
//////////////////////////////////////////////////

// 一次性读取多个字节
// 其内部实际上依赖ReatByte()来逐个字节的读取再一并返回
// 虽然每次可能足够读取（最后一次除外），但为了避免无效的读取操作，因此需要一次性的判断是否足够读取
func (bs *BitString) ReadBytes(size int) ([]byte, error) {
	if bs.bitsRemainingForRead() < size*8 {
		return nil, errors.New("not enough bits to read in bitstring")
	}
	res := make([]byte, size)
	for i := 0; i < size; i++ {
		value, err := bs.ReadUint(8)
		if err != nil {
			return nil, err
		}
		res[i] = byte(value)
	}
	return res, nil
}

// 读取一个字节
// 注：所有的读取操作，开始读取的位置都是游标cursor当前指明的位置
func (bs *BitString) ReadByte() (byte, error) {
	res, err := bs.ReadUint(8)
	if err != nil {
		// 注: byte的零值是0，[]byte的零值是nil
		return 0, nil
	}
	return byte(res), nil
}

// 读出长度为bitLen的位串，并将其转为uint64类型返回
// 读取开始的位置是cursor指定的位置
// 具体做法：构建一个uint64值，逐个bit位取出的同时将其设置在该uint64的对应位上
// 1. 当bitLen不足64位时，该uint64数值的高位自然的保留着0
// 2. 当bitLen大于64位时，该情况不允许！
func (bs *BitString) ReadUint(bitLen int) (uint64, error) {
	// 判断是否要取的bit位数超过最大允许位数
	if bitLen > 64 {
		return 0, errors.New("too much bits to read beyond uint64")
	}
	// 判断剩下可读的bit位数是否足够读取
	if bs.bitsRemainingForRead() < bitLen {
		return 0, errors.New("not enough bits to read in bitstring")
	}
	if bitLen == 0 {
		return 0, nil
	}
	var res uint64 = 0
	for i := bitLen - 1; i >= 0; i-- {
		if bs.readBit() {
			res |= 1 << i
		}
	}
	return res, nil
}

// 读取长度为bitLen的位串，并将其转为big.Int类型返回
// 注：bitLen的长度可以超过64位(big.Int在形式上类似字符串，因此足以表达任意数值)
// 具体做法：逐个位取出，并以字符串形式缀连，然后基于此构建big.Int
// TODO：使用场景
func (bs *BitString) ReadBigUint(bitLen int) (*big.Int, error) {
	// 判断剩下可读的bit位数是否足够读取
	if bs.bitsRemainingForRead() < bitLen {
		return nil, errors.New("not enough bits to read in bitstring")
	}
	if bitLen == 0 {
		return new(big.Int), nil
	}
	str := ""
	for i := 0; i < bitLen; i++ {
		if bs.readBit() {
			str += "1"
		} else {
			str += "0"
		}
	}

	bigNum := new(big.Int)
	// big.Int是个形式上的字符串，第二个参数是解析的进制
	bigNum.SetString(str, 2)
	return bigNum, nil
}

// 读出长度为bitLen的位串，并将其视为int64类型的值返回
// 重点是要判断该bitLen的第一个bit位是1还是0，以此来判定int64值的构建流程
func (bs *BitString) ReadInt(bitLen int) (int64, error) {
	// 判断是否要取的bit位数超过最大允许位数
	if bitLen > 64 {
		return 0, errors.New("too much bits to read beyond uint64")
	}
	// 判断剩下可读的bit位数是否足够读取
	if bs.bitsRemainingForRead() < bitLen {
		return 0, errors.New("not enough bits to read in bitstring")
	}
	if bitLen == 0 {
		return 0, nil
	}
	// 当只取出一个bit，且要将其视为int64的话：
	// 1. 若该bit为1，则为-1
	// 2. 若该bit为0，则为0
	if bitLen == 1 {
		if bs.readBit() {
			return -1, nil
		}
		return 0, nil
	}
	// 若为负数：
	// 1. 先取出除去符号位之外剩余bit位表示的正数-positiveValue(uint64类型)
	// 2. 再计算剩余bit位全满且刚好溢出的值-overflowValue
	// 3. positive 减去 overflowValue 的值，就为这个位串所表示的负数
	if bs.readBit() {
		positiveValue, err := bs.ReadUint(bitLen - 1)
		if err != nil {
			return 0, err
		}
		// 数值1默认是int类型，因此需要刻意转换成uint64，才能和同为uint64类型的positiveValue相减
		overflowValue := uint64(1) << (bitLen - 1)
		return int64(positiveValue - overflowValue), nil
	}
	// 若为正数：
	// 直接取出就好，不过需要将取出的uint64类型转换成int64类型
	positiveValue, err := bs.ReadUint(bitLen - 1)
	if err != nil {
		return 0, err
	}
	return int64(positiveValue), nil
}

// 读取长度为bitLen的位串，并将其转为big.Int类型返回
// 注：bitLen的长度可以超过64位(big.Int在形式上类似字符串，因此足以表达任意数值)
// 具体做法：逐个位取出，并以字符串形式缀连，然后基于此构建big.Int
// TODO：使用场景
// ReadBigInt与ReadBigUint的区别：前者会将读取的内容视为有符号大数
func (bs *BitString) ReadBigInt(bitLen int) (*big.Int, error) {
	// 判断剩下可读的bit位数是否足够读取
	if bs.bitsRemainingForRead() < bitLen {
		return nil, errors.New("not enough bits to read in bitstring")
	}
	if bitLen == 0 {
		return new(big.Int), nil
	}
	// 当只取出一个bit，且要将其视为int64的话：
	// 1. 若该bit为1，则为-1，不过需要转换成big.Int类型
	// 2. 若该bit为0，则为0，不过需要转换成big.Int类型
	if bitLen == 1 {
		if bs.readBit() {
			return big.NewInt(-1), nil
		}
		return big.NewInt(0), nil
	}
	// 若为负大数：
	// 1. 先取出除去符号位之外剩余bit位表示的正大数-positiveBigValue
	// 2.
	// 3.
	if bs.readBit() {
		positiveBigValue, err := bs.ReadBigUint(bitLen - 1)
		if err != nil {
			return nil, err
		}
		// TODO
		b := big.NewInt(2)
		nb := b.Exp(b, big.NewInt(int64(bitLen-1)), nil)
		return positiveBigValue.Sub(positiveBigValue, nb), nil
	}
	// 若为正大数：
	// 直接借助ReadBigUint取出就好，错误处理也会自然的在其中处理并返回
	return bs.ReadBigUint(bitLen - 1)
}

// 一次性读取n个bit位
// 开始读取的位置，是当前cursor指定的位置
// 若当前BitString的可读取bit位不足n，则直接返回空的BitString对象
func (bs *BitString) ReadBits(n int) (BitString, error) {
	bitString := NewBitString(n)
	for i := 0; i < n; i++ {
		bit, err := bs.ReadBit()
		if err != nil {
			return BitString{}, err
		}
		err = bitString.WriteBit(bit)
		if err != nil {
			return BitString{}, err
		}
	}
	return bitString, nil
}

// 逐个自动取出游标cursor指定的位置的bit（用布尔量表示1或0）
func (bs *BitString) ReadBit() (bool, error) {
	if bs.bitsRemainingForRead() < 1 {
		return false, errors.New("not enough bits to read in bitstring")
	}

	return bs.readBit(), nil
}

// 调整下一次的读取位置
// 1. 当n为正数，则将读取游标cursor右移n个位置(注意别越界)
// 2. 当n为负数，则将读取游标cursor左移n个位置(注意别越界)
func (bs *BitString) ReadSkip(n int) error {
	if n > 0 && bs.bitsRemainingForRead() < 1 || n < 0 && (bs.cursor+n) < 0 {
		return errors.New("not enough bits to skip in bitstring")
	}

	bs.cursor += n
	return nil
}

// 在ReadUint的基础上再将读取游标后退至读取之前的位置
func (bs *BitString) ReadUintAndBackward(bitLen int) (uint64, error) {
	res, err := bs.ReadUint(bitLen)
	if err != nil {
		return 0, err
	}
	bs.cursor -= bitLen
	return res, nil
}

// 读取游标重置
func (bs *BitString) ReadReset() {
	bs.cursor = 0
}

//////////////////////////////////////////////////
///											   ///
///                   Utils                    ///
///                                            ///
//////////////////////////////////////////////////

// 获取当前可写游标的位置（即目前已经使用的bit位的长度）
// TODO：思考是否需要Public
func (bs *BitString) GetWriteCursor() int {
	return bs.len
}

// 获取当前可读游标的位置
// TODO：思考是否需要Public
func (bs *BitString) GetReadCursor() int {
	return bs.cursor
}

// 获取当前可供写入bits的字节slice
// TODO：思考是否需要Public
func (bs *BitString) GetBuffer() []byte {
	return bs.buf
}

// 复制一份该BitSting结构体对象
func (bs *BitString) Copy() BitString {
	buf := make([]byte, len(bs.buf))
	copy(buf, bs.buf)
	return BitString{
		buf: buf,
		cap: bs.cap,
		len: bs.len,
	}
}

// 将该BitString所有的bit位打印出来
func (bs *BitString) FullBits() string {
	// TODO：思考下能否不要strings包
	str := strings.Builder{}
	for _, item := range bs.buf {
		str.WriteString(fmt.Sprintf("%08b", item))
	}
	return str.String()
}

// 将BitString已经使用的bit位打印出来
// TODO：思考是否需要以8位为单位，间隔打印
func (bs *BitString) UsedBits() string {
	buf := strings.Builder{}
	for i, item := range bs.buf {
		if (i+1)*8 <= bs.len {
			buf.WriteString(fmt.Sprintf("%08b", item))
		} else if (i)*8 > bs.len {
			break
		} else {
			str := fmt.Sprintf("%08b", item)
			for j := 0; buf.Len() < bs.len; j++ {
				buf.WriteByte(str[j])
			}
		}

	}
	return buf.String()
}

// TODO
func (bs *BitString) ToFiftHex() string {
	if bs.len%4 == 0 {
		str := strings.ToUpper(hex.EncodeToString(bs.buf[0 : (bs.len+7)/8]))
		if bs.len%8 == 0 {
			return str
		}
		return str[0 : len(str)-1]
	}
	temp := bs.Copy()
	temp.WriteBit(true)
	for temp.len%4 != 0 {
		temp.WriteBit(false)
	}
	hex := temp.ToFiftHex()
	return hex + "_"
}

// 将当前位串以TL的方式序列化写入到给定的Cell中
// TODO：此处的tag应该是有作用的吧？
func (bs BitString) MarshalTL(cell *Cell, tag string) error {
	err := cell.bits.WriteBitString(bs)
	if err != nil {
		return err
	}
	return nil
}

// 从给定Cell中以TL的方式反序列化读出(tag指定的长度的)位串来给当前位串赋值
func (bs *BitString) UnmarshalTL(cell *Cell, tag string) error {
	bitLen, err := decodeBitStringTag(tag)
	if err != nil {
		return err
	}
	s, err := cell.bits.ReadBits(bitLen)
	if err != nil {
		return err
	}
	// Temp:留意更改值的方式
	*bs = s
	return nil
}

// 暂未移植
func (s *BitString) SetTopUppedArray(arr []byte, fulfilledBytes bool) error {
	s.cap = len(arr) * 8
	s.buf = make([]byte, len(arr))
	copy(s.buf, arr)
	s.len = s.cap

	if fulfilledBytes || s.cap == 0 {
		return nil
	}

	foundEndBit := false

	for i := 0; i < 7; i++ {
		s.len--
		if s.getBit(s.len) {
			foundEndBit = true
			err := s.Off(s.len) //todo: check
			if err != nil {
				return err
			}
			break
		}
	}
	if !foundEndBit {
		return errors.New("incorrect topUppedArray")
	}
	return nil
}

// 暂未移植
func (s *BitString) GetTopUppedArray() ([]byte, error) {
	ret := s.Copy()
	tu := int(math.Ceil(float64(ret.GetWriteCursor())/8))*8 - ret.GetWriteCursor()
	fmt.Println(tu)

	if tu > 0 {
		tu = tu - 1
		err := ret.WriteBit(true)
		if err != nil {
			return nil, err
		}
		for tu > 0 {
			tu = tu - 1
			err := ret.WriteBit(false)
			if err != nil {
				return nil, err
			}
		}
	}
	// Notice: 和-8相与（注意-8的补码表示形式），就是将最后三个bit位清零
	// Question: 这和`(ret.len + 7) / 8`有什么区别？
	ret.buf = ret.buf[0 : ((ret.len+7)&-8)/8]
	return ret.buf, nil
}

//////////////////////////////////////////////////
///											   ///
///                   Helper                   ///
///                                            ///
//////////////////////////////////////////////////

// 检查给定的索引位置是否在容量的允许范围内
// 注: 必须要排除index==cap，因为当buf是整数个字节，且On(index)写入该buf的最后一个bit时，buf[index/8]会溢出
func (bs *BitString) checkValid(index int) error {
	if index >= bs.cap {
		return errors.New("BitString overflow")
	}
	return nil
}

// 取出BitString给定索引位置的值（用布尔量表示1或0）
func (bs *BitString) getBit(index int) bool {
	return (bs.buf[index/8] & (1 << (7 - (index % 8)))) > 0
}

// 逐个自动取出游标cursor指定的位置的bit（用布尔量表示1或0）
func (bs *BitString) readBit() bool {
	bit := bs.getBit(bs.cursor)
	bs.cursor++
	return bit
}

// 计算在已经写入的bits位中，还有多少个bits位可以读取
func (bs *BitString) bitsRemainingForRead() int {
	return bs.len - bs.cursor
}

// 计算在容量cap所允许bits位中，还有多少个bits位可以写入
func (bs *BitString) bitsRemainingForWrite() int {
	return bs.cap - bs.len
}

// 位串扩容
func (bs *BitString) grow(bitLen int) {
	bs.buf = append(bs.buf, make([]byte, (bitLen+7)/8)...)
	bs.cap += bitLen
}

// 解码位串的标签
// 即将`tlb:"8bits"`中的bit位数8解析出来
func decodeBitStringTag(tag string) (int, error) {
	var bitLen int
	if tag == "" {
		return 0, nil
	}
	_, err := fmt.Sscanf(tag, "%dbits", &bitLen)
	if err != nil {
		return 0, err
	}
	return bitLen, nil
}
