// scripts/minikube_setup.ps1
# PowerShell script for Windows setup
Write-Host "Setting up Minikube environment..."

# Check if Minikube is installed
if (!(Get-Command minikube -ErrorAction SilentlyContinue)) {
    Write-Host "Installing Minikube..."
    choco install minikube
}

# Start Minikube
minikube start --cpus 2 --memory 4096

# Enable required addons
minikube addons enable ingress
minikube addons enable metrics-server

# Set Docker environment to use Minikube's Docker daemon
Write-Host "Configuring Docker environment..."
minikube docker-env | Invoke-Expression

# Build the Docker image
Write-Host "Building Docker image..."
docker build -t tmf632-service:latest .

# Apply Kubernetes configurations
Write-Host "Applying Kubernetes configurations..."
kubectl apply -f deployment/k8s/

# Wait for pods to be ready
Write-Host "Waiting for pods to be ready..."
kubectl wait --for=condition=ready pod -l app=tmf632-service --timeout=300s

Write-Host "Setup complete! Use 'minikube service tmf632-service --url' to get the service URL"
