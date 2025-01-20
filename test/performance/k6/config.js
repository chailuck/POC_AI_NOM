// test/performance/k6/config.js
export const config = {
    baseUrl: __ENV.BASE_URL || 'http://localhost:8080',
    stages: [
        { duration: '1m', target: 10 },  // Ramp up to 10 users
        { duration: '3m', target: 10 },  // Stay at 10 users
        { duration: '1m', target: 20 },  // Ramp up to 20 users
        { duration: '3m', target: 20 },  // Stay at 20 users
        { duration: '1m', target: 0 },   // Ramp down to 0 users
    ],
    thresholds: {
        http_req_duration: ['p(95)<500'], // 95% of requests should be below 500ms
        http_req_failed: ['rate<0.01'],   // Less than 1% of requests should fail
    },
};
