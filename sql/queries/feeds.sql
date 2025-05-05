-- name: AddFeed :one
INSERT INTO feeds (user_id, url, name)
VALUES (
    $1,
    $2,
    $3
)
RETURNING *;

-- name: GetFeeds :many
SELECT
    f.name AS feed_name,
    u.name AS user_name,
    f.url
FROM feeds f
JOIN users u ON f.user_id = u.id;