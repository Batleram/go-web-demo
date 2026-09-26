package core

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"demo/core/templates"
)

func RegisterRoutes() *fiber.App {
	router := fiber.New()

	router.Static("/statics", "./statics")

	router.Get("/", indexHandler)
	router.Get("index.html", indexHandler)

	return router
}

func indexHandler(c *fiber.Ctx) error {
	c.Set("Content-type", "text/html")
	return templates.Index(nil).Render(context.Background(), c.Response().BodyWriter())
}
