# TON-Go（Taiki）

> TON的Go语言自定义实现和优化。TON是Telegram内部用C++编写的一套公链。拥有新颖高效的技术设计（多链multi-blockchain、异构heterogeneous、动态分片dynamic-sharding、POS+PBFT、智能合约smartcontract、紧密耦合tightly-coupled system）

> 参考TON白皮书、Etherum、BTCD等来实现。

## Usage
> 见 "./cmd/taiki/readme.md"

## TODO
- 主G等待中断信号（如ETCD）再退出 -ok

> 纯命令行的cli工具，无需如此。本身就应该执行完就退出；而是应该在主项目下，类似etcd那样。-ok

- 子命令的嵌入（如Ethereum中chainmd.go）-ok

- debug和metrics的嵌入（用于debug和指标统计）
- 不支持的命令时，不是报错而是给出友好提示


## ROADMAP

已经完成：
- 区块链的工作量证明POW
- 把区块链存放到bolt数据库里面，实现命令行接口CLI
- 链上交易，首先实现的是coinbase的交易，然后实现了未花费交易输出的查找从而能得到地址的余额，最后实现地址之间的币发送交易。此时没有实现交易池- ，所以一个区块只能包括一个交易。
- 实现了区块链中的钱包，钱包存储了一对秘钥，用公钥导出了地址，此时有了正真意义上的地址。最后实现了交易的签名。

接下来：
- 将数据库由`github.com/boltdb/bolt`替换为`leveldb`
- 代码结构和实现重新设计 -ok
- 基于`gopkg.in/urfave/cli.v1`将Taiki支持的命令getbalance等植入（已经升级cli.v1->cli/v2）-ok



## references 
[临时性参考简易区块链的实现](https://github.com/zyjblockchain/A_golang_blockchain/blob/master/CLI.go)

## NOTICE
- 源码的注释，要像ETCD那么详尽
- TON核心概念（boc, shardChain, masterChain, 虚拟的workChain, 虚拟的accountChain，IHR等）源码的实现，可以部分参考ton-c++中validator部分或其它部分