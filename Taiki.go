package main

import (
	// "debug"
	// "fmt"
	"os"
	// "runtime"
	// "runtime/debug"
	// "runtime/pprof"
)

var (
	cfg *config
)

func main() {
	if err := taikiMain(nil); err != nil {
		os.Exit(1)
	}
}

func taikiMain(serverChan chan<- *server) error {
	// 加载配置文件（配置文件&命令行参数）
	cfg, _ := loadConfig()

	// 加载数据库
	db, _ := loadDatabase(cfg)
	defer func() {
		db.Close()
		log.Info("[2/2]database closed.")
	}()

	interrupt := interruptListeners()
	server, _ := newServer(cfg.Listeners, db, interrupt)
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
