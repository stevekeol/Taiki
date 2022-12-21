package boc

import (
	"Taiki/logger"
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"hash/crc32"
)

var log = logger.Log
var ErrCircularReference = errors.New("circular reference is not allowed")

var (
	reachBocMagicPrefix   = []byte{0xb5, 0xee, 0x9c, 0x72}
	leanBocMagicPrefix    = []byte{0x68, 0xff, 0x65, 0xf3}
	leanBocMagicPrefixCRC = []byte{0xac, 0xc3, 0xa7, 0x28}
)

/////////////////////////////////////////////////////////////
///											              ///
///                       core                            ///
///                                                       ///
/////////////////////////////////////////////////////////////

// “细胞袋”的序列化
// 以给定Cell为根，将所有引用到的Cells序列化为一个[]byte
// TODO
func SerializeBoC(c *Cell, idx bool, hasCrc32 bool, cacheBits bool, flags uint) ([]byte, error) {
	// 将当前cell作为root
	cells, indexMap, err := topoSort(c)
	if err != nil {
		return nil, err
	}
	// 计算表示Cell个数的数值所需要的字节数
	// TODO：需改名为CellsCountByte
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
	// TODO：需改名为CellsLengthByte
	_, offsetByte := GetBitsAndBytesNeed(uint(fullSize))

	// 构建序列化的位串
	// TODO: 下面都省略了err的处理
	bitStr := NewBitString((CellBits + 32*4 + 32*3) * len(cells))
	// 1. 4字节-magicPrefix
	err = bitStr.WriteBytes([]byte(reachBocMagicPrefix))
	// 2. 1字节-flags集合
	//    1bit的hasIdx
	//    1bit的hasCrc32
	//    1bit的cacheBits
	//    2bits的flags
	//    3bits的sizeBytes(即byteNeed，后期需在源码中更正)
	err = bitStr.WriteBits([]bool{idx, hasCrc32, cacheBits})
	err = bitStr.WriteUint(uint64(flags), 2)
	err = bitStr.WriteInt(int64(byteNeed), 3)
	// 3. 1字节-描述Cell内容长度的数值
	err = bitStr.WriteInt(int64(offsetByte), 8)
	// 4. byteNeed字节-CellsNum(Cells数量)
	err = bitStr.WriteUint(uint64(len(cells)), int(byteNeed*8))
	// 5. byteNeed字节-rootsNum(Cell根的数量，此处默认都是一个Cell根)
	err = bitStr.WriteUint(1, int(byteNeed*8))
	// 6. byteNeed字节-absentNum(?，此处默认为0)
	err = bitStr.WriteUint(0, int(byteNeed*8))
	// 7. offsetByte字节-Cell内容的长度的数值
	err = bitStr.WriteUint(uint64(fullSize), int(offsetByte*8))
	// 8. byteNeed字节-Cell根编号的值(由于默认就一个CellRoot，因此直接默认该编号的值为0)
	err = bitStr.WriteUint(0, int(byteNeed*8))

	if idx {
		for i := range cells {
			// 9. offsetByte字节-每个Cell的所占字节的长度 * Cell个数
			err = bitStr.WriteUint(uint64(sizeIndex[i]), int(offsetByte*8))
			if err != nil {
				return nil, err
			}
		}
	}

	// 这里才是将每一个Cell（当然是以boc格式编排的）写入序列化位串中
	for _, cell := range cells {
		// 注: 每一个Cell都是通过buildBocSchema构建并写入bitStr中的，
		// 因此最终的boc中纯粹CellData那一部分，其实是多个结构相同的Cell子结构缀连而成的
		bocSchema, err := buildBocSchema(cell, indexMap, int(byteNeed))
		// 10. fullSize个字节 - 纯粹的Cell编码内容
		err = bitStr.WriteBytes(bocSchema)
		if err != nil {
			return nil, err
		}
	}

	// 10.5 将上述内容编码后(?)
	bitBytes, err := bitStr.GetTopUppedArray()
	if err != nil {
		return nil, err
	}

	//
	if hasCrc32 {
		checksum := make([]byte, 4)
		binary.LittleEndian.PutUint32(checksum, crc32.Checksum(bitBytes, crc32.MakeTable(crc32.Castagnoli)))
		// 11. 4字节 - 将crc32.CheckSum后缀在其后
		bitBytes = append(bitBytes, checksum...)
	}
	return bitBytes, nil
}

// "细胞袋"的反序列化
// 将给定的[]byte类型的细胞袋反序列化为一组Cell(其实就只是Cell根这一个元素，但形式上是个数组)
func DeSerializeBoC(boc []byte) ([]*Cell, error) {
	// 先将[]byte类型的boc全部解析成bocInfo结构体以备用
	bocInfo, err := parseBoc(boc)
	if err != nil {
		return nil, err
	}
	cellsData := bocInfo.cellsData
	cellsArray := make([]*Cell, 0)
	refsArray := make([][]int, 0)

	// 逐个取出Cell及其对应的引用Cells
	for i := 0; i < int(bocInfo.cellsNum); i++ {
		cell, refs, remaining, err := deserializeCellData(cellsData, bocInfo.sizeBytes)
		if err != nil {
			return nil, err
		}
		cellsData = remaining
		cellsArray = append(cellsArray, cell)
		refsArray = append(refsArray, refs)
	}

	// 逆序逐个取出引用Cell
	for i := int(bocInfo.cellsNum - 1); i >= 0; i-- {
		refCells := refsArray[i]
		if len(refCells) > 4 {
			return nil, ErrCellRefsOverflow
		}
		for refIndex, refCell := range refCells {
			// TODO
			if refIndex < i {
				return nil, errors.New("topological order is broken")
			}
			if refCell >= len(cellsArray) {
				return nil, errors.New("index out of range for boc deseriailization")
			}
			// TODO
			cellsArray[i].refs[refIndex] = cellsArray[refCell]
		}
	}

	rootCells := make([]*Cell, 0)
	for _, item := range bocInfo.rootList {
		rootCells = append(rootCells, cellsArray[item])
	}
	return rootCells, nil
}

// "细胞袋"的反序列化
// 注: 不同的是，不是将[]byte类型的序列化数据，而是base64编码的数据解析成Cell(Cell树)
func DeserializeBocBase64(boc string) ([]*Cell, error) {
	bocData, err := base64.StdEncoding.DecodeString(boc)
	if err != nil {
		return nil, err
	}
	return DeSerializeBoC(bocData)
}

/////////////////////////////////////////////////////////////
///											              ///
///                      Helper                           ///
///                                                       ///
/////////////////////////////////////////////////////////////

// 拓扑排序：汇集所有的Cell（以当前Cell为根，所有递归引用的Cell），以及其哈希和位置的映射关系
// 返回：
// 1.所有引用的Cells的Slice；
// 2.以每个Cell的Hash为键，其在该Slice中的位置为值的Map；
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

type bocInfo struct {
	idx            bool   //
	hasCrc32       bool   //
	hasCacheBits   bool   //
	flags          int    //
	sizeBytes      int    //
	cellsNum       uint   //
	rootsNum       uint   //
	absentNum      uint   //
	totalCellsSize uint   //
	rootList       []uint //
	index          []uint //
	cellsData      []byte //
}

var crcTable = crc32.MakeTable(crc32.Castagnoli)

// 将细胞袋的Header解析出来
// 其实就是将[]byte类型的boc全部解析成一个可直接读取的结构体对象
// TODO: 还需要精简提炼
func parseBoc(boc []byte) (*bocInfo, error) {
	fmt.Printf("len(boc): %v\n", len(boc)) // 7782

	// 边界条件
	if len(boc) < 4+1 {
		return nil, errors.New("not enough bytes for magic prefix")
	}

	// 0.每个boc最后4字节是crc32.Checksum的值（类似校验；可开启）（此处是计算,以备校验）
	checkSum := crc32.Checksum(boc[0:len(boc)-4], crcTable)
	fmt.Printf("last4bytes: %v\n", boc[len(boc)-4:]) // [246 159 218 144]
	fmt.Printf("checkSum: %v\n", checkSum)           // 694736544

	// 1.取出boc中的magicPrefix
	var prefix = boc[0:4]
	// 裁剪掉boc中的magicPrefix
	boc = boc[4:]

	fmt.Printf("prefix: %v\n", prefix) // [181 238 156 114]

	var (
		hasIdx       bool
		hasCrc32     bool
		hasCacheBits bool
		flags        int
		sizeBytes    int
	)

	// 2.处理第一个字节（即flags字节）：取出各种标志位
	if bytes.Equal(prefix, reachBocMagicPrefix) {
		var flagsByte = boc[0]
		fmt.Printf("flagsByte: %v\n", flagsByte) // 2
		hasIdx = (flagsByte & 128) > 0
		hasCrc32 = (flagsByte & 64) > 0
		hasCacheBits = (flagsByte & 32) > 0
		flags = int((flagsByte&16)*2 + (flagsByte & 8))
		// 即最后三个Bit位一起表示sizeBytes
		sizeBytes = int(flagsByte % 8)
	} else if bytes.Equal(prefix, leanBocMagicPrefix) {
		hasIdx = true
		hasCrc32 = false
		hasCacheBits = false
		flags = 0
		sizeBytes = int(boc[0])
	} else if bytes.Equal(prefix, leanBocMagicPrefixCRC) {
		hasIdx = true
		hasCrc32 = true
		hasCacheBits = false
		flags = 0
		sizeBytes = int(boc[0])
	} else {
		return nil, errors.New("unknown magic prefix")
	}

	fmt.Printf("sizeBytes: %v\n", sizeBytes) // 2

	// 裁剪掉boc中的flags
	boc = boc[1:]
	if len(boc) < 1+5*sizeBytes {
		return nil, errors.New("not enough bytes for encoding cells counters")
	}

	offsetBytes := int(boc[0])
	fmt.Printf("offsetBytes: %v\n", offsetBytes) // 2

	// 裁剪掉boc中的offset
	boc = boc[1:]

	cellsNum := readNBytesUIntFromArray(sizeBytes, boc)
	// 裁剪掉boc中的cellsNum
	boc = boc[sizeBytes:]
	rootsNum := readNBytesUIntFromArray(sizeBytes, boc)
	// 裁剪掉boc中的rootsNum
	boc = boc[sizeBytes:]
	absentNum := readNBytesUIntFromArray(sizeBytes, boc)
	// 裁剪掉boc中的absentNum
	boc = boc[sizeBytes:]
	totalCellsSize := readNBytesUIntFromArray(offsetBytes, boc)
	// 裁剪掉boc中的totCellsSize
	boc = boc[offsetBytes:]

	fmt.Printf("cellsNum: %v\n", cellsNum)             // 976
	fmt.Printf("rootsNum: %v\n", rootsNum)             // 1
	fmt.Printf("absentNum: %v\n", absentNum)           // 0
	fmt.Printf("totalCellsSize: %v\n", totalCellsSize) // 7766

	fmt.Printf("len(boc): %v\n", len(boc)) // 7768

	if len(boc) < int(rootsNum)*sizeBytes {
		return nil, errors.New("not enough bytes for encoding root cells hashes")
	}

	// Roots
	rootList := make([]uint, 0)
	for i := 0; i < int(rootsNum); i++ {
		rootList = append(rootList, readNBytesUIntFromArray(sizeBytes, boc))
		// 裁剪掉boc中所有的描述CellRootHash的字节
		boc = boc[sizeBytes:]
	}

	fmt.Printf("rootList: %v\n", rootList)

	// Index
	index := make([]uint, 0, cellsNum)
	if hasIdx {
		if len(boc) < offsetBytes*int(cellsNum) {
			return nil, errors.New("not enough bytes for index encoding")
		}
		for i := 0; i < int(cellsNum); i++ {
			val := readNBytesUIntFromArray(offsetBytes, boc)
			if hasCacheBits {
				val /= 2
			}
			index = append(index, val)
			// 裁剪掉boc中所有的描述CellIndex字节
			boc = boc[offsetBytes:]
		}
	}
	fmt.Printf("index: %v\n", index)

	// Cells
	if len(boc) < int(totalCellsSize) {
		return nil, errors.New("not enough bytes for cells data")
	}

	cellsData := boc[0:totalCellsSize]
	// 裁剪掉boc中所有的描述Cell内容数据的字节
	boc = boc[totalCellsSize:]

	if hasCrc32 {
		if len(boc) < 4 {
			return nil, errors.New("not enough bytes for crc32c hashsum")
		}
		if binary.LittleEndian.Uint32(boc[0:4]) != checkSum {
			return nil, errors.New("crc32c hashsum mismatch")
		}
		// 裁剪掉boc中最后的checkSum的4字节
		boc = boc[4:]
	}

	// 必须干干净净的结束
	if len(boc) > 0 {
		return nil, errors.New("too much bytes in provided boc")
	}

	return &bocInfo{
		hasIdx,
		hasCrc32,
		hasCacheBits,
		flags,
		sizeBytes,
		cellsNum,
		rootsNum,
		absentNum,
		totalCellsSize,
		rootList,
		index,
		cellsData,
	}, nil
}

// 从多个字节中组合读出一个uint数字
func readBytesAsUint(n int, arr []byte) uint {
	var res uint = 0
	for i := 0; i < n; i++ {
		res *= 256
		res += uint(arr[i])
	}
	return res
}
