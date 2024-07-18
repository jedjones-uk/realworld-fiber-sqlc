-- comment.sql

-- name: CreateComment :one
WITH article_id_cte AS (
    SELECT a.id
    FROM articles a
    WHERE a.slug = $1
), insert_comment AS (
    INSERT INTO comments (body, user_id, article_id)
        SELECT $2, $3, id
        FROM article_id_cte
        RETURNING id, body, user_id, article_id, created_at
)
SELECT
    c.id AS comment_id,
    c.body AS comment_body,
    c.article_id AS comment_article_id,
    c.created_at AS comment_created_at,
    c.updated_at AS comment_updated_at,
    u.username AS user_username,
    u.image AS user_image,
    u.bio AS user_bio
FROM comments c
         JOIN articles a ON c.article_id = a.id
         JOIN users u ON c.user_id = u.id
WHERE c.id = (SELECT id FROM insert_comment);
