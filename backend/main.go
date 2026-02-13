package main

import (
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"

	"github.com/subkeep/backend/config"
	"github.com/subkeep/backend/database"
	"github.com/subkeep/backend/handlers"
	"github.com/subkeep/backend/models"
	"github.com/subkeep/backend/repositories"
	"github.com/subkeep/backend/routes"
	"github.com/subkeep/backend/services"
	"github.com/subkeep/backend/utils"
)

func main() {
	// Load configuration.
	cfg := config.Load()
	setupLogger(cfg)

	slog.Info("starting SubKeep backend",
		"env", cfg.Server.Env,
		"port", cfg.Server.Port,
	)

	// Connect to database.
	_, err := database.Connect(&cfg.Database)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer func() {
		if closeErr := database.Close(); closeErr != nil {
			slog.Error("failed to close database connection", "error", closeErr)
		}
	}()

	// Initialize Fiber app.
	app := fiber.New(fiber.Config{
		AppName:               "SubKeep API v1",
		ReadTimeout:           15 * time.Second,
		WriteTimeout:          15 * time.Second,
		IdleTimeout:           60 * time.Second,
		BodyLimit:             4 * 1024 * 1024, // 4MB
		DisableStartupMessage: cfg.IsProduction(),
		ErrorHandler:          errorHandler,
	})

	// Global middleware.
	app.Use(recover.New(recover.Config{
		EnableStackTrace: cfg.IsDevelopment(),
	}))
	app.Use(requestid.New())
	app.Use(fiberlogger.New(fiberlogger.Config{
		Format:     "${time} | ${status} | ${latency} | ${ip} | ${method} | ${path} | ${error}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Asia/Seoul",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Join(cfg.CORS.AllowedOrigins, ","),
		AllowMethods:     strings.Join(cfg.CORS.AllowedMethods, ","),
		AllowHeaders:     strings.Join(cfg.CORS.AllowedHeaders, ","),
		AllowCredentials: true,
		MaxAge:           86400,
	}))

	// Security headers.
	app.Use(helmet.New())

	// Global rate limiter: limit per IP per minute.
	app.Use(limiter.New(limiter.Config{
		Max:               cfg.RateLimit.RequestsPerMinute,
		Expiration:        1 * time.Minute,
		LimitReached: func(c *fiber.Ctx) error {
			slog.Warn("rate limit exceeded",
				"ip", c.IP(),
				"path", c.Path(),
			)
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":   "Too Many Requests",
				"message": "요청이 너무 많습니다. 잠시 후 다시 시도해주세요.",
			})
		},
		KeyGenerator: func(c *fiber.Ctx) string {
			// Cloudflare real IP (CF-Connecting-IP > X-Real-IP > RemoteIP)
			if cfIP := c.Get("CF-Connecting-IP"); cfIP != "" {
				return cfIP
			}
			if realIP := c.Get("X-Real-IP"); realIP != "" {
				return realIP
			}
			return c.IP()
		},
		SkipFailedRequests:     false,
		SkipSuccessfulRequests: false,
	}))

	// Auto-migrate database tables.
	db := database.DB
	if err := models.AutoMigrateAll(db); err != nil {
		slog.Error("failed to auto-migrate database", "error", err)
		os.Exit(1)
	}
	if err := models.SeedDefaultCategories(db); err != nil {
		slog.Error("failed to seed default categories", "error", err)
		os.Exit(1)
	}

	// Initialize repositories.
	userRepo := repositories.NewUserRepository(db)
	subRepo := repositories.NewSubscriptionRepository(db)
	catRepo := repositories.NewCategoryRepository(db)
	folderRepo := repositories.NewFolderRepository(db)
	shareGroupRepo := repositories.NewShareGroupRepository(db)
	subShareRepo := repositories.NewSubscriptionShareRepository(db)

	// Initialize services.
	authService := services.NewAuthService(userRepo, cfg.JWT)
	oauthService := services.NewOAuthService(cfg.OAuth)
	subService := services.NewSubscriptionService(subRepo)
	dashboardService := services.NewDashboardService(subRepo, subShareRepo)
	simService := services.NewSimulationService(subRepo, subShareRepo)
	calendarService := services.NewCalendarService(subRepo, subShareRepo)
	catService := services.NewCategoryService(catRepo)
	folderService := services.NewFolderService(folderRepo)
	shareGroupService := services.NewShareGroupService(shareGroupRepo)
	subShareService := services.NewSubscriptionShareService(subShareRepo, subRepo, shareGroupRepo)
	reportService := services.NewReportService(subRepo, subShareRepo)

	// Initialize handlers.
	authHandler := handlers.NewAuthHandler(authService, oauthService)
	subHandler := handlers.NewSubscriptionHandler(subService)
	dashboardHandler := handlers.NewDashboardHandler(dashboardService)
	simHandler := handlers.NewSimulationHandler(simService)
	calendarHandler := handlers.NewCalendarHandler(calendarService)
	catHandler := handlers.NewCategoryHandler(catService)
	folderHandler := handlers.NewFolderHandler(folderService)
	shareGroupHandler := handlers.NewShareGroupHandler(shareGroupService)
	subShareHandler := handlers.NewSubscriptionShareHandler(subShareService)
	reportHandler := handlers.NewReportHandler(reportService)

	// Health check endpoint.
	app.Get("/health", func(c *fiber.Ctx) error {
		dbStatus := "up"
		if dbErr := database.HealthCheck(); dbErr != nil {
			dbStatus = "down"
		}

		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "subkeep-api",
			"version": "1.0.0",
			"database": fiber.Map{
				"status": dbStatus,
			},
		})
	})

	// Register route handlers.
	routes.SetupRoutes(app, &routes.Handlers{
		Auth:              authHandler,
		Subscription:      subHandler,
		Dashboard:         dashboardHandler,
		Simulation:        simHandler,
		Calendar:          calendarHandler,
		Category:          catHandler,
		Folder:            folderHandler,
		ShareGroup:        shareGroupHandler,
		SubscriptionShare: subShareHandler,
		Report:            reportHandler,
		AuthService:       authService,
	})

	// Graceful shutdown.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		addr := cfg.Server.Host + ":" + cfg.Server.Port
		slog.Info("server listening", "address", addr)
		if listenErr := app.Listen(addr); listenErr != nil {
			slog.Error("server failed to start", "error", listenErr)
			os.Exit(1)
		}
	}()

	sig := <-quit
	slog.Info("shutting down server", "signal", sig.String())

	shutdownTimeout := 10 * time.Second
	if err := app.ShutdownWithTimeout(shutdownTimeout); err != nil {
		slog.Error("server forced shutdown", "error", err)
	}

	slog.Info("server stopped gracefully")
}

// errorHandler is the custom Fiber error handler.
func errorHandler(c *fiber.Ctx, err error) error {
	// Handle AppError.
	if appErr, ok := err.(*utils.AppError); ok {
		return utils.Error(c, appErr)
	}

	// Handle Fiber errors.
	if fiberErr, ok := err.(*fiber.Error); ok {
		return utils.ErrorWithStatus(c, fiberErr.Code, fiberErr.Message)
	}

	// Unexpected errors.
	slog.Error("unhandled error",
		"error", err.Error(),
		"path", c.Path(),
		"method", c.Method(),
	)
	return utils.Error(c, utils.ErrInternal(""))
}

// setupLogger configures the global slog logger.
func setupLogger(cfg *config.Config) {
	var level slog.Level
	switch strings.ToLower(cfg.Log.Level) {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn", "warning":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: cfg.IsDevelopment(),
	}

	var handler slog.Handler
	if cfg.Log.Format == "json" || cfg.IsProduction() {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)
}
