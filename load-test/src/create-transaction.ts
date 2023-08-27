import http, { get } from "k6/http";

const host = __ENV.API_HOST ? __ENV.API_HOST : "http://host.docker.internal";

export const options = {
  vus: 20,
  duration: "5m",
};

// list of available user-ids
const userIds = [
  "2TWlPQ2AhstX9PtJ5UTOE6xQ7Ga",
  "2TWlPVmPjhonQ2DpOFt09O990th",
  "2TWlPdYWFbP3iXPMIKVRmdZ3ozC",
  "2TWlPkjd1YcRDCBDQk11nygVDpe",
  "2TWlPw3U64nhEtzeazP5ELd7q4c",
  "2TWlQ2JWBetv9MUkFsLPd4zhLa4",
  "2TWlQ83iEHKIOtvMCfSnBwM2sEB",
  "2TWlQJsEhRIW8XpeGGz2u75phWN",
  "2TWlQRnYNr5ViFi0wLTatYImXz5",
  "2TWlQYYZVIC566T3XFVldQckPsB",
  "2TWlQcFLvqE67A2qbZpWAdSJiZL",
  "2TWlQjcMepYBrRJRRzwgXWLA0gX",
  "2TWlQvGoGcHxUa4iod28bMsfW7e",
  "2TWlR5WFe7VUVMSxhgpEmfqlAAX",
  "2TWlR8ifwogFmktTwET0Eb2s4PE",
  "2TWlRFfCwsZo903aTO7xSRCYQIU",
  "2TWlRRtrVofuy7C1ZzcWIICVEME",
  "2TWlRZBGQlbQe34dVNU3GhQKshe",
  "2TWlRedfruxmFYvYLJur7oGesXY",
  "2TWlRnR3C0hopS7NkcwyjIOq5Kd",
];

export default function () {
  // pick user id at random
  const userId = userIds[Math.floor(Math.random() * userIds.length)];
  const params = {
    headers: {
      "Content-Type": "application/json",
      "x-user-id": userId,
    },
  };

  const payload = {
    amount: Math.floor(Math.random() * 1000),
    description: `testing-${Date.now()}`,
  };

  // get balance
  const getBalanceRes = http.get(`${host}/balance`, params);
  const out = getBalanceRes.json();

  if (out.payload.amount < payload.amount) {
    // do topup
    const topupPayload = {
      user_id: userId,
      description: `reason #${Date.now()}`,
      amount: Math.floor(Math.random() * 10) * 10000
    }
    http.post(`${host}/balance/add`, JSON.stringify(topupPayload), params)
  }

  // create transaction
  http.post(`${host}/transaction`, JSON.stringify(payload), params);
}
