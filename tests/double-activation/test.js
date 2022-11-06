import http from "k6/http";
import { check } from "k6";

const singleId = 2;
export const options = {
  stages: [
    { target: 200, duration: "30s" },
    { target: 0, duration: "30s" },
    { target: 10000, duration: "1s" },
  ],
};

export default function () {
  const result = http.get(
    `http://127.0.0.1:3500/v1.0/actors/fake-actor-type/${singleId}/method/Lock`
  );
  check(result, {
    "lock response status code is 2xx":
      result.status >= 200 && result.status < 300,
  });
}

export function teardown(_) {
  const result = http.post("http://127.0.0.1:3500/v1.0/shutdown");
  check(result, {
    "shutdown response status code is 2xx":
      result.status >= 200 && result.status < 300,
  });
}
