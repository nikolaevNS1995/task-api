-- refresh_tokens
DROP INDEX IF EXISTS idx_refresh_tokens_expires_at;
DROP INDEX IF EXISTS idx_refresh_tokens_user_id;
DROP INDEX IF EXISTS idx_refresh_tokens_token;

-- tasks_tags
DROP INDEX IF EXISTS idx_tasks_tags_tag_id;
DROP INDEX IF EXISTS idx_tasks_tags_task_id;
DROP INDEX IF EXISTS idx_tasks_tags_pair;

-- tags
DROP INDEX IF EXISTS idx_tags_name;
DROP INDEX IF EXISTS idx_tags_id;

-- comments
DROP INDEX IF EXISTS idx_comments_author;
DROP INDEX IF EXISTS idx_comments_task_id;
DROP INDEX IF EXISTS idx_comments_id;

-- tasks
DROP INDEX IF EXISTS idx_tasks_user_id;
DROP INDEX IF EXISTS idx_tasks_id;

-- users
DROP INDEX IF EXISTS idx_users_id;
DROP INDEX IF EXISTS idx_users_email;