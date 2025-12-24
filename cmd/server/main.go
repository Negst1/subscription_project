package main

import (
	"fmt"
	"subscribe_project/internal/config"
	"subscribe_project/internal/handlers"
	"subscribe_project/internal/middleware"
	"subscribe_project/internal/repository"
	"subscribe_project/internal/services"
	"subscribe_project/pkg/logger"

	_ "subscribe_project/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

// @title Subscription Service API
// @version 1.0
// @description API для управления подписками
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
// @schemes http

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @accept json
// @produce json
func main() {
	logger.InitLogger("subscription-service")

	logger.Log.Info("Starting subscription service...")

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Log.WithError(err).Fatal("Failed to load configuration")
	}

	logger.Log.WithField("db", cfg.DBName).Info("Connecting to database...")
	db, err := sqlx.Connect("postgres", cfg.GetDBConnectionString())
	if err != nil {
		logger.Log.WithError(err).Fatal("Failed to connect to database")
	}
	defer db.Close()

	logger.Log.Info("Database connection established")

	repo := repository.NewSubscriptionRepository(db)
	logger.Log.Info("Repository initialized")

	svc := services.NewSubscriptionService(repo)
	logger.Log.Info("Service initialized")

	handler := handlers.NewSubscriptionHandler(svc)
	logger.Log.Info("Handlers initialized")

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			logger.Log.WithFields(logrus.Fields{
				"error":  err.Error(),
				"method": c.Method(),
				"path":   c.Path(),
				"ip":     c.IP(),
			}).Error("Unhandled error in request")

			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal server error",
			})
		},
	})

	app.Use(middleware.LoggerMiddleware())
	logger.Log.Info("Middleware registered")

	setupRoutes(app, handler)
	logger.Log.WithField("port", cfg.ServerPort).Info("Routes registered")

	logger.Log.WithField("port", cfg.ServerPort).Info("Starting server...")

	if err := app.Listen(fmt.Sprintf(":%s", cfg.ServerPort)); err != nil {
		logger.Log.WithError(err).Fatal("Failed to start server")
	}
}

func setupRoutes(app *fiber.App, handler *handlers.SubscriptionHandler) {
	logger.Log.Info("Setting up routes...")

	api := app.Group("/api")
	app.Get("/swagger/*", swagger.HandlerDefault)
	app.Get("/swagger/doc.json", func(c *fiber.Ctx) error {
		return c.SendFile("./docs/swagger.json")
	})
	logger.Log.Info("Swagger UI registered at /swagger/index.html")

	api.Post("/subscriptions", handler.CreateSubscription)
	logger.Log.Info("Registered POST /api/subscriptions")

	api.Get("/subscriptions/:id", handler.GetSubscription)
	logger.Log.Info("Registered GET /api/subscriptions/:id")

	api.Put("/subscriptions/:id", handler.UpdateSubscription)
	api.Delete("/subscriptions/:id", handler.DeleteSubscription)
	api.Get("/subscriptions", handler.ListSubscriptions)
	api.Post("/summary", handler.GetSummary)

	app.Get("/health", func(c *fiber.Ctx) error {
		logger.Log.Debug("Health check requested")
		return c.JSON(fiber.Map{"status": "ok"})
	})

	logger.Log.Info("All routes registered")
}
