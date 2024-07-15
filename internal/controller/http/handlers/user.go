package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type UpdateProfileParams struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	Image    string `json:"image"`
	Bio      string `json:"bio"`
}

func getUserIDFromToken(c *fiber.Ctx) int64 {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	id := claims["id"].(int64)
	return id
}

func (h *HandlerBase) GetUser(c *fiber.Ctx) error {
	var id int64
	id = getUserIDFromToken(c)
	user, err := h.Queries.GetUser(c.Context(), id)
	if err != nil {
		return err
	}

	return c.Status(200).JSON(user)
}

//func UpdateProfile(c *fiber.Ctx) error {
//	var id int64
//	//	TODO get user id from token
//
//	var params UpdateProfileParams
//	if err := c.BodyParser(&params); err != nil {
//		return err
//	}
//
//	updateProfileParams := sqlc.UpdateProfileParams{
//		UserID: pgtype.Int8{Int64: id},
//		Bio:    pgtype.Text{String: params.Bio},
//		Image:  pgtype.Text{String: params.Image},
//	}
//
//	updateUserParams := sqlc.UpdateUserParams{
//		ID:       id,
//		Email:    params.Email,
//		Username: params.Username,
//		Password: params.Password,
//	}
//
//	if err := dto.UpdateProfileTX(dto.DB, &updateProfileParams, &updateUserParams); err != nil {
//		return err
//	}
//
//	return c.Status(200).JSON(fiber.Map{"message": "Profile updated"})
//}
