package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/subkeep/backend/handlers"
	"github.com/subkeep/backend/middleware"
	"github.com/subkeep/backend/services"
)

// Handlers holds all handler instances used by the router.
type Handlers struct {
	Auth        *handlers.AuthHandler
	AuthService *services.AuthService
}

// SetupRoutes configures all API routes on the Fiber app.
func SetupRoutes(app *fiber.App, h *Handlers) {
	api := app.Group("/api/v1")

	// Auth routes (public).
	auth := api.Group("/auth")
	auth.Post("/:provider/callback", h.Auth.OAuthCallback)
	auth.Post("/refresh", h.Auth.RefreshToken)

	// Auth routes (protected).
	authProtected := auth.Group("", middleware.AuthMiddleware(h.AuthService))
	authProtected.Post("/logout", h.Auth.Logout)
	authProtected.Get("/me", h.Auth.GetMe)

	// Protected routes placeholder for future resource endpoints.
	_ = api.Group("", middleware.AuthMiddleware(h.AuthService))
}
