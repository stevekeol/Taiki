package boc

import (
	"Taiki/logger"
	"errors"
	"math"
	"math/bits"
)

var log = logger.Log

var ErrCircularReference = errors.New("circular reference is not allowed")

func SerializeBoc(c *Cell, idx bool, hasCrc32 bool, cacheBits bool, flags uint) ([]byte, error) {
	// 将当前cell作为root
	cells, indexMap, err := topoSort(c)
	if err != nil {
		return nil, err
	}
	// 计算表示Cell个数的数值，所需要的位数
	bitNeed := bits.Len(uint(len(cells)))
	// 计算表示Cell个数的数值，所需要的字节数
	byteNeed := int(math.Max(math.Ceil(float64(bitNeed)/8), 1))
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

	if seen.Has(hash) {
		return nil, ErrCircularReference
	}

	res := make([]*Cell, 0)
	res = append(res, cell)

	for _, ref := range c.Refs() {
		ress, err := topoSortImpl(ref, append(seen, hash))
		if err != nil {
			return nil, err
		}
		res = append(res, ress...)
	}
	return res, nil
}
