import http from "k6/http";

export const options = {
  vus: 10,
  duration: "10s",
};

export default function () {
  const payload = JSON.stringify({
    amount: 10,
    description: `testing-${Date.now()}`,
  });
  const params = {
    headers: {
      "Content-Type": "application/json",
      "x-user-id": "2TWlPQ2AhstX9PtJ5UTOE6xQ7Ga",
    },
  };
  http.post("http://host.docker.internal:80/transaction", payload, params);
}
