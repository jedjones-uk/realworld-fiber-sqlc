package main

import (
	"realworld-fiber-sqlc/internal/app"
)

func main() {
	realWorld := app.New()
	realWorld.Listen(":3000")
}
