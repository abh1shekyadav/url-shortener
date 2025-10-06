package controllers

import (
	"context"
	"net/http"

	"github.com/abh1shekyadav/url-shortener/repositories"
	"github.com/abh1shekyadav/url-shortener/services"
	"github.com/gin-gonic/gin"
)

var urlService = services.NewURLService(repositories.NewPostgresURLRepository())

func ShortenURL(c *gin.Context) {
	var req struct {
		URL string `json:"url" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	code, err := urlService.ShortenURL(context.Background(), req.URL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to shorten URL"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"short_url": "http://localhost:8080/" + code})
}

func ResolveURL(c *gin.Context) {
	code := c.Param("code")
	original, err := urlService.ResolveURL(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}
	c.Redirect(http.StatusFound, original)
}

func GetStats(c *gin.Context) {
	code := c.Param("code")
	data, err := urlService.GetStats(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"short_code":   data.ShortCode,
		"original_url": data.OriginalURL,
		"click_count":  data.ClickCount,
	})

}
