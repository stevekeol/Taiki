# README

## Usage

```bash
# 1. 进入合适目录并编译
cd ./Taiki/cmd/taiki && go build -o Taiki

# 2. 测试是否编译成功
./Taiki

# Usage:
#  createblockchain --address ADDRESS 创建一条链并且该地址会得到狗头金
#  createwallet                      创建一个钱包，里面放着一对秘钥
#  getbalance       --address ADDRESS 得到该地址的余额
#  listaddresses                     罗列钱包中所有的地址
#  printchain                        打印链
#  reindexutxo                       重构UTXO集合
#  send             --from FROM --to TO --amount AMOUNT 地址from发送amount的币给地址to

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
send             --from FROM --to TO --amount AMOUNT 

# 8. 查询账户余额
./Taiki getbalance       --address ADDRESS

# 9. 打印整个链
./Taiki printchain
```

## TODO
- 参考aptos的命令行工具输出

```bash
stevekeol@linux:~/Code/BlockChain-Projects/Aptos/aptos-cli-0.3.4-Ubuntu-x86_64$ ./aptos -h
aptos 0.3.4
Aptos Labs <opensource@aptoslabs.com>
Command Line Interface (CLI) for developing and interacting with the Aptos blockchain

USAGE:
    aptos <SUBCOMMAND>

OPTIONS:
    -h, --help       Print help information
    -V, --version    Print version information

SUBCOMMANDS:
    account       Tool for interacting with accounts
    config        Tool for interacting with configuration of the Aptos CLI tool
    genesis       Tool for setting up an Aptos chain Genesis transaction
    governance    Tool for on-chain governance
    help          Print this message or the help of the given subcommand(s)
    info          Show build information about the CLI
    init          Tool to initialize current directory for the aptos tool
    key           Tool for generating, inspecting, and interacting with keys
    move          Tool for Move related operations
    node          Tool for operations related to nodes
    stake         Tool for manipulating stake
stevekeol@linux:~/Code/BlockChain-Projects/Aptos/aptos-cli-0.3.4-Ubuntu-x86_64$ ./aptos stake -h
aptos-stake 0.3.4
Tool for manipulating stake

USAGE:
    aptos stake <SUBCOMMAND>

OPTIONS:
    -h, --help       Print help information
    -V, --version    Print version information

SUBCOMMANDS:
    add-stake                 Stake coins to the stake pool
    help                      Print this message or the help of the given subcommand(s)
    increase-lockup           Increase lockup of all staked coins in the stake pool
    initialize-stake-owner    Initialize stake owner
    set-delegated-voter       Delegate voting capability from the stake owner to another account
    set-operator              Delegate operator capability from the stake owner to another
                                  account
    unlock-stake              Unlock staked coins
    withdraw-stake            Withdraw unlocked staked coins

```
