-- name: CreateFeedFollow :one
WITH inserted_feed_follows AS (INSERT INTO feed_follows (user_id,feed_id,name)
VALUES(
    $1,
    $2,
    $3
)
ON CONFLICT (user_id,feed_id) DO UPDATE
SET name = EXCLUDED.name, updated_at = NOW()
RETURNING *
)

SELECT inserted_feed_follows.*,users.name as user_name FROM inserted_feed_follows
INNER JOIN users on inserted_feed_follows.user_id = users.id;

-- name: GetFeedFollowsForUser :many
SELECT feed_follows.*,
users.name as user_name
FROM feed_follows
INNER JOIN users ON feed_follows.user_id = users.id
WHERE feed_follows.user_id = $1;

-- name: DeleteFeedFollow :one
DELETE FROM feed_follows
WHERE user_id = $1 AND feed_id = $2
RETURNING id;