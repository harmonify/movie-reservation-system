---

excalidraw-plugin: parsed
tags: [excalidraw]

---
==⚠  Switch to EXCALIDRAW VIEW in the MORE OPTIONS menu of this document. ⚠== You can decompress Drawing data with the command palette: 'Decompress current Excalidraw file'. For more info check in plugin settings under 'Saving'


# Excalidraw Data

## Text Elements
Consistency-Availability-Partition Tolerance
(CAP) Theorem ^Pqks1HwA

CAP theorem (also known as Brewer's theorem) states that in a distributed system, it is impossible to simultaneously guarantee all three of the following properties:

Consistency (C):
Every read receives the most recent write (or an error if the write isn't completed yet).
All replicas of data are in sync, and clients always see the same data regardless of which node they query.

Availability (A):
Every request (read or write) receives a response, even if some nodes are down or unreachable.
The system remains operational, and users never experience downtime.

Partition Tolerance (P):
The system continues to operate despite network partitions (where nodes cannot communicate with each other due to network failures). ^c071ZlQx

partition tolerance is non-negotiable in distributed system because network error is inevitable ^VVKWqbgJ

Partition tolerance ^Puh13LFk

Node Failure: One or more nodes in the cluster crash or go offline.
Network Failure: Communication links between nodes are disrupted (e.g., due to hardware issues or high latency). ^HALYeBNz

CP Systems
(Consistency + Partition Tolerance) ^Ayjl4KOW

AP Systems
(Availability + Partition Tolerance) ^MaJAkA3x

eventual consistency ^AO4Phr9C

Consistency model
https://jepsen.io/consistency ^NuZ29Oq1

Prioritize consistency during partitions but sacrifice availability.
Behavior: If nodes can't communicate, they won't process requests to avoid inconsistency.
Use Case: Banking systems, where it’s critical to maintain a consistent view of account balances. ^PTJ9k8SI

Prioritize availability during partitions but may return stale or inconsistent data.
Behavior: The system serves requests even if some nodes are out of sync.
Use Case: Social media feeds or IoT systems, where availability is more critical than strict consistency. ^gJizST11

Examples ^a6nvHYEP

Imagine a distributed e-commerce system with multiple nodes handling user requests.

Scenario: Network Partition
The network splits into two partitions:
Partition A: Contains Node 1 and Node 2.
Partition B: Contains Node 3 and Node 4.
CP System:

If a user adds an item to their cart in Partition A, the system will block the operation until it confirms with Partition B.
Guarantees consistency but sacrifices availability for some users.
AP System:

The user can add an item to their cart in Partition A, even if Partition B is unreachable.
Availability is preserved, but the cart in Partition B may temporarily show outdated data.
CA System:

The system can't handle the partition and stops processing requests until connectivity is restored.
Neither consistency nor availability is maintained during the partition. ^zTLpeu7s

Which to choose? ^7y05KDlO

Choose CP when correctness is critical.
Choose AP when availability is essential, and temporary inconsistencies are acceptable.
Never choose CA: (Almost) every system is a distributed system
-- Chas Emerick ^baW8ih6X

%%
## Drawing
```compressed-json
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
```
%%