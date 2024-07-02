package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load()
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("Port Not Bound ..failed to connect from env")

	}

	fmt.Println("PORT :", port)

	router := chi.NewRouter()

	server := &http.Server{

		Handler: router,
		Addr:    ":" + port,
	}

	err := server.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
}
