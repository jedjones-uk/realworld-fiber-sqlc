// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: article.sql

package sqlc

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createArticle = `-- name: CreateArticle :one
WITH inserted_article AS (
    INSERT INTO articles (slug, title, description, body, author_id)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, slug, title, description, body, created_at, updated_at, favorites_count
),
     inserted_tags AS (
         INSERT INTO tags (tag)
             SELECT unnest($6::text[])
             ON CONFLICT (tag) DO NOTHING
             RETURNING id, tag
     ),
     tag_ids AS (
         SELECT id
         FROM tags
         WHERE tag = ANY ($6::text[])
     ),
     inserted_article_tags AS (
         INSERT INTO article_tags (article_id, tag_id)
             SELECT inserted_article.id, tag_ids.id
             FROM inserted_article, tag_ids
     )
SELECT
    inserted_article.slug,
    inserted_article.title,
    inserted_article.description,
    inserted_article.body,
    array_agg(tag.tag) AS tag_list,
    to_char(inserted_article.created_at, 'YYYY-MM-DD"T"HH24:MI:SS.MSZ') AS created_at,
    to_char(inserted_article.updated_at, 'YYYY-MM-DD"T"HH24:MI:SS.MSZ') AS updated_at,
    false AS favorited,
    inserted_article.favorites_count as favorites_count
FROM
    inserted_article
        JOIN
    article_tags ON inserted_article.id = article_tags.article_id
        JOIN
    tags AS tag ON article_tags.tag_id = tag.id
GROUP BY
    inserted_article.id
`

type CreateArticleParams struct {
	Slug        string      `json:"slug"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Body        string      `json:"body"`
	AuthorID    pgtype.Int8 `json:"authorId"`
	Column6     []string    `json:"column6"`
}

type CreateArticleRow struct {
	Slug           string      `json:"slug"`
	Title          string      `json:"title"`
	Description    string      `json:"description"`
	Body           string      `json:"body"`
	TagList        interface{} `json:"tagList"`
	CreatedAt      string      `json:"createdAt"`
	UpdatedAt      string      `json:"updatedAt"`
	Favorited      bool        `json:"favorited"`
	FavoritesCount int32       `json:"favoritesCount"`
}

func (q *Queries) CreateArticle(ctx context.Context, arg CreateArticleParams) (CreateArticleRow, error) {
	row := q.db.QueryRow(ctx, createArticle,
		arg.Slug,
		arg.Title,
		arg.Description,
		arg.Body,
		arg.AuthorID,
		arg.Column6,
	)
	var i CreateArticleRow
	err := row.Scan(
		&i.Slug,
		&i.Title,
		&i.Description,
		&i.Body,
		&i.TagList,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Favorited,
		&i.FavoritesCount,
	)
	return i, err
}

const deleteArticle = `-- name: DeleteArticle :one
DELETE FROM articles
WHERE slug = $1 and author_id = $2
RETURNING id, slug, title, description, body, created_at, updated_at, favorites_count, author_id
`

type DeleteArticleParams struct {
	Slug     string      `json:"slug"`
	AuthorID pgtype.Int8 `json:"authorId"`
}

func (q *Queries) DeleteArticle(ctx context.Context, arg DeleteArticleParams) (Article, error) {
	row := q.db.QueryRow(ctx, deleteArticle, arg.Slug, arg.AuthorID)
	var i Article
	err := row.Scan(
		&i.ID,
		&i.Slug,
		&i.Title,
		&i.Description,
		&i.Body,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.FavoritesCount,
		&i.AuthorID,
	)
	return i, err
}

const favoriteArticle = `-- name: FavoriteArticle :one
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
RETURNING id, slug, title, description, body, created_at, updated_at, favorites_count, author_id
`

type FavoriteArticleParams struct {
	Slug   string `json:"slug"`
	UserID int64  `json:"userId"`
}

func (q *Queries) FavoriteArticle(ctx context.Context, arg FavoriteArticleParams) (Article, error) {
	row := q.db.QueryRow(ctx, favoriteArticle, arg.Slug, arg.UserID)
	var i Article
	err := row.Scan(
		&i.ID,
		&i.Slug,
		&i.Title,
		&i.Description,
		&i.Body,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.FavoritesCount,
		&i.AuthorID,
	)
	return i, err
}

const getArticle = `-- name: GetArticle :one

SELECT
    a.slug,
    a.title,
    a.description,
    a.body,
    a.created_at,
    a.updated_at,
    a.favorites_count,
    u.username AS username,
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
    a.id, u.id
`

type GetArticleRow struct {
	Slug           string           `json:"slug"`
	Title          string           `json:"title"`
	Description    string           `json:"description"`
	Body           string           `json:"body"`
	CreatedAt      pgtype.Timestamp `json:"createdAt"`
	UpdatedAt      pgtype.Timestamp `json:"updatedAt"`
	FavoritesCount int32            `json:"favoritesCount"`
	Username       string           `json:"username"`
	TagList        interface{}      `json:"tagList"`
	Favorited      bool             `json:"favorited"`
}

// article.sql
func (q *Queries) GetArticle(ctx context.Context, slug string) (GetArticleRow, error) {
	row := q.db.QueryRow(ctx, getArticle, slug)
	var i GetArticleRow
	err := row.Scan(
		&i.Slug,
		&i.Title,
		&i.Description,
		&i.Body,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.FavoritesCount,
		&i.Username,
		&i.TagList,
		&i.Favorited,
	)
	return i, err
}

const getTags = `-- name: GetTags :many
SELECT tag FROM tags
`

func (q *Queries) GetTags(ctx context.Context) ([]string, error) {
	rows, err := q.db.Query(ctx, getTags)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []string{}
	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return nil, err
		}
		items = append(items, tag)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const unfavoriteArticle = `-- name: UnfavoriteArticle :one
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
RETURNING id, slug, title, description, body, created_at, updated_at, favorites_count, author_id
`

type UnfavoriteArticleParams struct {
	Slug   string `json:"slug"`
	UserID int64  `json:"userId"`
}

func (q *Queries) UnfavoriteArticle(ctx context.Context, arg UnfavoriteArticleParams) (Article, error) {
	row := q.db.QueryRow(ctx, unfavoriteArticle, arg.Slug, arg.UserID)
	var i Article
	err := row.Scan(
		&i.ID,
		&i.Slug,
		&i.Title,
		&i.Description,
		&i.Body,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.FavoritesCount,
		&i.AuthorID,
	)
	return i, err
}

const updateArticle = `-- name: UpdateArticle :one
WITH updated_article AS (
    UPDATE articles
        SET slug        = CASE WHEN $3::text IS NOT NULL AND $3::text <> '' THEN $3::text ELSE slug END,
            title       = CASE WHEN $4::text IS NOT NULL AND $4::text <> '' THEN $4::text ELSE title END,
            description = CASE WHEN $5::text IS NOT NULL AND $5::text <> '' THEN $5::text ELSE description END,
            body        = CASE WHEN $6::text IS NOT NULL AND $6::text <> '' THEN $6::text ELSE body END,
            updated_at  = CURRENT_TIMESTAMP
        WHERE slug = $1 and author_id = $2
        RETURNING id, slug, title, description, body, created_at, updated_at, favorites_count, author_id
)
SELECT
    ua.slug,
    ua.title,
    ua.description,
    ua.body,
    to_char(ua.created_at, 'YYYY-MM-DD"T"HH24:MI:SS.MSZ') AS created_at,
    to_char(ua.updated_at, 'YYYY-MM-DD"T"HH24:MI:SS.MSZ') AS updated_at,
    ua.favorites_count AS favorites_count,
    u.username,
    (CASE WHEN EXISTS (SELECT 1 FROM favorites f WHERE f.article_id = ua.id) THEN TRUE ELSE FALSE END) AS favorited,
    ARRAY_AGG(t.tag) AS tagList
FROM
    updated_article ua
        JOIN
    users u ON ua.author_id = u.id
        LEFT JOIN
    article_tags at ON ua.id = at.article_id
        LEFT JOIN
    tags t ON at.tag_id = t.id
GROUP BY
    ua.id, u.id
`

type UpdateArticleParams struct {
	Slug        string      `json:"slug"`
	AuthorID    pgtype.Int8 `json:"authorId"`
	Slug_2      string      `json:"slug2"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Body        string      `json:"body"`
}

type UpdateArticleRow struct {
	Slug           string      `json:"slug"`
	Title          string      `json:"title"`
	Description    string      `json:"description"`
	Body           string      `json:"body"`
	CreatedAt      string      `json:"createdAt"`
	UpdatedAt      string      `json:"updatedAt"`
	FavoritesCount int32       `json:"favoritesCount"`
	Username       string      `json:"username"`
	Favorited      bool        `json:"favorited"`
	Taglist        interface{} `json:"taglist"`
}

func (q *Queries) UpdateArticle(ctx context.Context, arg UpdateArticleParams) (UpdateArticleRow, error) {
	row := q.db.QueryRow(ctx, updateArticle,
		arg.Slug,
		arg.AuthorID,
		arg.Slug_2,
		arg.Title,
		arg.Description,
		arg.Body,
	)
	var i UpdateArticleRow
	err := row.Scan(
		&i.Slug,
		&i.Title,
		&i.Description,
		&i.Body,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.FavoritesCount,
		&i.Username,
		&i.Favorited,
		&i.Taglist,
	)
	return i, err
}
