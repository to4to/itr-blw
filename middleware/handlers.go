package middleware

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/to4to/itr-blw/models"
)

// response format
type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

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

// HandlerCreateITR handles the creation of a new employee record.
// It decodes the request body to extract employee information, inserts the employee into the database,
// and responds with the ID of the newly created employee.
func HandlerCreateITR(w http.ResponseWriter, r *http.Request) {
	// Create a connection to the database
	db := createConnection()

	// Defer closing the database connection until the function returns
	defer db.Close()

	// Initialize a variable to hold the decoded employee data
	var employee models.Employee

	// Decode the request body into the employee struct
	err := json.NewDecoder(r.Body).Decode(&employee)

	// Handle any decoding errors
	if err != nil {
		log.Fatalf("Unable to decode the request body. %v", err)
	}

	// Insert the employee into the database and get the insert ID
	insertID := insertEmployee(employee)

	// Prepare the response with the insert ID and success message
	res := response{
		ID:      insertID,
		Message: "Employee created successfully",
	}

	// Encode the response as JSON and write it to the response writer
	json.NewEncoder(w).Encode(res)
}

// /////////////////////////////////////////////////////////////
func HandlerGetITR(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	// convert the id type from string to int
	employee_id, err := strconv.Atoi(params["id"])
	//parsedUUID, err := uuid.Parse(uuidStr)
	if err == nil {
		http.Error(w, "Invalid UUID format", http.StatusBadRequest)
		return

		// call the getStock function with stock id to retrieve a single stock
		employee, err := getemployee(employee_id)

		if err != nil {
			http.Error(w, "Unable to get employee", http.StatusBadRequest)
			return
		}

		// send the response
		json.NewEncoder(w).Encode(employee)
	}
}

// //////////////////////////////////////////////////////
// HandlerGetAllITR is an HTTP handler function that retrieves all ITR (Income Tax Return) information.
// It calls getAllITR function to fetch the data and encodes the result as JSON to the response writer.
func HandlerGetAllITR(w http.ResponseWriter, r *http.Request) {
	employees, err := getAllITR()

	if err != nil {
		log.Fatalf("Unable to get all ITR Info. %v", err)
	}

	json.NewEncoder(w).Encode(employees)
}

///////////////////////////////////////////////////////////////

// HandlerUpdateITR updates an employee record based on the provided ID in the request URL.
func HandlerUpdateITR(w http.ResponseWriter, r *http.Request) {
	// Extract parameters from the request URL
	params := mux.Vars(r)

	// Convert the ID parameter to an integer
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("Unable to convert the string into int. %v", err)
	}

	// Initialize a variable to hold the employee data
	var stock models.Employee

	// Decode the request body into the stock variable
	err = json.NewDecoder(r.Body).Decode(&stock)
	if err != nil {
		log.Fatalf("Unable to decode the request body. %v", err)
	}

	// Update the employee record in the database
	updatedRows := updateEmployee(int64(id), stock)

	// Prepare a success message with the number of rows affected
	msg := fmt.Sprintf("Employee updated successfully. Total rows/record affected %v", updatedRows)

	// Create a response object with the updated employee ID and message
	res := response{
		ID:      int64(id),
		Message: msg,
	}

	// Encode the response object and send it back in the HTTP response
	json.NewEncoder(w).Encode(res)
}

// /////////////////////////////////////////////////////////////////
// HandlerDeleteITR handles the HTTP DELETE request to delete a specific resource.
func HandlerDeleteITR(w http.ResponseWriter, r *http.Request) {
    // Extract the parameters from the request URL
    params := mux.Vars(r)

    // Convert the ID parameter to an integer
    id, err := strconv.Atoi(params["id"])
    if err != nil {
        log.Fatalf("Unable to convert the string into int. %v", err)
    }

    // Delete the resource with the specified ID
    deletedRows := deleteITR(int64(id))

    // Prepare the success message with the number of rows affected
    msg := fmt.Sprintf("Stock Deleted successfully. Total rows/record affected %v", deletedRows)

    // Create a response object with the deleted resource ID and message
    res := response{
        ID:      int64(id),
        Message: msg,
    }

    // Send the response back to the client in JSON format
    json.NewEncoder(w).Encode(res)
}

    // Send the response back to the client in JSON format
    json.NewEncoder(w).Encode(res)
}
