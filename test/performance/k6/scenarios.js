// test/performance/k6/scenarios.js
import { check, sleep } from 'k6';
import http from 'k6/http';
import { config } from './config.js';

export let options = {
    stages: config.stages,
    thresholds: config.thresholds,
};

const baseUrl = config.baseUrl;

export function setup() {
    // Create test data
    const individual = {
        id: 'PERF_TEST_001',
        title: 'Mr.',
        givenName: 'Performance',
        familyName: 'Test',
        maritalStatus: 'S',
        gender: 'M',
        nationality: 'USA',
    };

    const res = http.post(
        `${baseUrl}/tmf-api/partyManagement/v4/individual`,
        JSON.stringify(individual),
        { headers: { 'Content-Type': 'application/json' } }
    );

    check(res, {
        'setup created individual': (r) => r.status === 201,
    });

    return { individualId: individual.id };
}

export default function(data) {
    // Get individual
    const getRes = http.get(
        `${baseUrl}/tmf-api/partyManagement/v4/individual/${data.individualId}`
    );

    check(getRes, {
        'get individual status is 200': (r) => r.status === 200,
        'get individual has correct id': (r) => r.json('id') === data.individualId,
    });

    sleep(1);

    // List individuals
    const listRes = http.get(
        `${baseUrl}/tmf-api/partyManagement/v4/individual`
    );

    check(listRes, {
        'list individuals status is 200': (r) => r.status === 200,
        'list individuals returns array': (r) => Array.isArray(r.json()),
    });

    sleep(1);
}

export function teardown(data) {
    // Cleanup test data
    const delRes = http.del(
        `${baseUrl}/tmf-api/partyManagement/v4/individual/${data.individualId}`
    );

    check(delRes, {
        'teardown deleted individual': (r) => r.status === 204,
    });
}
