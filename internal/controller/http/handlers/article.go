package handlers

import (
	"fmt"
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

type Article struct {
	Slug           string   `json:"slug"`
	Title          string   `json:"title"`
	Description    string   `json:"description"`
	Body           string   `json:"body"`
	TagList        []string `json:"tagList"`
	CreatedAt      string   `json:"createdAt"`
	UpdatedAt      string   `json:"updatedAt"`
	Favorited      bool     `json:"favorited"`
	FavoritesCount int32    `json:"favoritesCount"`
	Author         Author   `json:"author"`
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

	author, err := h.Queries.GetUserProfile(c.Context(), &arg)
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

	article, err := h.Queries.CreateArticle(c.Context(), &sqlc.CreateArticleParams{
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

	author, err := h.Queries.GetUserProfileById(c.Context(), &sqlc.GetUserProfileByIdParams{
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

	article, err := h.Queries.UpdateArticle(c.Context(), &sqlc.UpdateArticleParams{
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

	author, err := h.Queries.GetUserProfile(c.Context(), &sqlc.GetUserProfileParams{
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

	_, err := h.Queries.DeleteArticle(c.Context(), &sqlc.DeleteArticleParams{
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

	article, err := h.Queries.FavoriteArticle(c.Context(), &sqlc.FavoriteArticleParams{
		Slug:   slug,
		UserID: userId,
	})
	if err != nil {
		return err
	}

	author, err := h.Queries.GetUserProfileById(c.Context(), &sqlc.GetUserProfileByIdParams{
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

	article, err := h.Queries.UnfavoriteArticle(c.Context(), &sqlc.UnfavoriteArticleParams{
		Slug:   slug,
		UserID: userId,
	})
	if err != nil {
		return err
	}

	author, err := h.Queries.GetUserProfileById(c.Context(), &sqlc.GetUserProfileByIdParams{
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

	userID := userIDFromToken(c)
	userIDPG := &pgtype.Int8{}
	if userID != 0 {
		_ = userIDPG.Scan(userID)
	}

	tag := c.Query("tag")
	tagPG := &pgtype.Text{}
	if tag != "" {
		_ = tagPG.Scan(tag)
	}

	author := c.Query("author")
	authorPG := &pgtype.Text{}
	if author != "" {
		_ = authorPG.Scan(author)
	}

	favorited := c.Query("favorited")
	favoritedPG := &pgtype.Text{}
	if favorited != "" {
		_ = favoritedPG.Scan(favorited)
	}

	limit := c.Query("limit", "20")
	limitPG := &pgtype.Int4{}
	_ = limitPG.Scan(limit)

	offset := c.Query("offset", "0")
	offsetPG := &pgtype.Int4{}
	_ = offsetPG.Scan(offset)

	var params sqlc.ListArticlesParams
	params = sqlc.ListArticlesParams{
		Tag:         *tagPG,
		Author:      *authorPG,
		FavoritedBy: *favoritedPG,
		Limitt:      *limitPG,
		Offsett:     *offsetPG,
		UserID:      *userIDPG,
	}

	articlesData, err := h.Queries.ListArticles(c.Context(), &params)
	if err != nil {
		return err
	}

	var articles []Article
	cnt := 0
	for _, article := range articlesData {
		cnt += 1
		articles = append(articles, Article{
			Slug:        article.Slug,
			Title:       article.Title,
			Description: article.Description,
			Body:        article.Body,
			TagList:     formTagList(article.TagList),
			CreatedAt:   article.CreatedAt.Time.Format("2006-01-02T15:04:05.000Z"),
			UpdatedAt:   article.UpdatedAt.Time.Format("2006-01-02T15:04:05.000Z"),

			Author: Author{
				Username:  article.AuthorUsername,
				Bio:       article.AuthorBio.String,
				Image:     article.AuthorImage.String,
				Following: article.Following.(bool),
			},
		})
	}

	return c.Status(200).JSON(fiber.Map{"articles": articles, "articlesCount": cnt})

}

func (h *HandlerBase) FeedArticles(c *fiber.Ctx) error {
	userID := userIDFromToken(c)
	fmt.Println("userID", userID)
	if userID == 0 {
		return c.Status(401).JSON(fiber.Map{"message": "Unauthorized"})
	}
	userIDPG := &pgtype.Int4{}
	_ = userIDPG.Scan(userID)

	limit := c.Query("limit", "20")
	limitPG := &pgtype.Int4{}
	_ = limitPG.Scan(limit)

	offset := c.Query("offset", "0")
	offsetPG := &pgtype.Int4{}
	_ = offsetPG.Scan(offset)

	fmt.Println("userIDPG", userIDPG)

	articlesData, err := h.Queries.FeedArticles(c.Context(), &sqlc.FeedArticlesParams{
		Limitt:  limitPG.Int32,
		Offsett: offsetPG.Int32,
		UserID:  *userIDPG,
	})
	if err != nil {
		return err
	}

	var articles []Article
	for _, article := range articlesData {
		articles = append(articles, Article{
			Slug:        article.Slug,
			Title:       article.Title,
			Description: article.Description,
			Body:        article.Body,
			TagList:     formTagList(article.TagList),

			CreatedAt:      article.CreatedAt,
			UpdatedAt:      article.UpdatedAt,
			Favorited:      article.Favorited.(bool),
			FavoritesCount: article.FavoritesCount,
			Author: Author{
				Username:  article.Username.String,
				Bio:       article.Bio.String,
				Image:     article.Image.String,
				Following: article.Following.(bool),
			},
		})
	}

	return c.Status(200).JSON(fiber.Map{"articles": articles})

}
