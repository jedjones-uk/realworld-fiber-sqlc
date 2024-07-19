package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"realworld-fiber-sqlc/pkg/hash"
	"realworld-fiber-sqlc/usecase/dto/sqlc"
)

type UserUPDReq struct {
	User struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
		Image    string `json:"image"`
		Bio      string `json:"bio"`
	} `json:"user"`
}

type CurrentUserResp struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Bio      string `json:"bio"`
	Image    string `json:"image"`
	Token    string `json:"token"`
}

func (h *HandlerBase) CurrentUser(c *fiber.Ctx) error {
	var id int64
	id = userIDFromToken(c)

	user, err := h.Queries.GetUser(c.Context(), id)
	if err != nil {
		return err
	}

	token := c.Locals("user").(*jwt.Token).Raw

	return c.Status(200).JSON(fiber.Map{"user": User{
		Email:    user.Email,
		Username: user.Username,
		Bio:      user.Bio.String,
		Image:    user.Image.String,
		Token:    token,
	}})
}

func (h *HandlerBase) UpdateProfile(c *fiber.Ctx) error {
	var id int64
	id = userIDFromToken(c)

	var params UserUPDReq
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	hashed, err := hash.HashPassword(params.User.Password)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	user, err := h.Queries.UpdateUser(c.Context(), &sqlc.UpdateUserParams{
		ID:       id,
		Email:    params.User.Email,
		Username: params.User.Username,
		Password: hashed,
		Bio:      params.User.Bio,
		Image:    params.User.Image,
	})
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	token := c.Locals("user").(*jwt.Token).Raw

	return c.Status(200).JSON(
		fiber.Map{
			"user": User{
				Email:    user.Email,
				Username: user.Username,
				Bio:      user.Bio.String,
				Image:    user.Image.String,
				Token:    token,
			},
		})
}
