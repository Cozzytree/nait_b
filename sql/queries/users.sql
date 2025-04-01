-- name: CreateUser :one
INSERT INTO users (name, email, auth_id, avatar, provider)
VALUES ($1, $2, $3, $4, $5)
RETURNING id;

-- name: GetUser :one
SELECT * FROM users
WHERE auth_id = $1;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;
