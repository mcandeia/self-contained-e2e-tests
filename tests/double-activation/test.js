import http from "k6/http";
import { check } from "k6";
import exec from "k6/execution";

export const options = {
  stages: [{ target: 1, duration: "1s" }],
};
const DAPR_ADDRESS = "http://127.0.0.1:3500/v1.0";

function callActorMethod(id, method) {
  return http.post(
    `${DAPR_ADDRESS}/actors/fake-actor-type/${id}/method/${method}`,
    JSON.stringify({})
  );
}
export default function () {
  const result = callActorMethod(exec.scenario.iterationInTest, "Inc");
  check(result, {
    "lock response status code is 2xx":
      result.status >= 200 && result.status < 300,
  });
}

export function teardown(_) {
  const getCurrentValue = callActorMethod("teardown", "Get");
  console.log(getCurrentValue);

  check(getCurrentValue, {
    "inc should be equal to total iterations completed":
      +getCurrentValue.body > 0,
  });

  const shutdownResult = http.post(`${DAPR_ADDRESS}/v1.0/shutdown`);
  check(shutdownResult, {
    "shutdown response status code is 2xx":
      shutdownResult.status >= 200 && shutdownResult.status < 300,
  });
}
