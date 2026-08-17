import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  stages: [
    { duration: '10m', target: 50 }, // Mantém carga constante por muito tempo
  ],
};

const BASE_URL = __ENV.BASE_URL || 'http://host.docker.internal:8080';

export default function () {
  const email = `soak_${__VU}_${Date.now()}@example.com`;
  http.post(`${BASE_URL}/api/v1/auth/register`, JSON.stringify({ name: `Soak`, email, password: 'pw' }), { headers: { 'Content-Type': 'application/json' } });
  sleep(1);
}
