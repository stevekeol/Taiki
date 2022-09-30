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
./Taiki createblockchain -address=1Lp5qqtq6xH148tLMKZkzyRDwjxSPoAEea
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