CREATE TABLE IF NOT EXISTS indicator_btc_1h PARTITION OF indicator_btc FOR VALUES IN ('1h');
CREATE TABLE IF NOT EXISTS indicator_btc_4h PARTITION OF indicator_btc FOR VALUES IN ('4h');
CREATE TABLE IF NOT EXISTS indicator_btc_1d PARTITION OF indicator_btc FOR VALUES IN ('1d');

CREATE TABLE IF NOT EXISTS indicator_eth_1h PARTITION OF indicator_eth FOR VALUES IN ('1h');
CREATE TABLE IF NOT EXISTS indicator_eth_4h PARTITION OF indicator_eth FOR VALUES IN ('4h');
CREATE TABLE IF NOT EXISTS indicator_eth_1d PARTITION OF indicator_eth FOR VALUES IN ('1d');

CREATE TABLE IF NOT EXISTS indicator_sol_1h PARTITION OF indicator_sol FOR VALUES IN ('1h');
CREATE TABLE IF NOT EXISTS indicator_sol_4h PARTITION OF indicator_sol FOR VALUES IN ('4h');
CREATE TABLE IF NOT EXISTS indicator_sol_1d PARTITION OF indicator_sol FOR VALUES IN ('1d');

CREATE TABLE IF NOT EXISTS indicator_bnb_1h PARTITION OF indicator_bnb FOR VALUES IN ('1h');
CREATE TABLE IF NOT EXISTS indicator_bnb_4h PARTITION OF indicator_bnb FOR VALUES IN ('4h');
CREATE TABLE IF NOT EXISTS indicator_bnb_1d PARTITION OF indicator_bnb FOR VALUES IN ('1d');

CREATE TABLE IF NOT EXISTS indicator_trump_1h PARTITION OF indicator_trump FOR VALUES IN ('1h');
CREATE TABLE IF NOT EXISTS indicator_trump_4h PARTITION OF indicator_trump FOR VALUES IN ('4h');
CREATE TABLE IF NOT EXISTS indicator_trump_1d PARTITION OF indicator_trump FOR VALUES IN ('1d');

