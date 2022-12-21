package boc

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"math"
)

// TODO：思考为何1023而不是1024?
const CellBits = 1023
const CellMaxRefs = 4
const CellTreeDepthLimit = 10000

var ErrCellDepthLimit = errors.New("Depth Limit of Cell")
var ErrCellRefsOverflow = errors.New("Too much refs of Cell")
var ErrCellRefsShortage = errors.New("Shortage refs of Cell")

// Taiki所有数据结构(不管是传输还是存储)的序列化的最小单元
// 其底层是依赖的BitString；其上层是BoC；
// TODO：add capacity checking
type Cell struct {
	bits      BitString
	refs      [CellMaxRefs]*Cell
	refCursor int  // 引用Cell的游标
	isExotic  bool // 是否是Exotic Cell的标记
}

func NewCell() *Cell {
	return &Cell{
		// TODO: 每一个Cell所对应的底层Buffer都是分配满的吗？
		bits:     NewBitString(CellBits),
		refs:     [4]*Cell{},
		isExotic: false,
	}
}

// Exotic Cell
// 除了`Ordinary`（也叫simple/data）Cell，其它类型的Cells，都叫`Exotic`Cells。有时也会出现在区块Block或其它数据结构中
// 两者的区分表示是：前者的第一个字节不超过4；后者大于等于5
func NewCellExotic() *Cell {
	return &Cell{
		bits:     NewBitString(CellBits),
		refs:     [4]*Cell{},
		isExotic: true,
	}
}

/////////////////////////////////////////////////////////////
///											              ///
///                       core                            ///
///                                                       ///
/////////////////////////////////////////////////////////////

// 利用给定私钥对当前Cell进行签名，并返回签名(其本质是一个位串)
func (c *Cell) Sign(prvKey ed25519.PrivateKey) (BitString, error) {
	hash, err := c.Hash()
	if err != nil {
		//TODO: 应该返回BitString{}，还是nil?
		return BitString{}, err
	}
	// ed25519的私钥大小为64字节，公钥大小为32字节，签名为64字节
	bs := NewBitString(512)
	err = bs.WriteBytes(ed25519.Sign(prvKey, hash[:]))
	return bs, err
}

// 计算当前Cell的Hash
// 返回结果是[]byte(所有Hash默认都如此，以便操作)
func (c *Cell) Hash() ([]byte, error) {
	return hashCell(c)
}

// 将该Cell哈希结果以字符串编码并返回
func (c *Cell) HashString() (string, error) {
	hash, err := hashCell(c)
	if err != nil {
		log.Error("Cell hash error")
	}
	return hex.EncodeToString(hash), err
}

// 将当前的Cell编码为16进制字符串
// 注：递归的将自身，以及所有引用的Cells全部编码为16进制字符串
func (c *Cell) String() string {
	iter := CellTreeDepthLimit
	return c.string("", &iter)
}

// 将该Cell序列化为Boc
func (c *Cell) Boc() ([]byte, error) {
	return SerializeBoC(c, false, false, false, 0)
}

// 将该Cell序列化为Boc并编码为字符串
func (c *Cell) BocString() (string, error) {
	boc, err := SerializeBoC(c, false, false, false, 0)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(boc), nil
}

// 将该Cell序列化为Boc并编码为base64
func (c *Cell) BocBase64() (string, error) {
	boc, err := SerializeBoC(c, false, false, false, 0)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(boc), nil
}

// 将当前Cell以TL的方式序列化写入到给定的Cell中
// TODO：此处的tag应该是有作用的吧？
func (c Cell) MarshalTL(cell *Cell, tag string) error {
	*cell = c
	return nil
}

// 从给定Cell中以Cell的方式反序列化读出(tag指定的长度的)Cell来给当前Cell赋值
// TODO
func (c *Cell) UnMarshalTL(cell *Cell, tag string) error {
	*c = *cell
	return nil
}

/////////////////////////////////////////////////////////////
///											              ///
///                       Utils                           ///
///                                                       ///
/////////////////////////////////////////////////////////////

// 判断该Cell是否为特种Cell
func (c *Cell) IsExotic() bool {
	return c.isExotic
}

// 当前Cell的游标重置
// 注：游标包括: 位串读取游标，以及引用读取游标
// 注：是递归重置(即将所有引用的Cell的游标都重置)
func (c *Cell) ReadReset() {
	c.readReset(make(map[*Cell]struct{}))
}

/////////////////////////////////////////////////////////////
///											              ///
///                       bits                            ///
///    为了精简Cell挂载的方法，其bits上定义的方法不再重复定义    ///
///                                                       ///
/// 1.  bits.GetWriteCursor(): 获取已经写入的bit位的长度      ///
/// 2.  bits.BitsAvailableForWrite: 还可以写入的bit位的长度
/// 3.  bits.BitsAvailableForRead
/// 4.  bits.WriteBytes:
/// 5.  bits.ReadUint
/// 6.  bits.ReadUintAndBackward
/// 7.  bits.WriteUint
/// 8.  bits.WriteInt
/// 9.  bits.SetTopUppedArray
/// 10. bits.Buffer
/// 11. bits.ReadSkip
/// 12. bits.ReadBits
/// 13. bits.ReadBit
/// 14. bits.WriteBitString
/// 15. bits.ReadInt
/// 16. bits.ReadBytes
/// 17. bits.ReadBigUint
/// 18. ReadRemainingBits
/////////////////////////////////////////////////////////////

// 获取当前Cell的底层的位串已使用的长度
func (c *Cell) BitLen() int {
	return c.bits.GetWriteCursor()
}

/////////////////////////////////////////////////////////////
///											              ///
///                       refs                            ///
///                                                       ///
/////////////////////////////////////////////////////////////

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

// 尝试向当前Cell新增一个空白的Cell引用
// 注：一般在调用NewRef之后，紧接着就是对该新建的Cell的填充处理
func (c *Cell) NewRef() (*Cell, error) {
	cell := NewCell()
	return cell, c.AddRef(cell)
}

// 向当前Cell增加一个外部Cell的引用
// 注：AddRef常用于添加一个已经存在的Cell
// 注：Cell只允许最多4个Refs；且已提前分配好空间
func (c *Cell) AddRef(cell *Cell) error {
	for i := range c.refs {
		if c.refs[i] == nil {
			c.refs[i] = cell
			return nil
		}
	}
	return ErrCellRefsOverflow
}

// 取出当前Cell中下一个外部引用的Cell
func (c *Cell) NextRef() (*Cell, error) {
	if c.refCursor > 3 {
		return nil, ErrCellRefsShortage
	}
	ref := c.refs[c.refCursor]
	if ref != nil {
		c.refCursor++
		ref.ReadReset()
		return ref, nil
	}
	return nil, ErrCellRefsShortage
}

/////////////////////////////////////////////////////////////
///											              ///
///                      Helper                           ///
///                                                       ///
/////////////////////////////////////////////////////////////

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
	res := buildBocSchemaWithoutRefs(c)

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
	for _, ref := range c.Refs() {
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
	if c.isExotic {
		flag = 8
	}

	// 也是通过第一个字节，除了取出引用Cell的个数，也能判断Cell的类型（Ordinary/Exotic等）
	// Notice: 第一个字节最简情况下就是引用Cell的个数；此处不是最简情况
	res[0] = byte(c.RefsCount() + flag)
	// 每个Cell布局中的第2个字节：前7位为数据位长度/8向下取整的值，最后一位为能否整除的标志（有余数为1）
	// TODO：第2个字节在Hash化的过程中的表示方法？
	res[1] = byte((c.BitLen()+7)/8 + c.BitLen()/8)
	copy(res[2:], c.bits.GetBuffer())

	if c.BitLen()%8 != 0 {
		// TODO: 再次思考为何如此设计
		res[len(res)-1] |= 1 << (7 - c.BitLen()%8)
	}
	return res
}

// 将给定的符合BoC编码的单个Cell内容，反序列化为Cell结构
// TODO
func deserializeCellData(cellData []byte, referenceIndexSize int) (*Cell, []int, []byte, error) {
	if len(cellData) < 2 {
		return nil, nil, nil, errors.New("not enough bytes to encode cell descriptors")
	}

	// 注: 每一个Cell都是通过buildBocSchema构建并写入bitStr中的，
	// 因此最终的boc中纯粹CellData那一部分，其实是多个结构相同的Cell子结构缀连而成的
	// 此处，只是取出第一个Cell子结构
	d1 := cellData[0]
	d2 := cellData[1]
	cellData = cellData[2:]

	// 取出每一个Cell内容中的元数据
	// 1. 是Ordinary/Exotic类型的Cell
	// 2. 该Cell引用的Cells数量
	// 3.
	isExotic := (d1 & 8) > 0
	refNum := int(d1 % 8)
	dataBytesSize := int(math.Ceil(float64(d2) / 2))
	fullfilledBytes := !((d2 % 2) > 0)
	withHashes := (d1 & 0b10000) != 0
	levelMask := d1 >> 5

	// 构造Cell以填充
	var cell *Cell
	if isExotic {
		cell = NewCellExotic()
	} else {
		cell = NewCell()
	}
	var refs = make([]int, 0)

	if withHashes {
		maskBits := int(math.Ceil(math.Log2(float64(levelMask) + 1)))
		hashesNum := maskBits + 1
		offset := hashesNum*hashSize + hashesNum*depthSize
		cellData = cellData[offset:]
	}

	if len(cellData) < dataBytesSize+referenceIndexSize*refNum {
		return nil, nil, nil, errors.New("not enough bytes to encode cell data")
	}

	//
	err := cell.bits.SetTopUppedArray(cellData[0:dataBytesSize], fullfilledBytes)
	if err != nil {
		return nil, nil, nil, err
	}
	cellData = cellData[dataBytesSize:]

	for i := 0; i < refNum; i++ {
		refs = append(refs, int(readBytesAsUint(referenceIndexSize, cellData)))
		cellData = cellData[referenceIndexSize:]
	}

	return cell, refs, cellData, nil
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
func (c *Cell) getBuffer() []byte {
	return c.bits.GetBuffer()
}

// 递归的对所有引用的Cell重置游标
func (c *Cell) readReset(seen map[*Cell]struct{}) {
	if _, has := seen[c]; has {
		return
	}
	seen[c] = struct{}{}
	c.bits.ReadReset()
	c.refCursor = 0
	for _, ref := range c.Refs() {
		ref.readReset(seen)
	}
	return
}

// 将当前的Cell编码为16进制字符串
// 注：递归的将自身，以及所有引用的Cells全部编码为16进制字符串
func (c *Cell) string(ident string, iterLimit *int) string {
	s := ident + "x{" + c.bits.ToFiftHex() + "}\n"
	if *iterLimit == 0 {
		return s
	}
	*iterLimit -= 1
	for _, ref := range c.Refs() {
		s += ref.string(ident+" ", iterLimit)
	}
	return s
}
