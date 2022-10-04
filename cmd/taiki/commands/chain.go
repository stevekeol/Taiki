package commands

import (
	"strconv"

	"Taiki/blockchain"
	"Taiki/logger"
	"Taiki/pow"
	"Taiki/wallet"

	"github.com/urfave/cli/v2"
)

var (
	log = logger.Log
)

var (
	CreateBlockchainCommand = &cli.Command{
		Name:      "createBlockChain",
		Action:    createBlockChain,
		Usage:     "...",
		ArgsUsage: "<...>",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "address"},
		},
		Description: `创建一条链并且该地址会得到打块激励`,
	}

	PrintChainCommand = &cli.Command{
		Name:        "printChain",
		Usage:       "...",
		ArgsUsage:   "<...>",
		Flags:       nil,
		Description: `打印链（闹着玩儿）`,
	}
)

func createBlockChain(ctx *cli.Context) error {
	address := ctx.String("address")
	log.Debug("address get from command-line successed", "address", address)
	if !wallet.ValidateAddress(address) {
		log.Error("ERROR: Address is not valid")
	}

	bc := blockchain.CreateBlockchain(address)
	defer bc.Db().Close()

	// UTXOSet := utxo.UTXOSet{bc}
	// UTXOSet.Reindex()
	log.Info("createBlockchain done")
	return nil
}

func printChain(ctx *cli.Context) error {
	//实例化一条链
	bc := blockchain.NewBlockchain() //因为已经有了链，不会重新创建链，所以接收的address设置为空
	defer bc.Db().Close()

	//这里需要用到迭代区块链的思想
	//创建一个迭代器
	bci := bc.Iterator()

	for {

		block := bci.Next() //从顶端区块向前面的区块迭代

		// fmt.Printf("------======= 区块 %x ============\n", block.Hash)
		// fmt.Printf("时间戳:%v\n", block.Timestamp)
		// fmt.Printf("PrevHash:%x\n", block.PrevBlockHash)
		log.Info("BlockInfo", "Hash", block.Hash,
			"Timestamp", block.Timestamp,
			"PrevBlockHash", block.PrevBlockHash)

		//fmt.Printf("Data:%s\n",block.Data)
		//fmt.Printf("Hash:%x\n",block.Hash)
		//验证当前区块的pow
		pow := pow.NewProofOfWork(block)
		boolen := pow.Validate()
		// fmt.Printf("POW is %s\n", strconv.FormatBool(boolen))

		log.Info("POW", "pow", strconv.FormatBool(boolen))

		for _, tx := range block.Transactions {
			transaction := (*tx).String()
			log.Info("transaction", "tx", transaction)
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return nil
}
