import { Client, StatusOK } from 'k6/net/grpc';
import { check, } from 'k6';


export const options = {
    vus: 1,
    duration: "1m",
};

const host = __ENV.BROKER_HOST
    ? __ENV.BROKER_HOST
    : "http://host.docker.internal:4000";


const client = new Client();
client.load([], 'broker.proto')

export default function () {
    client.connect(host, {
        plaintext: true
    })

    try {
        const data = { /** payloads */ }
        const response = client.invoke('broker.Broker/GetQueue', data)

        check(response, {
            'status is OK': (r) => r && r.status === StatusOK
        })
    } catch (e) {
        console.log("caught error")
        console.log(e)
    }
    client.close()
}