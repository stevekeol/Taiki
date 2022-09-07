# TON-Go（Taiki）

> TON的Go语言自定义实现和优化。TON是Telegram内部用C++编写的一套公链。拥有新颖高效的技术设计。

> 参考TON白皮书和Etherum来实现。

## references 
[临时性参考简易区块链的实现](https://github.com/zyjblockchain/A_golang_blockchain/blob/master/CLI.go)

## ROADMAP

已经完成：
- 区块链的工作量证明POW
- 把区块链存放到bolt数据库里面，实现命令行接口CLI
- 链上交易，首先实现的是coinbase的交易，然后实现了未花费交易输出的查找从而能得到地址的余额，最后实现地址之间的币发送交易。此时没有实现交易池- ，所以一个区块只能包括一个交易。
- 实现了区块链中的钱包，钱包存储了一对秘钥，用公钥导出了地址，此时有了正真意义上的地址。最后实现了交易的签名。

接下来：
- 将数据库由`github.com/boltdb/bolt`替换为`leveldb`
- 代码结构和实现重新设计
