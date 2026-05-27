package main

import (
	"log"
	"os"

	"github.com/bernaddwiki/koda-b7-weekly10/internal/config"
	"github.com/bernaddwiki/koda-b7-weekly10/internal/router"

	_ "github.com/bernaddwiki/koda-b7-weekly10/docs"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// @title E-Wallet API
// @version 1.0
// @description API documentation for E-Wallet Backend
// @host localhost:8080
// @BasePath /
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("failed load env")
	}

	db, err := config.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	redisClient := config.NewRedis()
	defer redisClient.Close()

	r := gin.Default()

	router.SetupRouter(r, db, redisClient)

	port := os.Getenv("APP_PORT")

	log.Println("server running at port", port)

	r.Run(":" + port)
}
