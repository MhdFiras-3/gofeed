-- +goose Up
ALTER TABLE feeds DROP COLUMN name;
ALTER TABLE feed_follows ADD COLUMN name TEXT NOT NULL;

-- +goose Down
ALTER TABLE feed_follows DROP COLUMN name;
ALTER TABLE feeds ADD COLUMN name TEXT NOT NULL;