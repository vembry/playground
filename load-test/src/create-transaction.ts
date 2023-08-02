import http from "k6/http";

export default function () {
  console.log(Date.now().toLocaleString());
  http.post("http://0.0.0.0:8080/transaction");
}
