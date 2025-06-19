# E-commerce Payment Management Microservice

This is a Go microservice implementing clean architecture principles for managing e-commerce payments with Stripe.

## Architecture

The service follows clean architecture principles with the following layers:

- Domain Layer: Core business logic and entities
- Use Case Layer: Application business rules
- Interface Layer: External interfaces (HTTP, gRPC)
- Infrastructure Layer: External services and implementations

## Project Structure

```
.
├── cmd/
│   └── api/            # Application entry point
├── internal/
│   ├── domain/         # Enterprise business rules
│   ├── usecase/        # Application business rules
│   ├── infra/          # IO Infrasctructure (Redis, Kafka, Mongo)
│   ├── lib/            # Third party package wrappers (Stripe, JWT)
│   ├── worker/         # Background workers
│   └── delivery/       # Delivery mechanisms
└── go.mod             # Go module file
```

## Getting Started

1. Install dependencies:

```bash
go mod download
```

2. Run the service:

```bash
go run cmd/api/main.go
```

## API Endpoints

- POST /api/payment/start - Trigger Stripe checkout session
- POST /api/payment/confirm - Confirm Stripe checkout session
