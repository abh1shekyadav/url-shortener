package main

import (
	"github.com/abh1shekyadav/url-shortener/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	routes.RegisterRoutes(r)
	r.Run(":8080")
}
