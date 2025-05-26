package main

import (
	"log"

	_ "github.com/XSidik/go-fiber-template/docs"
	"github.com/XSidik/go-fiber-template/internal/config"
	"github.com/XSidik/go-fiber-template/internal/database"
	"github.com/XSidik/go-fiber-template/internal/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

// @title go-fiber template API
// @version 1.0
// @description This is an API template developed using the GoFiber framework in Golang, with PostgreSQL as the database, JWT for authentication, and Redis for caching.

// @contact.name go-fiber template API Support
// @contact.github https://github.com/XSidik

// @host localhost:3000
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name bearer
func main() {
	config := config.GetConfig()

	app := fiber.New()

	app.Get("/swagger/*", swagger.HandlerDefault) // default

	database.Connect(config)
	database.Migrate(database.DB)

	// Run all seeders
	if err := database.RunSeeders(database.DB); err != nil {
		log.Fatal(err)
	}

	routes.SetupRoutes(app)

	app.Listen(":3000")
}
