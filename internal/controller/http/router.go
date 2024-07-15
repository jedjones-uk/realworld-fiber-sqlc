package http

import (
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"realworld-fiber-sqlc/internal/controller/http/handlers"
	"realworld-fiber-sqlc/usecase/dto/sqlc"
)

func NewRouter(app *fiber.App, dbQueries *sqlc.Queries) {

	handlerBase := handlers.NewHandlerQ(dbQueries)

	api := app.Group("/api")
	//profilesRoute := api.Group("/profile")
	users := api.Group("/users")
	user := api.Group("/user")
	//articlesRoute := api.Group("articles")
	//commentsRoute := articlesRoute.Group("/:slug/comments")

	//auth

	users.Post("/login", handlerBase.Login)
	users.Post("/", handlerBase.Register)

	user.Get("/", handlerBase.GetUser)
	//user.Put("/", handlers.UpdateProfile)
	//
	//profilesRoute.Get("/:username")
	//
	//articlesRoute.Get("/")
	//articlesRoute.Get("/:slug")
	//
	//app.Get("api/tags")

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			Key: []byte("secret"),
		},
	}))

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
