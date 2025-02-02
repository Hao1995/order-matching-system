# Order Matching System

## Introduction
Users can create or cancel orders and the matching mechanism prioritizes better prices first, followed by earlier creation times.

For example, the following orders have already created.
Orders
```
4,  sell,   $24.0,  10
3,  buy,    $23.0,  19
2,  buy,    $22.9,  10
1,  sell,   $23.1,  5
```

If we insert a sell order with $23.0 and 5 quantities.
Orders
```
5,  sell,   $23.0,  5
4,  sell,   $24.0,  10
3,  buy,    $23.0,  19
2,  buy,    $22.9,  5
1,  sell,   $23.1,  5
```

Both prices of no.2 and no.3 are matched with no.5, but the no.2 order is earlier than no.3
The sell order would be matched with $22.9 price and 5 quantities .
Orders
```
4. sell	    $24.0	10
3. buy      $23.0	19
2. buy	    $22.9	5
1. sell	    $23.1	5
```
Transactions
```
1, 2, 5, $22.9, 5
```

If we insert a buy order with $23.1 and 10 quantities.
Orders
```
6,  buy,   $23.1,  10
4,  sell,   $24.0,  10
3,  buy,    $23.0,  19
2,  buy,    $22.9,  5
1,  sell,   $23.1,  5
```

The no.1 is the first matched order, therefore the remaining quantities of the no.6 order is 5 quantities.
Orders
```
6,  buy,    $23.1,  5
4,  sell,   $24.0,  10
3,  buy,    $23.0,  19
2,  buy,    $22.9,  5
```
Transactions
```
1, 1, 6, $23.1, 5
1, 2, 5, $22.9, 5
```

## Project Layout
Project: "github.com/Hao1995/order-matching-system"

```
.
├── cmd
│   ├── api
│   │   └── order
│   │       ├── main.go
│   │       └── Dockerfile
│   └── worker
│       ├── matching_engine
│       │   ├── main.go
│       │   └── Dockerfile
│       └── matching_persister
│           ├── main.go
│           └── Dockerfile
├── internal
│   ├── common
│   │   └── models
│   │       └── events
│   ├── api
│   │   └── order
│   └── worker
│       ├── matching_engine
│       └── matching_persister
├── docs
│   ├── system_architecture.md
│   └── class_uml.md
├── go.mod
├── README.md
└── docker-compose.yaml
```

# How to use it?
```
make up
make down
```

# Kafka Test
Consume Order events, add `--group MATCHING_WORKER` for consuming messages of matching_worker
```
docker exec -it kafka kafka-console-consumer --bootstrap-server localhost:9092 --topic AAPL_ORDER
```

Consume Matching events
```
docker exec -it kafka kafka-console-consumer --bootstrap-server localhost:9092 --topic AAPL_MATCHING
```

