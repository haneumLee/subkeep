package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/subkeep/backend/handlers"
	"github.com/subkeep/backend/middleware"
	"github.com/subkeep/backend/services"
)

// Handlers holds all handler instances used by the router.
type Handlers struct {
	Auth         *handlers.AuthHandler
	Subscription *handlers.SubscriptionHandler
	Dashboard    *handlers.DashboardHandler
	Simulation   *handlers.SimulationHandler
	Category     *handlers.CategoryHandler
	ShareGroup   *handlers.ShareGroupHandler
	AuthService  *services.AuthService
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

	// Protected routes.
	protected := api.Group("", middleware.AuthMiddleware(h.AuthService))

	// Subscription routes.
	subs := protected.Group("/subscriptions")
	subs.Get("/", h.Subscription.GetAll)
	subs.Post("/", h.Subscription.Create)
	subs.Get("/:id", h.Subscription.GetByID)
	subs.Put("/:id", h.Subscription.Update)
	subs.Delete("/:id", h.Subscription.Delete)
	subs.Patch("/:id/satisfaction", h.Subscription.UpdateSatisfaction)

	// Dashboard routes.
	dashboard := protected.Group("/dashboard")
	dashboard.Get("/summary", h.Dashboard.GetSummary)
	dashboard.Get("/recommendations", h.Dashboard.GetRecommendations)

	// Simulation routes.
	simulation := protected.Group("/simulation")
	simulation.Post("/cancel", h.Simulation.SimulateCancel)
	simulation.Post("/add", h.Simulation.SimulateAdd)
	simulation.Post("/apply", h.Simulation.ApplySimulation)

	// Category routes.
	categories := protected.Group("/categories")
	categories.Get("/", h.Category.GetAll)
	categories.Post("/", h.Category.Create)
	categories.Put("/:id", h.Category.Update)
	categories.Delete("/:id", h.Category.Delete)

	// Share group routes.
	shareGroups := protected.Group("/share-groups")
	shareGroups.Get("/", h.ShareGroup.GetAll)
	shareGroups.Get("/:id", h.ShareGroup.GetByID)
	shareGroups.Post("/", h.ShareGroup.Create)
	shareGroups.Put("/:id", h.ShareGroup.Update)
	shareGroups.Delete("/:id", h.ShareGroup.Delete)
}
