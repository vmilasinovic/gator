-- name: AddFeed :one
INSERT INTO feeds (user_id, url, name)
VALUES (
    $1,
    $2,
    $3
)
RETURNING *;

-- name: GetFeedID :one

SELECT id
FROM feeds
WHERE feeds.url = $1;

-- name: GetFeeds :many
SELECT
    f.name AS feed_name,
    u.name AS user_name,
    f.url
FROM feeds f
JOIN users u ON f.user_id = u.id;