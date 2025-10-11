package controllers

import (
	"log"
	"net/http"

	"github.com/abh1shekyadav/url-shortener/repositories"
	"github.com/abh1shekyadav/url-shortener/services"
	"github.com/gin-gonic/gin"
)

var (
	baseRepo   = repositories.NewPostgresURLRepository()
	cacheRepo  = repositories.NewRedisURLRepository(baseRepo)
	urlService = services.NewURLService(cacheRepo)
)

func ShortenURL(c *gin.Context) {
	var req struct {
		URL string `json:"url" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[CONTROLLER] Invalid request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request payload"})
		return
	}

	ctx := c.Request.Context()
	log.Printf("[CONTROLLER] Shorten request for URL: %s", req.URL)

	code, data, err := urlService.ShortenURL(ctx, req.URL)
	if err != nil {
		log.Printf("[CONTROLLER] Error shortening URL: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to shorten URL"})
		return
	}
	if data == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Unexpected nil data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"short_code":   code,
			"short_url":    "http://localhost:8080/" + code,
			"original_url": data.OriginalURL,
			"expires_at":   data.ExpiresAt,
		},
	})
}

func ResolveURL(c *gin.Context) {
	ctx := c.Request.Context()
	code := c.Param("code")

	log.Printf("[CONTROLLER] Resolving shortcode: %s", code)

	original, err := urlService.ResolveURL(ctx, code)
	if err != nil || original == "" {
		log.Printf("[CONTROLLER] Failed to resolve shortcode: %s, err: %v", code, err)
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "URL not found or expired"})
		return
	}

	c.Redirect(http.StatusFound, original)
}

func GetStats(c *gin.Context) {
	ctx := c.Request.Context()
	code := c.Param("code")

	log.Printf("[CONTROLLER] Fetching stats for: %s", code)

	data, err := urlService.GetStats(ctx, code)
	if err != nil || data == nil {
		log.Printf("[CONTROLLER] Error fetching stats for %s: %v", code, err)
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "URL not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"short_code":   data.ShortCode,
			"original_url": data.OriginalURL,
			"click_count":  data.ClickCount,
			"expires_at":   data.ExpiresAt,
		},
	})
}
