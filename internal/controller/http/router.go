package http

import (
	"github.com/gofiber/fiber/v2"
	"realworld-fiber-sqlc/internal/controller/http/handlers"
	"realworld-fiber-sqlc/pkg/middleware"
	"realworld-fiber-sqlc/usecase/dto/sqlc"
)

func SetupRoutes(app *fiber.App, dbQueries *sqlc.Queries) {

	handlerBase := handlers.NewHandlerQ(dbQueries)

	api := app.Group("/api")
	//profilesRoute := api.Group("/profile")
	users := api.Group("/users")
	user := api.Group("/user")
	//articlesRoute := api.Group("articles")
	//commentsRoute := articlesRoute.Group("/:slug/comments")

	users.Post("/login", handlerBase.Login)
	users.Post("/", handlerBase.Register)

	user.Get("/", middleware.Protected(), handlerBase.CurrentUser)
	user.Put("/", middleware.Protected(), handlerBase.UpdateProfile)
	//
	//profilesRoute.Get("/:username")
	//
	//articlesRoute.Get("/")
	//articlesRoute.Get("/:slug")
	//
	//app.Get("api/tags")

	//user

	//	profilesRoute
	//
	//profilesRoute.Post("/:username/follow")
	//profilesRoute.Delete("/:username/follow")
	//
	////	articles

	//articlesRoute.Get("/feed")
	//
	//articlesRoute.Post("/")
	//articlesRoute.Put("/:slug")
	//articlesRoute.Delete("/:slug")
	//
	////comments
	//
	//commentsRoute.Post("/")
	//commentsRoute.Get("/")
	//commentsRoute.Delete("/:id")
	//
	////ffv
	//app.Get("/api/articles/:slug/favorite")
	//app.Delete("/api/articles/:slug/favorite")

}
