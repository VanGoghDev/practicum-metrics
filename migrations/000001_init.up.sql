CREATE TABLE IF NOT EXISTS gauges
(
    id      VARCHAR(200) PRIMARY KEY,
    delta   INT,
    g_value NUMERIC,
    g_type  VARCHAR(200)
)