package main

import (
	"github.com/XSidik/go-fiber-template/internal/config"
	"github.com/XSidik/go-fiber-template/internal/database"
	"github.com/XSidik/go-fiber-template/internal/models"
	"github.com/XSidik/go-fiber-template/internal/routes"
	"github.com/gofiber/fiber/v2"
)

func main() {
	config := config.GetConfig()

	app := fiber.New()

	database.Connect(config)
	models.Migrate(database.DB)

	routes.SetupRoutes(app)

	app.Listen(":3000")
}
