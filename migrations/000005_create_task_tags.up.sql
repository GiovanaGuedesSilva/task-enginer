CREATE TABLE task_tags (
    task_id INTEGER REFERENCES tasks(id),
    tag_id INTEGER REFERENCES tags(id),
    PRIMARY KEY (task_id, tag_id)
);