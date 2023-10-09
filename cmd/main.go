package main

import "github.com/gofiber/fiber/v2"

func main() {
	// TODO: handle environment variables

	app := fiber.New()

	// TODO: handle fiber endpoints on http/handler.go
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON("Testing...")
	})

	// TODO: handle container database connection
	// TODO: handle repository initialization

	// TODO: handle use cases instances

	app.Listen(":8001")
}
