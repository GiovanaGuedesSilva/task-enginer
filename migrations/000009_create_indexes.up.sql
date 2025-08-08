-- ==== Performance indexes

-- Indexes to optimize filtering by foreign keys
CREATE INDEX idx_tasks_project_id ON tasks(project_id);           -- For filtering tasks by project
CREATE INDEX idx_assignments_task_id ON assignments(task_id);     -- For filtering assignments by task
CREATE INDEX idx_comments_task_id ON comments(task_id);           -- For filtering comments by task
CREATE INDEX idx_task_tags_task_id ON task_tags(task_id);         -- For filtering task_tags by task
CREATE INDEX idx_task_tags_tag_id ON task_tags(tag_id);           -- For filtering task_tags by tag

-- Indexes on date/time columns for temporal queries
CREATE INDEX idx_tasks_created_at ON tasks(created_at);           -- For temporal queries on tasks.created_at
CREATE INDEX idx_tasks_updated_at ON tasks(updated_at);           -- For temporal queries on tasks.updated_at
CREATE INDEX idx_tasks_due_date ON tasks(due_date);               -- For temporal queries on tasks.due_date
CREATE INDEX idx_projects_created_at ON projects(created_at);     -- For temporal queries on projects.created_at
CREATE INDEX idx_projects_updated_at ON projects(updated_at);     -- For temporal queries on projects.updated_at
CREATE INDEX idx_assignments_assigned_at ON assignments(assigned_at); -- For temporal queries on assignments.assigned_at

-- Indexes on status and priority fields for faster filtering
CREATE INDEX idx_tasks_status ON tasks(status);                   -- For filtering tasks by status
CREATE INDEX idx_tasks_priority ON tasks(priority);               -- For filtering tasks by priority

-- ==== Full-text search indexes

-- Full-text index on task title and description
CREATE INDEX idx_tasks_fulltext ON tasks USING gin(
  to_tsvector('portuguese', title || ' ' || COALESCE(description, '')) -- For full-text search on tasks
);

-- Full-text index on project name and description
CREATE INDEX idx_projects_fulltext ON projects USING gin(
  to_tsvector('portuguese', name || ' ' || COALESCE(description, '')) -- For full-text search on projects
);

-- Full-text index on comment content
CREATE INDEX idx_comments_fulltext ON comments USING gin( -- For full-text search on comments
  to_tsvector('portuguese', content)
);

-- Full-text index on tag name
CREATE INDEX idx_tags_fulltext ON tags USING gin( -- For full-text search on tags
  to_tsvector('portuguese', name)
);

-- ==== Composite indexes

CREATE INDEX idx_tasks_project_status ON tasks(project_id, status);         -- For filtering tasks by project and status
CREATE INDEX idx_tasks_project_priority ON tasks(project_id, priority);     -- For filtering tasks by project and priority
CREATE INDEX idx_tasks_status_due_date ON tasks(status, due_date);          -- For filtering tasks by status and due date
CREATE INDEX idx_tasks_project_created_at ON tasks(project_id, created_at DESC); -- For filtering recent tasks by project
CREATE INDEX idx_assignments_user_assigned_at ON assignments(user_id, assigned_at DESC); -- For filtering recent assignments by user
CREATE INDEX idx_metrics_project_created_at ON metrics(project_id, created_at DESC); -- For filtering recent metrics by project

-- ==== Partial indexes

CREATE INDEX idx_tasks_pending ON tasks(project_id, priority, due_date) 
WHERE status = 'pending'; -- For active (pending) tasks

CREATE INDEX idx_tasks_with_due_date ON tasks(project_id, due_date) 
WHERE due_date IS NOT NULL; -- For tasks with deadlines

CREATE INDEX idx_tasks_high_priority ON tasks(project_id, due_date) 
WHERE priority = 'high'; -- For high priority tasks

-- ==== Functional indexes

CREATE INDEX idx_projects_name_lower ON projects(LOWER(name));    -- Case-insensitive search on project names
CREATE INDEX idx_tasks_title_lower ON tasks(LOWER(title));        -- Case-insensitive search on task titles
CREATE INDEX idx_tasks_due_year ON tasks(EXTRACT(YEAR FROM due_date)); -- For year-based queries on due_date

-- ==== Unique indexes

CREATE UNIQUE INDEX idx_assignments_unique_task_user ON assignments(task_id, user_id); -- Prevent duplicate assignments per task/user
CREATE UNIQUE INDEX idx_task_tags_unique_task_tag ON task_tags(task_id, tag_id);      -- Prevent duplicate tags per task
