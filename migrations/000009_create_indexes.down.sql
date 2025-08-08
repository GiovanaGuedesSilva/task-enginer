-- Performance indexes

DROP INDEX IF EXISTS idx_tasks_project_id;
DROP INDEX IF EXISTS idx_assignments_task_id;
DROP INDEX IF EXISTS idx_comments_task_id;
DROP INDEX IF EXISTS idx_task_tags_task_id;
DROP INDEX IF EXISTS idx_task_tags_tag_id;

DROP INDEX IF EXISTS idx_tasks_created_at;
DROP INDEX IF EXISTS idx_tasks_updated_at;
DROP INDEX IF EXISTS idx_tasks_due_date;
DROP INDEX IF EXISTS idx_projects_created_at;
DROP INDEX IF EXISTS idx_projects_updated_at;
DROP INDEX IF EXISTS idx_assignments_assigned_at;

DROP INDEX IF EXISTS idx_tasks_status;
DROP INDEX IF EXISTS idx_tasks_priority;

-- Fulltext indexes

DROP INDEX IF EXISTS idx_tasks_fulltext;
DROP INDEX IF EXISTS idx_projects_fulltext;
DROP INDEX IF EXISTS idx_comments_fulltext;
DROP INDEX IF EXISTS idx_tags_fulltext;

-- Composite indexes

DROP INDEX IF EXISTS idx_tasks_project_status;
DROP INDEX IF EXISTS idx_tasks_project_priority;
DROP INDEX IF EXISTS idx_tasks_status_due_date;
DROP INDEX IF EXISTS idx_tasks_project_created_at;
DROP INDEX IF EXISTS idx_assignments_user_assigned_at;
DROP INDEX IF EXISTS idx_metrics_project_created_at;

-- Partial indexes

DROP INDEX IF EXISTS idx_tasks_pending;
DROP INDEX IF EXISTS idx_tasks_with_due_date;
DROP INDEX IF EXISTS idx_tasks_high_priority;

-- Functional indexes

DROP INDEX IF EXISTS idx_projects_name_lower;
DROP INDEX IF EXISTS idx_tasks_title_lower;
DROP INDEX IF EXISTS idx_tasks_due_year;

-- Unique indexes

DROP INDEX IF EXISTS idx_assignments_unique_task_user;
DROP INDEX IF EXISTS idx_task_tags_unique_task_tag;
