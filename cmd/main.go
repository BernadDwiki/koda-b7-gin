package main

import (
	"log"
	"os"

	"github.com/bernaddwiki/koda-b7-weekly10/internal/config"
	"github.com/bernaddwiki/koda-b7-weekly10/internal/router"
	"github.com/joho/godotenv"
)

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

	r := router.SetupRouter(db)

	port := os.Getenv("APP_PORT")

	log.Println("server running at port", port)

	r.Run(":" + port)
}
