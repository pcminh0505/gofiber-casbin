package routes

import "github.com/gofiber/fiber/v2"

// NotFoundRoute describes 404 Error route.
func NotFoundRoute(app *fiber.App) {
	// Register new special route.
	app.Use(
		// Anonymous function.
		func(c *fiber.Ctx) error {
			// Return HTTP 404 status and JSON response.
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   true,
				"message": "API Endpoint is unavailable",
			})
		},
	)
}
