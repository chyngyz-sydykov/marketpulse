CREATE UNIQUE INDEX idx_indicator_btc_1h_timestamp ON indicator_btc_1h(timestamp);
CREATE UNIQUE INDEX idx_indicator_btc_4h_timestamp ON indicator_btc_4h(timestamp);
CREATE UNIQUE INDEX idx_indicator_btc_1d_timestamp ON indicator_btc_1d(timestamp);
CREATE UNIQUE INDEX idx_indicator_btc_4w_timestamp ON indicator_btc_1w(timestamp);
CREATE UNIQUE INDEX idx_indicator_btc_1m_timestamp ON indicator_btc_1m(timestamp);

CREATE UNIQUE INDEX idx_indicator_eth_1h_timestamp ON indicator_eth_1h(timestamp);
CREATE UNIQUE INDEX idx_indicator_eth_4h_timestamp ON indicator_eth_4h(timestamp);
CREATE UNIQUE INDEX idx_indicator_eth_1d_timestamp ON indicator_eth_1d(timestamp);
CREATE UNIQUE INDEX idx_indicator_eth_4w_timestamp ON indicator_eth_1w(timestamp);
CREATE UNIQUE INDEX idx_indicator_eth_1m_timestamp ON indicator_eth_1m(timestamp);

CREATE UNIQUE INDEX idx_indicator_sol_1h_timestamp ON indicator_sol_1h(timestamp);
CREATE UNIQUE INDEX idx_indicator_sol_4h_timestamp ON indicator_sol_4h(timestamp);
CREATE UNIQUE INDEX idx_indicator_sol_1d_timestamp ON indicator_sol_1d(timestamp);
CREATE UNIQUE INDEX idx_indicator_sol_4w_timestamp ON indicator_sol_1w(timestamp);
CREATE UNIQUE INDEX idx_indicator_sol_1m_timestamp ON indicator_sol_1m(timestamp);

CREATE UNIQUE INDEX idx_indicator_bnb_1h_timestamp ON indicator_bnb_1h(timestamp);
CREATE UNIQUE INDEX idx_indicator_bnb_4h_timestamp ON indicator_bnb_4h(timestamp);
CREATE UNIQUE INDEX idx_indicator_bnb_1d_timestamp ON indicator_bnb_1d(timestamp);
CREATE UNIQUE INDEX idx_indicator_bnb_4w_timestamp ON indicator_bnb_1w(timestamp);
CREATE UNIQUE INDEX idx_indicator_bnb_1m_timestamp ON indicator_bnb_1m(timestamp);

CREATE UNIQUE INDEX idx_indicator_trump_1h_timestamp ON indicator_trump_1h(timestamp);
CREATE UNIQUE INDEX idx_indicator_trump_4h_timestamp ON indicator_trump_4h(timestamp);
CREATE UNIQUE INDEX idx_indicator_trump_1d_timestamp ON indicator_trump_1d(timestamp);
CREATE UNIQUE INDEX idx_indicator_trump_4w_timestamp ON indicator_trump_1w(timestamp);
CREATE UNIQUE INDEX idx_indicator_trump_1m_timestamp ON indicator_trump_1m(timestamp);