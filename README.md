# Ecommerce Microservices Platform

A Go-based e-commerce platform built as a set of microservices with an API Gateway, shared infrastructure packages, and service-specific domains. The system is designed for modularity, scalability, and independent deployment of core business capabilities.

## Table of Contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Services](#services)
- [Shared Packages](#shared-packages)
- [Service-to-Service Communication](#service-to-service-communication)
- [Data & Infrastructure](#data--infrastructure)
- [Repository Structure](#repository-structure)
- [Configuration](#configuration)
- [Local Development](#local-development)
- [Kubernetes Deployment](#kubernetes-deployment)
- [Security: Gateway-Only Access](#security-gateway-only-access)
- [Operational Notes](#operational-notes)
- [Documentation Index](#documentation-index)

## Overview

This repository contains a multi-service e-commerce system. Each service encapsulates a bounded context and exposes APIs through the API Gateway. Shared libraries provide common infrastructure (logging, tracing, JWT, Redis, database helpers).

## Architecture

- **API Gateway** handles client traffic, authentication, routing, and cross-cutting concerns.
- **Domain Services** implement business capabilities (User, Product, Cart, Order).
- **Shared Libraries** in [pkg](pkg) standardize infrastructure and utilities.
- **Shared Proto** in [shared/proto](shared/proto) supports gRPC contracts.

High-level flow:

1. Client calls API Gateway
2. Gateway validates/authenticates and routes to target service
3. Services use shared packages for persistence, messaging, caching, and observability

## Services

- **ApiGateway**: Entry point for clients; routing, auth, and middleware.
- **UserService**: User registration, authentication, profiles.
- **ProductService**: Catalog, inventory, and product queries; caching support.
- **CartService**: Shopping cart lifecycle and persistence.
- **OrderService**: Order creation, status tracking, and payment workflow hooks.

## Shared Packages

Located in [pkg](pkg):

- **db**: Database initialization and helpers
- **jwt**: Token creation and validation
- **logger**: Structured logging with GORM integration
- **password**: Hashing and verification
- **rabbitmq**: Messaging utilities
- **redis**: Caching client helpers
- **tracer**: Tracing and error instrumentation

## Service-to-Service Communication

- gRPC contracts are stored in [shared/proto](shared/proto).
- The API Gateway routes external traffic to internal services.

## Data & Infrastructure

- **Databases**: Each service can own its data store.
- **Caching**: Redis helpers for fast reads and session data.
- **Tracing & Logging**: Shared tracing and logger packages.

## Repository Structure

```
.
├─ pkg/                 # shared infrastructure packages
├─ services/            # microservices
│  ├─ ApiGateway/
│  ├─ UserService/
│  ├─ ProductService/
│  ├─ CartService/
│  └─ OrderService/
├─ shared/              # shared protobuf definitions
├─ docker-compose.yaml  # top-level orchestration
└─ Makefile             # root build helpers
```

## Configuration

Each service has its own config in:

- [services/ApiGateway/config](services/ApiGateway/config)
- [services/UserService/config](services/UserService/config)
- [services/ProductService/config](services/ProductService/config)
- [services/CartService/config](services/CartService/config)
- [services/OrderService/config](services/OrderService/config)

Refer to each service’s README and config files for environment variables and defaults.

## Local Development

- Use the root [docker-compose.yaml](docker-compose.yaml) to bring up infrastructure.
- Each service includes its own Docker files under its [services/\*/docker](services) folder.
- Service binaries are defined in [services/\*/cmd](services) and can be run independently.

## Kubernetes Deployment

See [k8s/README.md](k8s/README.md) for production-oriented Kubernetes manifests and run steps.

## Security: Gateway-Only Access

Production best practice is to combine **network isolation** with **service-to-service authentication**.

### 1) Network isolation (primary control)

Only the API Gateway is exposed to the host. Internal services and databases do not publish host ports.

Where this is enforced:

- [services/ApiGateway/docker/docker-compose.yml](services/ApiGateway/docker/docker-compose.yml)
- [services/UserService/docker/docker-compose.yml](services/UserService/docker/docker-compose.yml)
- [services/ProductService/docker/docker-compose.yml](services/ProductService/docker/docker-compose.yml)
- [services/CartService/docker/docker-compose.yml](services/CartService/docker/docker-compose.yml)
- [services/OrderService/docker/docker-compose.yml](services/OrderService/docker/docker-compose.yml)

### 2) Internal token auth (defense in depth)

All gRPC servers require an internal token, and the API Gateway (and service-to-service clients) attach it to outbound calls.

Set the same token for all services and the gateway:

- Environment variable: INTERNAL_AUTH_TOKEN
- Configuration files:
  - [services/ApiGateway/config/config.go](services/ApiGateway/config/config.go)
  - [services/UserService/config/config.go](services/UserService/config/config.go)
  - [services/ProductService/config/config.go](services/ProductService/config/config.go)
  - [services/CartService/config/config.go](services/CartService/config/config.go)
  - [services/OrderService/config/config.go](services/OrderService/config/config.go)

Implementation details:

- Server interceptor: [pkg/grpcmiddleware/internal_auth.go](pkg/grpcmiddleware/internal_auth.go)
- Gateway client interceptor: [services/ApiGateway/internal/clients/grpc_clients.go](services/ApiGateway/internal/clients/grpc_clients.go)

### 3) mTLS (recommended for production)

For strong identity and encryption between services, use mTLS via a service mesh (e.g., Istio or Linkerd).
This complements the internal token check and prevents traffic spoofing.

## Operational Notes

- Logs are stored under [logs](logs) when enabled.
- Migrations exist under [services/\*/internal/migrations](services) where applicable.
- Temporary artifacts are under [tmp](tmp) and [services/\*/tmp](services).

## Documentation Index

- [AUTHORIZATION_LAYER_GUIDE.md](AUTHORIZATION_LAYER_GUIDE.md)
- [SERVICES_FLOW.md](SERVICES_FLOW.md)
- [k8s/README.md](k8s/README.md)
- API Gateway docs:
  - [services/ApiGateway/README.md](services/ApiGateway/README.md)
  - [services/ApiGateway/API_GATEWAY_GUIDE.md](services/ApiGateway/API_GATEWAY_GUIDE.md)
  - [services/ApiGateway/GRACEFUL_SHUTDOWN_README.md](services/ApiGateway/GRACEFUL_SHUTDOWN_README.md)
  - [services/ApiGateway/IMPLEMENTATION_SUMMARY.md](services/ApiGateway/IMPLEMENTATION_SUMMARY.md)
- Order Service docs:
  - [services/OrderService/README.md](services/OrderService/README.md)
