package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mitchellh/mapstructure"
	"realworld-fiber-sqlc/usecase/dto/sqlc"
	"strconv"
)

type CreateCommentReq struct {
	Comment struct {
		Body string `json:"body"`
	} `json:"comment"`
}

type CommentResp struct {
	Comment struct {
		ID        int64  `json:"id"`
		CreatedAt string `json:"createdAt"`
		UpdatedAt string `json:"updatedAt"`
		Body      string `json:"body"`
		Author    struct {
			Username  string `json:"username"`
			Bio       string `json:"bio"`
			Image     string `json:"image"`
			Following bool   `json:"following"`
		} `json:"author"`
	}
}

func (h *HandlerBase) CreateComment(c *fiber.Ctx) error {
	userID := userIDFromToken(c)
	fmt.Println(userID)
	if userID == 0 {
		return fiber.ErrUnauthorized
	}
	slug := c.Params("slug")

	var req CreateCommentReq
	if err := c.BodyParser(&req); err != nil {
		return err
	}

	comment, err := h.Queries.CreateComment(c.Context(), &sqlc.CreateCommentParams{
		Slug:   slug,
		Body:   req.Comment.Body,
		UserID: pgtype.Int8{Int64: userID, Valid: true},
	})
	if err != nil {
		return err
	}

	author, err := h.Queries.GetUserProfileById(c.Context(), &sqlc.GetUserProfileByIdParams{
		ID:         userID,
		FollowerID: userID,
	})
	if err != nil {
		return err
	}

	commentMap := make(map[string]interface{})
	authorMap := make(map[string]interface{})

	mapstructure.Decode(comment, &commentMap)
	mapstructure.Decode(author, &authorMap)

	commentMap["author"] = authorMap

	return c.JSON(fiber.Map{"comment": commentMap})

}

func (h *HandlerBase) DeleteComment(c *fiber.Ctx) error {
	userID := userIDFromToken(c)
	if userID == 0 {
		return fiber.ErrUnauthorized
	}

	commentID := c.Params("id")

	i, err := strconv.ParseInt(commentID, 10, 32)
	if err != nil {
		panic(err)
	}
	commentID32 := int32(i)

	err = h.Queries.DeleteComment(c.Context(), &sqlc.DeleteCommentParams{
		ID:     commentID32,
		UserID: pgtype.Int8{Int64: userID, Valid: true},
	})
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}

type Comment struct {
	ID        int32  `json:"id"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	Body      string `json:"body"`
	Author    Author `json:"author"`
}

func (h *HandlerBase) GetComments(c *fiber.Ctx) error {
	slug := c.Params("slug")

	commentsData, err := h.Queries.GetCommentsByArticleSlug(c.Context(), slug)
	if err != nil {
		return err
	}

	var comments []Comment
	for _, comment := range commentsData {
		comments = append(comments, Comment{
			ID:        comment.ID,
			CreatedAt: comment.CreatedAt,
			UpdatedAt: comment.UpdatedAt,
			Body:      comment.Body,
			Author: Author{
				Username:  comment.Username,
				Bio:       comment.Bio.String,
				Image:     comment.Image.String,
				Following: comment.Following,
			},
		})
	}

	return c.JSON(fiber.Map{"comments": comments})
}
