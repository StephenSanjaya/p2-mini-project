package routes

import (
	"os"
	"p2-mini-project/docs"
	"p2-mini-project/src/handler"
	"p2-mini-project/src/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func Routes(db *gorm.DB) {
	authService := handler.NewAuthService(db)
	carService := handler.NewCarService(db)
	adminService := handler.NewAdminService(db)
	userService := handler.NewUserService(db)

	r := gin.Default()
	r.Use(middleware.ErrorMiddleware)

	api := r.Group("/api/v1")
	{
		users := api.Group("/users")
		{
			users.POST("/register", authService.RegisterHandler)
			users.POST("/login", authService.LoginHandler)
		}
		authUsers := api.Group("/users")
		authUsers.Use(middleware.AuthMiddleware("user"))
		{
			authUsers.POST("/topup", userService.TopUp)
		}
		cars := api.Group("/cars")
		cars.Use(middleware.AuthMiddleware("user"))
		{
			cars.GET("", carService.GetAllCars)
			cars.GET("/:category_id", carService.GetAllCarsByCategory)
			cars.POST("/rental", carService.RentalCar)
			cars.POST("/pay/:rental_id", carService.PayRentalCar)
			cars.POST("/return/:rental_id", carService.ReturnRentalCar)
		}
		admin := api.Group("/admin/cars")
		admin.Use(middleware.AuthMiddleware("admin"))
		{
			admin.POST("", adminService.CreateNewCar)
			admin.PUT("/:car_id", adminService.UpdateCar)
			admin.DELETE("/:car_id", adminService.DeleteCar)
			admin.GET("/users", adminService.GetAllUsers)
			admin.GET("/rental-history", adminService.GetRentalHistory)
		}
	}

	docs.SwaggerInfo.BasePath = "/"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	r.Run(":" + port)
}
