CREATE TABLE IF NOT EXISTS indicator_btc (
    id BIGSERIAL,
    timeframe TEXT NOT NULL,    -- e.g., 1h, 8h, 1d, 4d, 1w, 1m
    timestamp TIMESTAMPTZ NOT NULL,
    sma NUMERIC NOT NULL,
    ema NUMERIC NOT NULL,
    std_dev NUMERIC NOT NULL,
    lower_bollinger NUMERIC NOT NULL,
    upper_bollinger NUMERIC NOT NULL,
    volatility NUMERIC NOT NULL,
    rsi NUMERIC NOT NULL,
    macd NUMERIC NOT NULL,
    macd_signal NUMERIC NOT NULL,
    data_timestamp TIMESTAMPTZ NOT NULL,
    PRIMARY KEY (id, timeframe)
) PARTITION BY LIST (timeframe);

CREATE TABLE IF NOT EXISTS indicator_bnb (
    id BIGSERIAL,
    timeframe TEXT NOT NULL,    -- e.g., 1h, 8h, 1d, 4d, 1w, 1m
    timestamp TIMESTAMPTZ NOT NULL,
    sma NUMERIC NOT NULL,
    ema NUMERIC NOT NULL,
    std_dev NUMERIC NOT NULL,
    lower_bollinger NUMERIC NOT NULL,
    upper_bollinger NUMERIC NOT NULL,
    volatility NUMERIC NOT NULL,
    rsi NUMERIC NOT NULL,
    macd NUMERIC NOT NULL,
    macd_signal NUMERIC NOT NULL,
    data_timestamp TIMESTAMPTZ NOT NULL,
    PRIMARY KEY (id, timeframe)
) PARTITION BY LIST (timeframe);

CREATE TABLE IF NOT EXISTS indicator_eth (
    id BIGSERIAL,
    timeframe TEXT NOT NULL,    -- e.g., 1h, 8h, 1d, 4d, 1w, 1m
    timestamp TIMESTAMPTZ NOT NULL,
    sma NUMERIC NOT NULL,
    ema NUMERIC NOT NULL,
    std_dev NUMERIC NOT NULL,
    lower_bollinger NUMERIC NOT NULL,
    upper_bollinger NUMERIC NOT NULL,
    volatility NUMERIC NOT NULL,
    rsi NUMERIC NOT NULL,
    macd NUMERIC NOT NULL,
    macd_signal NUMERIC NOT NULL,
    data_timestamp TIMESTAMPTZ NOT NULL,
    PRIMARY KEY (id, timeframe)
) PARTITION BY LIST (timeframe);

CREATE TABLE IF NOT EXISTS indicator_sol (
    id BIGSERIAL,
    timeframe TEXT NOT NULL,    -- e.g., 1h, 8h, 1d, 4d, 1w, 1m
    timestamp TIMESTAMPTZ NOT NULL,
    sma NUMERIC NOT NULL,
    ema NUMERIC NOT NULL,
    std_dev NUMERIC NOT NULL,
    lower_bollinger NUMERIC NOT NULL,
    upper_bollinger NUMERIC NOT NULL,
    volatility NUMERIC NOT NULL,
    rsi NUMERIC NOT NULL,
    macd NUMERIC NOT NULL,
    macd_signal NUMERIC NOT NULL,
    data_timestamp TIMESTAMPTZ NOT NULL,
    PRIMARY KEY (id, timeframe)
) PARTITION BY LIST (timeframe);

CREATE TABLE IF NOT EXISTS indicator_trump (
    id BIGSERIAL,
    timeframe TEXT NOT NULL,    -- e.g., 1h, 8h, 1d, 4d, 1w, 1m
    timestamp TIMESTAMPTZ NOT NULL,
    sma NUMERIC NOT NULL,
    ema NUMERIC NOT NULL,
    std_dev NUMERIC NOT NULL,
    lower_bollinger NUMERIC NOT NULL,
    upper_bollinger NUMERIC NOT NULL,
    volatility NUMERIC NOT NULL,
    rsi NUMERIC NOT NULL,
    macd NUMERIC NOT NULL,
    macd_signal NUMERIC NOT NULL,
    data_timestamp TIMESTAMPTZ NOT NULL,
    PRIMARY KEY (id, timeframe)
) PARTITION BY LIST (timeframe);