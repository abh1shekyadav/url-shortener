package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/abh1shekyadav/url-shortener/config"
	"github.com/abh1shekyadav/url-shortener/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	// Connect to Postgres and Redis
	config.ConnectDB()
	config.ConnectRedis()

	// Create Gin router
	r := gin.Default()

	// Register existing application routes
	routes.RegisterRoutes(r)

	// Add health check endpoint
	r.GET("/healthz", func(c *gin.Context) {
		if err := config.DB.Ping(c.Request.Context()); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "db unavailable"})
			return
		}
		if err := config.RedisClient.Ping(c.Request.Context()).Err(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "redis unavailable"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Create HTTP server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// Start server in a separate goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server listen failed: %v\n", err)
		}
	}()

	log.Println("Server running on http://localhost:8080")

	// Wait for termination signal (SIGINT or SIGTERM)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown signal received, shutting down server...")

	// Create context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown server gracefully
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
