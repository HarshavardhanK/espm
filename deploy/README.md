# ESPM Deployment Guide

## Local Development Setup

### Prerequisites
- Docker and Docker Compose
- Go 1.21+
- PostgreSQL client (optional, for direct database access)
- Redis CLI (optional, for cache inspection)

### Starting the Environment
```bash
cd deploy/docker
docker-compose up -d
```

This will start:
- PostgreSQL database on port 5432
- PgAdmin web interface on port 5050 (credentials: admin@espm.local / admin123)

### Verifying the Setup
```bash
# Check if containers are running
docker ps

# Test database connection
go run cmd/eventstore-test/main.go
```

## Preparing for Load Testing

Before running load tests:

1. **Expose API Endpoints**
   - Implement or mock the Command API (port 8080)
   - Implement or mock the Query API (port 8081)
   - Ensure PostgreSQL and Redis are running

2. **Configure Connection Pooling**
   - Adjust database connection pool settings based on expected load
   - Configure Redis connection settings appropriately

3. **Set Up Monitoring**
   - Use Docker stats for container monitoring
   - Set up PostgreSQL query monitoring
   - Monitor Redis cache hit/miss rates

4. **Prepare Artillery Tests**
   - Create test scenarios as outlined in `.load_test.md`
   - Adjust concurrency levels to match your local machine capacity

## Production Deployment Preparation

For production deployment:

1. **Resource Requirements**
   - Minimum 2 CPU cores per service
   - 2GB RAM per service instance
   - 10GB disk space for PostgreSQL
   - 2GB memory for Redis

2. **Security Considerations**
   - Configure proper authentication for all services
   - Use network isolation for internal services
   - Implement TLS for all API endpoints
   - Set up proper database credentials

3. **Scaling Recommendations**
   - Horizontally scale Command and Query APIs
   - Consider PostgreSQL read replicas for high-read workloads
   - Use Redis cluster for high cache throughput

4. **Monitoring Setup**
   - Implement health checks for all services
   - Set up application metrics with Prometheus
   - Configure alerting for performance degradation

The architecture is designed to be cloud-native and will work well in Kubernetes environments.