// README.md additions
## Local Development

### Prerequisites
- Go 1.21 or later
- Docker Desktop
- PowerShell

### Setting Up Development Environment

1. Run the development setup script:
```powershell
.\scripts\dev_setup.ps1
```

2. Start the development server:
```powershell
.\scripts\run_dev.ps1
```

### Running Tests

To run all tests (unit, integration, and performance):
```powershell
.\scripts\test_all.ps1
```

### API Documentation

Generate Swagger documentation:
```powershell
make generate-docs
```

Access the Swagger UI at: http://localhost:8080/swagger/index.html

### Testing with curl

1. Create an individual:
```powershell
$headers = @{
    "Content-Type" = "application/json"
}

$body = @{
    id = "PATY000001"
    title = "Mr."
    givenName = "John"
    familyName = "Doe"
    maritalStatus = "S"
    gender = "M"
    nationality = "USA"
    contactMedium = @(
        @{
            type = "PhoneContactMedium"
            mediumType = "homePhoneNumber"
            preferred = $true
            phoneNumber = "1234567890"
        }
    )
} | ConvertTo-Json

Invoke-RestMethod -Method Post -Uri "http://localhost:8080/tmf-api/partyManagement/v4/individual" -Headers $headers -Body $body
```

2. Get an individual:
```powershell
Invoke-RestMethod -Method Get -Uri "http://localhost:8080/tmf-api/partyManagement/v4/individual/PATY000001"
```

3. Update an individual:
```powershell
$updateBody = @{
    title = "Dr."
    givenName = "John"
    familyName = "Doe"
} | ConvertTo-Json

Invoke-RestMethod -Method Put -Uri "http://localhost:8080/tmf-api/partyManagement/v4/individual/PATY000001" -Headers $headers -Body $updateBody
```

4. Delete an individual:
```powershell
Invoke-RestMethod -Method Delete -Uri "http://localhost:8080/tmf-api/partyManagement/v4/individual/PATY000001"
```

### Monitoring

Access monitoring dashboards:
- Prometheus: http://localhost:9090
- Grafana: http://localhost:3000

Default Grafana credentials:
- Username: admin
- Password: admin

### Troubleshooting

1. If the application fails to start:
   - Check if PostgreSQL is running: `docker ps`
   - Check application logs: `docker compose logs app`
   - Verify database connection: `docker compose exec db psql -U postgres -d tmf632db -c "\l"`

2. If tests fail:
   - Check test logs: `go test -v ./... 2>&1 | Tee-Object -FilePath test.log`
   - Verify database connection for integration tests
   - Check if required services are running

3. If Swagger documentation is not accessible:
   - Regenerate documentation: `make generate-docs`
   - Verify server is running
   - Check application logs for errors


# Complete Deployment Guide for TMF632 Service

## Windows Deployment with Minikube

### 1. Prerequisites Installation

```powershell
# Install Chocolatey if not already installed
Set-ExecutionPolicy Bypass -Scope Process -Force
[System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072
iex ((New-Object System.Net.WebClient).DownloadString('https://chocolatey.org/install.ps1'))

# Install required tools
choco install -y kubernetes-cli
choco install -y minikube
choco install -y docker-desktop
```

### 2. Environment Setup

```powershell
# Start Minikube
minikube start --cpus=2 --memory=4096mb --driver=hyperv

# Enable required addons
minikube addons enable ingress
minikube addons enable metrics-server

# Configure Docker to use Minikube's Docker daemon
& minikube -p minikube docker-env --shell powershell | Invoke-Expression
```

### 3. Database Setup

```powershell
# Apply database secrets and config
kubectl apply -f deployment/k8s/secret.yaml
kubectl apply -f deployment/k8s/configmap.yaml

# Deploy PostgreSQL
kubectl apply -f deployment/k8s/postgresql.yaml

# Wait for PostgreSQL to be ready
kubectl wait --for=condition=ready pod -l app=tmf632-postgresql --timeout=300s
```

### 4. Application Deployment

```powershell
# Build the application image
docker build -t tmf632-service:latest .

# Deploy the application
kubectl apply -f deployment/k8s/deployment.yaml
kubectl apply -f deployment/k8s/service.yaml
kubectl apply -f deployment/k8s/ingress.yaml

# Wait for the application to be ready
kubectl wait --for=condition=ready pod -l app=tmf632-service --timeout=300s
```

### 5. Monitoring Setup

```powershell
# Deploy Prometheus and Grafana
kubectl apply -f deployment/k8s/prometheus.yaml
kubectl apply -f deployment/k8s/grafana.yaml

# Create monitoring dashboards
kubectl apply -f deployment/k8s/grafana-dashboards-configmap.yaml
```

## Testing the Deployment

### 1. Get Service URL

```powershell
$SERVICE_URL = minikube service tmf632-service --url
```

### 2. Test API Endpoints

```powershell
# Variables for testing
$HEADERS = @{
    "Content-Type" = "application/json"
}

# Create Individual
$CREATE_BODY = @{
    id = "TEST001"
    title = "Mr."
    givenName = "Test"
    familyName = "User"
    maritalStatus = "S"
    gender = "M"
    nationality = "USA"
    contactMedium = @(
        @{
            type = "PhoneContactMedium"
            mediumType = "homePhoneNumber"
            preferred = $true
            phoneNumber = "1234567890"
        }
    )
} | ConvertTo-Json

$RESPONSE = Invoke-RestMethod `
    -Method Post `
    -Uri "$SERVICE_URL/tmf-api/partyManagement/v4/individual" `
    -Headers $HEADERS `
    -Body $CREATE_BODY

Write-Host "Created Individual: $($RESPONSE.id)"

# Get Individual
$GET_RESPONSE = Invoke-RestMethod `
    -Method Get `
    -Uri "$SERVICE_URL/tmf-api/partyManagement/v4/individual/$($RESPONSE.id)"

Write-Host "Retrieved Individual: $($GET_RESPONSE.givenName) $($GET_RESPONSE.familyName)"
```

## Monitoring and Maintenance

### 1. Access Monitoring Dashboards

```powershell
# Get Prometheus URL
Write-Host "Prometheus URL: $(minikube service prometheus --url)"

# Get Grafana URL
Write-Host "Grafana URL: $(minikube service grafana --url)"
```

### 2. View Application Logs

```powershell
# Get pod name
$POD_NAME = kubectl get pods -l app=tmf632-service -o jsonpath="{.items[0].metadata.name}"

# View logs
kubectl logs $POD_NAME
```

### 3. Check System Status

```powershell
# Check all resources
kubectl get all

# Check specific component status
kubectl describe deployment tmf632-service
kubectl describe service tmf632-service
kubectl describe ingress tmf632-ingress
```

## Common Issues and Solutions

### 1. Database Connection Issues

If the application can't connect to the database:

```powershell
# Check database pod status
kubectl get pods -l app=tmf632-postgresql

# Check database logs
kubectl logs -l app=tmf632-postgresql

# Verify database secrets
kubectl describe secret tmf632-secret

# Test database connection from application pod
$APP_POD = kubectl get pods -l app=tmf632-service -o jsonpath="{.items[0].metadata.name}"
kubectl exec -it $APP_POD -- sh -c "nc -zv tmf632-postgresql 5432"
```

### 2. Service Unavailable

If the service is not accessible:

```powershell
# Check if pods are running
kubectl get pods

# Check service endpoints
kubectl get endpoints tmf632-service

# Verify ingress configuration
kubectl describe ingress tmf632-ingress

# Check service logs
kubectl logs -l app=tmf632-service
```

### 3. Performance Issues

If the service is running slowly:

```powershell
# Check resource usage
kubectl top pods
kubectl top nodes

# View detailed metrics in Prometheus/Grafana
Write-Host "Access Grafana at: $(minikube service grafana --url)"
```

## Cleanup

### 1. Remove Application

```powershell
# Delete all resources
kubectl delete -f deployment/k8s/

# Verify resources are removed
kubectl get all
```

### 2. Stop Minikube

```powershell
# Stop Minikube
minikube stop

# Optional: Delete Minikube cluster
minikube delete
```

## Additional Resources

### API Documentation
- Swagger UI: `$SERVICE_URL/swagger/index.html`
- Postman Collection: `test/api/postman/TMF632_Collection.json`

### Monitoring
- Prometheus Metrics: `$SERVICE_URL/metrics`
- Grafana Dashboards: See `deployment/docker/grafana/dashboards/`

### Source Code
- Main Application: `cmd/server/main.go`
- Database Models: `internal/models/`
- API Handlers: `internal/handlers/`
- Configuration: `internal/config/`

### Testing
- Unit Tests: `test/`
- Performance Tests: `test/performance/k6/`
- Integration Tests: `test/integration/`

# Performance Testing and Monitoring Guide

## Performance Testing with k6

### 1. Install k6

```powershell
# Install k6 using Chocolatey
choco install k6
```

### 2. Run Performance Tests

```powershell
# Get service URL
$SERVICE_URL = minikube service tmf632-service --url

# Run basic performance test
k6 run ./test/performance/k6/scenarios.js

# Run with custom configuration
k6 run --vus 10 --duration 30s ./test/performance/k6/scenarios.js

# Run with environment variables
$env:BASE_URL="$SERVICE_URL"; k6 run ./test/performance/k6/scenarios.js
```

### 3. Performance Test Scenarios

The test suite includes:

1. Base Load Test:
```powershell
k6 run ./test/performance/k6/base_load.js
```
- Simulates normal usage patterns
- Ramps up to 10 concurrent users
- Runs for 5 minutes

2. Stress Test:
```powershell
k6 run ./test/performance/k6/stress_test.js
```
- Tests system under high load
- Ramps up to 50 concurrent users
- Runs for 10 minutes

3. Soak Test:
```powershell
k6 run ./test/performance/k6/soak_test.js
```
- Tests system stability over time
- Maintains 20 concurrent users
- Runs for 1 hour

## Monitoring Setup

### 1. Deploy Monitoring Stack

```powershell
# Apply monitoring configurations
kubectl apply -f deployment/k8s/prometheus.yaml
kubectl apply -f deployment/k8s/grafana.yaml

# Wait for pods to be ready
kubectl wait --for=condition=ready pod -l app=prometheus --timeout=300s
kubectl wait --for=condition=ready pod -l app=grafana --timeout=300s
```

### 2. Access Monitoring Dashboards

```powershell
# Get Prometheus URL
$PROMETHEUS_URL = minikube service prometheus --url
Write-Host "Prometheus URL: $PROMETHEUS_URL"

# Get Grafana URL
$GRAFANA_URL = minikube service grafana --url
Write-Host "Grafana URL: $GRAFANA_URL"
```

### 3. Configure Grafana

1. Login to Grafana:
   - URL: $GRAFANA_URL
   - Default credentials:
     - Username: admin
     - Password: admin

2. Add Prometheus Data Source:
   - Configuration > Data Sources > Add data source
   - Select Prometheus
   - URL: http://prometheus:9090
   - Click "Save & Test"

3. Import Dashboards:
   - Create > Import
   - Upload JSON file: `deployment/docker/grafana/dashboards/tmf632.json`

## Metrics and Alerts

### Key Metrics Monitored

1. Application Metrics:
   - Request rate
   - Response time (p50, p95, p99)
   - Error rate
   - Success rate by endpoint

2. Database Metrics:
   - Query response time
   - Connection pool utilization
   - Active connections
   - Query cache hit rate

3. System Metrics:
   - CPU usage
   - Memory usage
   - Network I/O
   - Disk I/O

### Alert Rules

1. High Response Time:
```yaml
- alert: HighResponseTime
  expr: histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m])) > 0.5
  for: 5m
  labels:
    severity: warning
  annotations:
    summary: High response time (instance {{ $labels.instance }})
    description: 95th percentile of response time is above 500ms
```

2. High Error Rate:
```yaml
- alert: HighErrorRate
  expr: rate(http_requests_total{status=~"5.."}[5m]) / rate(http_requests_total[5m]) > 0.05
  for: 5m
  labels:
    severity: critical
  annotations:
    summary: High error rate (instance {{ $labels.instance }})
    description: Error rate is above 5% for more than 5 minutes
```

3. Database Connection Issues:
```yaml
- alert: DatabaseConnectionIssues
  expr: pg_stat_activity_count{state="active"} > 100
  for: 5m
  labels:
    severity: warning
  annotations:
    summary: High number of active database connections
    description: More than 100 active database connections
```

## Performance Optimization Tips

1. Database Optimization:
```sql
-- Add indexes for common queries
CREATE INDEX idx_individual_search ON individuals (given_name, family_name);

-- Analyze table statistics
ANALYZE individuals;

-- Monitor slow queries
SELECT * FROM pg_stat_activity WHERE state = 'active' AND query_start < now() - interval '1 minute';
```

2. Application Configuration:
```yaml
# Optimize connection pool
DB_MAX_OPEN_CONNS: 25
DB_MAX_IDLE_CONNS: 5
DB_CONN_MAX_LIFETIME: 5m

# Enable response compression
ENABLE_COMPRESSION: true
COMPRESSION_LEVEL: 5

# Cache configuration
CACHE_TTL: 5m
CACHE_SIZE: 1000
```

3. Kubernetes Resources:
```yaml
resources:
  requests:
    cpu: 200m
    memory: 256Mi
  limits:
    cpu: 500m
    memory: 512Mi
```

## Troubleshooting Performance Issues

1. High Response Time:
```powershell
# Check application logs for slow requests
kubectl logs -l app=tmf632-service | Select-String "duration"

# Monitor database performance
kubectl exec -it $(kubectl get pod -l app=tmf632-postgresql -o jsonpath='{.items[0].metadata.name}') -- psql -U postgres -c "SELECT * FROM pg_stat_activity;"
```

2. Memory Issues:
```powershell
# Check memory usage
kubectl top pods

# View detailed memory metrics
kubectl exec -it $(kubectl get pod -l app=tmf632-service -o jsonpath='{.items[0].metadata.name}') -- go tool pprof -http=:8080 /debug/pprof/heap
```

3. CPU Issues:
```powershell
# Check CPU usage
kubectl top pods --containers

# Profile CPU usage
kubectl exec -it $(kubectl get pod -l app=tmf632-service -o jsonpath='{.items[0].metadata.name}') -- go tool pprof -http=:8080 /debug/pprof/profile
```

## Regular Maintenance Tasks

1. Database Maintenance:
```sql
-- Vacuum database
VACUUM ANALYZE;

-- Update statistics
ANALYZE individuals;
ANALYZE contact_media;
ANALYZE external_references;
```

2. Log Rotation:
```yaml
# logrotate configuration
/var/log/tmf632/*.log {
    daily
    rotate 7
    compress
    delaycompress
    missingok
    notifempty
    create 0640 tmf632 tmf632
}
```

3. Backup Strategy:
```powershell
# Backup database
kubectl exec -it $(kubectl get pod -l app=tmf632-postgresql -o jsonpath='{.items[0].metadata.name}') -- pg_dump -U postgres tmf632db > backup.sql

# Restore database
kubectl exec -it $(kubectl get pod -l app=tmf632-postgresql -o jsonpath='{.items[0].metadata.name}') -- psql -U postgres tmf632db < backup.sql
```
