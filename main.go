package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/to4to/itr-blw/router"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found")
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable not set")
	}

	r := router.Router()

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
}

// package main

// import (
// 	"fmt"

// 	"log"
// 	"net/http"

// 	"github.com/to4to/itr-blw/router"
// )

// func main() {
//     r := router.Router()
//     // Create a new router instance using the Router function

//     // fs := http.FileServer(http.Dir("build"))
//     // Create a file server handler to serve static files from the "build" directory
//     // http.Handle("/", fs)
//     // Handle requests to the root path "/" by serving static files using the file server

//     fmt.Println("Starting server on the port 8080...")
//     // Print a message indicating that the server is starting on port 8080

//     log.Fatal(http.ListenAndServe(":8080", r))
//     // Start the HTTP server on port 8080 using the router instance 'r' and log any errors
// }
