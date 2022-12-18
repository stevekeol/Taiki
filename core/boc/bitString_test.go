package boc

import (
	"bytes"
	"fmt"
	"testing"
)

// 测试到的方法：NewBitString(), WriteUint(), WriteBit(), On(), Off(), checkValid()
func TestWriteUintAndUtils(t *testing.T) {
	bs := NewBitString(8 * 3)

	bs.WriteUint(17, 5)
	if !bytes.Equal(bs.buf, []byte{136, 0, 0}) {
		t.Error(`WriteUint(value, bitLen): wrong when bitLen is enough for value`)
	}

	// 验证高位是否会被正确截断(当给定的bitLen不足以表示该value时)
	bs.WriteUint(17, 4)
	if !bytes.Equal(bs.buf, []byte{136, 128, 0}) {
		t.Error(`WriteUint(value, bitLen): wrong when bitLen is not enough for value`)
	}

	// 验证以位串形式打印的结果(所有数据位)
	if bs.FullBits() != "100010001000000000000000" {
		t.Error(`FullBits(): wrong`)
	}

	// 验证以位串形式打印的结果（已经使用的数据位）
	if bs.UsedBits() != "100010001" {
		t.Error(`UsedBits(): wrong`)
	}

	bs.WriteInt(-17, 4)
	if bs.UsedBits() != "1000100011111" {
		t.Error(`WriteInt(): wrong when bitLen is not enough for value`)
	}

	bs1 := NewBitString(8)
	bs1.WriteUint(255, 8)
	fmt.Println(bs.FullBits())
	bs.WriteBitString(bs1)
	fmt.Println(bs1.FullBits())
	fmt.Println(bs.UsedBits())
	fmt.Println(bs.FullBits())
	fmt.Println(bs.FullBits())
	bs2 := NewBitString(4)
	bs2.WriteUint(15, 4)
	bs.WriteBitString(bs2)
	fmt.Println(bs.FullBits())
	fmt.Println(bs.ToFiftHex()) // 疑问：最后一位为什么是C而非8？
}

// func TestWriteByte
// func TestWriteBitString(t *testing.T) {
// 	bs1 :=
// }
