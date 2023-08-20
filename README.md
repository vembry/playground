# credits-management

## context
1. This is an *experiment* of "How to handle credits ledger"
2. credit = balance = available funds

## scenario
1. setup starting balance at 10K
2. create multiple transaction entries that will exhaust the balance
   1. we're gonna use `grafana/k6` to replicate huge intakes
   2. entries will consist of both credit addition and deduction.

## expected outcome
1. remaining credits shouldnt be negative
2. final balance on ledger should always be `in` > `out`
2. nothing went brokey

## learning points
1. learn how to make credit-ledger doesnt go kaboom when theres huge intake of transactions
2. learn how to utilize docker for service orchestration(?)
3. learn how to load test using `grafana/k6`
4. visualize the observability metrics from the experiment, using:
   1. grafana, to visualize
   2. cadvisor, to scrap metrics from the containers
   3. prometheus, to get metrics from `api` and `api-worker`

## setup
Run the following command which will start everything
> make start

## links
- [grafana](http://localhost:10000)
   - [dashboard](http://localhost:10000/d/f744f256-9d10-4e29-9dd0-c4ea6f8d7dd2/dashboard?orgId=1&var-container_name=api-worker&var-container_name=api&from=now-5m&to=now&refresh=5s) in works