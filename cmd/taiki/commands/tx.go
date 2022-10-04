package commands

import (
	"encoding/hex"

	"Taiki/blockchain"
	"Taiki/transaction"
	"Taiki/utxo"
	"Taiki/wallet"

	"github.com/urfave/cli/v2"
)

var (
	SendValueCommand = &cli.Command{
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
)

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
