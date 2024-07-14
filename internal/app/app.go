package app

import (
	"github.com/gofiber/fiber/v2"
)

func Run() {
	// logger
	//l := logger.New("debug")

	//	http server
	app := fiber.New()
	app.Listen(":4000")
}
