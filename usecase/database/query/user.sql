-- user.sql

-- name: GetUser :one
SELECT * FROM users WHERE id = $1;

-- name: UpdateUser :exec
UPDATE users SET email = $2, username = $3, password = $4 WHERE id = $1;

-- name: UpdateProfile :exec
UPDATE profiles SET bio = $2, image = $3 WHERE user_id = $1;
