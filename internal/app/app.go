package app

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"realworld-fiber-sqlc/pkg/logger"
)

func Run() {

	// logger
	l := logger.New("debug")

	// database
	connString := "host=localhost port=5432 user=postgres password=postgres dbname=realworld sslmode=disable"
	dbpool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		l.Fatal("Unable to connect to database", err)
	}
	defer dbpool.Close()

	//	http server
	app := fiber.New()
	err = app.Listen(":4000")
	if err != nil {
		return
	}
}
