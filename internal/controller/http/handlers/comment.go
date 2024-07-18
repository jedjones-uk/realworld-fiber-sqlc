package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgtype"
	"realworld-fiber-sqlc/usecase/dto/sqlc"
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
	fmt.Println(req)

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
	resp.Comment.ID = int64(comment.CommentID)
	resp.Comment.CreatedAt = comment.CommentCreatedAt.Time.String()
	resp.Comment.UpdatedAt = comment.CommentUpdatedAt.Time.String()
	resp.Comment.Body = comment.CommentBody
	resp.Comment.Author.Username = comment.UserUsername
	resp.Comment.Author.Bio = comment.UserBio.String
	resp.Comment.Author.Image = comment.UserImage.String
	resp.Comment.Author.Following = false

	return c.JSON(fiber.Map{"comment": resp.Comment})

}
