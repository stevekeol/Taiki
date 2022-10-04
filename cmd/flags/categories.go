// categories.go in package flags is used to category the command-line flags
package flags

// import "github.com/urfave/cli/v2"

const (
	TAIKICategory       = "TAIKI"
	WalletCategory      = "WALLET"
	BlockchainCategory  = "CHAIN"
	TransactionCategory = "TX"
	// LightCategory      = "LIGHT CLIENT"
	// DevCategory        = "DEVELOPER CHAIN"
	// EthashCategory     = "ETHASH"
	// TxPoolCategory     = "TRANSACTION POOL"
	// PerfCategory       = "PERFORMANCE TUNING"
	// AccountCategory = "ACCOUNT"
	// APICategory        = "API AND CONSOLE"
	// NetworkingCategory = "NETWORKING"
	// MinerCategory      = "MINER"
	// GasPriceCategory   = "GAS PRICE ORACLE"
	// VMCategory         = "VIRTUAL MACHINE"
	// LoggingCategory    = "LOGGING AND DEBUGGING"
	// MetricsCategory    = "METRICS AND STATS"
	// MiscCategory       = "MISC"
	// DeprecatedCategory = "ALIASED (deprecated)"
)

func init() {
	// cli.HelpFlag.(*cli.BoolFlag).Category = MiscCategory
	// cli.VersionFlag.(*cli.BoolFlag).Category = MiscCategory
}
