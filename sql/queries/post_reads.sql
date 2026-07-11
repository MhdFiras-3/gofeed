-- name: MarkPostRead :exec
INSERT INTO post_reads (user_id,post_id) VALUES (
    $1,
    $2
);

-- name: CheckPostRead :one
SELECT * from post_reads
WHERE user_id = $1 AND post_id = $2;

-- name: GetPostReadsPerUser :many
SELECT posts.* 
FROM post_reads
INNER JOIN posts ON post_reads.post_id = posts.id
WHERE post_reads.user_id = $1
ORDER BY post_reads.read_at DESC;
