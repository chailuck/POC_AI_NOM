// scripts/run_tests.sh
#!/bin/bash

echo "Running tests..."
go test -v -cover ./...

if [ $? -eq 0 ]; then
    echo "Tests passed successfully!"
    exit 0
else
    echo "Tests failed!"
    exit 1
fi