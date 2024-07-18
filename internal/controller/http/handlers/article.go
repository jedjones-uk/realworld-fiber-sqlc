package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mitchellh/mapstructure"
	"realworld-fiber-sqlc/usecase/dto/sqlc"
	"strings"
)

type CreateArticleReq struct {
	Article struct {
		Title       string   `json:"title"`
		Description string   `json:"description"`
		Body        string   `json:"body"`
		TagList     []string `json:"tagList"`
	} `json:"article"`
}

func formTagList(tagList interface{}) []string {
	tags := []string{}
	if tagList != nil {
		tagList, ok := tagList.([]interface{})
		if !ok {
			return tags
		}

		for _, tag := range tagList {
			if tag != nil {
				tagStr, ok := tag.(string)
				if !ok {
					return tags
				}
				tags = append(tags, tagStr)
			}
		}
	}
	return tags
}

func transformString(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	return s
}

func (h *HandlerBase) GetArticle(c *fiber.Ctx) error {
	slug := c.Params("slug")
	article, err := h.Queries.GetArticle(c.Context(), slug)
	if err != nil {
		return err
	}

	var arg sqlc.GetUserProfileParams
	userID := userIDFromToken(c)
	if userID == 0 {
		arg = sqlc.GetUserProfileParams{
			Username: article.Username,
		}
	} else {
		arg = sqlc.GetUserProfileParams{
			Username:   article.Username,
			FollowerID: userID,
		}
	}

	author, err := h.Queries.GetUserProfile(c.Context(), arg)
	if err != nil {
		return err
	}

	tagList := formTagList(article.TagList)

	articleMap := make(map[string]interface{})
	authorMap := make(map[string]interface{})

	mapstructure.Decode(article, &articleMap)
	mapstructure.Decode(author, &authorMap)

	articleMap["author"] = authorMap
	articleMap["tagList"] = tagList

	delete(articleMap, "username")

	return c.Status(200).JSON(fiber.Map{"article": articleMap})
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
		Column6:     req.Article.TagList,
	})
	if err != nil {
		return err
	}

	author, err := h.Queries.GetUserProfileById(c.Context(), sqlc.GetUserProfileByIdParams{
		ID:         authorId,
		FollowerID: authorId,
	})
	if err != nil {
		return err
	}

	articleMap := make(map[string]interface{})
	authorMap := make(map[string]interface{})

	mapstructure.Decode(article, &articleMap)
	mapstructure.Decode(author, &authorMap)

	articleMap["author"] = authorMap

	return c.Status(200).JSON(
		fiber.Map{
			"article": articleMap,
		})
}

func (h *HandlerBase) UpdateArticle(c *fiber.Ctx) error {
	userId := userIDFromToken(c)
	if userId == 0 {
		return c.Status(401).JSON(fiber.Map{"message": "Unauthorized"})
	}

	slug := c.Params("slug")

	var req CreateArticleReq
	if err := c.BodyParser(&req); err != nil {
		return err
	}

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

	author, err := h.Queries.GetUserProfile(c.Context(), sqlc.GetUserProfileParams{
		Username:   article.Username,
		FollowerID: userId,
	})
	if err != nil {
		return err
	}

	articleMap := make(map[string]interface{})
	authorMap := make(map[string]interface{})

	mapstructure.Decode(article, &articleMap)
	mapstructure.Decode(author, &authorMap)

	articleMap["author"] = authorMap
	delete(articleMap, "username")

	return c.Status(200).JSON(
		fiber.Map{
			"article": articleMap,
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

	author, err := h.Queries.GetUserProfileById(c.Context(), sqlc.GetUserProfileByIdParams{
		ID:         userId,
		FollowerID: userId,
	})
	if err != nil {
		return err
	}

	articleMap := make(map[string]interface{})
	authorMap := make(map[string]interface{})

	mapstructure.Decode(article, &articleMap)
	mapstructure.Decode(author, &authorMap)

	articleMap["author"] = authorMap
	delete(articleMap, "authorId")

	return c.Status(200).JSON(fiber.Map{"article": articleMap})
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

	author, err := h.Queries.GetUserProfileById(c.Context(), sqlc.GetUserProfileByIdParams{
		ID:         userId,
		FollowerID: userId,
	})
	if err != nil {
		return err
	}

	articleMap := make(map[string]interface{})
	authorMap := make(map[string]interface{})

	mapstructure.Decode(article, &articleMap)
	mapstructure.Decode(author, &authorMap)

	articleMap["author"] = authorMap
	delete(articleMap, "authorId")

	return c.Status(200).JSON(fiber.Map{"article": articleMap})

}

func (h *HandlerBase) GetTags(c *fiber.Ctx) error {
	tags, err := h.Queries.GetTags(c.Context())
	if err != nil {
		return err
	}

	return c.Status(200).JSON(fiber.Map{"tags": tags})
}

func (h *HandlerBase) GetArticles(c *fiber.Ctx) error {
	//tag := c.Query("tag")
	//author := c.Query("author")
	//favorited := c.Query("favorited")
	//limit := c.Query("limit", "20")
	//offset := c.Query("offset", "0")
	return nil

}
