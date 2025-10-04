package controllers

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

type URLData struct {
	OriginalURL string `json:"original_url"`
	ClickCount  int    `json:"click_count"`
}

var (
	urlStore  = make(map[string]*URLData)
	idCounter = 1
	mu        sync.Mutex
)

func ShortenURL(c *gin.Context) {
	var req struct {
		URL string `json:"url" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	mu.Lock()
	code := strconv.Itoa(idCounter)
	idCounter++
	urlStore[code] = &URLData{OriginalURL: req.URL, ClickCount: 0}
	mu.Unlock()
	c.JSON(http.StatusOK, gin.H{"short_url": "http://localhost:8080/" + code})
}

func ResolveURL(c *gin.Context) {
	code := c.Param("code")
	mu.Lock()
	data, exists := urlStore[code]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}
	data.ClickCount++
	mu.Unlock()
	c.Redirect(http.StatusFound, data.OriginalURL)
}

func GetStats(c *gin.Context) {
	code := c.Param("code")
	mu.Lock()
	data, exists := urlStore[code]
	mu.Unlock()
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"short_code":   code,
		"original_url": data.OriginalURL,
		"click_count":  data.ClickCount,
	})

}
