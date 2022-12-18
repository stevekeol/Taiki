package boc

import (
	"Taiki/logger"
	"encoding/binary"
	"errors"
	"hash/crc32"
)

var log = logger.Log
var ErrCircularReference = errors.New("circular reference is not allowed")

var (
	reachBocMagicPrefix   = []byte{0xb5, 0xee, 0x9c, 0x72}
	leanBocMagicPrefix    = []byte{0x68, 0xff, 0x65, 0xf3}
	leanBocMagicPrefixCRC = []byte{0xac, 0xc3, 0xa7, 0x28}
)

// “细胞袋”的序列化
// 以给定Cell为根，将所有引用到的Cells序列化为一个[]byte
func SerializeBoc(c *Cell, idx bool, hasCrc32 bool, cacheBits bool, flags uint) ([]byte, error) {
	// 将当前cell作为root
	cells, indexMap, err := topoSort(c)
	if err != nil {
		return nil, err
	}
	// 计算表示Cell个数的数值所需要的字节数
	_, byteNeed := GetBitsAndBytesNeed(uint(len(cells)))

	// 计算所有Cell所占的长度
	fullSize := 0
	// 计算并存储所有Cell的大小（注：累和的形式，揣摩）
	sizeIndex := make([]uint, 0)
	for _, cell := range cells {
		sizeIndex = append(sizeIndex, uint(fullSize))
		bocSchema, err := buildBocSchema(cell, indexMap, int(byteNeed))
		if err != nil {
			return nil, err
		}
		fullSize += len(bocSchema)
	}
	// 计算表示所有Cell所占长度的数值，所需要的字节数
	_, offsetByte := GetBitsAndBytesNeed(uint(fullSize))

	// 构建序列化的位串
	// Notice: ?
	// Notice: 下面都省略了err的处理
	bitStr := NewBitString((CellBits + 32*4 + 32*3) * len(cells))
	err = bitStr.WriteBytes([]byte(reachBocMagicPrefix))
	err = bitStr.WriteBits([]bool{idx, hasCrc32, cacheBits})
	err = bitStr.WriteUint(uint64(flags), 2)
	err = bitStr.WriteInt(int64(byteNeed), 3)
	err = bitStr.WriteInt(int64(offsetByte), 8)
	err = bitStr.WriteUint(uint64(len(cells)), int(byteNeed*8))
	err = bitStr.WriteUint(1, int(byteNeed*8))
	err = bitStr.WriteUint(0, int(byteNeed*8))
	err = bitStr.WriteUint(uint64(fullSize), int(offsetByte*8))
	err = bitStr.WriteUint(0, int(byteNeed*8))

	if idx {
		for i := range cells {
			err = bitStr.WriteUint(uint64(sizeIndex[i]), int(offsetByte*8))
			if err != nil {
				return nil, err
			}
		}
	}

	// 这里才是将每一个Cell（当然是以boc格式编排的）写入序列化位串中
	for _, cell := range cells {
		bocSchema, err := buildBocSchema(cell, indexMap, int(byteNeed))
		err = bitStr.WriteBytes(bocSchema)
		if err != nil {
			return nil, err
		}
	}

	// ？
	bitBytes, err := bitStr.GetTopUppedArray()
	if err != nil {
		return nil, err
	}

	if hasCrc32 {
		checksum := make([]byte, 4)
		binary.LittleEndian.PutUint32(checksum, crc32.Checksum(bitBytes, crc32.MakeTable(crc32.Castagnoli)))

		bitBytes = append(bitBytes, checksum...)
	}
	return bitBytes, nil
}

//-----------------------------------------------------------

// 拓扑排序：汇集所有的Cell（以当前Cell为根，所有递归引用的Cell），以及其哈希和位置的映射关系
// 返回：1.所有引用的Cells的Slice；2.以每个Cell的Hash为键，其在该Slice中的位置为值的Map；
func topoSort(c *Cell) ([]*Cell, map[string]int, error) {
	res, err := topoSortImpl(c, []string{})
	if err != nil {
		return nil, nil, err
	}

	indexMap := make(map[string]int)

	for index, cell := range res {
		hashStr, err := cell.HashString()
		if err != nil {
			return nil, nil, err
		}
		indexMap[hashStr] = index
	}
	return res, indexMap, nil
}

// 以该Cell为根，然后"深度优先"的遍历所有引用的Cell，将其逐个压入到一个串行的[]*Cell中
func topoSortImpl(c *Cell, seen []string) ([]*Cell, error) {
	hash, err := c.HashString()
	if err != nil {
		return nil, err
	}

	if contains(seen, hash) {
		return nil, ErrCircularReference
	}

	res := make([]*Cell, 0)
	res = append(res, c)

	for _, ref := range c.Refs() {
		ress, err := topoSortImpl(ref, append(seen, hash))
		if err != nil {
			return nil, err
		}
		res = append(res, ress...)
	}
	return res, nil
}

// 对给定Cell构建成Boc的表示格式：给定Cell的标准布局在前，后面依次是所有"直接引用"的Cells的内部索引
// indexMap: 对给定Cell的引用Cells，逐个取出并计算出其Hash后，然后从中取出对应的位置索引，每个占8字节
// byteNeed: 用来计算位置索引在给定的8字节中占最后面的多少个字节
func buildBocSchema(c *Cell, indexMap map[string]int, byteNeed int) ([]byte, error) {
	// 先将给定Cell（除去引用的Cells）
	res := buildBocSchemaWithoutRefs(c)
	for _, ref := range c.Refs() {
		hash, err := ref.HashString()
		if err != nil {
			return nil, err
		}
		b := make([]byte, 8)
		// Notice: Cell内部引用使用的是内部索引，而不是其HashString
		// 但同时仍用8个字节：为了...
		binary.BigEndian.PutUint64(b, uint64(indexMap[hash]))
		// Notice: 明显占不完8个字节；所以Cell的索引位置的数值靠这8字节的右边(低位)存储
		res = append(res, b[8-byteNeed:]...)
	}
	return res, nil
}
