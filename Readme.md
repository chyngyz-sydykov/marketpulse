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
run migration `docker run --rm -v $(pwd)/migrations:/migrations --network=host migrate/migrate -path=/migrations -database "postgres://postgres:password@localhost:5402/marketpulse_db?sslmode=disable" up`

revert migration `docker run --rm -v $(pwd)/migrations:/migrations --network=host migrate/migrate -path=/migrations -database "postgres://postgres:password@localhost:5402/marketpulse_db?sslmode=disable" down 1`

force migration `docker run --rm -v $(pwd)/migrations:/migrations --network=host migrate/migrate -path=/migrations -database "postgres://postgres:password@localhost:5402/marketpulse_db?sslmode=disable" force 1`


create migration `docker run --rm -v $(pwd)/migrations:/migrations --network=host migrate/migrate create -ext sql -dir /migrations -seq add_column_phone`

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