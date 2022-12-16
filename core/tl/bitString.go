package tl

import (
	"errors"
	"fmt"
	"strings"
)

var ErrBitStingOverflow = errors.New("BitString overflow")

// BitString提供位操作
type BitString struct {
	buf    []byte // 容纳这些bit位
	cap    int    // 总容量
	len    int    // 实际使用的bit位长
	cursor int    // 接下来可以操作的位置
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
///                   Utils                    ///
///                                            ///
//////////////////////////////////////////////////

// 获取当前可写游标的位置（即目前已经使用的bit位的长度）
func (bs *BitString) GetWriteCursor() int {
	return bs.len
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

//////////////////////////////////////////////////
///											   ///
///                   Write                    ///
///                                            ///
//////////////////////////////////////////////////

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

//////////////////////////////////////////////////
///											   ///
///                   Helper                   ///
///                                            ///
//////////////////////////////////////////////////

// 检查给定的索引位置是否在容量的允许范围内
func (bs *BitString) checkValid(index int) error {
	if index > bs.cap {
		return ErrBitStingOverflow
	}
	return nil
}
