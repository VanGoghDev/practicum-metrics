CREATE TABLE IF NOT EXISTS gauges
(
    id      CHAR(50) PRIMARY KEY,
    delta   INTEGER,
    g_value FLOAT,
    g_type  CHAR(50)
)