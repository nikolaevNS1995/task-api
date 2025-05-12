CREATE SCHEMA IF NOT EXISTS tasks;

CREATE TABLE IF NOT EXISTS tasks.tasks
(
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT NOT NULL,
    description TEXT,
    status TEXT NOT NULL DEFAULT 'todo',
    created_by uuid REFERENCES users.users(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);