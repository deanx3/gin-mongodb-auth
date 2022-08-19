package main

import (
	"log"
	"time"

	"github.com/deanx3/gin-mongodb-auth/controllers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	//Connecting to DB
	err := mgm.SetDefaultConfig(nil, "efficon", options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Printf(err.Error())
	}

	//starting server
	startGinServer()

}

func startGinServer() {
	AuthController := controllers.NewAuthController()
	HomeController := controllers.NewHomeController()

	server := gin.Default()
	router := server.Group("/api")

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost"},
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		// AllowOriginFunc: func(origin string) bool {
		// 	return origin == "https://github.com"
		// },
		MaxAge: 12 * time.Hour,
	}))

	AuthController.AuthRoutes(router)
	HomeController.HomeRoutes(router)
	server.Run()
}
