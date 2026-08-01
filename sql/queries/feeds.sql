-- name: CreateFeed :one
INSERT INTO feeds (url)
VALUES (
    $1
)
ON CONFLICT (url) DO UPDATE
SET updated_at = NOW()
RETURNING *;

-- name: GetAllFeeds :many
SELECT * FROM feeds;

-- name: GetFeedByUrl :one
SELECT * FROM feeds
WHERE url = $1;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET last_fetched_at = NOW(),
    updated_at = NOW()
WHERE id = $1;

-- name: GetFeedByID :one
SELECT * FROM feeds
WHERE id = $1;

-- name: GetNextFeedToFetch :one
SELECT url,id FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT $1;