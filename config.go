package main

const (
	defaultConfigFilename = "taiki.conf"
	defaultDataDirname    = "data"
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
	LogDir     string `long:"logdir" description:"Directory to log output."`
}

// loadConfig initializes and parses the config using a config file and command line options.
func loadConfig() (*config, error) {
	cfg := config{
		ConfigFile: defaultConfigFilename,
		DataDir:    defaultDataDirname,
		LogDir:     defaultLogDirname,
	}

	log.Debug("loadConfig done.")

	//@TODO after a serices of config procudure
	return &cfg, nil
}
