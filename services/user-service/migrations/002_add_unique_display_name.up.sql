CREATE UNIQUE INDEX IF NOT EXISTS unique_display_name ON profiles (display_name) WHERE display_name != '';
