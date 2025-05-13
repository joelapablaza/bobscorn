package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"bobscorn/internal/database"
	"bobscorn/internal/handler"
	cornservice "bobscorn/internal/service"
	"bobscorn/internal/storage"
	utils "bobscorn/pkg"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "bob")
	dbName := getEnv("DB_NAME", "corn_db")
	dbSSLMode := getEnv("DB_SSLMODE", "disable")

	connStr := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbName, dbSSLMode)

	db, err := database.NewPostgreSQL(connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize schema (create table if not exists)
	if err := database.InitializeSchema(db); err != nil {
		log.Fatalf("Failed to initialize database schema: %v", err)
	}
	log.Println("Database schema initialized successfully.")

	// Dependency Injection
	rateLimitStorage := storage.NewPostgresRateLimitStorage(db)
	cornSvc := cornservice.NewCornService(rateLimitStorage)
	cornHandler := handler.NewCornHandler(cornSvc)

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "POST, GET, OPTIONS",
	}))

	app.Post("/buy", cornHandler.BuyCorn)

	// Start background cleanup task
	cleanupIntervalMinutesStr := getEnv("CLEANUP_INTERVAL_MINUTES", "5")
	cleanupIntervalMinutes, err := strconv.Atoi(cleanupIntervalMinutesStr)
	if err != nil {
		log.Printf("Warning: Invalid CLEANUP_INTERVAL_MINUTES value '%s', defaulting to 5 minutes. Error: %v", cleanupIntervalMinutesStr, err)
		cleanupIntervalMinutes = 5
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go utils.StartCleanupTask(ctx, rateLimitStorage, time.Duration(cleanupIntervalMinutes)*time.Minute)

	// Graceful shutdown
	// Wait for interrupt signal to gracefully shutdown the server
	// Block until a signal is received
	// Signal cleanup goroutine to stop
	// Give Fiber some time to shutdown gracefully
	go func() {
		if err := app.Listen(":8000"); err != nil {
			log.Fatalf("Error starting Fiber server: %v", err)
		}
	}()
	log.Println("Bob's Corn API started on port 3000")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	cancel()

	if err := app.ShutdownWithTimeout(5 * time.Second); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
