CREATE TABLE events (
    id SERIAL PRIMARY KEY,
    task_id INTEGER REFERENCES tasks(id),
    event_type VARCHAR(100) NOT NULL,
    payload JSONB,
    created_at TIMESTAMP DEFAULT NOW()
);