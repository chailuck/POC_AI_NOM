// scripts/setup_dev.sh
#!/bin/bash

echo "Setting up development environment..."

# Install dependencies
go mod download

# Create test database
PGPASSWORD=password psql -h localhost -U postgres -c "CREATE DATABASE tmf632db_test;"

# Run migrations
migrate -path db/migrations -database "postgresql://postgres:password@localhost:5432/tmf632db_test?sslmode=disable" up

echo "Development environment setup complete!"
