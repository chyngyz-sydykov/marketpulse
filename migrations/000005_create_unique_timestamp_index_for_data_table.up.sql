CREATE UNIQUE INDEX idx_data_btc_1h_timestamp ON data_btc_1h(timestamp);
CREATE UNIQUE INDEX idx_data_btc_4h_timestamp ON data_btc_4h(timestamp);
CREATE UNIQUE INDEX idx_data_btc_1d_timestamp ON data_btc_1d(timestamp);

CREATE UNIQUE INDEX idx_data_eth_1h_timestamp ON data_eth_1h(timestamp);
CREATE UNIQUE INDEX idx_data_eth_4h_timestamp ON data_eth_4h(timestamp);
CREATE UNIQUE INDEX idx_data_eth_1d_timestamp ON data_eth_1d(timestamp);

CREATE UNIQUE INDEX idx_data_sol_1h_timestamp ON data_sol_1h(timestamp);
CREATE UNIQUE INDEX idx_data_sol_4h_timestamp ON data_sol_4h(timestamp);
CREATE UNIQUE INDEX idx_data_sol_1d_timestamp ON data_sol_1d(timestamp);

CREATE UNIQUE INDEX idx_data_bnb_1h_timestamp ON data_bnb_1h(timestamp);
CREATE UNIQUE INDEX idx_data_bnb_4h_timestamp ON data_bnb_4h(timestamp);
CREATE UNIQUE INDEX idx_data_bnb_1d_timestamp ON data_bnb_1d(timestamp);

CREATE UNIQUE INDEX idx_data_trump_1h_timestamp ON data_trump_1h(timestamp);
CREATE UNIQUE INDEX idx_data_trump_4h_timestamp ON data_trump_4h(timestamp);
CREATE UNIQUE INDEX idx_data_trump_1d_timestamp ON data_trump_1d(timestamp);