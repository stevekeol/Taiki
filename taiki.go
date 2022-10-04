package main

import (
	// "debug"
	// "fmt"
	"os"
	// "runtime"
	// "runtime/debug"
	// "runtime/pprof"
)

func main() {
	if err := taikiMain(nil); err != nil {
		os.Exit(1)
	}
}

func taikiMain(serverChan chan<- *server) error {
	// 加载配置文件（配置文件&命令行参数）
	// cfg, _ := loadConfig()
	// db, err := loadBlockDB()
	interrupt := interruptListeners()

	// server, err := newServer(cfg.Listeners, db)
	Listeners := []string{}
	server, _ := newServer(Listeners, interrupt)
	server.Start()
	defer func() {
		server.Stop()
	}()

	if serverChan != nil {
		serverChan <- server
	}

	<-interrupt
	return nil
}
