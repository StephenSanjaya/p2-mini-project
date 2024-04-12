package routes

import (
	"os"
	"p2-mini-project/src/handler"
	"p2-mini-project/src/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Routes(db *gorm.DB) {
	authService := handler.NewAuthService(db)

	r := gin.Default()
	r.Use(middleware.ErrorMiddleware)

	api := r.Group("/api/v1")
	{
		users := api.Group("/users")
		{
			users.POST("/register", authService.RegisterHandler)
			users.POST("/login", authService.LoginHandler)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	r.Run(":" + port)
}
