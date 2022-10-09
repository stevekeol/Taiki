# taiki

```golang title="整个taiki应用的测试示例"
package taiki_test

import (
	"Taiki/logger"
)

var log = logger.Log

func TaikiDemo() {
	stack, err := node.New(&node.Config{})
	defer stack.Close()

	// Create and register a simple network Lifecycle.
	service := new(SampleLifecycle)
	stack.RegisterLifecycle(service)

	// Boot up the entire protocol stack, do a restart and terminate
	if err := stack.Start(); err != nil {
		log.Fatalf("Failed to start the protocol stack: %v", err)
	}
	if err := stack.Close(); err != nil {
		log.Fatalf("Failed to stop the protocol stack: %v", err)
	}
}

```

## Bootstrap

`./taiki.go`文件

1. 加载配置（默认配置文件+命令行参数）
2. 加载相应数据库
3. 运行核心Host服务器