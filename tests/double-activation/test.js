import http from "k6/http";
import { check } from "k6";

export const options = {
  discardResponseBodies: true,
  thresholds: {
    // During the whole test execution none error should occur.
    http_req_failed: ["rate==0"],
  },
  scenarios: {
    doubleActivation: {
      executor: "ramping-vus",
      startVUs: 0,
      stages: [
        { duration: "5s", target: 1000 },
        { duration: "2s", target: 10000 },
        { duration: "5s", target: 5000 },
      ],
      gracefulRampDown: "0s",
    },
  },
};

const DAPR_ADDRESS = `http://127.0.0.1:${__ENV.DAPR_HTTP_PORT}/v1.0`;

function callActorMethod(id, method) {
  return http.post(
    `${DAPR_ADDRESS}/actors/fake-actor-type/${id}/method/${method}`,
    JSON.stringify({})
  );
}
export default function () {
  const result = callActorMethod("exec", "Lock");
  check(result, {
    "lock response status code is 2xx":
      result.status >= 200 && result.status < 300,
  });
}

export function teardown(_) {
  const shutdownResult = http.post(`${DAPR_ADDRESS}/shutdown`);
  check(shutdownResult, {
    "shutdown response status code is 2xx":
      shutdownResult.status >= 200 && shutdownResult.status < 300,
  });
}
