CREATE TABLE assignments (
    id SERIAL PRIMARY KEY,
    task_id INTEGER REFERENCES tasks(id),
    user_id INTEGER NOT NULL,
    assigned_at TIMESTAMP DEFAULT NOW()
);