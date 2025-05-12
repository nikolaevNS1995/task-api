CREATE TABLE IF NOT EXISTS tasks.comments
(
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    task_id uuid NOT NULL,
    author_id uuid NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    FOREIGN KEY (task_id) REFERENCES tasks.tasks(id) ON DELETE CASCADE,
    FOREIGN KEY (author_id) REFERENCES users.users(id) ON DELETE CASCADE
);