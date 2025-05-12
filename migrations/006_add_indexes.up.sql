-- users
CREATE UNIQUE INDEX idx_users_email ON users.users(email);
CREATE UNIQUE INDEX idx_users_id ON users.users(id);

-- tasks
CREATE UNIQUE INDEX idx_tasks_id ON tasks.tasks(id);
CREATE INDEX idx_tasks_user_id ON tasks.tasks(created_by);

-- comments
CREATE UNIQUE INDEX idx_comments_id ON tasks.comments(id);
CREATE INDEX idx_comments_task_id ON tasks.comments(task_id);
CREATE INDEX idx_comments_author ON tasks.comments(author_id);

-- tags
CREATE UNIQUE INDEX idx_tags_id ON tasks.tags(id);
CREATE UNIQUE INDEX idx_tags_name ON tasks.tags(title);

-- tasks_tags
CREATE UNIQUE INDEX idx_tasks_tags_pair ON tasks.tasks_tags(task_id, tag_id);
CREATE INDEX idx_tasks_tags_task_id ON tasks.tasks_tags(task_id);
CREATE INDEX idx_tasks_tags_tag_id ON tasks.tasks_tags(tag_id);

-- refresh_tokens
CREATE UNIQUE INDEX idx_refresh_tokens_token ON users.refresh_tokens(token);
CREATE INDEX idx_refresh_tokens_user_id ON users.refresh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_expires_at ON users.refresh_tokens(expires_at);