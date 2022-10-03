package main

import (
	"encoding/hex"
	"strconv"

	"Taiki/base58"
	"Taiki/blockchain"
	"Taiki/pow"
	"Taiki/transaction"
	"Taiki/utxo"
	"Taiki/wallet"

	"github.com/urfave/cli/v2"
)

var (
	initCommand = &cli.Command{
		Name:      "init",
		Action:    initGenesis,
		Aliases:   []string{"i"},
		Usage:     "Bootstrap and initialize a new genesis block",
		ArgsUsage: "<genesisPath>",
		Flags:     nil,
		Description: `
The init command initializes a new genesis block and definition for the network.
This is a destructive action and changes the network in which you will be
participating.

It expects the genesis file as argument.`,
	}

	createWalletCommand = &cli.Command{
		Name:        "createWallet",
		Action:      createWallet,
		Usage:       "...",
		ArgsUsage:   "",
		Flags:       nil,
		Description: `创建一个钱包，里面放着一对秘钥`,
	}

	createBlockchainCommand = &cli.Command{
		Name:      "createBlockChain",
		Action:    createBlockChain,
		Usage:     "...",
		ArgsUsage: "<...>",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "address"},
		},
		Description: `创建一条链并且该地址会得到打块激励`,
	}

	listAddressesCommand = &cli.Command{
		Name:        "listAddresses",
		Action:      listAddresses,
		Usage:       "...",
		ArgsUsage:   "<...>",
		Flags:       nil,
		Description: `罗列钱包中所有的地址`,
	}
	sendValueCommand = &cli.Command{
		Name:      "send",
		Action:    sendValue,
		Usage:     "...",
		ArgsUsage: "<...>",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "from"},
			&cli.StringFlag{Name: "to"},
			&cli.UintFlag{Name: "amount"},
		},
		Description: `从地址from发送amount的币给地址to`,
	}
	getBalanceCommand = &cli.Command{
		Name:      "getBalance",
		Action:    getBalance,
		Usage:     "...",
		ArgsUsage: "<address>",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "address"},
		},
		Description: `获取该地址的余额`,
	}
	printChainCommand = &cli.Command{
		Name:        "printChain",
		Action:      printChain,
		Usage:       "...",
		ArgsUsage:   "<...>",
		Flags:       nil,
		Description: `打印链（闹着玩儿）`,
	}
)

func initGenesis(ctx *cli.Context) error {
	log.Debug("initGenesis cmd")
	return nil
}

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

func createWallet(ctx *cli.Context) error {
	wallets, _ := wallet.NewWallets()
	address := wallets.CreateWallet()
	wallets.SaveToFile()
	log.Info("new address created", "address", address)

	return nil
}

func listAddresses(ctx *cli.Context) error {
	wallets, err := wallet.NewWallets()
	if err != nil {
		log.Error("listAddresses", "err", err)
	}
	addresses := wallets.GetAddresses()
	for _, address := range addresses {
		log.Info("listAddresses", "address", address)
	}

	return nil
}

func sendValue(ctx *cli.Context) error {
	from, to, amount := ctx.String("from"), ctx.String("to"), ctx.Uint("amount")
	if !wallet.ValidateAddress(from) {
		log.Error("Address(from) is not valid")
	}
	if !wallet.ValidateAddress(to) {
		log.Error("Address(to) is not valid")
	}

	bc := blockchain.NewBlockchain()
	UTXOSet := utxo.UTXOSet{bc}
	defer bc.Db().Close()

	//tx := NewUTXOTransaction(from,to,amount,bc)
	////挖矿奖励的交易，把挖矿的奖励发送给矿工，这里的矿工默认为发送交易的地址
	//cbtx := transaction.NewCoinbaseTX(from,"")

	//挖出一个包含该交易的区块,此时区块还包含了-挖矿奖励的交易
	//bc.MineBlock([]*transaction.Transaction{cbtx,tx})
	tx := NewUTXOTransaction(from, to, int(amount), &UTXOSet)
	cbTx := transaction.NewCoinbaseTX(ctx.Value("from").(string), "")
	txs := []*transaction.Transaction{cbTx, tx}
	newBlock := bc.MineBlock(txs)
	UTXOSet.Update(newBlock)
	log.Info("sendValue successed.", "from", from, "to", to, "amount", amount)

	return nil
}

func getBalance(ctx *cli.Context) error {
	log.Debug("getBalance", "address", ctx.String("address"))

	address := ctx.String("address") // urfave/cli中获取命令行flag的标准做法
	if !wallet.ValidateAddress(address) {
		log.Error("address is not valid", "address", address)
		return nil //@TODO
	}
	bc := blockchain.NewBlockchain()
	UTXOSet := utxo.UTXOSet{bc}
	defer bc.Db().Close()

	balance := 0
	pubKeyHash := base58.Base58Decode([]byte(address))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4] //这里的4是校验位字节数，这里就不在其他包调过来了

	UTXOs := UTXOSet.FindUTXO(pubKeyHash)

	//遍历UTXOs中的交易输出out，得到输出字段out.Value,求出余额
	for _, out := range UTXOs {
		balance += out.Value
	}

	log.Info("getBalance", "address", address, "balance", balance)

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

//发送币操作,相当于创建一笔未花费输出交易
func NewUTXOTransaction(from, to string, amount int, UTXOSet *utxo.UTXOSet) *transaction.Transaction {
	var inputs []transaction.TXInput
	var outputs []transaction.TXOutput
	//validOutputs是一个存放要用到的未花费输出的交易/输出的map
	//acc,validOutputs := bc.FindSpendableOutputs(from,amount)
	wallets, err := wallet.NewWallets()
	if err != nil {
		log.Error("NewWallets error", "err", err)
	}
	_wallet := wallets.GetWallet(from)
	pubKeyHash := wallet.HashPubKey(_wallet.PublicKey)
	acc, validOutputs := UTXOSet.FindSpendableOutputs(pubKeyHash, amount)
	if acc < amount {
		log.Error("Not enough tokens", "acc", acc, "amount", amount)
	}
	//通过validOutputs里面的数据来放入建立一个输入列表
	for txid, outs := range validOutputs {
		//反序列化得到txID
		txID, err := hex.DecodeString(txid)
		if err != nil {
			log.Error("hex.DecodeString error", "err", err)
		}
		//遍历输出outs切片,得到TXInput里的Vout字段值
		for _, out := range outs {
			//input := transaction.TXInput{txID,out,from}
			input := transaction.TXInput{txID, out, nil, _wallet.PublicKey}
			inputs = append(inputs, input)
		}
	}
	//建立一个输出列表
	//outputs = append(outputs,transaction.TXOutput{amount,to})
	outputs = append(outputs, *transaction.NewTXOutput(amount, to))
	if acc > amount {
		//outputs = append(outputs,transaction.TXOutput{acc - amount,from}) //相当于找零
		outputs = append(outputs, *transaction.NewTXOutput(acc-amount, from)) //相当于找零
	}
	tx := transaction.Transaction{nil, inputs, outputs}
	//tx.SetID()
	tx.ID = tx.Hash()
	UTXOSet.Blockchain.SignTransaction(&tx, _wallet.PrivateKey)

	return &tx
}
