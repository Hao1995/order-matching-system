# Order Matching System

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
│   ├── api
│   │   └── order
│   │       ├── handlers
│   │       │   └── order.go
│   │       └── repositories
│   │           └── order_producer.go
│   └── worker
│       ├── matching_engine
│       │   ├── use_cases
│       │   │   ├── order_book_initializer.go
│       │   │   └── matching_engine.go
│       │   ├── repositories
│       │   │   ├── order_consumer.go
│       │   │   ├── order_repository.go
│       │   │   └── matching_producer.go
│       │   └── domains
│       │       ├── order_book.go
│       │       └── price_level.go
│       └── matching_persister
│           └── repository
│               ├── matching_repository.go
│               └── transaction.go
├── pkg
│   └── models
│       ├── events
│       │   ├── matching_event.go
│       │   └── order_event.go
│       └── order.go
├── docs
│   ├── system_architecture.md
│   └── class_uml.md
├── go.mod
├── README.md
└── docker-compose.yaml
```