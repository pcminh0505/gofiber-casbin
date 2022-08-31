package controllers

import (
	"errors"
	"fmt"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/pcminh0505/gofiber-casbin/api/models"
	"github.com/pcminh0505/gofiber-casbin/api/utils"
	"github.com/pcminh0505/gofiber-casbin/infras/database"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserInput struct {
	Username string
	Password string
	Name     string
	Email    string
	Role     string
}

type UpdatePasswordInput struct {
	CurrentPassword string
	NewPassword     string
}

// GetUsers godoc
// @Summary Get all users
// @Tags    users
// @Accept  json
// @Produce json
// @Success 200 {array}  models.User
// @Failure 401 {object} models.Response
// @Failure 404 {object} models.Response
// @Router  /admin/users/ [get]
func GetUsers(c *fiber.Ctx) error {
	var users []models.User
	if err := database.GetAdminDB().Find(&users).Error; err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error":   true,
			"message": "Internal Server Error",
		})
	}

	return c.JSON(users)
}

// GetUser godoc
// @Summary Get a user by ID
// @Tags    users
// @Accept  json
// @Produce json
// @Param   id  path     int true "User ID"
// @Success 200 {object} models.User
// @Failure 404 {object} models.Response
// @Router  /admin/users/{id} [get]
func GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.User
	err := database.GetAdminDB().First(&user, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"error":   true,
			"message": "User not found",
		})
	}

	return c.JSON(user)
}

// CreateUser godoc
// @Summary     Create new user
// @Description Create new user with username, password, name, email, and role
// @Tags        users
// @Param       data body UserInput true "Enter user's info"
// @Accept      json
// @Produce     json
// @Success     200 {object} models.Response
// @Failure     400 {object} models.Response
// @Router      /admin/users/ [post]
func CreateUser(e *casbin.Enforcer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Parse input from request body
		var data UserInput
		if err := c.BodyParser(&data); err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"error":   true,
				"message": "Invalid request params!",
			})
		}

		// If existed user is found, return error
		if count := database.GetAdminDB().
			Where(&models.User{Email: data.Email}).
			First(new(models.User)).
			RowsAffected; count > 0 {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"error":   true,
				"message": "Email is already registered",
			})
		}

		if count := database.GetAdminDB().
			Where(&models.User{Username: data.Username}).
			First(new(models.User)).
			RowsAffected; count > 0 {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"error":   true,
				"message": "Username is already registered",
			})
		}

		// Encrypt password and push to AdminDB
		password, _ := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)

		user := models.User{
			Username: data.Username,
			Password: string(password),
			Name:     data.Name,
			Email:    data.Email,
			Role:     data.Role,
		}

		if user.CreatedAt.IsZero() {
			user.CreatedAt = time.Now()
		}
		user.UpdatedAt = time.Now()

		// Write into user DB
		database.GetAdminDB().Create(&user)
		// Write into Casbin rule DB
		e.AddGroupingPolicy(fmt.Sprint(user.ID), user.Role)

		return c.JSON(fiber.Map{
			"error":   false,
			"message": "New user registered successfully!",
			"user":    user,
		})
	}
}

// UpdateUser godoc
// @Summary     Update user's information
// @Description Update username, name, email, and role
// @Tags        users
// @Param       id   path int       true "User ID"
// @Param       data body UserInput true "Enter user's info"
// @Accept      json
// @Produce     json
// @Success     200 {object} models.Response
// @Failure     400 {object} models.Response
// @Failure     404 {object} models.Response
// @Failure     500 {object} models.Response
// @Router      /admin/users/{id} [put]
func UpdateUser(e *casbin.Enforcer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Parse input from request body
		var data UserInput
		if err := c.BodyParser(&data); err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"error":   true,
				"message": "Invalid request params!",
			})
		}

		db := database.GetAdminDB()

		var user models.User
		id := c.Params("id")

		res := db.Transaction(func(tx *gorm.DB) error {
			// Update role
			_, err := e.UpdateGroupingPolicy([]string{id, user.Role}, []string{id, data.Role})
			if err != nil {
				return err
			}

			// Find and Update user's info
			if err := tx.First(&user, id).
				Updates(models.User{
					Name:     data.Name,
					Email:    data.Email,
					Username: data.Username,
					Role:     data.Role,
				}).Error; err != nil {
				return err
			}

			return nil
		})

		if res != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{
				"error":   true,
				"message": "Error updating userID: " + id,
			})
		}

		return c.JSON(fiber.Map{
			"error":   false,
			"message": "Update user successfully!",
			"user":    user,
		})
	}
}

// DeleteUser godoc
// @Summary Delete user
// @Tags    users
// @Accept  json
// @Produce json
// @Param   id  path     int true "User ID"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router  /admin/users/{id} [delete]
func DeleteUser(e *casbin.Enforcer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		db := database.GetAdminDB()

		res := db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Delete(models.User{}, id).Error; err != nil {
				return err
			}

			_, err := e.RemoveGroupingPolicy(id)
			if err != nil {
				return err
			}

			return nil
		})

		if res != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{
				"error":   true,
				"message": "Error when deleting userID: " + id,
			})
		}

		return c.JSON(fiber.Map{
			"error":   false,
			"message": "Delete user successfully!",
		})
	}
}

// UpdatePassword godoc
// @Summary     Update user's password
// @Description Update password
// @Tags        users
// @Param       id   path int                 true "User ID"
// @Param       data body UpdatePasswordInput true "Enter user's info"
// @Accept      json
// @Produce     json
// @Success     200 {object} models.Response
// @Failure     400 {object} models.Response
// @Failure     404 {object} models.Response
// @Failure     500 {object} models.Response
// @Router      /users/{id}/password [put]
func UpdatePassword(c *fiber.Ctx) error {

	// Parse input from request body
	var data UpdatePasswordInput
	if err := c.BodyParser(&data); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request params!",
		})
	}

	var user models.User
	id := c.Params("id")

	if utils.IsEmpty(data.CurrentPassword) || utils.IsEmpty(data.NewPassword) {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"error":   true,
			"message": "Cannot proceed with empty input!",
		})
	}

	if data.CurrentPassword == data.NewPassword {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"error":   true,
			"message": "Updated password must be different from the current one!",
		})
	}

	db := database.GetAdminDB()

	// If user is not found, return error
	if err := db.First(&user, id).Error; err != nil {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"error":   true,
			"message": "User not found!",
		})
	}

	// If current password is incorrect, return error
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.CurrentPassword)); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"error":   true,
			"message": "Incorrect password!",
		})
	}

	// Update password
	newPassword, _ := bcrypt.GenerateFromPassword([]byte(data.NewPassword), bcrypt.DefaultCost)
	if err := db.Model(&user).Update("password", string(newPassword)).Error; err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error":   true,
			"message": "Error updating password userID: " + id,
		})
	}

	fmt.Println(data)

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Update password successfully!",
		"user":    user,
	})
}
