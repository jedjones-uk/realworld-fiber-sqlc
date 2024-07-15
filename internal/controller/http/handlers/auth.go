package handlers

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"realworld-fiber-sqlc/pkg/hash"
	"realworld-fiber-sqlc/usecase/dto/sqlc"
	"time"
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

type UserObj struct {
	Email    string `json:"email"`
	Token    string `json:"token"`
	Username string `json:"username"`
	Bio      string `json:"bio"`
	Image    string `json:"image"`
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

	claims := jwt.MapClaims{
		"id":  user.ID,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	loginResponse := UserObj{
		Email:    user.Email,
		Token:    t,
		Username: user.Username,
		Bio:      user.Bio.String,
		Image:    user.Image.String,
	}

	jsonResp := fiber.Map{
		"user": loginResponse,
	}

	return c.Status(fiber.StatusOK).JSON(jsonResp)
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

	//fmt.Println(params.User.Username, params.User.Email, hashPass)

	id, err := h.Queries.CreateUser(context.Background(), sqlc.CreateUserParams{
		Email:    params.User.Email,
		Username: params.User.Username,
		Password: hashPass,
	})

	if err != nil {
		//log.Fatal(err)
		log.Printf("error: %v", err)
		return err
	}

	claims := jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		log.Printf("error: %v", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	loginResponse := UserObj{
		Email:    params.User.Email,
		Token:    t,
		Username: params.User.Username,
		Bio:      "",
		Image:    "",
	}

	jsonResp := fiber.Map{
		"user": loginResponse,
	}

	return c.Status(fiber.StatusOK).JSON(jsonResp)
}
