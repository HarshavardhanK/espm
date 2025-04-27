# ESPM Kubernetes Deployment

This directory contains Kubernetes manifests for deploying the Event Sourcing Platform (ESPM) and its monitoring stack.

## Directory Structure

```
kubernetes/
├── base/                    # Base Kubernetes manifests
│   ├── command-api/         # Command API service
│   ├── event-publisher/     # Event publisher service
│   ├── projections/         # Projection services
│   └── query-api/           # Query API service
├── monitoring/              # Monitoring stack
│   ├── prometheus/          # Prometheus configuration
│   ├── opentelemetry/       # OpenTelemetry collector
│   └── grafana/             # Grafana dashboards
└── overlays/                # Environment-specific overlays
    ├── development/         # Development environment
    └── production/          # Production environment
```

## Prerequisites

1. Kubernetes cluster (v1.21 or later)
2. kubectl configured to access the cluster
3. Helm 3 for deploying monitoring stack

## Deployment Steps

1. Create the namespace:
   ```bash
   kubectl create namespace espm
   ```

2. Deploy the monitoring stack:
   ```bash
   helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
   helm repo add grafana https://grafana.github.io/helm-charts
   helm repo update
   
   # Deploy Prometheus
   helm install prometheus prometheus-community/kube-prometheus-stack -n monitoring
   
   # Deploy OpenTelemetry
   helm repo add open-telemetry https://open-telemetry.github.io/opentelemetry-helm-charts
   helm install opentelemetry open-telemetry/opentelemetry-collector -n monitoring
   
   # Deploy Grafana
   helm install grafana grafana/grafana -n monitoring
   ```

3. Deploy the application:
   ```bash
   kubectl apply -k overlays/development  # For development
   # or
   kubectl apply -k overlays/production   # For production
   ```

## Monitoring Setup

### Prometheus Configuration

Prometheus is configured to scrape metrics from:
- Application services (command-api, event-publisher, etc.)
- Kubernetes components
- PostgreSQL database
- Redis cache

### OpenTelemetry Configuration

OpenTelemetry collector is configured to:
- Collect traces from all services
- Export to Jaeger for visualization
- Configure sampling rates
- Set up correlation IDs

### Grafana Dashboards

Pre-configured dashboards include:
- Service health and performance
- Event processing metrics
- Database performance
- Cache hit/miss rates
- Error rates and alerts

## Customization

1. Environment-specific configurations:
   - Edit `overlays/development` or `overlays/production`
   - Adjust resource limits and requests
   - Configure environment variables

2. Monitoring stack:
   - Modify Prometheus scrape configs
   - Adjust OpenTelemetry sampling
   - Customize Grafana dashboards

## Troubleshooting

1. Check pod status:
   ```bash
   kubectl get pods -n espm
   kubectl get pods -n monitoring
   ```

2. View logs:
   ```bash
   kubectl logs -n espm deployment/command-api
   kubectl logs -n monitoring deployment/prometheus
   ```

3. Access Grafana:
   ```bash
   kubectl port-forward -n monitoring svc/grafana 3000:80
   ```
   Then access http://localhost:3000 (default credentials: admin/admin)
