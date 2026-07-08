-- +goose Up
CREATE table feeds (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    url TEXT  UNIQUE NOT NULL,
    category TEXT NOT NULL DEFAULT 'General',
    last_fetched_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);





-- +goose Down
DROP TABLE feeds;