package commands

import (
	"Taiki/base58"
	"Taiki/blockchain"
	"Taiki/utxo"
	"Taiki/wallet"

	"github.com/urfave/cli/v2"
)

var (
	CreateWalletCommand = &cli.Command{
		Name:        "createWallet",
		Action:      createWallet,
		Usage:       "...",
		ArgsUsage:   "",
		Flags:       nil,
		Description: `创建一个钱包，里面放着一对秘钥`,
	}
	ListAddressesCommand = &cli.Command{
		Name:        "listAddresses",
		Action:      listAddresses,
		Usage:       "...",
		ArgsUsage:   "<...>",
		Flags:       nil,
		Description: `罗列钱包中所有的地址`,
	}
	GetBalanceCommand = &cli.Command{
		Name:      "getBalance",
		Action:    getBalance,
		Usage:     "...",
		ArgsUsage: "<address>",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "address"},
		},
		Description: `获取该地址的余额`,
	}
)

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
