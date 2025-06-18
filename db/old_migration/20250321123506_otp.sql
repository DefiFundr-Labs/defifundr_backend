-- +goose Up
-- Create otp_verifications table for managing one-time password verifications

-- Create an enum for OTP purpose
DROP TYPE IF EXISTS otp_purpose;

CREATE TYPE otp_purpose AS ENUM (
    'email_verification',
    'password_reset',
    'phone_verification',
    'account_recovery',
    'two_factor_auth',
    'login_confirmation'
);

-- Create the otp_verifications table
CREATE TABLE otp_verifications (
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

-- Indexes for performance
CREATE INDEX idx_otp_verifications_user_id ON otp_verifications(user_id);
CREATE INDEX idx_otp_verifications_purpose ON otp_verifications(purpose);
CREATE INDEX idx_otp_verifications_contact_method ON otp_verifications(contact_method);
CREATE INDEX idx_otp_verifications_created_at ON otp_verifications(created_at);

-- +goose Down
-- Drop the table
DROP TABLE IF EXISTS otp_verifications;

-- Drop the custom type
DROP TYPE IF EXISTS otp_purpose;

-- Optionally remove the UUID extension if no longer needed
-- DROP EXTENSION IF EXISTS "uuid-ossp";