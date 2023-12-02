# credits-management

## context
1. This is an *experiment* of "How to manage credits"
2. credit = balance = available funds

## scenario
1. prepared x amount as the transacting balance
2. check whether current active balance is enough
3. withdraw money

## expected outcome
1. remaining credits shouldnt be negative
2. accumulated balance on ledger should always be `in` > `out`
2. nothing went brokey

## learning points
1. learn how to make credits doesnt go kaboom when theres huge transaction intake
2. learn how to utilize docker for service orchestration
3. learn how to load test using `grafana/k6`, which we'll be using to run the scenario + replicate huge intake process
4. visualize the observability metrics from the experiment, using:
   1. grafana, to visualize
   2. cadvisor, to scrap metrics from the containers
   3. prometheus, to get metrics from `app-go` and `app-go-worker`

## setup
Run the following command which will start everything
> make start

## links
- [grafana](http://localhost:10000)
   - [dashboard](http://localhost:10000/d/f744f256-9d10-4e29-9dd0-c4ea6f8d7dd2/dashboard?orgId=1&var-container_name=app-go-worker&var-container_name=app-go&from=now-5m&to=now&refresh=5s) in works