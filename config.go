// TODO Config可能需要移植到单独的模块中（如etcd），赞以文件代替吧
package main

const (
	defaultConfigFilename = "taiki.conf"
	defaultDataDirname    = "data"
	defaultDbType         = "leveldb"
	defaultLogLevel       = "info"
	defaultLogDirname     = "logs"
	defaultLogFilename    = "taiki.log"
	defaultBlockMinSize   = 0
	defaultBlockMaxSize   = 750000 // @TODO
)

// config defines the configuration options for taiki
// NOTICE: reference from btcd/config.go
type config struct {
	ConfigFile string `short:"C" long:"configfile" description:"Path to configuration file"`
	CPUProfile string `long:"cpuprofile" description:"Write CPU profile to the specified file"`
	DataDir    string `short:"b" long:"datadir" description:"Directory to store data"`
	DbType     string `long:"dbtype" description:"memorydb, leveldb or rpcdb"`
	LogDir     string `long:"logdir" description:"Directory to log output."`
	Listeners  []string
}

// loadConfig initializes and parses the config using a config file and command line options.
func loadConfig() (*config, error) {
	cfg := config{
		ConfigFile: defaultConfigFilename,
		DataDir:    defaultDataDirname,
		DbType:     defaultDbType, // memorydb or leveldb
		LogDir:     defaultLogDirname,
		Listeners:  []string{},
	}

	log.Info("[1/3]config loaded.")

	//@TODO after a serices of config procudure
	return &cfg, nil
}
