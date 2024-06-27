CREATE TABLE IF NOT EXISTS metrics(
			name 	VARCHAR(200) PRIMARY KEY,
			g_type 	VARCHAR(200) NOT NULL,
			g_value DOUBLE PRECISION,
			delta 	bigint,
			UNIQUE(name, g_type),
			CHECK (g_value IS NOT NULL OR delta IS NOT NULL)
		)