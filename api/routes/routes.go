package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pcminh0505/gofiber-casbin/api/controllers"
	"github.com/pcminh0505/gofiber-casbin/infras/database"
	"github.com/pcminh0505/gofiber-casbin/middleware"
)

// Setup register all the route of the app
func Setup(app *fiber.App) {
	// Default route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World! Please go to /swagger for API documentation")
	})

	api := app.Group("/api")                                // Public route
	admin := api.Group("/admin", middleware.AuthorizeJWT()) // Admin route

	// Init Casbin for Role-based Authorization Control (RBAC)
	enforcer := database.Casbin()

	// Authentication Routes
	auth := api.Group("/auth")
	auth.Post("/login", controllers.Login)
	auth.Post("/logout", controllers.Logout)
	auth.Post("/register", controllers.CreateUser(enforcer)) // Backup for dev env - Delete when deploy

	// Users route
	// Public
	user := api.Group("/users", middleware.AuthorizeJWT())
	user.Put("/:id/password", middleware.AuthorizeCasbin(enforcer), controllers.UpdatePassword) // Update password

	// Admin
	adminUser := admin.Group("/users", middleware.AuthorizeJWT())
	adminUser.Get("/", middleware.AuthorizeCasbin(enforcer), controllers.GetUsers)
	adminUser.Get("/:id", middleware.AuthorizeCasbin(enforcer), controllers.GetUser)
	adminUser.Post("/", middleware.AuthorizeCasbin(enforcer), controllers.CreateUser(enforcer))
	adminUser.Put("/:id", middleware.AuthorizeCasbin(enforcer), controllers.UpdateUser(enforcer))
	adminUser.Delete("/:id", middleware.AuthorizeCasbin(enforcer), controllers.DeleteUser(enforcer))
}
