SELECT create_hypertable('forecasts', by_range('ts'),migrate_data => true);

ALTER TABLE forecasts SET (
    timescaledb.compress,
    timescaledb.compress_segmentby = 'location_id'
    );

SELECT add_compression_policy('forecasts', INTERVAL '7 days');