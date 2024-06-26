CREATE TABLE IF NOT EXISTS metrics(
			name 	VARCHAR(200) PRIMARY KEY,
			g_type 	VARCHAR(200) NOT NULL,
			g_value DOUBLE PRECISION NOT NULL,
			delta 	bigint NOT NULL,
			UNIQUE(name, g_type)
		)