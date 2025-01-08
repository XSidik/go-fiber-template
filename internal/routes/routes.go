package routes

import (
	"github.com/XSidik/go-fiber-template/internal/controllers"
	"github.com/XSidik/go-fiber-template/internal/middleware"
	"github.com/XSidik/go-fiber-template/internal/models"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(models.APIResponse{
			Status:  true,
			Message: "Hello, this is an example API template of golang with fiber",
		})
	})

	// Auth routes
	auth := app.Group("api/v1/auth")
	auth.Post("/register", controllers.Register)
	auth.Post("/login", controllers.Login)
	app.Use(middleware.AuthMiddleware)
	auth.Post("/logout", controllers.Logout)
	auth.Get("/refresh", controllers.RefreshToken)
}
