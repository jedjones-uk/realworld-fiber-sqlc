package handlers

import (
	"github.com/gofiber/fiber/v2"
	"realworld-fiber-sqlc/usecase/dto/sqlc"
)

type Profile struct {
	Username  string `json:"username"`
	Bio       string `json:"bio"`
	Image     string `json:"image"`
	Following bool   `json:"following"`
}

type Author struct {
	Username  string `json:"username"`
	Bio       string `json:"bio"`
	Image     string `json:"image"`
	Following bool   `json:"following"`
}

func (h *HandlerBase) GetProfile(c *fiber.Ctx) error {
	username := c.Params("username")
	userID := userIDFromToken(c)

	var arg sqlc.GetUserProfileParams
	if userID == 0 {
		arg = sqlc.GetUserProfileParams{
			Username: username,
		}
	} else {
		arg = sqlc.GetUserProfileParams{
			Username:   username,
			FollowerID: userID,
		}
	}

	profile, err := h.Queries.GetUserProfile(c.Context(), &arg)
	if err != nil {
		return err
	}

	return c.Status(200).JSON(fiber.Map{"profile": profile})

}

func (h *HandlerBase) Follow(c *fiber.Ctx) error {
	username := c.Params("username")
	userID := userIDFromToken(c)

	err := h.Queries.FollowUser(c.Context(), &sqlc.FollowUserParams{
		Username:   username,
		FollowerID: userID,
	})
	if err != nil {
		return err
	}

	profile, err := h.Queries.GetUserProfile(c.Context(), &sqlc.GetUserProfileParams{
		Username:   username,
		FollowerID: userID,
	})
	if err != nil {
		return err
	}
	profile.Following = true

	return c.Status(200).JSON(fiber.Map{
		"profile": Profile{
			Username:  profile.Username,
			Bio:       profile.Bio.String,
			Image:     profile.Image.String,
			Following: true,
		},
	})
}

func (h *HandlerBase) Unfollow(c *fiber.Ctx) error {
	username := c.Params("username")
	userID := userIDFromToken(c)

	err := h.Queries.UnfollowUser(c.Context(), &sqlc.UnfollowUserParams{
		Username:   username,
		FollowerID: userID,
	})
	if err != nil {
		return err
	}

	profile, err := h.Queries.GetUserProfile(c.Context(), &sqlc.GetUserProfileParams{
		Username:   username,
		FollowerID: userID,
	})
	if err != nil {
		return err
	}

	return c.Status(200).JSON(fiber.Map{
		"profile": Profile{
			Username:  profile.Username,
			Bio:       profile.Bio.String,
			Image:     profile.Image.String,
			Following: false,
		},
	})
}
