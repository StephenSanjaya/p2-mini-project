package main

import (
	"p2-mini-project/src/config"
	"p2-mini-project/src/routes"
)

// @title           Mini Project - Rental Car
// @version         1.0
// @description     This is a rental car api docs

// @contact.name   stephen
// @contact.email  stephen@email.com

// @host      localhost:8081
// @BasePath  /api/v1
func main() {
	db := config.GetConnection()

	routes.Routes(db)
}
