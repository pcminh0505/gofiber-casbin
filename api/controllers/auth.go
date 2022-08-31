package controllers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/pcminh0505/gofiber-casbin/api/models"
	"github.com/pcminh0505/gofiber-casbin/api/utils"
	"github.com/pcminh0505/gofiber-casbin/infras/database"
	"golang.org/x/crypto/bcrypt"
)

type AuthInput struct {
	Identity string
	Password string
}

// Login godoc
// @Summary     Login a user
// @Description Login with username/email and password, return a cookie
// @Tags        auth
// @Param       data body AuthInput true "Login with Username and Password"
// @Accept      json
// @Produce     json
// @Success     200 {object} models.Response
// @Failure     400 {object} models.Response
// @Failure     404 {object} models.Response
// @Router      /auth/login [post]
func Login(c *fiber.Ctx) error {
	// Parse input from request body
	var data AuthInput
	if err := c.BodyParser(&data); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request params!",
		})
	}

	if utils.IsEmpty(data.Identity) || utils.IsEmpty(data.Password) {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"error":   true,
			"message": "Cannot login with empty input!",
		})
	}

	var user models.User
	// If user is not found, return error
	if res := database.GetAdminDB().Where(
		&models.User{Email: data.Identity}).Or(
		&models.User{Username: data.Identity},
	).First(&user); res.RowsAffected <= 0 {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"error":   true,
			"message": "User not found!",
		})
	}

	// If password is incorrect, return error
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password)); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"error":   true,
			"message": "Incorrect password!",
		})
	}

	// Create JWT token with userID.
	token, err := utils.GenerateJWT(string(c.Request().Host()), user.ID)

	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error":   true,
			"message": "Internal Server Error",
		})
	}

	// Create cookie
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 1),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(user)
}

// Logout godoc
// @Summary     Logout a user
// @Description Logout by overriding cookie expired time
// @Tags        auth
// @Accept      json
// @Produce     json
// @Success     200 {object} models.Response
// @Router      /auth/logout [post]
func Logout(c *fiber.Ctx) error {
	// Override expired time of cookie
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Logout successfully!",
	})
}
