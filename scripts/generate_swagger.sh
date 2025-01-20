// scripts/generate_swagger.sh
#!/bin/bash

# Install swag if not present
if ! command -v swag &> /dev/null
then
    echo "Installing swag..."
    GO111MODULE=on go get -u github.com/swaggo/swag/cmd/swag
fi

# Generate swagger documentation
echo "Generating Swagger documentation..."
swag init -g cmd/server/main.go -o api/swagger

echo "Swagger documentation generated successfully!"
