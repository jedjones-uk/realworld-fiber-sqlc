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

	//articlesRoute := api.Group("articles")
	//commentsRoute := articlesRoute.Group("/:slug/comments")
	users := api.Group("/users")

	users.Post("/login", handlerBase.Login)
	users.Post("/", handlerBase.Register)

	user := api.Group("/user")

	user.Get("/", middleware.Protected(), handlerBase.CurrentUser)
	user.Put("/", middleware.Protected(), handlerBase.UpdateProfile)
	//

	profilesRoute := api.Group("/profiles")
	profilesRoute.Get("/:username", handlerBase.GetProfile)
	//
	//articlesRoute.Get("/")
	//articlesRoute.Get("/:slug")
	//
	//app.Get("api/tags")

	//user

	//profilesRoute
	//
	profilesRoute.Post("/:username/follow", middleware.Protected(), handlerBase.Follow)
	profilesRoute.Delete("/:username/follow", middleware.Protected(), handlerBase.Unfollow)
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
