CREATE DOMAIN job_status AS TEXT 
CHECK (VALUE IN ('pending', 'processing', 'completed', 'failed'));

CREATE TABLE jobs ( 
    id UUID PRIMARY KEY,
    status job_status NOT NULL DEFAULT 'pending',
    retries INT NOT NULL DEFAULT 0,
    error TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);