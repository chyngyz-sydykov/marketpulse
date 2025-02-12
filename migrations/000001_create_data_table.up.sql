CREATE TABLE IF NOT EXISTS data_btc (
    id BIGSERIAL,
    symbol TEXT NOT NULL,       -- e.g., BTC, ETH
    timeframe TEXT NOT NULL,    -- e.g., 1h, 8h, 1d, 4d, 1w, 1m
    timestamp TIMESTAMPTZ NOT NULL,
    open NUMERIC NOT NULL,
    high NUMERIC NOT NULL,
    low NUMERIC NOT NULL,
    close NUMERIC NOT NULL,
    volume NUMERIC NOT NULL,
    PRIMARY KEY (id, timeframe)
) PARTITION BY LIST (timeframe);


CREATE TABLE IF NOT EXISTS data_eth (
    id BIGSERIAL,
    symbol TEXT NOT NULL,       -- e.g., BTC, ETH
    timeframe TEXT NOT NULL,    -- e.g., 1h, 8h, 1d, 4d, 1w, 1m
    timestamp TIMESTAMPTZ NOT NULL,
    open NUMERIC NOT NULL,
    high NUMERIC NOT NULL,
    low NUMERIC NOT NULL,
    close NUMERIC NOT NULL,
    volume NUMERIC NOT NULL,
    PRIMARY KEY (id, timeframe)
) PARTITION BY LIST (timeframe);



CREATE TABLE IF NOT EXISTS data_sol (
    id BIGSERIAL,
    symbol TEXT NOT NULL,       -- e.g., BTC, ETH
    timeframe TEXT NOT NULL,    -- e.g., 1h, 8h, 1d, 4d, 1w, 1m
    timestamp TIMESTAMPTZ NOT NULL,
    open NUMERIC NOT NULL,
    high NUMERIC NOT NULL,
    low NUMERIC NOT NULL,
    close NUMERIC NOT NULL,
    volume NUMERIC NOT NULL,
    PRIMARY KEY (id, timeframe)
) PARTITION BY LIST (timeframe);


CREATE TABLE IF NOT EXISTS data_bnb (
    id BIGSERIAL,
    symbol TEXT NOT NULL,       -- e.g., BTC, ETH
    timeframe TEXT NOT NULL,    -- e.g., 1h, 8h, 1d, 4d, 1w, 1m
    timestamp TIMESTAMPTZ NOT NULL,
    open NUMERIC NOT NULL,
    high NUMERIC NOT NULL,
    low NUMERIC NOT NULL,
    close NUMERIC NOT NULL,
    volume NUMERIC NOT NULL,
    PRIMARY KEY (id, timeframe)
) PARTITION BY LIST (timeframe);


CREATE TABLE IF NOT EXISTS data_trump (
    id BIGSERIAL,
    symbol TEXT NOT NULL,       -- e.g., BTC, ETH
    timeframe TEXT NOT NULL,    -- e.g., 1h, 8h, 1d, 4d, 1w, 1m
    timestamp TIMESTAMPTZ NOT NULL,
    open NUMERIC NOT NULL,
    high NUMERIC NOT NULL,
    low NUMERIC NOT NULL,
    close NUMERIC NOT NULL,
    volume NUMERIC NOT NULL,
    PRIMARY KEY (id, timeframe)
) PARTITION BY LIST (timeframe);