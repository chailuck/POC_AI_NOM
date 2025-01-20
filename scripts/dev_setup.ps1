// scripts/dev_setup.ps1
# PowerShell script for local development setup
Write-Host "Setting up local development environment..."

# Check if Go is installed
if (!(Get-Command go -ErrorAction SilentlyContinue)) {
    Write-Host "Installing Go..."
    choco install golang
}

# Check if Docker Desktop is installed
if (!(Get-Command docker -ErrorAction SilentlyContinue)) {
    Write-Host "Installing Docker Desktop..."
    choco install docker-desktop
}

# Install development tools
Write-Host "Installing development tools..."
go install github.com/cosmtrek/air@latest
go install github.com/swaggo/swag/cmd/swag@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Create development database
Write-Host "Setting up development database..."
docker compose up -d db

# Install dependencies
Write-Host "Installing project dependencies..."
go mod download
go mod tidy

Write-Host "Development environment setup complete!"
