# struct-tree

> 核心重要模块之间的依赖&内嵌关系的设计 （参考ETCD）

---------------------------------------------------

0. `ServiceTag`定义了该地址支持的服务集合
```golang title="以单一的一个01序列来表示所支持的服务集合"
type ServiceTag uint64
const (
	SFNodeNetwork ServiceTag = 1 << iota
	SFNodeBloom
	SFNodeWitness
	...
)
```
1. `NetAddress`定义了网络层的节点信息：最后一次可见时间，支持的服务，地址，端口
```golang title="etcd/wire/netaddressv2.go"
	Timestamp time.Time
	Services ServiceFlag // 该地址支持哪些服务
	Addr net.Addr // 地址
	Port uint16
```

2. `KnownAddress`通过跟踪已知的网络地址，来确定这些地址的可行性
```golang title="btcd/addrmgr/knownaddress.go"
type KnownAddress struct {
	mtx         sync.RWMutex // na and lastattempt
	na          *NetAddress
	srcAddr     *NetAddress
	attempts    int
	lastattempt time.Time
	lastsuccess time.Time
	tried       bool
	refs        int // reference count of new buckets
}
```
	- `KnownAddress`提供了内部的方法`chance`（即仅供`addrmgr`使用）：基于attempts、lastattempt经过一定自定义的规则来计算该已知地址的`被选中可能性`

	- `KnownAddress`也提供了内部的方法`isBad`来供`addrmgr`使用，经过一些规则来判断该地址是否很糟糕

3. `AddrManager`在btc网络层上提供一个并发安全的地址管理器，用来缓存潜在可用的节点

```golang title="btcd/addrmgr/addmanager.go"
type AddrManager struct {
	mtx            sync.RWMutex
	// 存储peers的json文件（每隔10分钟就将peers转存到该文件）
	peersFile      string 
	// 查询给定Host Address对应的NetAddress地址
	lookupFunc     func(string) ([]net.IP, error) 
	rand           *rand.Rand
	key            [32]byte
	addrIndex      map[string]*KnownAddress // address key to ka for all addrs.
	addrNew        [newBucketCount]map[string]*KnownAddress
	addrTried      [triedBucketCount]*list.List
	started        int32 // 根据atomic.AddInt32()来判断是否开启或关闭
	shutdown       int32
	wg             sync.WaitGroup
	quit           chan struct{}
	nTried         int
	nNew           int
	lamtx          sync.Mutex
	// key是ip+port，value是{NetAddress, score}，该score用来判断优先级
	localAddresses map[string]*localAddress
	version        int
}
```
- 分配了`newBucketCount`个桶，通过哈希后取余来落桶；每个桶可以有`newBucketSize`个地址

---------------------------------------------------

0. `ConnReq`即`连接请求`结构体，发向网络地址的连接请求
```golang title="btcd/connmgr/connmanager.go"
type ConnReq struct {
	id uint64 // 该连接请求的识别符（只能使用atomic原子化读写）（从0自增）

	Addr      net.Addr
	Permanent bool // 为true时会在连接断开的时候重试

	conn       net.Conn
	state      ConnState // 有pending,failing,canceled,sdtablished,disconnected五种状态
	stateMtx   sync.RWMutex
	retryCount uint32
}

	// 补充: net.Addr
	type Addr interface {
		Network() string // 网络协议名称（tcp/udp...）
		String() string  // 地址字符串 ("192.0.2.1:25", "[2001:db8::1]:80")
	}

	// 补充：net.Conn
	type Conn interface {
		Read(b []byte) (n int, err error)   // 从conn中读取数据
		Write(b []byte) (n int, err error)  // 向conn中写入数据
		Close() error						// 关闭conn
		LocalAddr() Addr                    // 返回本地网络地址
		RemoteAddr() Addr                   // 返回远程网络地址
		SetReadDeadline(t time.Time) error  // 设置未来Read仍可调用的截止日期
		SetWriteDeadline(t time.Time) error // 设置未来Write仍可调用的截止日期
		SetDeadline(t time.Time) error      // SetReadDeadline+SetWriteDeadline
	}
```

1. `ConnManager`提供一个管理器来处理网络连接
```golang title="btde/connmgr/connmanager.go"
type ConnManager struct {
	// 这些变量照样需要atomic原子化操作
	connReqCount uint64 // 统计`ConnReq`的个数
	start        int32
	stop         int32

	cfg            Config
	wg             sync.WaitGroup
	failedAttempts uint64
	// 作为connmgr统一的请求入口，内部分请求类型而作不同操作
	// registerPending，handleConnected，handleDisconnected，handleFailed
	requests       chan interface{} 
	quit           chan struct{}
}

// 补充： Config
type Config struct {
	// 定义了connmgr负责的一系列监听对象（该对象接受连接）（当connmgr关闭时，它们也需要关闭）
	Listeners []net.Listener
	// 当一个入站连接被accpted时被调用（调用者需要记住关闭该连接）
	OnAccept func(net.Conn)
	// 出站网络连接的维护数量（默认8）
	TargetOutbound uint32
	// 再试的时间间隔（默认5s，最长5min）
	RetryDuration time.Duration
	// 当出站连接建立时触发的回调方法
	OnConnection func(*ConnReq, net.Conn)
	// 当出站连接断开时触发的回调方法
	OnDisconnection func(*ConnReq)
	// 获取要进行网络连接的地址
	GetNewAddress func() (net.Addr, error)
	// 拨号以和给定地址建立连接（不能为空，代码中挂载的net.DialTimeout）
	Dial func(net.Addr) (net.Conn, error)
}
```

2. `connmgr`即网络连接管理器，处理一些通用性的连接问题：维护一组出站连接，sourcing peer（找出节点），禁止的连接，限制最大连接数，tor查找等

> 该connmgr能接收来自source节点（或一组给定addresses）的连接请求，向其拨号，并在建立连接的时候通知调用者；主要的使用目的是，维护一个活跃连接池以此形成P2P网络。此外还提供：

- 连接或断开时的通知
- 处理来自source的新地址的失败与重试
- 只连接到指定的地址
- 永久连接（具有不断增加的回退重试计时器的）
- 断开或移除一个已经建立的连接

---------------------------------------------------

`netsync`即网络同步模块，实现了并发安全的区块同步协议

1. `SyncManager`:和已经连接的对等节点通信来执行初始化的下载，维持链，维护未确认的交易池（in sync）；宣告新块上链等。当前，SyncManager简单的选择一个对等节点，并从其下载所有区块
```golang title="btcd/netsync/manager.go"
type SyncManager struct {
	peerNotifier   PeerNotifier
	started        int32
	shutdown       int32
	// 1. 核心链
	chain          *blockchain.BlockChain
	// 2. 交易池
	txMemPool      *mempool.TxPool
	// 3. 链参数
	chainParams    *chaincfg.Params
	// 区块进度日志
	progressLogger *blockProgressLogger
	msgChan        chan interface{}
	wg             sync.WaitGroup
	quit           chan struct{}

	// These fields should only be accessed from the blockHandler thread
	rejectedTxns     map[chainhash.Hash]struct{}
	requestedTxns    map[chainhash.Hash]struct{}
	requestedBlocks  map[chainhash.Hash]struct{}
	syncPeer         *peerpkg.Peer //存储选出的最佳同步节点
	peerStates       map[*peerpkg.Peer]*peerSyncState
	lastProgressTime time.Time

	// The following fields are used for headers-first mode.
	headersFirstMode bool
	headerList       *list.List
	startHeader      *list.Element
	nextCheckpoint   *chaincfg.Checkpoint

	// An optional fee estimator.
	feeEstimator *mempool.FeeEstimator
}
```
- Start方法会选择节点来同步区块；一旦处于同步过程中，SyncManager会处理进来的区块和header notifications，也会中继`新块的宣告`到别的节点;
- `startSync()`会从候选节点中选出最佳节点来下载区块并同步


---------------------------------------------------

1. 位于main包中的`server`用于处理peers之间往返的通信
```golang title="btcd/server.go"
type server struct {
	// 以下变量必须基于atomic原子化操作
	bytesReceived uint64 // 启动后收到的所有节点的字节总数
	bytesSent     uint64 // 启动后发向所有节点的字节总数
	started       int32  // 启动标志
	shutdown      int32  // 关闭标志
	shutdownSched int32  // 是否计划一段时间后关闭的标志
	startupTime   int64  // 启动时间

	chainParams          *chaincfg.Params
	addrManager          *addrmgr.AddrManager // 地址管理器
	connManager          *connmgr.ConnManager // 连接管理器
	sigCache             *txscript.SigCache
	hashCache            *txscript.HashCache
	rpcServer            *rpcServer
	syncManager          *netsync.SyncManager
	chain                *blockchain.BlockChain
	txMemPool            *mempool.TxPool
	cpuMiner             *cpuminer.CPUMiner
	modifyRebroadcastInv chan interface{}
	newPeers             chan *serverPeer
	donePeers            chan *serverPeer
	banPeers             chan *serverPeer
	query                chan interface{}
	relayInv             chan relayMsg
	broadcast            chan broadcastMsg
	peerHeightsUpdate    chan updatePeerHeightsMsg
	wg                   sync.WaitGroup
	quit                 chan struct{}
	nat                  NAT
	db                   database.DB
	timeSource           blockchain.MedianTimeSource
	services             wire.ServiceFlag

	// The following fields are used for optional indexes.  They will be nil
	// if the associated index is not enabled.  These fields are set during
	// initial creation of the server and never changed afterwards, so they
	// do not need to be protected for concurrent access.
	txIndex   *indexers.TxIndex
	addrIndex *indexers.AddrIndex
	cfIndex   *indexers.CfIndex

	// The fee estimator keeps track of how long transactions are left in
	// the mempool before they are mined into blocks.
	feeEstimator *mempool.FeeEstimator

	// cfCheckptCaches stores a cached slice of filter headers for cfcheckpt
	// messages for each filter type.
	cfCheckptCaches    map[wire.FilterType][]cfHeaderKV
	cfCheckptCachesMtx sync.RWMutex

	// agentBlacklist is a list of blacklisted substrings by which to filter
	// user agents.
	agentBlacklist []string

	// agentWhitelist is a list of whitelisted user agent substrings, no
	// whitelisting will be applied if the list is empty or nil.
	agentWhitelist []string
}
```


```golang title="etcd/peer/peer.go"
type Peer struct {
	// The following variables must only be used atomically.
	bytesReceived uint64
	bytesSent     uint64
	lastRecv      int64
	lastSend      int64
	connected     int32
	disconnect    int32

	conn net.Conn

	// These fields are set at creation time and never modified, so they are
	// safe to read from concurrently without a mutex.
	addr    string
	cfg     Config
	inbound bool

	flagsMtx             sync.Mutex // protects the peer flags below
	na                   *wire.NetAddressV2
	id                   int32
	userAgent            string
	services             wire.ServiceFlag
	versionKnown         bool
	advertisedProtoVer   uint32 // protocol version advertised by remote
	protocolVersion      uint32 // negotiated protocol version
	sendHeadersPreferred bool   // peer sent a sendheaders message
	verAckReceived       bool
	witnessEnabled       bool
	sendAddrV2           bool

	wireEncoding wire.MessageEncoding

	knownInventory     lru.Cache
	prevGetBlocksMtx   sync.Mutex
	prevGetBlocksBegin *chainhash.Hash
	prevGetBlocksStop  *chainhash.Hash
	prevGetHdrsMtx     sync.Mutex
	prevGetHdrsBegin   *chainhash.Hash
	prevGetHdrsStop    *chainhash.Hash

	// These fields keep track of statistics for the peer and are protected
	// by the statsMtx mutex.
	statsMtx           sync.RWMutex
	timeOffset         int64
	timeConnected      time.Time
	startingHeight     int32
	lastBlock          int32
	lastAnnouncedBlock *chainhash.Hash
	lastPingNonce      uint64    // Set to nonce if we have a pending ping.
	lastPingTime       time.Time // Time we sent last ping.
	lastPingMicros     int64     // Time for last ping to return.

	stallControl  chan stallControlMsg
	outputQueue   chan outMsg
	sendQueue     chan outMsg
	sendDoneQueue chan struct{}
	outputInvChan chan *wire.InvVect
	inQuit        chan struct{}
	queueQuit     chan struct{}
	outQuit       chan struct{}
	quit          chan struct{}
}
```