package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

// Swagger describes Swagger API documentation route.
func Swagger(app *fiber.App) {
	swag := app.Group("/swagger")

	// Swagger document
	swag.Get("*", swagger.New(swagger.Config{
		URL:         "/swagger/doc.json",
		DeepLinking: false,
		// Expand ("list") or Collapse ("none") tag groups by default
		DocExpansion: "list",
		// Prefill OAuth ClientId on Authorize popup
		OAuth: &swagger.OAuthConfig{
			AppName:  "OAuth Provider",
			ClientId: "21bb4edc-05a7-4afc-86f1-2e151e4ba6e2",
		},
		// Ability to change OAuth2 redirect uri location
		OAuth2RedirectUrl: "http://localhost:8000/swagger/oauth2-redirect.html",
	}))
}
