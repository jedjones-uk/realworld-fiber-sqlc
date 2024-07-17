package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgtype"
	"realworld-fiber-sqlc/usecase/dto/sqlc"
	"strings"
)

type SingleArticleResp struct {
	Slug           string   `json:"slug"`
	Title          string   `json:"title"`
	Description    string   `json:"description"`
	Body           string   `json:"body"`
	TagList        []string `json:"tagList"`
	CreatedAt      string   `json:"createdAt"`
	UpdatedAt      string   `json:"updatedAt"`
	Favorited      bool     `json:"favorited"`
	FavoritesCount int      `json:"favoritesCount"`
	Author         struct {
		Username  string `json:"username"`
		Bio       string `json:"bio"`
		Image     string `json:"image"`
		Following bool   `json:"following"`
	} `json:"author"`
}

type CreateArticleReq struct {
	Article struct {
		Title       string   `json:"title"`
		Description string   `json:"description"`
		Body        string   `json:"body"`
		TagList     []string `json:"tagList"`
	} `json:"article"`
}

func (h *HandlerBase) GetArticle(c *fiber.Ctx) error {
	slug := c.Params("slug")
	article, err := h.Queries.GetArticle(c.Context(), slug)
	if err != nil {
		return err
	}

	userID := userIDFromToken(c)
	var arg sqlc.GetUserProfileParams
	if userID == 0 {
		arg = sqlc.GetUserProfileParams{
			Username: article.AuthorUsername,
		}
	} else {
		arg = sqlc.GetUserProfileParams{
			Username:   article.AuthorUsername,
			FollowerID: userID,
		}
	}

	author, err := h.Queries.GetUserProfile(c.Context(), arg)
	if err != nil {
		return err
	}

	// Initialize an empty slice of strings for tags
	tags := []string{}
	if article.TagList != nil {
		tagList, ok := article.TagList.([]interface{})
		if !ok {
			return fmt.Errorf("failed to convert tagList to []interface{}")
		}

		for _, tag := range tagList {
			if tag != nil {
				tagStr, ok := tag.(string)
				if !ok {
					return fmt.Errorf("failed to convert tag to string")
				}
				tags = append(tags, tagStr)
			}
		}
	}
	fmt.Println(tags)

	var resp SingleArticleResp
	resp.Slug = article.Slug
	resp.Title = article.Title
	resp.Description = article.Description
	resp.Body = article.Body
	resp.TagList = tags
	resp.CreatedAt = article.CreatedAt.Time.Format("2006-01-02T15:04:05.000Z07:00")
	resp.UpdatedAt = article.UpdatedAt.Time.Format("2006-01-02T15:04:05.000Z07:00")
	resp.Favorited = article.Favorited
	resp.FavoritesCount = int(article.FavoritesCount)
	resp.Author.Username = author.Username
	resp.Author.Bio = author.Bio.String
	resp.Author.Image = author.Image.String
	resp.Author.Following = author.Following

	return c.Status(200).JSON(fiber.Map{"article": resp})
}

func transformString(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	return s
}

func (h *HandlerBase) CreateArticle(c *fiber.Ctx) error {
	var req CreateArticleReq
	if err := c.BodyParser(&req); err != nil {
		return err
	}

	authorId := userIDFromToken(c)
	if authorId == 0 {
		return c.Status(401).JSON(fiber.Map{"message": "Unauthorized"})
	}

	slug := transformString(req.Article.Title)

	article, err := h.Queries.CreateArticle(c.Context(), sqlc.CreateArticleParams{
		Title:       req.Article.Title,
		Slug:        slug,
		Description: req.Article.Description,
		Body:        req.Article.Body,
		AuthorID:    pgtype.Int8{Int64: authorId},
	})
	if err != nil {
		return err
	}

	profile, err := h.Queries.GetProfileById(c.Context(), sqlc.GetProfileByIdParams{
		ID: authorId,
	})
	if err != nil {
		return err
	}

	// TODO taglist

	resp := SingleArticleResp{
		Slug:           article.Slug,
		Title:          article.Title,
		Description:    article.Description,
		Body:           article.Body,
		TagList:        []string{},
		CreatedAt:      article.CreatedAt.Time.Format("2006-01-02T15:04:05.000Z07:00"),
		UpdatedAt:      article.UpdatedAt.Time.Format("2006-01-02T15:04:05.000Z07:00"),
		Favorited:      false,
		FavoritesCount: 0,
	}

	resp.Author.Username = profile.Username
	resp.Author.Bio = profile.Bio.String
	resp.Author.Image = profile.Image.String
	resp.Author.Following = profile.Following

	return c.Status(200).JSON(fiber.Map{
		"article": resp,
	})
}

func (h *HandlerBase) UpdateArticle(c *fiber.Ctx) error {
	userId := userIDFromToken(c)
	fmt.Println("userId", userId)
	if userId == 0 {
		return c.Status(401).JSON(fiber.Map{"message": "Unauthorized"})
	}

	slug := c.Params("slug")

	var req CreateArticleReq
	if err := c.BodyParser(&req); err != nil {
		return err
	}
	fmt.Println("req", req.Article.Body)

	newSlug := transformString(req.Article.Title)

	article, err := h.Queries.UpdateArticle(c.Context(), sqlc.UpdateArticleParams{
		Slug:        slug,
		Title:       req.Article.Title,
		Description: req.Article.Description,
		Body:        req.Article.Body,
		AuthorID:    pgtype.Int8{Int64: userId, Valid: true},
		Slug_2:      newSlug,
	})
	if err != nil {
		return err
	}

	profile, err := h.Queries.GetProfileById(c.Context(), sqlc.GetProfileByIdParams{
		ID: userId,
	})
	if err != nil {
		return err
	}

	resp := SingleArticleResp{
		Slug:           article.Slug,
		Title:          article.Title,
		Description:    article.Description,
		Body:           article.Body,
		TagList:        []string{},
		CreatedAt:      article.CreatedAt.Time.Format("2006-01-02T15:04:05.000Z07:00"),
		UpdatedAt:      article.UpdatedAt.Time.Format("2006-01-02T15:04:05.000Z07:00"),
		Favorited:      false,
		FavoritesCount: 0,
	}
	resp.Author.Username = profile.Username
	resp.Author.Bio = profile.Bio.String
	resp.Author.Image = profile.Image.String
	resp.Author.Following = profile.Following

	return c.Status(200).JSON(fiber.Map{
		"article": resp,
	})

}

func (h *HandlerBase) DeleteArticle(c *fiber.Ctx) error {
	userId := userIDFromToken(c)
	if userId == 0 {
		return c.Status(401).JSON(fiber.Map{"message": "Unauthorized"})
	}

	slug := c.Params("slug")

	_, err := h.Queries.DeleteArticle(c.Context(), sqlc.DeleteArticleParams{
		Slug:     slug,
		AuthorID: pgtype.Int8{Int64: userId, Valid: true},
	})
	if err != nil {
		return err
	}

	return c.Status(200).JSON(fiber.Map{})
}

func (h *HandlerBase) FavoriteArticle(c *fiber.Ctx) error {
	userId := userIDFromToken(c)
	if userId == 0 {
		return c.Status(401).JSON(fiber.Map{"message": "Unauthorized"})
	}

	slug := c.Params("slug")

	article, err := h.Queries.FavoriteArticle(c.Context(), sqlc.FavoriteArticleParams{
		Slug:   slug,
		UserID: userId,
	})
	if err != nil {
		return err
	}

	profile, err := h.Queries.GetProfileById(c.Context(), sqlc.GetProfileByIdParams{
		ID: userId,
	})
	if err != nil {
		return err
	}

	resp := SingleArticleResp{
		Slug:           slug,
		Title:          article.Title,
		Description:    article.Description,
		Body:           article.Body,
		TagList:        []string{},
		CreatedAt:      article.CreatedAt.Time.Format("2006-01-02T15:04:05.000Z07:00"),
		UpdatedAt:      article.UpdatedAt.Time.Format("2006-01-02T15:04:05.000Z07:00"),
		Favorited:      true,
		FavoritesCount: int(article.FavoritesCount),
	}
	resp.Author.Username = profile.Username
	resp.Author.Bio = profile.Bio.String
	resp.Author.Image = profile.Image.String
	resp.Author.Following = profile.Following

	return c.Status(200).JSON(fiber.Map{"article": resp})
}

func (h *HandlerBase) UnfavoriteArticle(c *fiber.Ctx) error {
	userId := userIDFromToken(c)
	if userId == 0 {
		return c.Status(401).JSON(fiber.Map{"message": "Unauthorized"})
	}

	slug := c.Params("slug")

	article, err := h.Queries.UnfavoriteArticle(c.Context(), sqlc.UnfavoriteArticleParams{
		Slug:   slug,
		UserID: userId,
	})
	if err != nil {
		return err
	}

	profile, err := h.Queries.GetProfileById(c.Context(), sqlc.GetProfileByIdParams{
		ID: userId,
	})
	if err != nil {
		return err
	}

	resp := SingleArticleResp{
		Slug:           slug,
		Title:          article.Title,
		Description:    article.Description,
		Body:           article.Body,
		TagList:        []string{},
		CreatedAt:      article.CreatedAt.Time.Format("2006-01-02T15:04:05.000Z07:00"),
		UpdatedAt:      article.UpdatedAt.Time.Format("2006-01-02T15:04:05.000Z07:00"),
		Favorited:      true,
		FavoritesCount: int(article.FavoritesCount),
	}
	resp.Author.Username = profile.Username
	resp.Author.Bio = profile.Bio.String
	resp.Author.Image = profile.Image.String
	resp.Author.Following = profile.Following

	return c.Status(200).JSON(fiber.Map{"article": resp})

}
