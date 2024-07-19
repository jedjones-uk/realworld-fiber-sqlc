package handlers

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"net/mail"
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

func validEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
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
	h.Logger.Info("Login handler")
	var params LoginParams
	if err := c.BodyParser(&params); err != nil {
		h.Logger.Error("error parsing body: %v", err)
		return c.Status(422).JSON(fiber.Map{"errors": fiber.Map{
			"body": []string{"Invalid body"},
		}})
	}

	hashPass, err := hash.HashPassword(params.User.Password)
	if err != nil {
		h.Logger.Error("error hashing password: %v", err)
		return c.SendStatus(500)
	}

	user, err := h.Queries.GetUserByEmail(c.Context(), params.User.Email)
	if err != nil {
		h.Logger.Error("error getting user by email: %v", err)
		return c.SendStatus(500)
	}

	if !hash.CheckPasswordHash(params.User.Password, hashPass) {
		return c.Status(422).JSON(fiber.Map{"errors": fiber.Map{
			"body": []string{"Invalid email or password"},
		}})
	}

	token, err := jwt2.GenerateToken(user.ID)
	if err != nil {
		h.Logger.Error("error generating token: %v", err)
		return c.SendStatus(500)
	}

	resp := generateUser(user, token)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"user": resp})
}

func (h *HandlerBase) Register(c *fiber.Ctx) error {
	h.Logger.Info("Register handler")
	var params RegisterParamsT
	if err := c.BodyParser(&params); err != nil {
		h.Logger.Error("error parsing body: %v", err)
		return c.Status(422).JSON(fiber.Map{"errors": fiber.Map{
			"body": []string{"Invalid body"},
		}})
	}

	if !validEmail(params.User.Email) {
		h.Logger.Error("invalid email", params.User.Email)
		return c.Status(422).JSON(fiber.Map{"errors": fiber.Map{
			"body": []string{"Invalid email"},
		}})

	}

	hashPass, err := hash.HashPassword(params.User.Password)
	if err != nil {
		h.Logger.Error("error hashing password: %v", err)
		return c.SendStatus(500)
	}

	user, err := h.Queries.CreateUser(context.Background(), &sqlc.CreateUserParams{
		Email:    params.User.Email,
		Username: params.User.Username,
		Password: hashPass,
	})
	if err != nil {
		h.Logger.Error("error creating user: %v", err)
		return c.SendStatus(500)
	}

	token, err := jwt2.GenerateToken(user.ID)
	if err != nil {
		h.Logger.Error("error generating token: %v", err)
		return c.SendStatus(500)
	}

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
