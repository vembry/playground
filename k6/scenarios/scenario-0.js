import http from "k6/http";

export const options = {
  vus: 1,
  // duration: "1m",

  iterations: 10,
};

const host = __ENV.API_HOST
  ? __ENV.API_HOST
  : "http://host.docker.internal:8080";

  
// list of available user-ids
const balanceIds = [
  "2TeSprhp2cN6nEIcayZsjjvnlsK",
  "2TeSppzLaJaxldlzMOkqYO37vqw",
  "2TeSprCLB0tF6HsJT9eCb1IBht0",
  "2TeSps1ECSPx2IRrWMKEd6oHvSJ",
  "2TeSpqL5Cq3t0rNJ4RDcR4ky029",
  "2TeSpo6Okj6BedA62PUO3nRYADm",
  "2TeSprwUG0oLfgGbIIObipgf5be",
  "2TeSpqi0vPnYFf82cdBCE4Gaxjx",
  "2TeSppFijPxu2bIKpRLhLy3ye4j",
  "2TeSpnNHCn0Yq0vJddsv0dcSviy",
  "2TeSpsBbTYE5ZBe3zpO1HDRf44o",
  // "2TeSprmZUc7HGpqlm7bUi6RB9lu",
  // "2TeSpqKkgW9pz15PaqR8ZIS0Wjh",
  // "2TeSpqfWmsmNLffPB4fhErgYRiC",
  // "2TeSpnH3d9Hr9zPyNpoBrzdd8fx",
  // "2TeSpsz6jeis3vJHCoSfFrnGPVl",
  // "2TeSptEF7KEwZ9cdlowOH4HSmbR",
  // "2TeSpsvjaufyDq78BwcFFZ6ibfi",
  // "2TeSpsFLvMYhPU3CpC21wgyOcS5",
  // "2TeSpmhKLpqgNXTNOvh3JPLkbjy",
  // "2TWlPQ2AhstX9PtJ5UTOE6xQ7Ga",
  // "2TWlPVmPjhonQ2DpOFt09O990th",
  // "2TWlPdYWFbP3iXPMIKVRmdZ3ozC",
  // "2TWlPkjd1YcRDCBDQk11nygVDpe",
  // "2TWlPw3U64nhEtzeazP5ELd7q4c",
  // "2TWlQ2JWBetv9MUkFsLPd4zhLa4",
  // "2TWlQ83iEHKIOtvMCfSnBwM2sEB",
  // "2TWlQJsEhRIW8XpeGGz2u75phWN",
  // "2TWlQRnYNr5ViFi0wLTatYImXz5",
  // "2TWlQYYZVIC566T3XFVldQckPsB",
  // "2TWlQcFLvqE67A2qbZpWAdSJiZL",
  // "2TWlQjcMepYBrRJRRzwgXWLA0gX",
  // "2TWlQvGoGcHxUa4iod28bMsfW7e",
  // "2TWlR5WFe7VUVMSxhgpEmfqlAAX",
  // "2TWlR8ifwogFmktTwET0Eb2s4PE",
  // "2TWlRFfCwsZo903aTO7xSRCYQIU",
  // "2TWlRRtrVofuy7C1ZzcWIICVEME",
  // "2TWlRZBGQlbQe34dVNU3GhQKshe",
  // "2TWlRedfruxmFYvYLJur7oGesXY",
  // "2TWlRnR3C0hopS7NkcwyjIOq5Kd",
];

const requestParams = {
  headers: {
    "Content-Type": "application/json",
  },
};

export default function () {
  // pick balance id at random
  const balanceId = balanceIds[Math.floor(Math.random() * balanceIds.length)];

  // get balance
  const getBalanceRes = http.get(`${host}/balance/${balanceId}`);
  const out = getBalanceRes.json();

  const trxAmount = Math.floor(Math.random() * 1000);

  if (out["object"]["amount"] < trxAmount) {
    const topupAmount = Math.floor(Math.random() * 10) * 10000 + 1000;

    // do topup
    const topupPayload = {
      amount: topupAmount,
    };
    http.post(
      `${host}/balance/${balanceId}/deposit`,
      JSON.stringify(topupPayload),
      requestParams
    );
  }

  // let balanceIdTo = balanceIds[Math.floor(Math.random() * balanceIds.length)];
  // while (balanceIdTo == balanceId) {
  //   balanceIdTo = balanceIds[Math.floor(Math.random() * balanceIds.length)];
  // }

  // create withdrawal
  http.post(
    `${host}/balance/${balanceId}/withdraw`,
    JSON.stringify({
      amount: trxAmount,
    }),
    requestParams
  );
}
