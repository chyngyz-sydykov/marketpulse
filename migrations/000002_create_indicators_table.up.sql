CREATE TABLE indicators (
    id SERIAL PRIMARY KEY,
    symbol TEXT NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL,
    sma NUMERIC NOT NULL,
    ema NUMERIC NOT NULL,
    std_dev NUMERIC NOT NULL,
    volatility NUMERIC NOT NULL,
    rsi NUMERIC NOT NULL,
    macd NUMERIC NOT NULL,
    macd_signal NUMERIC NOT NULL,
    lower_bollinger NUMERIC NOT NULL,
    upper_bollinger NUMERIC NOT NULL
);