-- +goose Up
-- Authentication tables
CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    refresh_token VARCHAR(1024),
    user_agent TEXT,
    client_ip VARCHAR(45),
    last_used_at TIMESTAMPTZ DEFAULT NOW(),
    web_oauth_client_id TEXT,
    oauth_access_token TEXT,
    oauth_id_token TEXT,
    user_login_type VARCHAR(100),
    mfa_verified BOOLEAN DEFAULT FALSE,
    is_blocked BOOLEAN DEFAULT FALSE,
    expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE user_devices (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    device_token VARCHAR(255),
    platform VARCHAR(50),
    device_type VARCHAR(100),
    device_model VARCHAR(100),
    os_name VARCHAR(50),
    os_version VARCHAR(50),
    push_notification_token VARCHAR(255),
    is_active BOOLEAN DEFAULT TRUE,
    is_verified BOOLEAN DEFAULT FALSE,
    last_used_at TIMESTAMPTZ DEFAULT NOW(),
    app_version VARCHAR(50),
    client_ip VARCHAR(45),
    expires_at TIMESTAMPTZ,
    is_revoked BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE security_events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    company_id UUID REFERENCES companies(id) ON DELETE SET NULL,
    event_type VARCHAR(100) NOT NULL,
    severity VARCHAR(50) NOT NULL,
    ip_address VARCHAR(45),
    user_agent TEXT,
    metadata JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

DROP TYPE IF EXISTS otp_purpose;

CREATE TYPE otp_purpose AS ENUM (
    'email_verification',
    'password_reset',
    'phone_verification',
    'account_recovery',
    'two_factor_auth',
    'login_confirmation'
);

-- Create the otp table
CREATE TABLE otp (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    
    -- OTP Details
    otp_code VARCHAR(10) NOT NULL,
    hashed_otp VARCHAR(255) NOT NULL,
    
    -- Verification Context
    purpose otp_purpose NOT NULL,
    contact_method VARCHAR(255),
    
    -- Tracking and Limits
    attempts_made INTEGER NOT NULL DEFAULT 0,
    max_attempts INTEGER NOT NULL DEFAULT 5,
    is_verified BOOLEAN NOT NULL DEFAULT false,
    
    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL,
    verified_at TIMESTAMPTZ,
    
    -- Additional Metadata
    ip_address INET,
    user_agent VARCHAR(500),
    device_id UUID,
    
    -- Constraints
    CONSTRAINT unique_unverified_otp UNIQUE (user_id, purpose, is_verified),
    CONSTRAINT max_attempts_check CHECK (attempts_made <= max_attempts)
);

-- Create indexes
CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);
CREATE INDEX idx_user_devices_user_id ON user_devices(user_id);
CREATE INDEX idx_user_devices_device_token ON user_devices(device_token);
CREATE INDEX idx_security_events_user_id ON security_events(user_id);
CREATE INDEX idx_security_events_company_id ON security_events(company_id);
CREATE INDEX idx_security_events_event_type ON security_events(event_type);
CREATE INDEX idx_security_events_created_at ON security_events(created_at);
CREATE INDEX idx_otp_user_id ON otp(user_id);
CREATE INDEX idx_otp_purpose ON otp(purpose);
CREATE INDEX idx_otp_contact_method ON otp(contact_method);
CREATE INDEX idx_otp_created_at ON otp(created_at);

-- +goose Down
DROP TABLE IF EXISTS security_events CASCADE;
DROP TABLE IF EXISTS user_devices CASCADE;
DROP TABLE IF EXISTS sessions CASCADE;
DROP TABLE IF EXISTS otp;
DROP TYPE IF EXISTS otp_purpose;