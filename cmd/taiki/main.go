// Taiki is the command-line client for Taiki Blockchain.
package main

import (
	"fmt"
	"os"
	"runtime"

	"Taiki/cmd/taiki/commands"
	"Taiki/cmd/taiki/flags" // 注意该使用方法（包名和文件名一样，且中间路径即为中间文件夹名）
	"Taiki/logger"
	"github.com/urfave/cli/v2"
)

const (
	clientIdentifier = "Taiki"
)

var (
	log = logger.Log
	app *cli.App
)

func init() { app = NewDefaultApp() }

func main() {
	if err := app.Run(os.Args); err != nil {
		// debug.Exit()
		// console.Stdin.Close()
		log.Error("somehting is wrong in whole app")
		os.Exit(1)
	}
}

func NewDefaultApp() *cli.App {
	return &cli.App{
		Name:  "Taiki",
		Usage: "the Taiki command line interface",
		Authors: []*cli.Author{{
			Name:  "stevekeol",
			Email: "stevekeol.x@gmail.com",
		}},
		Commands: []*cli.Command{
			commands.InitCommand,
			commands.CreateWalletCommand,
			commands.CreateBlockchainCommand,
			commands.ListAddressesCommand,
			commands.TransferCommand,
			commands.GetBalanceCommand,
			commands.PrintChainCommand,
		},
		Flags: []cli.Flag{
			flags.Address,
			flags.From,
			flags.To,
			flags.Amount,
		},
		Before:               beforeHandler,
		Action:               appHandler,
		After:                afterHandler,
		EnableBashCompletion: true, // 似乎没用
	}
}

// 创建节点前的前置工作
func beforeHandler(ctx *cli.Context) error {
	log.Info("Job before Taiki.Run()")
	runtime.GOMAXPROCS(runtime.NumCPU())
	return nil
}

// 将要挂载在cli-app上的内核工作
// appHandler is the root action for the Taiki command, creates a node configuration,
// loads the keystore, init the node, then creates and starts the node and node services
func appHandler(ctx *cli.Context) error {
	if args := ctx.Args().Slice(); len(args) > 0 {
		return fmt.Errorf("failed to read command argument: %q", args[0])
	}

	prepare(ctx)
	// TODO 创建钱包地址，创建节点，开启监听服务等
	return nil
}

// 创建节点后的后续工作
func afterHandler(ctx *cli.Context) error {
	log.Info("Job after Taiki.Run()")
	return nil
}

// --------------- Helpers ----------------

func prepare(ctx *cli.Context) error {
	fmt.Println("prepare context with flag from command-line")

	switch {
	case ctx.IsSet(flags.Address.Name):
		log.Info("Starting create a blockchain")
	}

	return nil

	// // Start metrics export if enabled
	// utils.SetupMetrics(ctx)

	// // Start system runtime metrics collection
	// go metrics.CollectProcessMetrics(3 * time.Second)
}
