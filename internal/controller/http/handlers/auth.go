package handlers

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"realworld-fiber-sqlc/pkg/hash"
	jwt2 "realworld-fiber-sqlc/pkg/jwt"
	"realworld-fiber-sqlc/usecase/dto/sqlc"
	"strconv"
)

type RegisterParamsT struct {
	User RegisterParams `json:"user"`
}

type RegisterParams struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginParams struct {
	User struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	} `json:"user"`
}

type User struct {
	Email    string `json:"email"`
	Token    string `json:"token"`
	Username string `json:"username"`
	Bio      string `json:"bio"`
	Image    string `json:"image"`
}

func userIDFromToken(c *fiber.Ctx) int64 {
	token, ok := c.Locals("user").(*jwt.Token)
	if !ok {
		return 0
	}
	claims := token.Claims.(jwt.MapClaims)
	subj, err := claims.GetSubject()
	if err != nil {
		return 0
	}

	id, _ := strconv.ParseInt(subj, 10, 64)
	return id
}

func (h *HandlerBase) Login(c *fiber.Ctx) error {
	var params LoginParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	hashPass, err := hash.HashPassword(params.User.Password)
	if err != nil {
		return err
	}

	user, err := h.Queries.GetUserByEmail(c.Context(), params.User.Email)
	if err != nil {
		log.Printf("error: %v", err)
		return err
	}

	if !hash.CheckPasswordHash(params.User.Password, hashPass) {
		return c.Status(401).JSON(fiber.Map{"message": "Invalid password"})
	}

	token, _ := jwt2.GenerateToken(user.ID)

	resp := generateUser(user, token)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"user": resp})
}

func (h *HandlerBase) Register(c *fiber.Ctx) error {
	var params RegisterParamsT
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	hashPass, err := hash.HashPassword(params.User.Password)
	if err != nil {
		return err
	}

	user, err := h.Queries.CreateUser(context.Background(), &sqlc.CreateUserParams{
		Email:    params.User.Email,
		Username: params.User.Username,
		Password: hashPass,
	})

	if err != nil {
		return err
	}

	token, _ := jwt2.GenerateToken(user.ID)

	resp := generateUser(user, token)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"user": resp})
}

func generateUser(userDB sqlc.User, token string) *User {
	return &User{
		Email:    userDB.Email,
		Username: userDB.Username,
		Bio:      userDB.Bio.String,
		Image:    userDB.Image.String,
		Token:    token,
	}

}
