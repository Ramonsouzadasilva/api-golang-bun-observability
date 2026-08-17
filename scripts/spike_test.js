import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  stages: [
    { duration: '10s', target: 500 }, // Rápido aumento para 500
    { duration: '30s', target: 500 }, // Mantém o pico
    { duration: '10s', target: 0 },   // Cai rapidamente
  ],
};

const BASE_URL = __ENV.BASE_URL || 'http://host.docker.internal:8080';

export default function () {
  const email = `spike_${__VU}_${Date.now()}@example.com`;
  const password = 'password123';

  http.post(`${BASE_URL}/api/v1/auth/register`, JSON.stringify({ name: `Spike`, email, password }), { headers: { 'Content-Type': 'application/json' } });
  http.post(`${BASE_URL}/api/v1/auth/login`, JSON.stringify({ email, password }), { headers: { 'Content-Type': 'application/json' } });
  sleep(1);
}
