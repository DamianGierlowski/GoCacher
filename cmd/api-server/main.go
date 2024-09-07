package main

import (
	"GoCacher/internal/Controller"
	"github.com/gin-gonic/gin"
)

func main() {
	var router = gin.Default()
	router.POST("/insert", Controller.Insert)
	router.POST("/fetch", Controller.Fetch)
	router.Run("0.0.0.0:8080") // Bind to 0.0.0.0 to listen to all network interfaces
}
