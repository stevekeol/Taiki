# log

> 基于[log15](https://github.com/inconshreveable/log15)的简易封装（暂时不封装了，直接使用）



## Usage
```
package main

import (
	log "github.com/inconshreveable/log15"
	"os"
)

func main() {

	// all loggers can have key/value context
	srvlog := log.New("module", "app/server")

	// child loggers with inherited context
	connlog := srvlog.New("raddr", "172.1.1.1")

	// lazy evaluation
	connlog.Debug("ping remote", "latency")

	connlog.Info("commandLine", "args", os.Args[1:])

	// all log messages can have key/value context
	srvlog.Warn("abnormal conn rate", "rate", 0.500, "low", 0.100, "high", 0.800)

	// flexible configuration
	// srvlog.SetHandler(log.MultiHandler(
	// 	log.StreamHandler(os.Stderr, log.LogfmtFormat()),
	// 	log.LvlFilterHandler(
	// 		log.LvlError,
	// 		log.Must.FileHandler("errors.json", log.JSONFormat()))))

	connlog.Error("ping remote", "rate", 1.5, "low", 2.5)
	connlog.Crit("ping remote", "rate", 1.5, "low", 2.5)

	// 打印文件名和行号
	h := log.CallerFileHandler(log.StdoutHandler)
	log.Root().SetHandler(h)
	log.Error("open file", "err", "err123")

	//将符合日志级别的record写入指定的文件
	h1 := log.MultiHandler(
		log.LvlFilterHandler(log.LvlError, log.Must.FileHandler("./service.json", log.JsonFormat())),
		// log.MatchFilterHandler("pkg", "app/rpc", log.StdoutHandler()),
	)
	log.Root().SetHandler(h1)
	log.Error("open file", "err", "err123")

}