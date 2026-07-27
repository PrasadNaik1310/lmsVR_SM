import http from 'k6/http';
import { check } from 'k6';

export const options = {
  stages: [
    { duration: '2m', target: 50 },   // Warm up
    { duration: '2m', target: 100 },  // Normal load
    { duration: '2m', target: 200 },  // Heavy load
    { duration: '2m', target: 300 },  // Stress
    { duration: '2m', target: 500 },  // Breaking point
    { duration: '1m', target: 0 },    // Cool down
  ],

  thresholds: {
    http_req_failed: ['rate<0.01'], // <1% failures
    http_req_duration: ['p(95)<1000'], // 95% under 1 second
  },
};

export default function () {
  const payload = JSON.stringify({
    useremail: 'prasad@lms.com',
    password: 'prasad',
  });

  const params = {
    headers: {
      'Content-Type': 'application/json',
    },
    timeout: '30s',
  };

  const res = http.post(
    'https://lmsvr-sm.onrender.com/lms/auth/login',
    payload,
    params
  );

  check(res, {
    'status is 200': (r) => r.status === 200,
  });
}