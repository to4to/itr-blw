package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("PORT Connection Failed")
	}

}
