-- Add foreign key constraints to indicator tables
ALTER TABLE indicator_btc_1h
ADD CONSTRAINT fk_indicator_btc_1h_data_timestamp
FOREIGN KEY (data_timestamp) REFERENCES data_btc_1h (timestamp) ON DELETE CASCADE;

ALTER TABLE indicator_btc_4h
ADD CONSTRAINT fk_indicator_btc_4h_data_timestamp
FOREIGN KEY (data_timestamp) REFERENCES data_btc_4h (timestamp) ON DELETE CASCADE;

ALTER TABLE indicator_btc_1d
ADD CONSTRAINT fk_indicator_btc_1d_data_timestamp
FOREIGN KEY (data_timestamp) REFERENCES data_btc_1d (timestamp) ON DELETE CASCADE;

ALTER TABLE indicator_btc_1w
ADD CONSTRAINT fk_indicator_btc_1w_data_timestamp
FOREIGN KEY (data_timestamp) REFERENCES data_btc_1w (timestamp) ON DELETE CASCADE;

ALTER TABLE indicator_btc_1m
ADD CONSTRAINT fk_indicator_btc_1m_data_timestamp
FOREIGN KEY (data_timestamp) REFERENCES data_btc_1m (timestamp) ON DELETE CASCADE;

ALTER TABLE indicator_bnb_1h
ADD CONSTRAINT fk_indicator_bnb_1h_data_timestamp
FOREIGN KEY (data_timestamp) REFERENCES data_bnb_1h (timestamp) ON DELETE CASCADE;

ALTER TABLE indicator_bnb_4h
ADD CONSTRAINT fk_indicator_bnb_4h_data_timestamp
FOREIGN KEY (data_timestamp) REFERENCES data_bnb_4h (timestamp) ON DELETE CASCADE;

ALTER TABLE indicator_bnb_1d
ADD CONSTRAINT fk_indicator_bnb_1d_data_timestamp
FOREIGN KEY (data_timestamp) REFERENCES data_bnb_1d (timestamp) ON DELETE CASCADE;

ALTER TABLE indicator_bnb_1w
ADD CONSTRAINT fk_indicator_bnb_1w_data_timestamp
FOREIGN KEY (data_timestamp) REFERENCES data_bnb_1w (timestamp) ON DELETE CASCADE;

ALTER TABLE indicator_bnb_1m
ADD CONSTRAINT fk_indicator_bnb_1m_data_timestamp
FOREIGN KEY (data_timestamp) REFERENCES data_bnb_1m (timestamp) ON DELETE CASCADE;

ALTER TABLE indicator_sol_1h
ADD CONSTRAINT fk_indicator_sol_1h_data_timestamp
FOREIGN KEY (data_timestamp) REFERENCES data_sol_1h (timestamp) ON DELETE CASCADE;

ALTER TABLE indicator_sol_4h
ADD CONSTRAINT fk_indicator_sol_4h_data_timestamp
FOREIGN KEY (data_timestamp) REFERENCES data_sol_4h (timestamp) ON DELETE CASCADE;

ALTER TABLE indicator_sol_1d
ADD CONSTRAINT fk_indicator_sol_1d_data_timestamp
FOREIGN KEY (data_timestamp) REFERENCES data_sol_1d (timestamp) ON DELETE CASCADE;

ALTER TABLE indicator_sol_1w
ADD CONSTRAINT fk_indicator_sol_1w_data_timestamp
FOREIGN KEY (data_timestamp) REFERENCES data_sol_1w (timestamp) ON DELETE CASCADE;

ALTER TABLE indicator_sol_1m
ADD CONSTRAINT fk_indicator_sol_1m_data_timestamp
FOREIGN KEY (data_timestamp) REFERENCES data_sol_1m (timestamp) ON DELETE CASCADE;

ALTER TABLE indicator_eth_1h
ADD CONSTRAINT fk_indicator_eth_1h_data_timestamp
FOREIGN KEY (data_timestamp) REFERENCES data_eth_1h (timestamp) ON DELETE CASCADE;

ALTER TABLE indicator_eth_4h
ADD CONSTRAINT fk_indicator_eth_4h_data_timestamp
FOREIGN KEY (data_timestamp) REFERENCES data_eth_4h (timestamp) ON DELETE CASCADE;

ALTER TABLE indicator_eth_1d
ADD CONSTRAINT fk_indicator_eth_1d_data_timestamp
FOREIGN KEY (data_timestamp) REFERENCES data_eth_1d (timestamp) ON DELETE CASCADE;

ALTER TABLE indicator_eth_1w
ADD CONSTRAINT fk_indicator_eth_1w_data_timestamp
FOREIGN KEY (data_timestamp) REFERENCES data_eth_1w (timestamp) ON DELETE CASCADE;

ALTER TABLE indicator_eth_1m
ADD CONSTRAINT fk_indicator_eth_1m_data_timestamp
FOREIGN KEY (data_timestamp) REFERENCES data_eth_1m (timestamp) ON DELETE CASCADE;

ALTER TABLE indicator_trump_1h
ADD CONSTRAINT fk_indicator_trump_1h_data_timestamp
FOREIGN KEY (data_timestamp) REFERENCES data_trump_1h (timestamp) ON DELETE CASCADE;

ALTER TABLE indicator_trump_4h
ADD CONSTRAINT fk_indicator_trump_4h_data_timestamp
FOREIGN KEY (data_timestamp) REFERENCES data_trump_4h (timestamp) ON DELETE CASCADE;

ALTER TABLE indicator_trump_1d
ADD CONSTRAINT fk_indicator_trump_1d_data_timestamp
FOREIGN KEY (data_timestamp) REFERENCES data_trump_1d (timestamp) ON DELETE CASCADE;

ALTER TABLE indicator_trump_1w
ADD CONSTRAINT fk_indicator_trump_1w_data_timestamp
FOREIGN KEY (data_timestamp) REFERENCES data_trump_1w (timestamp) ON DELETE CASCADE;

ALTER TABLE indicator_trump_1m
ADD CONSTRAINT fk_indicator_trump_1m_data_timestamp
FOREIGN KEY (data_timestamp) REFERENCES data_trump_1m (timestamp) ON DELETE CASCADE;