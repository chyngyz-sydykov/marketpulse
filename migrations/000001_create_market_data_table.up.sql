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

CREATE INDEX idx_data_btc_timeframe ON data_btc (timeframe);
CREATE INDEX idx_data_btc_timestamp ON data_btc (timestamp DESC);

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

CREATE INDEX idx_data_eth_timeframe ON data_eth (timeframe);
CREATE INDEX idx_data_eth_timestamp ON data_eth (timestamp DESC);


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

CREATE INDEX idx_data_sol_timeframe ON data_sol (timeframe);
CREATE INDEX idx_data_sol_timestamp ON data_sol (timestamp DESC);

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

CREATE INDEX idx_data_bnb_timeframe ON data_bnb (timeframe);
CREATE INDEX idx_data_bnb_timestamp ON data_bnb (timestamp DESC);

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

CREATE INDEX idx_data_trump_timeframe ON data_trump (timeframe);
CREATE INDEX idx_data_trump_timestamp ON data_trump (timestamp DESC);