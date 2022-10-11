# database

### 需要支持三种
- memory
- leveldb
- rpcdb

### leveldb

Definition-1: `ethdb/leveldb` -> `core/rawdb/database.go: NewLevelDBDatabase()` -> `node/node.go: OpenDatabase()`

Definition-2: `ethdb/leveldb` -> `core/rawdb/database.go: NewLevelDBDatabaseWithFreezer()` -> `node/node.go: OpenDatabaseWithFreezer()`

Usage-1: `les/client.go: New()` 、`les/server.go: NewLesServer()`、

Usage-2：`/cmd/geth/chaincmd.go: initGenesis()` 、`/cmd/utils/flags.go： MakeChainDatabase()`、`/eth/backend.go: New()` 、`/node/node.go: OpenDatabaseWithFreezer()`

> Freezer是用来将不可变链段挪到冷存储中






```golang title="node/node.go 根据config来切换数据库类型"
func (n *Node) OpenDatabase() {
	...
	if n.config.DataDir == "" {
		db = rawdb.NewMemoryDatabase()
	} else {
		db, err = rawdb.NewLevelDBDatabase(n.ResolvePath(name), cache, handles, namespace, readonly)
	}
	...
}
```
```golang title="核心结构体"
// Ethereum implements the Ethereum full node service.
type Ethereum struct {
	config *ethconfig.Config

	// Handlers
	txPool             *core.TxPool
	blockchain         *core.BlockChain
	handler            *handler
	ethDialCandidates  enode.Iterator
	snapDialCandidates enode.Iterator
	merger             *consensus.Merger

	// DB interfaces
	chainDb ethdb.Database // Block chain database

	eventMux       *event.TypeMux
	engine         consensus.Engine
	accountManager *accounts.Manager

	bloomRequests     chan chan *bloombits.Retrieval // Channel receiving bloom data retrieval requests
	bloomIndexer      *core.ChainIndexer             // Bloom indexer operating during block imports
	closeBloomHandler chan struct{}

	APIBackend *EthAPIBackend

	miner     *miner.Miner
	gasPrice  *big.Int
	etherbase common.Address

	networkID     uint64
	netRPCService *ethapi.NetAPI

	p2pServer *p2p.Server

	lock sync.RWMutex // Protects the variadic fields (e.g. gas price and etherbase)

	shutdownTracker *shutdowncheck.ShutdownTracker // Tracks if and when the node has shutdown ungracefully
}
```

## bolt -> leveldb

Open:

db, err := leveldb.OpenFile("path/to/db", nil)
db, err := bolt.Open("my.db", 0600, nil)

bolt.View()
bolt.Update()

