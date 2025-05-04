-- +goose Up
CREATE TABLE security_events (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    event_type TEXT NOT NULL,
    ip_address TEXT NOT NULL,
    user_agent TEXT,
    metadata JSONB,
    timestamp TIMESTAMP NOT NULL
);

CREATE INDEX idx_security_events_user_id ON security_events(user_id);
CREATE INDEX idx_security_events_type_time ON security_events(event_type, timestamp);

-- +goose Down
DROP TABLE security_events;
