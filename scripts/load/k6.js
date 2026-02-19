import http from "k6/http";
import { sleep } from "k6";

export const options = {
  stages: [
    { duration: "10s", target: 50 },
    { duration: "20s", target: 200 },
    { duration: "20s", target: 200 },
    { duration: "10s", target: 0 },
  ],
};

const base = __ENV.BASE_URL || "http://localhost:8080";

export default function () {

  const r = Math.random();
  if (r < 0.50) {
    http.get(`${base}/cpu?n=60000`);
  } else if (r < 0.80) {
    http.get(`${base}/io?delay=200`);
  } else if (r < 0.95) {
    http.get(`${base}/mem?mb=50&hold=100`);
  } else {
    http.get(`${base}/chaos?rate=0.3&mode=http500`);
  }
  sleep(0.05);
}
