package main

import (
	"fmt"
	"os"
	"runtime"

	"gopkg.in/urfave/cli.v1"
)

const (
	clientIdentifier = "Taiki"
)

var app = NewApp("", "the Taiki command line interface")

func init() {
	app.Action = Taiki
	app.Commands = []cli.Command{}
	app.Before = func(ctx *cli.Context) error {
		fmt.Println("prev action ...")
		runtime.GOMAXPROCS(runtime.NumCPU())
		return nil
	}
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

func NewApp(gitCommit, usage string) *cli.App {
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

func Taiki(ctx *cli.Context) error {
	// node := makeFullNode(ctx)
	// startNode(ctx, node)
	// node.Wait()
	fmt.Println("bootstrap a node")
	return nil
}
