-- init.sql
-- Create the metadata table
CREATE TABLE IF NOT EXISTS metadata_table (
                                              my_key TEXT PRIMARY KEY,
                                              my_value JSONB NOT NULL
);

-- Insert sample data
INSERT INTO metadata_table (my_key, my_value) VALUES
                                                  ('config1', '{"environment": "production", "max_connections": 100, "debug_mode": false}'::jsonb),
                                                  ('user_preferences', '{"theme": "dark", "notifications": true, "language": "en-US"}'::jsonb),
                                                  ('app_settings', '{"cache_ttl": 3600, "retry_attempts": 3, "timeout": 30}'::jsonb),
                                                  ('feature_flags', '{"new_ui": true, "beta_features": false, "maintenance_mode": false}'::jsonb),
                                                  ('api_config', '{"rate_limit": 1000, "allowed_origins": ["example.com"], "version": "v2"}'::jsonb);

-- Create an index on the my_key column for faster lookups
CREATE INDEX IF NOT EXISTS idx_metadata_key ON metadata_table(my_key);

-- Create a gin index on the jsonb column for faster JSON queries
CREATE INDEX IF NOT EXISTS idx_metadata_value ON metadata_table USING gin(my_value);