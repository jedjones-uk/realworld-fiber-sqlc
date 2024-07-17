package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"log"
	"realworld-fiber-sqlc/internal/controller/http"
	"realworld-fiber-sqlc/usecase/dto"
	"realworld-fiber-sqlc/usecase/dto/sqlc"
)

func Run() {

	// logger
	//connString := "host=localhost port=5432 user=postgres password=postgres dbname=realworld sslmode=disable"
	//maxRetries := 10
	//retryInterval := 5 * time.Second

	var err error
	pool, err := dto.NewPool()
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()
	dbQueries := sqlc.New(pool)

	log.Printf("Connected to the database")

	//handler := handlers.NewHandlerQ(dbQueries)

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, HEAD, PUT, PATCH, POST, DELETE",
	}))

	http.SetupRoutes(app, dbQueries)

	//app.Use(l)
	app.Listen(":3000")

}
