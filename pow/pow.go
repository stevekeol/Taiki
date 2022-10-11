package pow

import (
	"Taiki/block"
	"Taiki/logger"
	"Taiki/transaction"
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"math"
	"math/big"
	"strconv"
	"time"
)

//在实际的比特币区块链中，加入一个区块是非常困难的事情，其中运用得到的就是工作量证明

//创建一个工作量证明的结构体
type ProofOfWork struct {
	block  *block.Block //要证明的区块
	target *big.Int     //难度值
}

var log = logger.Log

//声明一个挖矿难度
const targetBits = 10

//实例化一个工作量证明
func NewProofOfWork(b *block.Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	pow := &ProofOfWork{b, target}
	return pow
}

//准备需要进行哈希的数据
func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.Header.PrevBlockHash,
			pow.block.HashTransactions(), //这里被修改，把之前的Data字段修改成交易字段的哈希
			[]byte(strconv.FormatInt(pow.block.Header.Timestamp, 10)),
			[]byte(strconv.FormatInt(targetBits, 10)),
			[]byte(strconv.FormatInt(int64(nonce), 10)),
		},
		[]byte{},
	)
	return data
}

//进行工作量证明,证明成功会返回随机数和区块哈希
func (pow *ProofOfWork) Run() (int, []byte) {
	nonce := 0
	var hash [32]byte
	var hashInt big.Int
	for nonce < math.MaxInt64 {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])

		//把哈希后的数据与难度值进行比较
		if hashInt.Cmp(pow.target) == -1 {
			// @TODO [32]byte转字符串
			log.Info("pow done", "hash", base64.URLEncoding.EncodeToString(hash[:]), "nonce", nonce)
			break
		} else {
			nonce++
		}
	}

	return nonce, hash[:]
}

//实例化一个区块    /更改data为transaction/
func NewBlock(transactions []*transaction.Transaction, prevBlockHash []byte) *block.Block {
	// block := &block.Block{time.Now().Unix(), transactions, prevBlockHash, []byte{}, 0}
	block := &block.Block{
		Header: block.BlockHeader{
			PrevBlockHash: prevBlockHash,
			MerkleRoot:    []byte{},
			Timestamp:     time.Now().Unix(),
			Nonce:         0,
		},
		Height:       uint32(0),
		Hash:         []byte{},
		Transactions: []*transaction.Transaction{},
	}
	// block.SetHash()

	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()
	block.Hash = hash
	block.Header.Nonce = nonce
	return block
}

//其他节点验证nonce是否正确
func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.prepareData(pow.block.Header.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(pow.target) == -1
	return isValid
}
