CREATE TABLE IF NOT EXISTS profiles (
    id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255) UNIQUE NOT NULL,
    display_name VARCHAR(255),
    bio TEXT,
    avatar_url VARCHAR(500),
    country VARCHAR(100),
    timezone VARCHAR(100),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE INDEX idx_profiles_user_id ON profiles(user_id);

CREATE TABLE IF NOT EXISTS user_settings (
    id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255) UNIQUE NOT NULL,
    language VARCHAR(10) DEFAULT 'en',
    theme VARCHAR(20) DEFAULT 'light',
    notifications_enabled BOOLEAN DEFAULT TRUE,
    privacy_level VARCHAR(20) DEFAULT 'public'
);

CREATE INDEX idx_user_settings_user_id ON user_settings(user_id);

