# credits-management

## context
1. credit = balance = available funds
2. This is an experiment of "How to handle credits ledger"

## scenario
1. setup starting balance at X
2. create multiple transaction entries that will exhaust the balance
   1. we're gonna use `grafana/k6` to replicate huge intakes
   2. entries will consist of both credit addition and deduction at random, but the general idea is to exhaust the funds

## expected outcome
1. remaining credits shouldnt be negative
2. nothing went brokey

## learning points
1. learn how to make credit ledger not goes kaboom when theres huge intake of transaction entries that will be add/deduct available credit/balance/funds (or whatever you wanna call it :D)
2. learn how to utilize docker for service orchestration(?)
3. learn how to load test using `grafana/k6`
4. visualize the observability metrics from the experiment


## setup

### cluster
> make start

### load test
`<work in progress>`