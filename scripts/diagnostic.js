import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  stages: [
    { duration: '10s', target: 50 },
    { duration: '20s', target: 150 },
    { duration: '20s', target: 300 },
    { duration: '10s', target: 0 },
  ],

  thresholds: {
    http_req_failed: ['rate<0.02'],
  },
};

const BASE_URL = __ENV.BASE_URL || 'http://host.docker.internal:8080';

export default function () {
  const email = `diagnostic_${__VU}_${Date.now()}@example.com`;

  const password = 'supersecret';

  // REGISTER
  const register = http.post(
    `${BASE_URL}/api/v1/auth/register`,
    JSON.stringify({
      name: `Diagnostic ${__VU}`,
      email,
      password,
    }),
    {
      headers: {
        'Content-Type': 'application/json',
      },
      tags: {
        endpoint: 'register',
      },
    },
  );

  check(register, {
    'register ok': (r) => r.status === 200 || r.status === 201,
  });

  // LOGIN
  const login = http.post(
    `${BASE_URL}/api/v1/auth/login`,
    JSON.stringify({
      email,
      password,
    }),
    {
      headers: {
        'Content-Type': 'application/json',
      },
      tags: {
        endpoint: 'login',
      },
    },
  );

  let token = null;

  try {
    token = login.json('access_token') || login.json('token');
  } catch {}

  check(login, {
    'login ok': (r) => r.status === 200,
  });

  // ME
  if (token) {
    const me = http.get(`${BASE_URL}/api/v1/users/me`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
      tags: {
        endpoint: 'me',
      },
    });

    check(me, {
      'me ok': (r) => r.status === 200,
    });
  }

  sleep(1);
}
