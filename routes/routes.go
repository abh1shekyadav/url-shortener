package routes

import (
	"github.com/abh1shekyadav/url-shortener/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Phase 1 under construction"})
	})
	r.POST("/shorten", controllers.ShortenURL)
	r.GET("/:code", controllers.ResolveURL)
}
