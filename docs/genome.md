# Genome

In order to build a highly scalable blockchain system, 
- for functionality scalable, it should support `hetergeneous` workchains (maybe having diferent rules);
- for performance scalable, it should support `dynamic sharding` (Splitting and Merging Shardchains when the load overweight or subside)
- when using sharding, becase of the delay confirm, it make sense only if the system is `tigtly-coupled`
- when using tight-coupled, it means there should be a `masterchain`, `IHR`(Instant Hypercube routing), `PoS+BFT(PBFT/RBFT)`, and so on.