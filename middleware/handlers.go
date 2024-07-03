package middleware

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

// createConnection establishes a connection to a PostgreSQL database using the POSTGRES_URL environment variable.
// It loads the database connection details from a .env file, opens a connection, and pings the database to ensure connectivity.
// If any errors occur during the process, it logs a fatal error or panics.
// It returns a pointer to the sql.DB database connection.
func createConnection() *sql.DB {
    // Load environment variables from .env file
    err := godotenv.Load(".env")
    if err != nil {
        log.Fatalf("Error loading .env file")
    }

    // Open a connection to the PostgreSQL database using the POSTGRES_URL environment variable
    db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))
    if err != nil {
        panic(err)
    }

    // Ping the database to ensure connectivity
    err = db.Ping()
    if err != nil {
        panic(err)
    }

    fmt.Println("Successfully connected!")

    return db
}

func GetITR(w http.ResponseWriter, r *http.Request) {}

func GetAllITR(w http.ResponseWriter, r *http.Request) {}

func CreateITR(w http.ResponseWriter, r *http.Request) {}

func UpdateITR(w http.ResponseWriter, r *http.Request) {}

func DeleteITR(w http.ResponseWriter, r *http.Request) {}
