CREATE TABLE IF NOT EXISTS forecasts
(
    ts TIMESTAMPTZ NOT NULL,
    type VARCHAR(255) NOT NULL,
    location_id VARCHAR(255) NOT NULL,
    value DECIMAL,
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (ts, type, location_id)
);