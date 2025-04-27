#!/bin/bash

#Comprehensive script for setting up the ESPM stack in a Kubernetes cluster

#Exit on error
set -e

#Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

#Function to print status messages
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

#Check prerequisites
check_prerequisites() {
    print_status "Checking prerequisites..."
    
    #Check kubectl
    if ! command -v kubectl &> /dev/null; then
        print_error "kubectl is not installed. Please install kubectl first."
        exit 1
    fi
    
    #Check helm
    if ! command -v helm &> /dev/null; then
        print_error "helm is not installed. Please install helm first."
        exit 1
    fi
    
    #Check Kubernetes cluster
    if ! kubectl cluster-info &> /dev/null; then
        print_error "Cannot connect to Kubernetes cluster. Please configure kubectl first."
        exit 1
    fi
    
    print_status "All prerequisites met."
}

#Create namespaces
create_namespaces() {
    print_status "Creating namespaces..."
    
    #Create espm namespace
    if ! kubectl get namespace espm &> /dev/null; then
        kubectl create namespace espm
        print_status "Created espm namespace"
    else
        print_warning "espm namespace already exists"
    fi
    
    #Create monitoring namespace
    if ! kubectl get namespace monitoring &> /dev/null; then
        kubectl create namespace monitoring
        print_status "Created monitoring namespace"
    else
        print_warning "monitoring namespace already exists"
    fi
}

#Deploy monitoring stack
deploy_monitoring() {
    print_status "Deploying monitoring stack..."
    
    #Add Helm repositories
    print_status "Adding Helm repositories..."
    helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
    helm repo add grafana https://grafana.github.io/helm-charts
    helm repo add open-telemetry https://open-telemetry.github.io/opentelemetry-helm-charts
    helm repo update
    
    #Deploy Prometheus
    if ! helm list -n monitoring | grep -q prometheus; then
        print_status "Deploying Prometheus..."
        helm install prometheus prometheus-community/kube-prometheus-stack \
            -n monitoring \
            --set prometheus.prometheusSpec.serviceMonitorSelectorNilUsesHelmValues=false \
            --set prometheus.prometheusSpec.podMonitorSelectorNilUsesHelmValues=false
    else
        print_warning "Prometheus is already deployed"
    fi
    
    #Deploy OpenTelemetry
    if ! helm list -n monitoring | grep -q opentelemetry; then
        print_status "Deploying OpenTelemetry..."
        helm install opentelemetry open-telemetry/opentelemetry-collector \
            -n monitoring \
            -f monitoring/opentelemetry/values.yaml
    else
        print_warning "OpenTelemetry is already deployed"
    fi
    
    #Deploy Grafana
    if ! helm list -n monitoring | grep -q grafana; then
        print_status "Deploying Grafana..."
        helm install grafana grafana/grafana \
            -n monitoring \
            --set persistence.enabled=true \
            --set persistence.size=10Gi \
            --set adminPassword=admin \
            --set service.type=LoadBalancer
    else
        print_warning "Grafana is already deployed"
    fi
}

#Deploy application
deploy_application() {
    print_status "Deploying application..."
    
    #Apply base configurations
    print_status "Applying base configurations..."
    kubectl apply -k base/
    
    #Apply environment-specific configurations
    if [ "$1" == "production" ]; then
        print_status "Applying production configurations..."
        kubectl apply -k overlays/production/
    else
        print_status "Applying development configurations..."
        kubectl apply -k overlays/development/
    fi
}

#Wait for deployments
wait_for_deployments() {
    print_status "Waiting for deployments to be ready..."
    
    #Wait for monitoring stack
    kubectl wait --for=condition=available --timeout=300s deployment/prometheus-operator -n monitoring
    kubectl wait --for=condition=available --timeout=300s deployment/opentelemetry-collector -n monitoring
    kubectl wait --for=condition=available --timeout=300s deployment/grafana -n monitoring
    
    #Wait for application deployments
    kubectl wait --for=condition=available --timeout=300s deployment/command-api -n espm
    kubectl wait --for=condition=available --timeout=300s deployment/event-publisher -n espm
    kubectl wait --for=condition=available --timeout=300s deployment/projections -n espm
    kubectl wait --for=condition=available --timeout=300s deployment/query-api -n espm
}

#Print access information
print_access_info() {
    print_status "Deployment completed successfully!"
    echo ""
    echo "Access Information:"
    echo "------------------"
    
    #Get Grafana URL
    GRAFANA_URL=$(kubectl get svc grafana -n monitoring -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
    if [ -z "$GRAFANA_URL" ]; then
        GRAFANA_URL=$(kubectl get svc grafana -n monitoring -o jsonpath='{.status.loadBalancer.ingress[0].hostname}')
    fi
    if [ -z "$GRAFANA_URL" ]; then
        print_warning "Grafana URL not available. You may need to use port-forwarding:"
        echo "kubectl port-forward -n monitoring svc/grafana 3000:80"
        echo "Access Grafana at http://localhost:3000 (admin/admin)"
    else
        echo "Grafana: http://${GRAFANA_URL}:3000 (admin/admin)"
    fi
    
    #Get Prometheus URL
    PROMETHEUS_URL=$(kubectl get svc prometheus-kube-prometheus-prometheus -n monitoring -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
    if [ -z "$PROMETHEUS_URL" ]; then
        PROMETHEUS_URL=$(kubectl get svc prometheus-kube-prometheus-prometheus -n monitoring -o jsonpath='{.status.loadBalancer.ingress[0].hostname}')
    fi
    if [ -z "$PROMETHEUS_URL" ]; then
        print_warning "Prometheus URL not available. You may need to use port-forwarding:"
        echo "kubectl port-forward -n monitoring svc/prometheus-kube-prometheus-prometheus 9090:9090"
        echo "Access Prometheus at http://localhost:9090"
    else
        echo "Prometheus: http://${PROMETHEUS_URL}:9090"
    fi
    
    #Get Jaeger URL
    JAEGER_URL=$(kubectl get svc jaeger-query -n monitoring -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
    if [ -z "$JAEGER_URL" ]; then
        JAEGER_URL=$(kubectl get svc jaeger-query -n monitoring -o jsonpath='{.status.loadBalancer.ingress[0].hostname}')
    fi
    if [ -z "$JAEGER_URL" ]; then
        print_warning "Jaeger URL not available. You may need to use port-forwarding:"
        echo "kubectl port-forward -n monitoring svc/jaeger-query 16686:16686"
        echo "Access Jaeger at http://localhost:16686"
    else
        echo "Jaeger: http://${JAEGER_URL}:16686"
    fi
}

#Main deployment process
main() {
    print_status "Starting ESPM deployment..."
    
    #Check environment (default to development)
    ENV=${1:-development}
    if [ "$ENV" != "development" ] && [ "$ENV" != "production" ]; then
        print_error "Invalid environment. Use 'development' or 'production'"
        exit 1
    fi
    
    check_prerequisites
    create_namespaces
    deploy_monitoring
    deploy_application "$ENV"
    wait_for_deployments
    print_access_info
}

# Execute main function
main "$@" 