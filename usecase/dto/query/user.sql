-- user.sql

-- name: GetUser :one
SELECT * FROM users WHERE id = $1;

-- name: UpdateUser :one
UPDATE users
SET
    email = CASE WHEN @email::text IS NOT NULL AND @email::text <> '' THEN @email::text ELSE email END,
    username = CASE WHEN @username::text IS NOT NULL AND @username::text <> '' THEN @username::text ELSE username END,
    password = CASE WHEN @password::text IS NOT NULL AND @password::text <> '' THEN @password::text ELSE password END,
    image = CASE WHEN @image::text IS NOT NULL AND @image::text <> '' THEN @image::text ELSE image END,
    bio = CASE WHEN @bio::text IS NOT NULL AND @bio::text <> '' THEN @bio::text ELSE bio END
WHERE
    id = $1
RETURNING *;





-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: CreateUser :one
INSERT INTO users (email, username, password) VALUES ($1, $2, $3)
RETURNING id;

