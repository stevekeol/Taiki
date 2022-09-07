// // Taiki is the command-line client for Taiki Blockchain.
// package main

// import (
// 	"fmt"
// 	"os"
// 	"runtime"
// 	"text/tabwriter"
// 	"time"

// 	"Taiki/block"
// 	"Taiki/blockchain"
// 	"gopkg.in/urfave/cli.v1"
// )

// const (
// 	clientIdentifier = "Taiki"
// )

// var app = NewDefaultApp("", "the Taiki command line interface")

// // cli app的初始化挂载工作
// func init() {
// 	app.Action = Taiki
// 	app.Commands = []cli.Command{}
// 	// 创建节点前的前置工作
// 	app.Before = func(ctx *cli.Context) error {
// 		fmt.Println("prev action ...")
// 		runtime.GOMAXPROCS(runtime.NumCPU())
// 		return nil
// 	}
// 	// 创建节点后的后续工作
// 	app.After = func(ctx *cli.Context) error {
// 		fmt.Println("post action ...")
// 		return nil
// 	}
// }

// func main() {
// 	if err := app.Run(os.Args); err != nil {
// 		// debug.Exit()
// 		// console.Stdin.Close()
// 		fmt.Println("somehting is wrong")
// 		os.Exit(1)
// 	}
// }

// // default cli app的创建工作
// func NewDefaultApp(gitCommit, usage string) *cli.App {
// 	app := cli.NewApp()
// 	app.Author = "stevekeol"
// 	app.Email = "stevekeol.x@gmial.com"
// 	app.Version = "0.1.0"
// 	if len(gitCommit) >= 9 {
// 		app.Version += "-" + gitCommit[:8]
// 	}
// 	app.Usage = usage
// 	return app
// }

// // 将要挂载在cli app上的内核工作
// func Taiki(ctx *cli.Context) error {
// 	TaikiDemo()
// 	return nil
// }

// func TaikiDemo() {
// 	fmt.Println("bootstrap a node")
// 	bc := blockchain.New()

// 	bc.AddBlock("Send 50.0 BTC to Minner01")
// 	time.Sleep(1 * time.Second) //延时记入下一区块，让时间戳不同
// 	bc.AddBlock("Send 25.0 BTC to Minner02")

// 	printBlocks(bc.Blocks())

// }

// func printBlocks(blocks []*block.Block) {
// 	const format = "%x\t %s\t %v\t %x\t \n"
// 	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
// 	// fmt.Fprintf(tw, format, "PrevBlockHash", "Data", "TimeStamp", "Hash123")
// 	// fmt.Fprintf(tw, format, "-----", "------", "-----", "----")
// 	for _, block := range blocks {
// 		fmt.Fprintf(tw, format, block.PrevBlockHash, block.Data, time.Unix(block.TimeStamp, 0), string(block.Hash))
// 	}
// 	tw.Flush() // calculate column widths and print table
// }

// Taiki is the command-line client for Taiki Blockchain.
package main

import (
	"fmt"
	"os"
	"runtime"

	CLI "Taiki/cli"
	"gopkg.in/urfave/cli.v1"
)

const (
	clientIdentifier = "Taiki"
)

var app = NewDefaultApp("", "the Taiki command line interface")

// cli app的初始化挂载工作
func init() {
	app.Action = Taiki
	app.Commands = []cli.Command{}
	// 创建节点前的前置工作
	app.Before = func(ctx *cli.Context) error {
		fmt.Println("prev action ...")
		runtime.GOMAXPROCS(runtime.NumCPU())
		return nil
	}
	// 创建节点后的后续工作
	app.After = func(ctx *cli.Context) error {
		fmt.Println("post action ...")
		return nil
	}
}

func main() {
	if err := app.Run(os.Args); err != nil {
		// debug.Exit()
		// console.Stdin.Close()
		fmt.Println("somehting is wrong")
		os.Exit(1)
	}
}

// default cli app的创建工作
func NewDefaultApp(gitCommit, usage string) *cli.App {
	app := cli.NewApp()
	app.Author = "stevekeol"
	app.Email = "stevekeol.x@gmial.com"
	app.Version = "0.1.0"
	if len(gitCommit) >= 9 {
		app.Version += "-" + gitCommit[:8]
	}
	app.Usage = usage
	return app
}

// 将要挂载在cli app上的内核工作
func Taiki(ctx *cli.Context) error {
	cli := CLI.CLI{}
	cli.Run()
	return nil
}
