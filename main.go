package main

import (
	"os"

	"github.com/Ris-Codes/eCommerce02/config"
	"github.com/Ris-Codes/eCommerce02/routes"
	"github.com/gin-gonic/gin"
)

func init() {
	config.LoadEnv()
	config.ConnectDB()
	R.LoadHTMLGlob("templates/*")
}

var R = gin.Default()

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	routes.AdminRoutes(R)
	routes.UserRoutes(R)
	R.Run(":" + port)
}
