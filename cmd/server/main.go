package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	// Load .env (no error if missing; handy for prod where env vars are injected)
	_ = godotenv.Load()

	// Structured logger
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer func() { _ = logger.Sync() }()

	// Config
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		logger.Fatal("DATABASE_URL not set")
	}

	migrationsPath := os.Getenv("MIGRATIONS_PATH") // e.g. "./migrations"
	if migrationsPath == "" {
		migrationsPath = "./migrations"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Run migrations at startup
	if err := runMigrations(dbURL, migrationsPath, logger); err != nil {
		logger.Fatal("migrations failed", zap.Error(err))
	}

	// DB Pool
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbpool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		logger.Fatal("failed to create db pool", zap.Error(err))
	}
	defer dbpool.Close()

	if err := dbpool.Ping(ctx); err != nil {
		logger.Fatal("database not reachable", zap.Error(err))
	}

	// Gin
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(ginZapLogger(logger)) // structured request logs

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	logger.Info("listening", zap.String("port", port))
	if err := r.Run(":" + port); err != nil {
		logger.Fatal("server exited", zap.Error(err))
	}
}

func runMigrations(dbURL, migrationsPath string, logger *zap.Logger) error {
	sourceURL := fmt.Sprintf("file://%s", migrationsPath)
	m, err := migrate.New(sourceURL, dbURL)
	if err != nil {
		return err
	}
	defer func() { _, _ = m.Close() }()

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	if errors.Is(err, migrate.ErrNoChange) {
		logger.Info("migrations up-to-date")
	} else {
		logger.Info("migrations applied")
	}
	return nil
}

// Minimal Gin middleware that logs requests with zap.
func ginZapLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method

		if raw != "" {
			path = path + "?" + raw
		}

		logger.Info("http_request",
			zap.Int("status", status),
			zap.String("method", method),
			zap.String("path", path),
			zap.String("client_ip", clientIP),
			zap.Duration("latency", latency),
		)
	}
}

