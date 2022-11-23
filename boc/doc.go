/*
Package boc implements bag of cell

BoC Overview:

1. TON区块链的块（Block）和状态（State）中的所有数据都表示为一个Cell的集合

2. TON区块链和TVM将所有永久存储的数据都表示为所谓的`细胞袋`

3. 每个Cell由最多1023个数据位和最多4个对其它Cell的引用组成

4. 不允许Cell循环引用（因此Cell通常被组织为CellTree，准确的说是CellDAG）

5. 任何抽象数据类型都可以表示/序列化为CellTree，其精确的表示方式即`TL-B scheme`

6. Cell的标准布局如下：
	1. 第一部分：占2字节的描述符。第一个字节是该Cell引用其它Cells的数量；第二个字节的前7位是数据位长度l/8向下取整的值；最后一位是是否整除的标志（能整除为0，不能为1）
	2. 第二部分：占l/8向上取整的字节数。8个bit位为一组，以大端序存入一个字节。如果不能被8整除，就将单个的二进制1和适量数量的二进制0(最多6个)附加在后面。
	3. 第三部分：r个引用，每个引用占32字节（即引用Cell的SHA256哈希值）

	因此，每个Cell的大小是： `CellRepr(c)` = `2 + ⌈l/8⌉ + 32 * r`。其中l是数据的位长度，r是引用的Cell数量

*/

package boc
