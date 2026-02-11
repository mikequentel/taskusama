package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// Config
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL not set")
	}

	// DB Pool
	ctx := context.Background()

	dbpool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer dbpool.Close()

	// Test DB
	if err := dbpool.Ping(ctx); err != nil {
		log.Fatal("DB not reachable:", err)
	}

	// Gin
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	log.Println("Listening on :8080")
	r.Run(":8080")
}

