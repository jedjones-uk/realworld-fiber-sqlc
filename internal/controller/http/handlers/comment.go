package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgtype"
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
	if userID == 0 {
		return fiber.ErrUnauthorized
	}
	slug := c.Params("slug")

	var req CreateCommentReq
	if err := c.BodyParser(&req); err != nil {
		return err
	}
	fmt.Println(req.Comment.Body, slug, userID)

	comment, err := h.Queries.CreateComment(c.Context(), sqlc.CreateCommentParams{
		Slug:   slug,
		Body:   req.Comment.Body,
		UserID: pgtype.Int8{Int64: userID, Valid: true},
	})
	if err != nil {
		fmt.Println(err)
		return err
	}

	resp := CommentResp{}
	resp.Comment.ID = int64(comment.ID)
	resp.Comment.CreatedAt = comment.CreatedAt.Time.Format("2006-01-02T15:04:05.000Z07:00")
	resp.Comment.UpdatedAt = comment.UpdatedAt.Time.Format("2006-01-02T15:04:05.000Z07:00")
	resp.Comment.Body = comment.Body

	profile, err := h.Queries.GetProfileById(c.Context(), sqlc.GetProfileByIdParams{
		ID:         userID,
		FollowerID: userID,
	})
	if err != nil {
		return err
	}

	resp.Comment.Author.Username = profile.Username
	resp.Comment.Author.Bio = profile.Bio.String
	resp.Comment.Author.Image = profile.Image.String
	resp.Comment.Author.Following = false

	return c.JSON(fiber.Map{"comment": resp.Comment})

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

	err = h.Queries.DeleteComment(c.Context(), sqlc.DeleteCommentParams{
		ID:     commentID32,
		UserID: pgtype.Int8{Int64: userID, Valid: true},
	})
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *HandlerBase) GetComments(c *fiber.Ctx) error {
	slug := c.Params("slug")

	comments, err := h.Queries.GetCommentsByArticleSlug(c.Context(), slug)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{"comments": comments})
}
