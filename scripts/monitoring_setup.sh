// scripts/monitoring_setup.sh
#!/bin/bash

echo "Setting up monitoring stack..."

# Install Prometheus operator
kubectl apply -f deployment/k8s/prometheus.yaml

# Install Grafana
kubectl apply -f deployment/k8s/grafana.yaml

# Wait for deployments to be ready
kubectl wait --for=condition=ready pod -l app=prometheus --timeout=300s
kubectl wait --for=condition=ready pod -l app=grafana --timeout=300s

echo "Monitoring stack setup complete!"