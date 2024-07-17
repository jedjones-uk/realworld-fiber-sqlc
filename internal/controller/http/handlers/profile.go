package handlers

import (
	"github.com/gofiber/fiber/v2"
	"realworld-fiber-sqlc/usecase/dto/sqlc"
)

type ProfileResp struct {
	Profile struct {
		Username  string `json:"username"`
		Bio       string `json:"bio"`
		Image     string `json:"image"`
		Following bool   `json:"following"`
	} `json:"profile"`
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

	profile, err := h.Queries.GetUserProfile(c.Context(), arg)
	if err != nil {
		return err
	}

	var resp ProfileResp
	resp.Profile.Username = profile.Username
	resp.Profile.Bio = profile.Bio.String
	resp.Profile.Image = profile.Image.String
	resp.Profile.Following = profile.Following

	return c.Status(200).JSON(resp)

}

func (h *HandlerBase) Follow(c *fiber.Ctx) error {
	username := c.Params("username")
	userID := userIDFromToken(c)

	err := h.Queries.FollowUser(c.Context(), sqlc.FollowUserParams{
		Username:   username,
		FollowerID: userID,
	})
	if err != nil {
		return err
	}

	profile, err := h.Queries.GetUserProfile(c.Context(), sqlc.GetUserProfileParams{
		Username:   username,
		FollowerID: userID,
	})
	if err != nil {
		return err
	}

	var resp ProfileResp
	resp.Profile.Username = profile.Username
	resp.Profile.Bio = profile.Bio.String
	resp.Profile.Image = profile.Image.String
	resp.Profile.Following = profile.Following

	return c.Status(200).JSON(resp)
}

func (h *HandlerBase) Unfollow(c *fiber.Ctx) error {
	username := c.Params("username")
	userID := userIDFromToken(c)

	err := h.Queries.UnfollowUser(c.Context(), sqlc.UnfollowUserParams{
		Username:   username,
		FollowerID: userID,
	})
	if err != nil {
		return err
	}

	profile, err := h.Queries.GetUserProfile(c.Context(), sqlc.GetUserProfileParams{
		Username:   username,
		FollowerID: userID,
	})
	if err != nil {
		return err
	}

	var resp ProfileResp
	resp.Profile.Username = profile.Username
	resp.Profile.Bio = profile.Bio.String
	resp.Profile.Image = profile.Image.String
	resp.Profile.Following = profile.Following

	return c.Status(200).JSON(resp)
}
