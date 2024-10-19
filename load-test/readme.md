# load-test
is a load/stress test executor. 

## `./sdk-go/loadtest` 
- contain the prototype of stress tester.
- has open-telemetry + otel related libraries integrated. 
    - it will publish otel traces to `monitoring-otel-collector`
    - currently it's pushes prometheus metrics yet, this need to work

## app tests

> go run ./load-test app

### scenario
1. prepared `x` amount as the transacting balance
2. call `get balance` to, check whether current active balance is enough
    1. when balance is not enough, then call `topup balance`
3. call `withdraw money`

### expected
1. remaining credits shouldnt be negative
2. accumulated balance on ledger should always be `in` > `out`
2. nothing went brokey

## broker tests

> go run ./load-test broker

### scenario
run `./broker` whole lifecycle

### expected
no crash