package middleware

import (
	"fmt"

	"github.com/casbin/casbin/v2"

	"github.com/gofiber/fiber/v2"
)

func AuthorizeCasbin(e *casbin.Enforcer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get current user/subject
		userID, ok := c.Locals("userID").(string)

		if userID == "" || !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   true,
				"message": "Current logged in user not found!",
			})
		}

		// Load policy from Database
		err := e.LoadPolicy()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   true,
				"message": "Failed to load Casbin policy!",
			})
		}

		// Casbin enforces policy
		accepted, err := e.Enforce(fmt.Sprint(userID), c.OriginalURL(), c.Method()) // id - url - method || 1 - /api/admin/users - GET

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   true,
				"message": "Error when authorizing user's accessibility",
			})
		}

		if !accepted {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error":   true,
				"message": "Unauthorized!",
			})
		}
		return c.Next()
	}
}
