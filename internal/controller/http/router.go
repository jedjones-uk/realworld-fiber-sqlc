package http

import (
	"github.com/gofiber/fiber/v2"
	"realworld-fiber-sqlc/internal/controller/http/handlers"
	"realworld-fiber-sqlc/pkg/logger"
	"realworld-fiber-sqlc/pkg/middleware"
	"realworld-fiber-sqlc/usecase/dto/sqlc"
)

func SetupRoutes(app *fiber.App, dbQueries *sqlc.Queries, l *logger.Logger) {

	handlerBase := handlers.NewHandlerQ(dbQueries, l)
	api := app.Group("/api")

	//auth
	users := api.Group("/users")
	users.Post("/login", handlerBase.Login)
	users.Post("/", handlerBase.Register)

	//user
	user := api.Group("/user")
	user.Get("/", middleware.Protected(), handlerBase.CurrentUser)
	user.Put("/", middleware.Protected(), handlerBase.UpdateProfile)

	//profiles
	profilesRoute := api.Group("/profiles")
	profilesRoute.Get("/:username", handlerBase.GetProfile)

	//follow
	profilesRoute.Post("/:username/follow", middleware.Protected(), handlerBase.Follow)
	profilesRoute.Delete("/:username/follow", middleware.Protected(), handlerBase.Unfollow)

	//	articles
	articlesRoute := api.Group("/articles")
	articlesRoute.Get("/:slug", handlerBase.GetArticle)
	articlesRoute.Get("/", handlerBase.GetArticles)
	articlesRoute.Get("/feed", middleware.Protected(), handlerBase.Feed)
	articlesRoute.Post("/", middleware.Protected(), handlerBase.CreateArticle)
	articlesRoute.Put("/:slug", middleware.Protected(), handlerBase.UpdateArticle)
	articlesRoute.Delete("/:slug", middleware.Protected(), handlerBase.DeleteArticle)

	//comments
	commentsRoute := articlesRoute.Group("/:slug/comments")
	commentsRoute.Post("/", middleware.Protected(), handlerBase.CreateComment)
	commentsRoute.Get("/", handlerBase.GetComments)
	commentsRoute.Delete("/:id", middleware.Protected(), handlerBase.DeleteComment)

	//ffv
	app.Post("/api/articles/:slug/favorite", middleware.Protected(), handlerBase.FavoriteArticle)
	app.Delete("/api/articles/:slug/favorite", middleware.Protected(), handlerBase.UnfavoriteArticle)

	app.Get("api/tags", handlerBase.GetTags)
}
