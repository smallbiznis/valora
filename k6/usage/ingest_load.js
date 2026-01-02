import http from 'k6/http';
import { BASE_URL } from '../config/env.js';
import { apiHeaders, requireApiKey } from '../helpers/auth.js';
import { usagePayload, ensureUsageEnv } from '../helpers/payloads.js';
import { assert2xx } from '../helpers/checks.js';

const vus = parseInt(__ENV.INGEST_VUS || '20', 10);
const duration = __ENV.INGEST_DURATION || '5m';

export const options = {
  scenarios: {
    ingest_load: {
      executor: 'constant-vus',
      vus,
      duration,
    },
  },
  thresholds: {
    'http_req_duration{name:usage_ingest_load}': ['p(95)<300'],
    'http_req_failed{name:usage_ingest_load}': ['rate<0.01'],
  },
};

export function setup() {
  requireApiKey();
  ensureUsageEnv();
}

export default function () {
  const payload = usagePayload();
  const res = http.post(`${BASE_URL}/api/usage`, payload, {
    headers: apiHeaders(),
    tags: { name: 'usage_ingest_load' },
  });

  assert2xx(res);
}
