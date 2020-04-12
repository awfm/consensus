# Alvalor Consensus

## Description

Alvalor consensus implements a byzantine fault-tolerant (BFT) consensus algorithm, based on HotStuff.

## Roadmap

### Milestone 1 - PoC (v0.0.1)

- event-driven consensus logic
- no liveness mechanism for leader
- no crypto-economic incentive scheme
- no cryptographic primitives
- no verification against graph

### Milestone 2 - MVP (v0.0.2)

- implement buffer component
- implement chain component
- add state extension check
- add finalization of vertices

### Milestone 3 - Cryptography (v0.0.3)

- implement signature component
- implement verification component
- add identity set for participants
- add signature creation & checking

### Milestone X - Randomness (v0.0.4)

- add entropy to vertices
- add threshold signatures

### Milestone X - Liveness (v0.1.0)

- add depth concept to vertices
- add timeout mechanism for leader

### Milestone X - Incentives (v0.2.0)

- add native economic token ledger
- add transaction fee distribution
- add slashing challenges

### Milestone X - Committee (v0.3.0)

- add checkpoints / epochs
- add staking / unstaking
- add validator token auctions
- add validator token buybacks

### Miscellaneous (to be decided)

- add stake delegation
