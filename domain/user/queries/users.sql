-- name: GetUser :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: CreateUser :exec
INSERT INTO users (id, name, email) VALUES ($1, $2, $3);

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;

-- name: UpdateUser :exec
UPDATE users SET name = $2 WHERE id = $1;

-- name: GetUsers :many
SELECT * FROM users;

