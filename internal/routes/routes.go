package routes

import (
	"github.com/XSidik/go-fiber-template/internal/controllers"
	"github.com/XSidik/go-fiber-template/internal/middleware"
	"github.com/XSidik/go-fiber-template/internal/models/dto"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(dto.APIResponse{
			Status:  true,
			Message: "Hello, this is an example API template of golang with fiber",
		})
	})

	// Auth routes
	auth := app.Group("api/v1/auth")
	auth.Post("/register", controllers.Register)
	auth.Post("/login", controllers.Login)
	auth.Get("/logout", middleware.AuthMiddleware, controllers.Logout)
	auth.Get("/refresh-token", middleware.AuthMiddleware, controllers.RefreshToken)
}
