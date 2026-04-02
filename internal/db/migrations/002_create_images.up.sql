CREATE DOMAIN image_size AS TEXT 
CHECK (VALUE IN ('thumbnail', 'medium', 'large'));

CREATE TABLE images (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    job_id UUID NOT NULL,
    url TEXT NOT NULL,
    size image_size NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (job_id) REFERENCES jobs(id)
);
