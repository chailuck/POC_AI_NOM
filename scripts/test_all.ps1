// scripts/test_all.ps1
# PowerShell script for running all tests
Write-Host "Running all tests..."

# Run unit tests
Write-Host "Running unit tests..."
go test -v ./...

# Run integration tests
Write-Host "Running integration tests..."
go test -v -tags=integration ./test/integration/...

# Run performance tests
Write-Host "Running performance tests..."
k6 run test/performance/k6/scenarios.js

Write-Host "All tests completed!"