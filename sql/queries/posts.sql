-- name: CreatePost :one
INSERT INTO posts (title,url,description,feed_id,published_at)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;

-- name: GetPostsForUser :many
SELECT posts.* from posts
INNER JOIN feed_follows ON feed_follows.feed_id = posts.feed_id
WHERE feed_follows.user_id = $1
ORDER BY posts.published_at DESC
LIMIT $2;