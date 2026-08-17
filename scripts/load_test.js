import http from 'k6/http';
import { check, sleep, group } from 'k6';

export const options = {
  stages: [
    { duration: '30s', target: 50 },
    { duration: '1m', target: 100 },
    { duration: '2m', target: 100 },
    { duration: '30s', target: 0 },
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'],
    http_req_failed: ['rate<0.01'],
  },
};

const BASE_URL = __ENV.BASE_URL || 'http://host.docker.internal:8080';

export default function () {
  const email = `load_${__VU}_${Date.now()}@example.com`;
  const password = 'supersecret';

  const resReg = http.post(`${BASE_URL}/api/v1/auth/register`, JSON.stringify({ name: `Load User`, email, password }), { headers: { 'Content-Type': 'application/json' }, tags: { endpoint: 'register' } });
  check(resReg, { 'register ok': (r) => r.status === 201 || r.status === 200 });

  const resLog = http.post(`${BASE_URL}/api/v1/auth/login`, JSON.stringify({ email, password }), { headers: { 'Content-Type': 'application/json' }, tags: { endpoint: 'login' } });
  check(resLog, { 'login ok': (r) => r.status === 200 });

  sleep(1);
}
