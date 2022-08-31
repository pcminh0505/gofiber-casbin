package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pcminh0505/gofiber-casbin/docs"
)

// SwaggerMiddleware provide configuration Swagger generated document.
func SwaggerMiddleware(a *fiber.App) {
	a.Use(func(c *fiber.Ctx) error {
		req := c.Request()

		// Config SwaggerInfo
		docs.SwaggerInfo.Title = "Admin Client REST API Documentation"
		docs.SwaggerInfo.Description = "API for administration dashboard"
		docs.SwaggerInfo.Host = string(req.Host()) // Dynamic Host based on URL
		docs.SwaggerInfo.Version = "1.0"           // API version control in .env in the future
		docs.SwaggerInfo.BasePath = "/api/"

		return c.Next()
	})
}
