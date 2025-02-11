CREATE TABLE indicator_aggregate (
    id SERIAL PRIMARY KEY,
    symbol TEXT NOT NULL,
    timeframe TEXT NOT NULL,    -- e.g., 1h, 8h, 1d, 4d, 1w, 1m
    timestamp TIMESTAMPTZ NOT NULL,
    sma NUMERIC NOT NULL,
    ema NUMERIC NOT NULL,
    std_dev NUMERIC NOT NULL,
    lower_bollinger NUMERIC NOT NULL,
    upper_bollinger NUMERIC NOT NULL
    
);

CREATE TABLE indicator_calculate (
    id SERIAL PRIMARY KEY,
    symbol TEXT NOT NULL,
    timeframe TEXT NOT NULL,    -- e.g., 1h, 8h, 1d, 4d, 1w, 1m
    timestamp TIMESTAMPTZ NOT NULL,
    volatility NUMERIC NOT NULL,
    rsi NUMERIC NOT NULL,
    macd NUMERIC NOT NULL,
    macd_signal NUMERIC NOT NULL
);