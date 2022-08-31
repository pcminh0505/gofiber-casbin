package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pcminh0505/gofiber-casbin/api/utils"
)

// AuthorizeJWT returns a middleware which secures all the private routes
func AuthorizeJWT() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		// Parse jwt token from cookie
		cookie := c.Cookies("jwt")

		publicKey := utils.LoadEcdsaPrivateKeyKey().PublicKey
		// Verify with public key
		token, err := jwt.ParseWithClaims(cookie, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
			return &publicKey, nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   true,
				"message": "Error when parsing JWT!",
			})
		}
		// Get userID inside cookie subject
		claims := token.Claims.(*jwt.RegisteredClaims)

		if token.Valid {
			if claims.ExpiresAt.Unix() < time.Now().Unix() {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error":   true,
					"message": "JWT Expired!",
				})
			}
		} else if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				// Malformed token -> Delete Cookie
				c.ClearCookie("jwt")
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error":   true,
					"message": "Missing or malformed JWT!",
				})
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				// Token is either expired or not active yet
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error":   true,
					"message": "Unauthorized token!",
				})
			} else {
				// Cannot handle -> Delete Cookie
				c.ClearCookie("jwt")
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"error":   true,
					"message": "Error when processing identity!",
				})
			}
		}

		// Store current userID into Fiber Context Locals
		c.Locals("userID", claims.Subject)
		return c.Next()
	}
}
