CREATE TABLE IF NOT EXISTS tasks.tags
(
    id uuid PRIMARY KEY  DEFAULT gen_random_uuid(),
    title TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);