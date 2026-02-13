package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	_ "modernc.org/sqlite"

	"github.com/sm2-cosign/backend/internal/config"
	"github.com/sm2-cosign/backend/internal/handler"
	"github.com/sm2-cosign/backend/internal/middleware"
	"github.com/sm2-cosign/backend/internal/repository"
	"github.com/sm2-cosign/backend/internal/service"
)

func main() {
	configPath := "config.yaml"
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}
	if err := config.Load(configPath); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if err := initDatabase(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer repository.CloseDB()

	if err := service.InitAdminUser(); err != nil {
		log.Printf("Warning: Failed to initialize admin user: %v", err)
	}

	app := fiber.New(fiber.Config{
		AppName:      "SM2 Co-Sign Server v1.0",
		ServerHeader: "SM2-CoSign",
	})

	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	setupRoutes(app)

	go func() {
		addr := fmt.Sprintf(":%d", config.AppConfig.Server.Port)
		if err := app.Listen(addr); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Printf("Server started on port %d", config.AppConfig.Server.Port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	if err := app.Shutdown(); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}
	log.Println("Server stopped")
}

func initDatabase() error {
	if err := repository.InitDB(config.AppConfig.Database.Path); err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}

	schemaPath := "scripts/schema.sql"
	schema, err := os.ReadFile(schemaPath)
	if err != nil {
		log.Printf("Warning: schema.sql not found at %s, skipping initialization", schemaPath)
		return nil
	}

	db := repository.GetDB()
	if _, err := db.Exec(string(schema)); err != nil {
		if err != sql.ErrTxDone {
			log.Printf("Schema execution warning: %v", err)
		}
	}

	log.Println("Database initialized successfully")
	return nil
}

func setupRoutes(app *fiber.App) {
	userHandler := handler.NewUserHandler()
	cosignHandler := handler.NewCosignHandler()
	adminHandler := handler.NewAdminHandler()

	api := app.Group("/api")
	api.Post("/register", userHandler.Register)
	api.Post("/login", userHandler.Login)
	api.Post("/logout", userHandler.Logout)

	authGroup := api.Group("", middleware.AuthMiddleware())
	authGroup.Get("/user/info", userHandler.GetUserInfo)
	authGroup.Post("/key/init", cosignHandler.KeyInit)
	authGroup.Post("/sign", cosignHandler.Sign)
	authGroup.Post("/decrypt", cosignHandler.Decrypt)

	mapi := app.Group("/mapi")
	mapi.Get("/health", adminHandler.Health)
	mapi.Get("/stats", adminHandler.Stats)
	mapi.Get("/users", adminHandler.ListUsers)
	mapi.Get("/users/:id", adminHandler.GetUser)
	mapi.Delete("/users/:id", adminHandler.DeleteUser)
	mapi.Put("/users/:id/status", adminHandler.UpdateUserStatus)
	mapi.Get("/keys", adminHandler.ListKeys)
	mapi.Delete("/keys/:id", adminHandler.DeleteKey)
	mapi.Get("/logs", adminHandler.ListLogs)
}
