import http from 'k6/http';
import { check, sleep, group } from 'k6';
import { Rate } from 'k6/metrics';

export const options = {
  stages: [
    { duration: '30s', target: 30 },
    { duration: '2m', target: 150 },
    { duration: '1m', target: 300 },
    { duration: '30s', target: 0 },
  ],
  thresholds: {
    http_req_duration: ['p(95)<500', 'p(99)<900'],
    http_req_failed: ['rate<0.02'],
  },
};

const BASE_URL = __ENV.BASE_URL || 'http://host.docker.internal:8080';
const errorRate = new Rate('errors');

export default function () {
  let accessToken = null;
  const email = `user${__VU}_${Date.now()}@example.com`;
  const password = 'supersecret';

  group('Register', () => {
    const res = http.post(
      `${BASE_URL}/api/v1/auth/register`,
      JSON.stringify({ name: `User ${__VU}`, email, password }),
      { headers: { 'Content-Type': 'application/json' } }
    );
    if (!check(res, { 'register status 201/200': (r) => r.status === 201 || r.status === 200 })) {
      errorRate.add(1);
    }
  });

  sleep(1);

  group('Login', () => {
    const res = http.post(
      `${BASE_URL}/api/v1/auth/login`,
      JSON.stringify({ email, password }),
      { headers: { 'Content-Type': 'application/json' } }
    );
    let body = null;
    try { body = res.json(); } catch {}
    const token = body?.access_token || body?.token;

    if (check(res, { 'login status 200': (r) => r.status === 200, 'has token': () => Boolean(token) })) {
      accessToken = token;
    } else {
      errorRate.add(1);
    }
  });

  sleep(1);

  group('Get Me', () => {
    if (!accessToken) return;
    const res = http.get(`${BASE_URL}/api/v1/users/me`, {
      headers: { Authorization: `Bearer ${accessToken}` },
    });
    if (!check(res, { 'me status 200': (r) => r.status === 200 })) {
      errorRate.add(1);
    }
  });

  sleep(2);
}
