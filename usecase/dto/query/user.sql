-- user.sql

-- name: GetUser :one
SELECT * FROM users WHERE id = $1;

-- name: UpdateUser :exec
UPDATE users SET email = $2, username = $3, password = $4 WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: CreateUser :one
INSERT INTO users (email, username, password) VALUES ($1, $2, $3)
RETURNING id;

