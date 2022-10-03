# README

## Usage

```bash
# 1. 进入合适目录并编译
cd ./Taiki/cmd/taiki && go build -o Taiki

# 2. 测试是否编译成功
./Taiki

# Usage:
#  createblockchain -address ADDRESS 创建一条链并且该地址会得到狗头金
#  createwallet                      创建一个钱包，里面放着一对秘钥
#  getbalance       -address ADDRESS 得到该地址的余额
#  listaddresses                     罗列钱包中所有的地址
#  printchain                        打印链
#  reindexutxo                       重构UTXO集合
#  send             -from FROM -to TO -amount AMOUNT 地址from发送amount的币给地址to

# 3. 创建钱包
./Taiki createwallet
# new address created                      address=1Lp5qqtq6xH148tLMKZkzyRDwjxSPoAEea

# 4. 创建一条链
./Taiki createblockchain --address 1Lp5qqtq6xH148tLMKZkzyRDwjxSPoAEea
# pow done  hash="AAhNl-lykDruIb3Ex5kMRqEu4R-yEukDHVqzKB9H17o=" nonce=82
# Blockchain already existed 

# 5. 再创建一个钱包地址
./Taiki createwallet
# new address created                      address=1AeoXGvco5QeGVazzuzTzXqePuXzKJXNWK

# 6. 罗列所有钱包地址
./Taiki listaddresses
# INFO[09-30|11:08:40] listAddresses                            address=1JwWqF4WTXTvb68WD1izns9PQq7aHvPVrk
# INFO[09-30|11:08:40] listAddresses                            address=14eA6EswuiuMGVXzpmwMxPJPR4qgR7bjRf
# INFO[09-30|11:08:40] listAddresses                            address=1Lp5qqtq6xH148tLMKZkzyRDwjxSPoAEea
# INFO[09-30|11:08:40] listAddresses                            address=1AeoXGvco5QeGVazzuzTzXqePuXzKJXNWK
# INFO[09-30|11:08:40] listAddresses                            address=12bx1JSy4Df7ZEN8GrrgVH7M2X7UB8M8YB
# INFO[09-30|11:08:40] listAddresses                            address=1GZCzzMfH5U8PB3PD9rhCNBkCkXHMwxEDz

# 7. 转账
send             -from FROM -to TO -amount AMOUNT 

# 8. 查询账户余额
./Taiki getbalance       -address ADDRESS

# 9. 打印整个链
./Taiki printchain
```

```
Usage:
  btcd [OPTIONS]

Application Options:
      --addcheckpoint=        Add a custom checkpoint.  Format: '<height>:<hash>'
  -a, --addpeer=              Add a peer to connect with at startup
      --addrindex             Maintain a full address-based transaction index which makes the searchrawtransactions RPC available
      --agentblacklist=       A comma separated list of user-agent substrings which will cause btcd to reject any peers whose user-agent contains any of the blacklisted
                              substrings.
      --agentwhitelist=       A comma separated list of user-agent substrings which will cause btcd to require all peers' user-agents to contain one of the whitelisted
                              substrings. The blacklist is applied before the blacklist, and an empty whitelist will allow all agents that do not fail the blacklist.
      --banduration=          How long to ban misbehaving peers.  Valid time units are {s, m, h}.  Minimum 1 second (default: 24h0m0s)
      --banthreshold=         Maximum allowed ban score before disconnecting and banning misbehaving peers. (default: 100)
      --blockmaxsize=         Maximum block size in bytes to be used when creating a block (default: 750000)
      --blockminsize=         Mininum block size in bytes to be used when creating a block
      --blockmaxweight=       Maximum block weight to be used when creating a block (default: 3000000)
      --blockminweight=       Mininum block weight to be used when creating a block
      --blockprioritysize=    Size in bytes for high-priority/low-fee transactions when creating a block (default: 50000)
      --blocksonly            Do not accept transactions from remote peers.
  -C, --configfile=           Path to configuration file (default: /home/stevekeol/.btcd/btcd.conf)
      --connect=              Connect only to the specified peers at startup
      --cpuprofile=           Write CPU profile to the specified file
  -b, --datadir=              Directory to store data (default: /home/stevekeol/.btcd/data)
      --dbtype=               Database backend to use for the Block Chain (default: ffldb)
  -d, --debuglevel=           Logging level for all subsystems {trace, debug, info, warn, error, critical} -- You may also specify
                              <subsystem>=<level>,<subsystem2>=<level>,... to set the log level for individual subsystems -- Use show to list available subsystems (default:
                              info)
      --dropaddrindex         Deletes the address-based transaction index from the database on start up and then exits.
      --dropcfindex           Deletes the index used for committed filtering (CF) support from the database on start up and then exits.
      --droptxindex           Deletes the hash-based transaction index from the database on start up and then exits.
      --externalip=           Add an ip to the list of local addresses we claim to listen on to peers
      --generate              Generate (mine) bitcoins using the CPU
      --limitfreerelay=       Limit relay of transactions with no transaction fee to the given amount in thousands of bytes per minute (default: 15)
      --listen=               Add an interface/port to listen for connections (default all interfaces port: 8333, testnet: 18333)
      --logdir=               Directory to log output. (default: /home/stevekeol/.btcd/logs)
      --maxorphantx=          Max number of orphan transactions to keep in memory (default: 100)
      --maxpeers=             Max number of inbound and outbound peers (default: 125)
      --miningaddr=           Add the specified payment address to the list of addresses to use for generated blocks -- At least one address is required if the generate
                              option is set
      --minrelaytxfee=        The minimum transaction fee in BTC/kB to be considered a non-zero fee. (default: 1e-05)
      --nobanning             Disable banning of misbehaving peers
      --nocfilters            Disable committed filtering (CF) support
      --nocheckpoints         Disable built-in checkpoints.  Don't do this unless you know what you're doing.
      --nodnsseed             Disable DNS seeding for peers
      --nolisten              Disable listening for incoming connections -- NOTE: Listening is automatically disabled if the --connect or --proxy options are used without
                              also specifying listen interfaces via --listen
      --noonion               Disable connecting to tor hidden services
      --nopeerbloomfilters    Disable bloom filtering support
      --norelaypriority       Do not require free or low-fee transactions to have high priority for relaying
      --nowinservice          Do not start as a background service on Windows -- NOTE: This flag only works on the command line, not in the config file
      --norpc                 Disable built-in RPC server -- NOTE: The RPC server is disabled by default if no rpcuser/rpcpass or rpclimituser/rpclimitpass is specified
      --nostalldetect         Disables the stall handler system for each peer, useful in simnet/regtest integration tests frameworks
      --notls                 Disable TLS for the RPC server -- NOTE: This is only allowed if the RPC server is bound to localhost
      --onion=                Connect to tor hidden services via SOCKS5 proxy (eg. 127.0.0.1:9050)
      --onionpass=            Password for onion proxy server
      --onionuser=            Username for onion proxy server
      --profile=              Enable HTTP profiling on given port -- NOTE port must be between 1024 and 65536
      --proxy=                Connect via SOCKS5 proxy (eg. 127.0.0.1:9050)
      --proxypass=            Password for proxy server
      --proxyuser=            Username for proxy server
      --regtest               Use the regression test network
      --rejectnonstd          Reject non-standard transactions regardless of the default settings for the active network.
      --rejectreplacement     Reject transactions that attempt to replace existing transactions within the mempool through the Replace-By-Fee (RBF) signaling policy.
      --relaynonstd           Relay non-standard transactions regardless of the default settings for the active network.
      --rpccert=              File containing the certificate file (default: /home/stevekeol/.btcd/rpc.cert)
      --rpckey=               File containing the certificate key (default: /home/stevekeol/.btcd/rpc.key)
      --rpclimitpass=         Password for limited RPC connections
      --rpclimituser=         Username for limited RPC connections
      --rpclisten=            Add an interface/port to listen for RPC connections (default port: 8334, testnet: 18334)
      --rpcmaxclients=        Max number of RPC clients for standard connections (default: 10)
      --rpcmaxconcurrentreqs= Max number of concurrent RPC requests that may be processed concurrently (default: 20)
      --rpcmaxwebsockets=     Max number of RPC websocket connections (default: 25)
      --rpcquirks             Mirror some JSON-RPC quirks of Bitcoin Core -- NOTE: Discouraged unless interoperability issues need to be worked around
  -P, --rpcpass=              Password for RPC connections
  -u, --rpcuser=              Username for RPC connections
      --sigcachemaxsize=      The maximum number of entries in the signature verification cache (default: 100000)
      --simnet                Use the simulation test network
      --signet                Use the signet test network
      --signetchallenge=      Connect to a custom signet network defined by this challenge instead of using the global default signet test network -- Can be specified
                              multiple times
      --signetseednode=       Specify a seed node for the signet network instead of using the global default signet network seed nodes
      --testnet               Use the test network
      --torisolation          Enable Tor stream isolation by randomizing user credentials for each connection.
      --trickleinterval=      Minimum time between attempts to send new inventory to a connected peer (default: 10s)
      --txindex               Maintain a full hash-based transaction index which makes all transactions available via the getrawtransaction RPC
      --uacomment=            Comment to add to the user agent -- See BIP 14 for more information.
      --upnp                  Use UPnP to map our listening port outside of NAT
  -V, --version               Display version information and exit
      --whitelist=            Add an IP network or IP that will not be banned. (eg. 192.168.1.0/24 or ::1)

```