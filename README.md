# Taiki

> 👋🏻 NOTE: Taiki is a very early work-in-progress. It's currently highly unstable and not very useful as it is.

🔭 Taiki is a simple implement & optimization of TON protocol by TaikiLab, and is a new layer2 blockchain with high performance & scalability.

## Genome

In order to build a highly scalable blockchain system, 
- for functionality scalable, it should support `hetergeneous` workchains (maybe having diferent rules);
- for performance scalable, it should support `dynamic sharding` (Splitting and Merging Shardchains when the load overweight or subside)
- when using sharding, becase of the delay confirm, it make sense only if the system is `tigtly-coupled`
- when using tight-coupled, it means there should be a `masterchain`, `IHR`(Instant Hypercube routing), `PoS+BFT(PBFT/RBFT)`, and so on.

So, The features of Taiki show below:
-  🌹 multi-chain
-  🦆 hetergeneous
-  🍓 smart-contract
-  💋 dynamic sharding
-  🍅 pos+pbft/rbft
-  ✍🏻️ tightly-coupled

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