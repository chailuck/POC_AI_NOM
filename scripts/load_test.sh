// scripts/load_test.sh
#!/bin/bash

echo "Running load tests..."

# Install k6 if not present
if ! command -v k6 &> /dev/null
then
    echo "Installing k6..."
    curl -L https://github.com/loadimpact/k6/releases/download/v0.33.0/k6-v0.33.0-linux-amd64.tar.gz -o k6.tar.gz
    tar xzf k6.tar.gz
    sudo mv k6-v0.33.0-linux-amd64/k6 /usr/local/bin/
    rm -rf k6-v0.33.0-linux-amd64 k6.tar.gz
fi

# Run load test
k6 run test/load/scenarios.js

// test/load/scenarios.js
import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
    stages: [
        { duration: '30s', target: 20 },
        { duration: '1m', target: 20 },
        { duration: '30s', target: 0 },
    ],
};

const BASE_URL = 'http://localhost:8080/tmf-api/partyManagement/v4';

export default function() {
    let response = http.get(`${BASE_URL}/individual`);
    check(response, {
        'is status 200': (r) => r.status === 200,
    });
    sleep(1);
}