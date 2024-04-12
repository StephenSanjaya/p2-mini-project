package routes

import (
	"os"
	"p2-mini-project/src/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Routes(db *gorm.DB) {
	r := gin.Default()
	r.Use(middleware.ErrorMiddleware)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	r.Run(":" + port)
}
