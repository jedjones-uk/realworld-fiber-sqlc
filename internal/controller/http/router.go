package http

import (
	"github.com/gofiber/fiber/v2"
	"realworld-fiber-sqlc/internal/controller/http/handlers"
)

func NewRouter(app *fiber.App) {

	api := app.Group("/api")

	//auth
	users := api.Group("users")
	users.Post("/login")
	users.Post("/")

	//user
	user := api.Group("user")
	user.Get("/", handlers.GetUser)
	user.Put("/", handlers.UpdateProfile)

	//	profilesRoute
	profilesRoute := api.Group("profilesRoute")
	profilesRoute.Get("/:username")
	profilesRoute.Post("/:username/follow")
	profilesRoute.Delete("/:username/follow")

	//	articles
	articlesRoute := api.Group("articles")
	articlesRoute.Get("/")
	articlesRoute.Get("/feed")
	articlesRoute.Get("/:slug")
	articlesRoute.Post("/")
	articlesRoute.Put("/:slug")
	articlesRoute.Delete("/:slug")

	//comments
	commentsRoute := articlesRoute.Group("/:slug/comments")
	commentsRoute.Post("/")
	commentsRoute.Get("/")
	commentsRoute.Delete("/:id")

	//fav
	app.Get("/api/articles/:slug/favorite")
	app.Delete("/api/articles/:slug/favorite")

	app.Get("api/tags")
}
