package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgtype"
	"realworld-fiber-sqlc/usecase/database"
	"realworld-fiber-sqlc/usecase/database/sqlc"
)

type UpdateProfileParams struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	Image    string `json:"image"`
	Bio      string `json:"bio"`
}

func GetUser(c *fiber.Ctx) error {
	var id int64
	//	TODO get user id from token
	db := sqlc.New(database.DB)
	user, err := db.GetUser(c.Context(), id)
	if err != nil {
		return err
	}

	return c.Status(200).JSON(user)
}

func UpdateProfile(c *fiber.Ctx) error {
	var id int64
	//	TODO get user id from token

	var params UpdateProfileParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	updateProfileParams := sqlc.UpdateProfileParams{
		UserID: pgtype.Int8{Int64: id},
		Bio:    pgtype.Text{String: params.Bio},
		Image:  pgtype.Text{String: params.Image},
	}

	updateUserParams := sqlc.UpdateUserParams{
		ID:       id,
		Email:    params.Email,
		Username: params.Username,
		Password: params.Password,
	}

	if err := database.UpdateProfileTX(database.DB, &updateProfileParams, &updateUserParams); err != nil {
		return err
	}

	return c.Status(200).JSON(fiber.Map{"message": "Profile updated"})
}
