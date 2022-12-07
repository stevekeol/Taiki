# Taiki

> 👋🏻 NOTE: Taiki is a very early work-in-progress. It's currently highly unstable and not very useful as it is.

🔭 Taiki is a simple implement & optimization of Taiki protocol by TaikiLab, and is a new layer1 blockchain with high performance & scalability.

## Features

> [Why choose these features](./docs/genome.md)

The features of Taiki show below:
-  🌹 multi-chain（多链）
-  🦆 hetergeneous（异构）
-  🍓 smart-contract（支持合约）
-  💋 dynamic sharding（动态分片）
-  🍅 pos+pbft/rbft（共识）
-  ✍🏻️ tightly-coupled（紧密耦合）

## Usage

1. Generate the Taiki binary
```bash
make Taiki
```

2. Review the Help options
```bash
cd ./bin && ./Taiki -h
```

3. Other usage (just for the raw test)
```bash
cd ./cmd/taiki && cat README.md
```

## Structure

## Roadmap
1. Primitives
- ADNL
	- Address
	- P2P Protocol(UDP over ADNL)
	- C/S Protocol(TCP over ADNL)
	- RLDP
	- Channel
	- Zero Channel(support for LiteClient)
	- TDHT
		- PING
		- STORE
		- FIND_NODE
		- FIND_VALUE
- Cell&BoC
- Account
- Transaction
- Message
- Block
- Masterchain
- Shardchain

2. Core Concept
- Validator
- Collector
- Dynamic Sharding
- HR/IHR

3. Core Functionality
- Message Transfering
- Transaction Handling
- Dynamic Sharding
- PoS+PBFT/RBFT
- ChainState Management
- HR/IHR Communication

4. Surrounding
- Wallet
- LiteClient
- MobileApp(ReactNative, ton-npm)

5. Diffcult
- VM (using Go binding instead)
