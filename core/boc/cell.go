package boc

import (
	"Taiki/core/tl"
	"Taiki/logger"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
)

var log = logger.Log

const CellBits = 1023
const CellMaxRefs = 4
const CellTreeDepthLimit = 10000

var ErrCellDepthLimit = errors.New("Depth Limit of Cell")

type Cell struct {
	bits      tl.BitString
	refs      [CellMaxRefs]*Cell
	refCursor int  // 引用Cell的游标
	isExotic  bool // 是否是Exotic Cell的标记
}

func NewCell() *Cell {
	return &Cell{
		bits: tl.NewBitString(CellBits),
		refs: [4]*Cell{},
	}
}

// Exotic Cell
// 除了`Ordinary`（也叫simple/data）Cell，其它类型的Cells，都叫`Exotic`Cells。有时也会出现在区块Block或其它数据结构中
// 两者的区分表示是：前者的第一个字节不超过4；后者大于等于5
func NewCellExotic() *Cell {
	return &Cell{
		bits: tl.NewBitString(CellBits),
		refs: [4]*Cell{},
		// isExotic: true,
	}
}

// 该Cell中外部引用的Cell的个数
func (c *Cell) RefsCount() int {
	count := 0
	for i := range c.refs {
		if c.refs[i] != nil {
			count++
		}
	}
	return count
}

// 获取Cell的外部引用
// NOTICE: 需要以slice的形式获取array，并返回给别处使用
func (c *Cell) Refs() []*Cell {
	res := make([]*Cell, 0, 4)
	for _, ref := range c.refs {
		if ref != nil {
			res = append(res, ref)
		}
	}
	return res
}

// 获取已经使用的Bit位的长度
func (c *Cell) BitLen() int {
	return c.bits.GetWriteCursor()
}

// 获取该Cell的Hash
func (c *Cell) Hash() ([]byte, error) {
	return hashCell(c)
}

// 将该Cell哈希结果以字符串编码
func (c *Cell) HashString() (string, error) {
	hash, err := hashCell(c)
	if err != nil {
		log.Error("Cell hash error")
	}
	return hex.EncodeToString(hash), err
}

// 将该Cell序列化为Boc
func (c *Cell) ToBoc() ([]byte, error) {
	return SerializeBoc(c, false, false, false, 0)
}

//------------------------------------------------------

// helper: Cell的哈希化
func hashCell(c *Cell) ([]byte, error) {
	content, err := buildHashSchema(c)
	if err != nil {
		return nil, err
	}
	hash := sha256.Sum256(content)
	return hash[:], nil
}

// helper: Cell哈希前构建成对应的[]byte表示
func buildHashSchema(c *Cell) ([]byte, error) {
	// 0.构建无引用Cell时的基本形式
	res := buildSchemaWithoutRefs(c)

	// 1.依次缀连上每一个引用Cell的最大深度（每个最大深度用2个字节表示）
	for _, ref := range c.Refs() {
		depthSchema := make([]byte, 2)
		depthLimit := CellTreeDepthLimit

		depth, err := getMaxDepth(ref, &depthLimit)
		if err != nil {
			return nil, err
		}
		binary.BigEndian.PutUint16(depthSchema, uint16(depth))
		res = append(res, depthSchema...)
	}

	// 2.再依次缀连上每一个引用Cell的256位哈希(每个哈希用32字节表示)
	for _, ref := range c.Refs {
		hash, err := ref.Hash()
		if err != nil {
			return nil, err
		}
		res = append(res, hash...)
	}
	return res, nil
}

// helper: 先构建没有外部引用的Cell对应的[]byte表示
// schema: 内容buffer，其前缀加两个字节（第一个表示外部引用的cell个数，外加-略；第二个字节表示数据位所占字节相关的信息）
// 		   内容buffer的最后一个字节有特殊处理
func buildBocSchemaWithoutRefs(c *Cell) []byte {
	// 注意这里+7的巧妙！（即哪怕该Cell的bits最后余1个bit，也会再多创建一个字节来容纳）
	// +2是为了在该Cell的前两个字节描述
	res := make([]byte, (c.BitLen()+7)/8+2)

	flag := 0
	if c.isExotic() {
		flag := 8
	}

	// 也是通过第一个字节，除了取出引用Cell的个数，也能判断Cell的类型（Ordinary/Exotic等）
	// Notice: 第一个字节最简情况下就是引用Cell的个数；此处不是最简情况
	res[0] = byte(c.RefsCount() + flag)
	// 每个Cell布局中的第2个字节：前7位为数据位长度/8向下取整的值，最后一位为能否整除的标志（有余数为1）
	// TODO：第2个字节在Hash化的过程中的表示方法？
	res[1] = byte((c.BitLen()+7)/8 + c.BitLen()/8)
	copy(res[2:], c.getBuffer())

	if cell.BitLen()%8 != 0 {
		// TODO: 再次思考为何如此设计
		res[len(res)-1] |= 1 << (7 - cell.BitLen()%8)
	}
	return res
}

// helper: 递归的计算出某个Cell引用的Cell的最大深度（即CellTree DAG的最大深度）
// Question: 在解析Cell时，最大深度的用处和处理是？
func getMaxDepth(c *Cell, iterCounter *int) (int, error) {
	// 每迭代一次getMaxDepth()就检测是否已经达到深度限制
	if *iterCounter == 0 {
		return 0, ErrCellDepthLimit
	}
	*iterCounter -= 1

	maxDepth := 0
	if c.RefsCount() > 0 {
		for _, ref := range c.Refs() {
			depth, err := getMaxDepth(ref, iterCounter)
			if err != nil {
				return 0, err
			}
			if depth > maxDepth {
				maxDepth = depth
			}
		}
		maxDepth++
	}
	return maxDepth, nil
}

// helper: 获取该Cell中真正存储数据的Bit的[]byte
func (c *Cell) getBuffer() {
	return c.bits.Buffer()
}
