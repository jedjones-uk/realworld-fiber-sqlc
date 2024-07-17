-- article.sql

-- name: GetArticle :one
SELECT
    a.slug,
    a.title,
    a.description,
    a.body,
    a.created_at,
    a.updated_at,
    a.favorites_count,
    u.username AS author_username,
    u.bio AS author_bio,
    u.image AS author_image,
    ARRAY_AGG(t.tag) AS tag_list,
    (CASE WHEN EXISTS (SELECT 1 FROM favorites f WHERE f.article_id = a.id) THEN TRUE ELSE FALSE END) AS favorited
FROM
    articles a
        JOIN users u ON a.author_id = u.id
        LEFT JOIN article_tags at ON a.id = at.article_id
        LEFT JOIN tags t ON at.tag_id = t.id
WHERE
    a.slug = $1
GROUP BY
    a.id, u.id;

-- name: CreateArticle :one
INSERT INTO articles (slug, title, description, body, author_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;


-- name: UpdateArticle :one
UPDATE articles
SET slug        = CASE WHEN @slug::text IS NOT NULL AND @slug::text <> '' THEN @slug::text ELSE slug END,
    title       = CASE WHEN @title::text IS NOT NULL AND @title::text <> '' THEN @title::text ELSE title END,
    description = CASE WHEN @description::text IS NOT NULL AND @description::text <> '' THEN @description::text ELSE description END,
    body        = CASE WHEN @body::text IS NOT NULL AND @body::text <> '' THEN @body::text ELSE body END,
    updated_at  = CURRENT_TIMESTAMP
WHERE slug = $1 and author_id = $2
RETURNING *;


-- name: DeleteArticle :one
DELETE FROM articles
WHERE slug = $1 and author_id = $2
RETURNING *;


-- name: FavoriteArticle :one
WITH article_id_cte AS (
    SELECT a.id
    FROM articles a
    WHERE a.slug = $1
), insert_favorite AS (
    INSERT INTO favorites (user_id, article_id)
        SELECT $2, a.id
        FROM article_id_cte a
        RETURNING article_id
)
UPDATE articles
SET favorites_count = favorites_count + 1
WHERE id = (SELECT article_id FROM insert_favorite)
RETURNING *;

-- name: UnfavoriteArticle :one
WITH article_id_cte AS (
    SELECT a.id
    FROM articles a
    WHERE a.slug = $1
), delete_favorite AS (
    DELETE FROM favorites
        WHERE user_id = $2 AND article_id = (SELECT id FROM article_id_cte)
        RETURNING article_id
)
UPDATE articles
SET favorites_count = GREATEST(favorites_count - 1, 0)
WHERE id = (SELECT article_id FROM delete_favorite)
RETURNING id, slug, title, description, body, created_at, updated_at, favorites_count, author_id;


