// scripts/run_dev.ps1
# PowerShell script for running the application in development mode
Write-Host "Starting development server..."

# Start the application using air for hot reload
air

// .air.toml
root = "."
tmp_dir = "tmp"

[build]
  cmd = "go build -o ./tmp/main ./cmd/server"
  bin = "tmp/main"
  delay = 1000
  exclude_dir = ["assets", "tmp", "vendor", "test"]
  include_ext = ["go", "tpl", "tmpl", "html"]
  exclude_regex = ["_test.go"]

[screen]
  clear_on_rebuild = true
