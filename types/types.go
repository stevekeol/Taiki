package types

type WorkchainId int32
type ShardId uint64
type AccountIdPrefix uint64
type BlockSeqno uint32

type Bits256 [256]BitArray
type BlockHash Bits256
type RootHash Bits256
type FileHash Bits256
type NodeIdShort Bits256 // compatible with adnl::AdnlNodeIdShort
// type StdSmcAddress Bits256; // masterchain / base workchain smart-contract addresses

type UnixTime uint32 // 思考 uint32合适，还是time.Time合适？
type LogicalTime uint64
type ValidatorWeight uint64

// type CatchainSeqno uint32

type ValidatorSessionId Bits256

const workchainIdNotYet WorkchainId = 1 << 31
const masterchainId WorkchainId = -1
const basechainId WorkchainId = 0
const workchainInvalid WorkchainId = 0x80000000

const shardIdAll ShardId = 1 << 63
const splitMergeDelay uint8 = 100      // prepare (delay) split/merge for 100 seconds
const splitMergeInterval uint8 = 100   // split/merge is enabled during 60 second interval
const minSplitMergeInterval uint8 = 30 // split/merge interval must be at least 30 seconds
const maxSplitMergeDelay uint8 = 1000  // end of split/merge interval must be at most 1000 seconds in the future

const (
	capIhrEnabled = 1 << iota // 开启了即时超立方体路由
	capCreateStatsEnabled
	capBounceMsgBody
	capReportVersion
	capSplitMergeTransactions
	capShortDequeue
)

// // 计算分片链前缀长度（63减去后面为0的长度即可）
// func ShardPrefixLen(shard ShardId) uint8 {}

// // 分片链Id字符串化（借助Buffer）
// func Shard2Str(shard ShardId) string {}

// 完整的分片链Id（由所属workchainId和shardId共同决定）
type ShardIdFull struct {
	Workchain WorkchainId
	Shard     ShardId
}

// 完整的账户Id前缀
type AccountIdPrefixFull struct {
	Workchain       WorkchainId
	AccountIdPrefix AccountIdPrefix
}

// 区块Id
type BlockId struct {
	Workchain WorkchainId
	Shard     ShardId
	Seqno     BlockSeqno
}

// 扩展的区块Id
type BlockIdExt struct {
	Id       BlockId
	RootHash RootHash
	FileHash FileHash
}

// 初始zero状态
type ZeroStateIdExt struct {
	Workchain WorkchainId
	RootHash  RootHash
	FileHash  FileHash
}

// type BlockStatus int8

// const (
// 	BlockNone BlockStatus = iota
// 	BlockPrevalidated
// 	BlockValidated
// 	BlockApplied
// )

type BufferSlice []byte // 思考：此处用[]byte，还是bytes.Buffer？(目前看来前者更合适)

type BlockSign struct {
	Node NodeIdShort
	Sign BufferSlice
}

type ReceivedBlock struct {
	Id   BlockIdExt
	Data BufferSlice
}

type BlockBroadcast struct {
	BlockId          BlockIdExt
	Signs            []BlockSign
	CatchainSeqno    CatchainSeqno
	ValidatorSetHash uint32
	Data             BufferSlice
	Proof            BufferSlice
}

type Ed25519PrivateKey struct {
	Prvkey Bits256
}

type Ed25519PublicKey struct {
	Pubkey Bits256
}

type BlockCandidate struct {
	Pubkey           Ed25519PublicKey
	Id               BlockIdExt
	Data             BufferSlice
	CollatedData     BufferSlice
	CollatedFileHash FileHash
}

type ValidatorDescr struct {
	PubKey Ed25519PublicKey
	Weight ValidatorWeight
	Addr   Bits256
}

type ValidatorSessionConfig struct {
	ProtoVersion         uint32
	CatchainIdleTimeout  float64
	CatchainMaxDeps      uint32
	RoundCandidates      uint32
	NextCandidateDelay   float64
	RoundAttemptDuration uint32
	MaxRoundAttempts     uint32
	MaxBlockSize         uint32
	MaxCollatedDataSize  uint32
	NewCatchainIds       bool
}

func NewValidatorSessionConfig() ValidatorSessionConfig {
	return ValidatorSessionConfig{
		ProtoVersion:         0,
		CatchainIdleTimeout:  16.0,
		CatchainMaxDeps:      4,
		RoundCandidates:      3,
		NextCandidateDelay:   2.0,
		RoundAttemptDuration: 16,
		MaxRoundAttempts:     4,
		MaxBlockSize:         4 << 20,
		MaxCollatedDataSize:  4 << 20,
		NewCatchainIds:       false,
	}
}
