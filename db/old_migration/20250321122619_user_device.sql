-- +goose Up
-- SQL in this section is executed when the migration is applied
CREATE TABLE user_device_tokens (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    -- Device Identification
    device_token VARCHAR NOT NULL UNIQUE,
    platform VARCHAR(50) NOT NULL, -- e.g., 'ios', 'android', 'web'
    device_type VARCHAR(100), -- e.g., 'smartphone', 'tablet', 'desktop'
    device_model VARCHAR(100), -- e.g., 'iPhone 13', 'Samsung Galaxy S21'
    
    -- Operating System Details
    os_name VARCHAR(50),
    os_version VARCHAR(50),
    
    -- Push Notification Tokens (if applicable)
    push_notification_token VARCHAR,
    
    -- Authentication and Security
    is_active BOOLEAN NOT NULL DEFAULT true,
    is_verified BOOLEAN NOT NULL DEFAULT false,
    
    -- Metadata
    last_used_at TIMESTAMPTZ,
    first_registered_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    
    -- Additional Context
    app_version VARCHAR(50),
    client_ip VARCHAR(45), -- IPv4 or IPv6
    
    -- Expiration and Management
    expires_at TIMESTAMPTZ,
    is_revoked BOOLEAN NOT NULL DEFAULT false
);

-- Indexes for performance and querying
CREATE INDEX idx_user_device_tokens_user_id ON user_device_tokens(user_id);
CREATE INDEX idx_user_device_tokens_platform ON user_device_tokens(platform);
CREATE INDEX idx_user_device_tokens_device_token ON user_device_tokens(device_token);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back
DROP TABLE IF EXISTS user_device_tokens;