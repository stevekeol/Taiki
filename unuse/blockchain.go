// NOTICE: 吃透其简易逻辑
//         将boltdb升级为leveldb
//         再参照btcd重新定义blockchain的结构体和功能
package blockchain

import (
	"Taiki/block"
	"Taiki/db"
	"Taiki/logger"
	"Taiki/pow"
	"Taiki/transaction"
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	// "github.com/boltdb/bolt"
	"os"
)

var log = logger.Log

const dbFile = "blockchain.db"
const blocksBucket = "blocks"
const genesisCoinbaseData = "💋 The Times 02/Dec/2018 Zhangzezhi born in Hangzhou"

type Blockchain struct {
	top []byte // 最新区块的hash
	db  *db.KeyValueStore
}

func (chain *Blockchain) Db() *db.KeyValueStore {
	return chain.db
}

// func txHandler(tx *bolt.Tx) ([]byte, error) {
// 	lastHash := tx.Bucket([]byte(blocksBucket)).Get([]byte("l"))
// 	return lastHash, nil
// }

//把区块添加进区块链,挖矿
func (chain *Blockchain) MineBlock(transactions []*transaction.Transaction) *block.Block {
	var lastHash []byte

	//在一笔交易被放入一个块之前进行验证
	for _, tx := range transactions {
		if chain.VerifyTransaction(tx) != true {
			log.Error("invalid transaction", "tx", tx)
		}
	}
	//只读的方式浏览数据库，获取当前区块链顶端区块的哈希，为加入下一区块做准备
	err := chain.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l")) //通过键"l"拿到区块链顶端区块哈希

		return nil
	})
	if err != nil {
		log.Error("chain.db.View", "err", err)
	}

	//prevBlock := chain.Blocks[len(chain.Blocks)-1]
	//求出新区块
	newBlock := pow.NewBlock(transactions, lastHash)
	// chain.Blocks = append(chain.Blocks,newBlock)
	//把新区块加入到数据库区块链中
	err = chain.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		err := b.Put(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			log.Error("tx.Bucket Put error", "err", err)
		}
		err = b.Put([]byte("l"), newBlock.Hash)
		chain.top = newBlock.Hash

		return nil
	})
	if err != nil {
		log.Error("chain.db.Update error", "err", err)
	}

	return newBlock
}

//创建创世区块  /修改/
func NewGenesisBlock(coinbase *transaction.Transaction) *block.Block {
	return pow.NewBlock([]*transaction.Transaction{coinbase}, []byte{})
}

//创建区块链数据库
func CreateBlockchain(address string) *Blockchain {
	var top []byte
	//此时的创世区块就要包含交易coinbaseTx
	cbtx := transaction.NewCoinbaseTX(address, genesisCoinbaseData)
	genesis := NewGenesisBlock(cbtx)

	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Error("bolt.Open error", "err", err)
	}
	//读写操作数据库
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		//查看名字为blocksBucket的Bucket是否存在
		if b != nil {
			log.Info("Blockchain already existed")
			os.Exit(1)
		}
		//否则，则重新创建
		b, err := tx.CreateBucket([]byte(blocksBucket))
		if err != nil {
			log.Error("tx.CreateBucket error", "err", err)
		}

		err = b.Put(genesis.Hash, genesis.Serialize()) //写入键值对，区块哈希对应序列化后的区块
		if err != nil {
			log.Error("tx.CreateBucket Put error", "err", err)
		}
		err = b.Put([]byte("l"), genesis.Hash) //"l"键对应区块链顶端区块的哈希
		if err != nil {
			log.Error("tx.CreateBucket Put1 error", "err", err)
		}
		top = genesis.Hash //指向最后一个区块，这里也就是创世区块
		return nil
	})
	if err != nil {
		log.Error("db.Update error", "err", err)
	}

	chain := Blockchain{top, db}

	return &chain
}

//实例化一个区块链,默认存储了创世区块 ,接收一个地址为挖矿奖励地址 /修改/
func NewBlockchain() *Blockchain {
	//return &Blockchain{[]*block.Block{NewGenesisBlock()}}
	var top []byte
	//打开一个数据库文件，如果文件不存在则创建该名字的文件
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Error("bolt.Open error", "err", err)
	}
	//读写操作数据库
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		//查看名字为blocksBucket的Bucket是否存在
		if b == nil {
			//不存在
			log.Warn("no blockchain, need build one")
			os.Exit(1)
		}
		//如果存在blocksBucket桶，也就是存在区块链
		//通过键"l"映射出顶端区块的Hash值
		top = b.Get([]byte("l"))

		return nil
	})
	if err != nil {
		log.Error("db.Update err", "err", err)
	}

	chain := Blockchain{top, db} //此时Blockchain结构体字段已经变成这样了
	return &chain
}

//分割线——————迭代器——————
type BlockchainIterator struct {
	currentHash []byte
	db          *db.KeyValueStore
}

//当需要遍历当前区块链时，创建一个此区块链的迭代器
func (chain *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{chain.top, chain.db}

	return bci
}

//迭代器的任务就是返回链中的下一个区块
func (i *BlockchainIterator) Next() *block.Block {
	var Block *block.Block

	//只读方式打开区块链数据库
	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		//获取数据库中当前区块哈希对应的被序列化后的区块
		encodeBlock := b.Get(i.currentHash)
		//反序列化，获得区块
		Block = block.DeserializeBlock(encodeBlock)

		return nil
	})
	if err != nil {
		log.Error("db.View err", "err", err)
	}

	//把迭代器中的当前区块哈希设置为上一区块的哈希，实现迭代的作用
	i.currentHash = Block.PrevBlockHash

	return Block

}

//通过找到未花费输出交易的集合，我们返回集合中的所有未花费交易的交易输出集合
func (chain *Blockchain) FindUTXO() map[string]transaction.TXOutputs {
	//var UTXOs []transaction.TXOutput
	UTXO := make(map[string]transaction.TXOutputs)
	//找到address地址下的未花费交易输出的交易的集合
	//unspentTransactions := chain.FindUnspentTransactions(pubKeyHash)
	//创建一个map，存储已经花费了的交易输出
	spentTXOs := make(map[string][]int)
	//因为要在链上遍历区块，所以要使用到迭代器
	bci := chain.Iterator()

	for {
		block := bci.Next() //迭代

		//遍历当前区块上的交易
		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID) //把交易ID转换成string类型，方便存入map中

			//标签
		Outputs:
			//遍历当前交易中的输出切片，取出交易输出
			for outIdx, out := range tx.Vout {
				//在已经花费了的交易输出map中，如果没有找到对应的交易输出，则表示当前交易的输出未花费
				//反之如下
				if spentTXOs[txID] != nil {
					//存在当前交易的输出中有已经花费的交易输出，
					//则我们遍历map中保存的该交易ID对应的输出的index
					//提示：(这里的已经花费的交易输出index其实就是输入TXInput结构体中的Vout字段)
					for _, spentOutIdx := range spentTXOs[txID] {
						//首先要清楚当前交易输出是一个切片，里面有很多输出，
						//如果map里存储的引用的输出和我们当前遍历到的输出index重合,则表示该输出被引用了
						if spentOutIdx == outIdx {
							continue Outputs //我们就继续遍历下一轮，找到未被引用的输出
						}
					}
				}
				outs := UTXO[txID]
				outs.Outputs = append(outs.Outputs, out)
				UTXO[txID] = outs
			}
			//判断是否为coinbase交易
			if tx.IsCoinbase() == false {
				//如果不是,则遍历当前交易的输入
				for _, in := range tx.Vin {
					inTxID := hex.EncodeToString(in.Txid)
					spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Vout)
				}
			}
		}
		//退出for循环的条件就是遍历到的创世区块后
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
	// //遍历交易集合得到交易，从交易中提取出输出字段Vout,从输出字段中提取出属于address的输出
	// for _,tx := range unspentTransactions {
	// 	for _, out := range tx.Vout {
	// 		if out.IsLockedWithKey(pubKeyHash) {
	// 			UTXOs = append(UTXOs,out)
	// 		}
	// 	}
	// }
	//返回未花费交易输出
	return UTXO
}

//通过交易ID找到一个交易
func (chain *Blockchain) FindTransaction(ID []byte) (transaction.Transaction, error) {
	bci := chain.Iterator()
	for {
		block := bci.Next()

		for _, tx := range block.Transactions {
			if bytes.Compare(tx.ID, ID) == 0 {
				return *tx, nil
			}
		}
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
	return transaction.Transaction{}, errors.New("Transaction is not found")
}

//对交易输入进行签名
func (chain *Blockchain) SignTransaction(tx *transaction.Transaction, privKey ecdsa.PrivateKey) {
	prevTXs := make(map[string]transaction.Transaction)
	for _, vin := range tx.Vin {
		prevTX, err := chain.FindTransaction(vin.Txid) //找到输入引用的输出所在的交易
		if err != nil {
			log.Error("FindTransaction err", "err", err)
		}
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}
	tx.Sign(privKey, prevTXs)
}

//验证交易
func (chain *Blockchain) VerifyTransaction(tx *transaction.Transaction) bool {
	if tx.IsCoinbase() {
		return true
	}
	prevTXs := make(map[string]transaction.Transaction)

	for _, vin := range tx.Vin {
		prevTX, err := chain.FindTransaction(vin.Txid)
		if err != nil {
			log.Error("FindTransaction err", "err", err)
		}
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}
	return tx.Verify(prevTXs) //验证签名
}
