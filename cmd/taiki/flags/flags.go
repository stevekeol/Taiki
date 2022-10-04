package flags

// references: /home/stevekeol/Code/BlockChain-Projects/Etherum/go-ethereum/cmd/utils/flags.go
import (
	// "fmt"
	// "path/filepath"
	"github.com/urfave/cli/v2"
)

var (
	//Demo Options
	Address = &cli.StringFlag{
		Name:     "address",
		Value:    "14eA6EswuiuMGVXzpmwMxPJPR4qgR7bjRf", // default value
		Usage:    "sendValue from=<FROM> to=<TO> amount=<AMOUNT>",
		Category: WalletCategory,
	}
	From = &cli.StringFlag{
		Name:     "from",
		Usage:    "sendValue from=<FROM> to=<TO> amount=<AMOUNT>",
		Category: WalletCategory,
	}
	To = &cli.StringFlag{
		Name:     "to",
		Usage:    "sendValue from=<FROM> to=<TO> amount=<AMOUNT>",
		Category: WalletCategory,
	}
	Amount = &cli.UintFlag{
		Name:     "amount",
		Usage:    "sendValue from=<FROM> to=<TO> amount=<AMOUNT>",
		Category: WalletCategory,
	}
)

// GroupFlags combines the given flag slices together and returns the merged one.
func GroupFlags(groups ...[]cli.Flag) []cli.Flag {
	var ret []cli.Flag
	for _, group := range groups {
		ret = append(ret, group...)
	}
	return ret
}
