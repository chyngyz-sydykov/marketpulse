# About the project

MarketPulse (Data Collection & Analysis)
Retrieves market data from Binance API every hour for selected cryptocurrencies (e.g., BTC/USDT, ETH/USDT).
Computes key statistics:
OHLC (Open, High, Low, Close)
Moving averages (SMA, EMA)
Standard deviation & volatility metrics
Technical indicators (RSI, MACD, Bollinger Bands)
Aggregates data over different timeframes (daily, weekly, monthly).

# Installation

 - clone the repo
 - install docker
 - copy `.env.dist` to `.env`
 - run `docker-compose up --build`
 - if everything is ok, please console

# Migration
run migration `make migrate-up`

revert migration `make migrate-down`

create migration ``make migrate-new name=add_column_phone`

# GRPC

the protobuf files are stored in different repo https://github.com/chyngyz-sydykov/crypto-bot-protoc and it is imported via following command.

generate grpc files `docker exec -it marketpulse bash -c ".scripts/generate_protoc.sh"`

check if the service is registered `grpcurl -plaintext localhost:50051 list`. you should see the following in the console `rating.RatingService` 

in order to communicate with the grpc, do following

1. create a local network `docker network create crypto-bot-network`
2. after running `docker-compose up` run `docker network inspect crypto-bot-network`
You should see a json with the list of containers ex:
```
"Containers": {
            "some_hash": {
                "Name": "some name","
            },
            "some_hash": {
                "Name": "some name",
            },
            ...
},
```
# Testing

On initial project setup, please manually create a database for tests. check the database name in env.test file. to run use following commands:

run tests `APP_ENV=test go test ./tests/`

run tests without cache `go test -count=1 ./tests/`

run tests within docker (preferred way) `docker exec -it marketpulse bash -c "go test -count=1 ./tests"`

run test coverage on local machine `docker exec -it marketpulse bash "scripts/coverage.sh"`
`go tool cover -html=coverage/filtered_coverage.out`

# Handy commands

To install new package

`go get package_name`

to clean up go.sum run

`go mod tidy`

to run test

running project via docker
`docker-compose up --build`
`docker-compose down`

`docker-compose logs -f`