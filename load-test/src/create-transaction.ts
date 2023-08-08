import http from "k6/http";

export const options = {
  vus: 10,
  duration: "10s",
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
  const payload = JSON.stringify({
    amount: 10,
    description: `testing-${Date.now()}`,
  });

  // pick user id at random
  const userId = userIds[Math.floor(Math.random() * userIds.length)];

  const params = {
    headers: {
      "Content-Type": "application/json",
      "x-user-id": userId,
    },
  };
  http.post("http://host.docker.internal:80/transaction", payload, params);
}
