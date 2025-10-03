package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var urlStore = make(map[string]string)
var idCounter = 1

func ShortenURL(c *gin.Context) {
	var req struct {
		URL string `json:"url" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	code := strconv.Itoa(idCounter)
	urlStore[code] = req.URL
	idCounter++
	c.JSON(http.StatusOK, gin.H{"short_url": "http://localhost:8080/" + code})
}

func ResolveURL(c *gin.Context) {
	code := c.Param("code")
	original, exists := urlStore[code]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}
	c.Redirect(http.StatusFound, original)
}
