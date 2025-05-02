-- name: CreateUser :one
INSERT INTO users (name)
VALUES (
    $1
)
RETURNING *;

-- name: GetUser :one
SELECT id, created_at, updated_at, name
FROM users
WHERE name = $1;

-- name: GetUsers :many
SELECT name
FROM users
ORDER BY name; 

-- name: ClearUsers :exec
DELETE
FROM users;