CREATE TABLE IF NOT EXISTS tasks.tasks_tags
(
    task_id uuid NOT NULL REFERENCES tasks.tasks(id) ON DELETE CASCADE,
    tag_id uuid NOT NULL REFERENCES tasks.tags(id) ON DELETE CASCADE,
    PRIMARY KEY (task_id, tag_id)
);