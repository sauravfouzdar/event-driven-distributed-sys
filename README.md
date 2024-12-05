# event-driven-distributed-sys
A massively scalable ecommerce application which follow event-driven-distributed-microservice architecture


```
ecommerce-app/
├── README.md
├── go.mod
├── go.sum
├── cmd/                     # Command-line applications for services
│   ├── orders/              # Order microservice entry point
│   │   └── main.go
│   ├── payments/            # Payment microservice entry point
│   │   └── main.go
│   ├── delivery/            # Delivery microservice entry point
│       └── main.go
├── pkg/                     # Shared reusable packages
│   ├── events/              # Event messaging abstractions
│   │   ├── producer.go
│   │   ├── consumer.go
│   │   └── broker.go        # Kafka or NATS JetStream integration
│   ├── models/              # Common domain models
│   │   ├── order.go
│   │   ├── payment.go
│   │   └── delivery.go
│   ├── logger/              # Centralized logging utilities
│   │   └── logger.go
│   └── config/              # Configuration utilities
│       └── config.go
├── services/                # Business logic for individual services
│   ├── orders/
│   │   ├── handler.go       # API handlers
│   │   ├── service.go       # Business logic
│   │   ├── repository.go    # Database interactions
│   │   └── events.go        # Event processing logic
│   ├── payments/
│   │   ├── handler.go
│   │   ├── service.go
│   │   ├── repository.go
│   │   └── events.go
│   └── delivery/
│       ├── handler.go
│       ├── service.go
│       ├── repository.go
│       └── events.go
├── docs/                    # API documentation
├── deployments/             # Deployment configurations (Docker, Kubernetes)
│   ├── docker-compose.yml
│   └── k8s/
│       ├── orders-deployment.yml
│       ├── payments-deployment.yml
│       └── delivery-deployment.yml
└── tests/                   # Integration and unit tests
    ├── orders_test.go
    ├── payments_test.go
    └── delivery_test.go

```
