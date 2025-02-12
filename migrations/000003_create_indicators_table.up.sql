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
    data_id bigint not null,
    data_timeframe TEXT NOT NULL,
    PRIMARY KEY (id, timeframe),
    FOREIGN KEY (data_id, data_timeframe) REFERENCES data_btc (id, timeframe) ON DELETE CASCADE
) PARTITION BY LIST (timeframe);

CREATE INDEX idx_indicator_btc_data_id_timeframe ON indicator_btc (data_id DESC, timeframe);

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
    data_id bigint not null,
    data_timeframe TEXT NOT NULL,
    PRIMARY KEY (id, timeframe),
    FOREIGN KEY (data_id, data_timeframe) REFERENCES data_bnb (id, timeframe) ON DELETE CASCADE
) PARTITION BY LIST (timeframe);

CREATE INDEX idx_indicator_bnb_data_id_timeframe ON indicator_bnb (data_id DESC, timeframe);

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
    data_id bigint not null,
    data_timeframe TEXT NOT NULL,
    PRIMARY KEY (id, timeframe),
    FOREIGN KEY (data_id, data_timeframe) REFERENCES data_eth (id, timeframe) ON DELETE CASCADE
) PARTITION BY LIST (timeframe);

CREATE INDEX idx_indicator_eth_data_id_timeframe ON indicator_eth (data_id DESC, timeframe);

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
    data_id bigint not null,
    data_timeframe TEXT NOT NULL,
    PRIMARY KEY (id, timeframe),
    FOREIGN KEY (data_id, data_timeframe) REFERENCES data_sol (id, timeframe) ON DELETE CASCADE
) PARTITION BY LIST (timeframe);

CREATE INDEX idx_indicator_sol_data_id_timeframe ON indicator_sol (data_id DESC, timeframe);

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
    data_id bigint not null,
    data_timeframe TEXT NOT NULL,
    PRIMARY KEY (id, timeframe),
    FOREIGN KEY (data_id, data_timeframe) REFERENCES data_trump (id, timeframe) ON DELETE CASCADE
) PARTITION BY LIST (timeframe);

CREATE INDEX idx_indicator_trump_data_id_timeframe ON indicator_trump (data_id DESC, timeframe);