import { Client, StatusOK } from 'k6/net/grpc';
import { check, } from 'k6';


export const options = {
    vus: 1,
    duration: "1m"
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
        // queue name will be autogenerated as queue_{N}
        // N will be 0 until 9
        const number = Math.floor(Math.random() * 10);
        const queue = `queue_${number}`

        // construct param
        const data = {
            queue_name: queue,
            payload: `${Date.now()}`, // mock payloads
        }

        // call rpc
        const response = client.invoke('broker.Broker/Enqueue', data)

        // handle response
        check(response, {
            'status is OK': (r) => r && r.status === StatusOK
        })
    } catch (e) {
        console.log("!catching error")
        console.log(e.message)
    }
    client.close()
}
