package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/subkeep/backend/handlers"
	"github.com/subkeep/backend/middleware"
	"github.com/subkeep/backend/services"
)

// Handlers holds all handler instances used by the router.
type Handlers struct {
	Auth              *handlers.AuthHandler
	Subscription      *handlers.SubscriptionHandler
	Dashboard         *handlers.DashboardHandler
	Simulation        *handlers.SimulationHandler
	Calendar          *handlers.CalendarHandler
	Category          *handlers.CategoryHandler
	Folder            *handlers.FolderHandler
	ShareGroup        *handlers.ShareGroupHandler
	SubscriptionShare *handlers.SubscriptionShareHandler
	Report            *handlers.ReportHandler
	AuthService       *services.AuthService
}

// SetupRoutes configures all API routes on the Fiber app.
func SetupRoutes(app *fiber.App, h *Handlers) {
	api := app.Group("/api/v1")

	// Auth routes (public).
	auth := api.Group("/auth")
	auth.Get("/dev-login", h.Auth.DevLogin)
	auth.Post("/refresh", h.Auth.RefreshToken)

	// Auth routes (protected) — per-route middleware to avoid blocking public OAuth routes.
	authMw := middleware.AuthMiddleware(h.AuthService)
	auth.Post("/logout", authMw, h.Auth.Logout)
	auth.Get("/me", authMw, h.Auth.GetMe)

	// OAuth routes (public) — /:provider must be last to avoid catching /me, /refresh, etc.
	auth.Get("/:provider", h.Auth.OAuthRedirect)
	auth.Post("/:provider/callback", h.Auth.OAuthCallback)

	// Protected routes.
	protected := api.Group("", middleware.AuthMiddleware(h.AuthService))

	// Subscription routes.
	subs := protected.Group("/subscriptions")
	subs.Get("/", h.Subscription.GetAll)
	subs.Post("/", h.Subscription.Create)
	subs.Get("/duplicates", h.Subscription.CheckDuplicates)
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
	simulation.Post("/combined", h.Simulation.SimulateCombined)
	simulation.Post("/apply", h.Simulation.ApplySimulation)
	simulation.Post("/undo", h.Simulation.UndoSimulation)

	// Calendar routes.
	calendar := protected.Group("/calendar")
	calendar.Get("/monthly", h.Calendar.GetMonthlyCalendar)
	calendar.Get("/daily", h.Calendar.GetDayDetail)
	calendar.Get("/upcoming", h.Calendar.GetUpcomingPayments)

	// Category routes.
	categories := protected.Group("/categories")
	categories.Get("/", h.Category.GetAll)
	categories.Post("/", h.Category.Create)
	categories.Put("/:id", h.Category.Update)
	categories.Delete("/:id", h.Category.Delete)

	// Folder routes.
	folders := protected.Group("/folders")
	folders.Get("/", h.Folder.GetAll)
	folders.Post("/", h.Folder.Create)
	folders.Put("/:id", h.Folder.Update)
	folders.Delete("/:id", h.Folder.Delete)

	// Share group routes.
	shareGroups := protected.Group("/share-groups")
	shareGroups.Get("/", h.ShareGroup.GetAll)
	shareGroups.Get("/:id", h.ShareGroup.GetByID)
	shareGroups.Post("/", h.ShareGroup.Create)
	shareGroups.Put("/:id", h.ShareGroup.Update)
	shareGroups.Delete("/:id", h.ShareGroup.Delete)

	// Report routes.
	reports := protected.Group("/reports")
	reports.Get("/overview", h.Report.GetOverview)

	// Subscription share routes.
	subs.Post("/:id/share", h.SubscriptionShare.Link)
	subs.Get("/:id/share", h.SubscriptionShare.GetBySubscription)

	subscriptionShares := protected.Group("/subscription-shares")
	subscriptionShares.Put("/:id", h.SubscriptionShare.Update)
	subscriptionShares.Delete("/:id", h.SubscriptionShare.Unlink)
}
