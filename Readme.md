# Event Sourcing Platform for Microservices (ESPM)

A cloud-native distributed event sourcing and CQRS platform built with Go and Kubernetes.

## Overview

ESPM is a production-ready implementation of event sourcing patterns for cloud-native environments. It provides a complete platform for building event-driven systems with separate write and read paths, full audit capability, and scalable projection rebuilding.

![Architecture Overview](docs/images/architecture-overview.png)

## Features

- **Event-Driven Architecture** - Complete implementation of the event sourcing pattern
- **CQRS** - Separate optimized read and write paths
- **Cloud-Native** - Designed for Kubernetes with StatefulSets, ConfigMaps, and HPAs
- **Complete Audit Trail** - Every state change is recorded as an immutable event
- **Scalable Projections** - Efficient read models with 40% faster rebuild capability
- **Observability** - Full OpenTelemetry, Prometheus, and Grafana integration
- **Resilience** - System can be rebuilt from event history

## Components

- **Command API** - Validates and processes commands, converting them to events
- **Event Store** - Immutable log of all events using PostgreSQL with JSONB
- **Event Publisher** - Distributes events to subscribers
- **Projection Services** - Build specialized read models from event streams
- **Query API** - Provides optimized read access to system state

## Getting Started

### Prerequisites

- Go 1.21+
- Docker and Docker Compose
- Kubernetes cluster (local or remote)
- kubectl and Helm

### Quick Start

1. Clone the repository
   ```
   git clone https://github.com/HarshavardhanK/espm.git
   cd espm
   ```

2. Start local development environment
   ```
   make dev-up
   ```

3. Run the sample application
   ```
   make run-example
   ```

4. Deploy to Kubernetes
   ```
   make deploy
   ```

## Documentation

- [Architecture Guide](docs/architecture.md)
- [Domain Modeling](docs/domain-modeling.md)
- [Kubernetes Deployment](docs/kubernetes.md)
- [Observability Setup](docs/observability.md)

## Use Cases

ESPM is ideal for:

- Systems requiring complete audit trails
- Applications with complex domain logic
- Microservices with different read and write requirements
- Systems needing point-in-time recovery capability
- Services requiring high scalability for reads

## Contributing

Contributions are welcome! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.