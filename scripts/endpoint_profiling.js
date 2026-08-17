import http from 'k6/http';
import { check, sleep, group } from 'k6';

export const options = {
  vus: 50,
  duration: '30s',
};

const BASE_URL = __ENV.BASE_URL || 'http://host.docker.internal:8080';

export default function () {
  const email = `prof_${__VU}_${Date.now()}@example.com`;
  const password = 'supersecret';

  group('Profile Register', () => {
    http.post(`${BASE_URL}/api/v1/auth/register`, JSON.stringify({ name: `Prof`, email, password }), { headers: { 'Content-Type': 'application/json' }, tags: { endpoint: 'profile_register' } });
  });

  let token = null;
  group('Profile Login', () => {
    const res = http.post(`${BASE_URL}/api/v1/auth/login`, JSON.stringify({ email, password }), { headers: { 'Content-Type': 'application/json' }, tags: { endpoint: 'profile_login' } });
    try { token = res.json('access_token'); } catch {}
  });

  group('Profile Me', () => {
    if (token) {
      http.get(`${BASE_URL}/api/v1/users/me`, { headers: { Authorization: `Bearer ${token}` }, tags: { endpoint: 'profile_me' } });
    }
  });

  sleep(1);
}
