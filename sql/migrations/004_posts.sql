-- +goose Up
CREATE table posts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT NOT NULL,
    url TEXT NOT NULL UNIQUE,
    description TEXT,
    feed_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    published_at TIMESTAMP,
    FOREIGN KEY (feed_id) REFERENCES feeds(id) ON DELETE CASCADE
);


-- +goose Down
DROP TABLE posts;