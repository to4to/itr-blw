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
	id, err := strconv.Atoi(params["id"])
	//parsedUUID, err := uuid.Parse(uuidStr)
	if err == nil {
		http.Error(w, "Invalid UUID format", http.StatusBadRequest)
		return
	}
	employee, err := getEmployee(int64(id))

	if err != nil {
		http.Error(w, "Unable to get employee", http.StatusBadRequest)
		return
	}

	// send the response
	json.NewEncoder(w).Encode(employee)
}

// //////////////////////////////////////////////////////
// HandlerGetAllITR is an HTTP handler function that retrieves all ITR (Income Tax Return) information.
// It calls getAllITR function to fetch the data and encodes the result as JSON to the response writer.
func HandlerGetAllITR(w http.ResponseWriter, r *http.Request) {
	employees, err := getAllEmployees()

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
	deletedRows := deleteEmployee(int64(id))

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

//------------------------- handler functions ----------------

// insertEmployee inserts a new employee record into the database and returns the generated employee_id.
func insertEmployee(employee models.Employee) int64 {
	// create a new database connection
	db := createConnection()
	defer db.Close()

	// SQL statement to insert employee details and return the generated employee_id
	sqlStatement := `
        INSERT INTO employees (name, joining_date, salary, pan_number, year, tax_income, deductions, designation)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        RETURNING employee_id
    `

	var id int64
	//created_at, updated_at,joining_date := time.Now(), time.Now(),time.Now()

	// Execute the SQL query and scan the generated employee_id into the id variable
	err := db.QueryRow(sqlStatement,
		// created_at,
		// updated_at,
		employee.Name,
		employee.JoiningDate,
		employee.Salary,
		employee.PanNumber,
		employee.Year,
		employee.TaxIncome,
		employee.Deductions,
		employee.Designation,
	).Scan(&id)

	// Handle any errors that occur during query execution
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	// Print a success message with the generated employee_id
	fmt.Printf("Inserted a single record with ID %v", id)

	return id
}

// getEmployee retrieves an employee from the database based on the provided employeeID.
// It returns the employee details if found, otherwise an error.
func getEmployee(employeeID int64) (models.Employee, error) {
	// Create a new database connection
	db := createConnection()
	defer db.Close() // Ensure the database connection is closed when the function exits

	var employee models.Employee

	// SQL query to select employee details based on employee_id
	sqlStatement := `SELECT * FROM employees WHERE employee_id=$1`

	// Execute the query and get a single row result
	row := db.QueryRow(sqlStatement, employeeID)

	// Scan the row data into the employee struct fields
	err := row.Scan(
		&employee.EmployeeID,
		&employee.CreatedAt,
		&employee.UpdatedAt,
		&employee.Name,
		&employee.JoiningDate,
		&employee.Salary,
		&employee.PanNumber,
		&employee.Year,
		&employee.TaxIncome,
		&employee.Deductions,
		&employee.Designation,
	)

	// Handle different scan outcomes
	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return employee, nil
	case nil:
		return employee, nil
	default:
		log.Fatalf("Unable to scan the row. %v", err)
	}

	return employee, err
}

// getAllEmployees retrieves all employees from the database.
// It returns a slice of models.Employee and an error if any.
func getAllEmployees() ([]models.Employee, error) {
	// Create a new database connection
	db := createConnection()
	defer db.Close()

	var employees []models.Employee

	// SQL statement to select all employees
	sqlStatement := `SELECT * FROM employees`

	// Execute the query and get the result set
	rows, err := db.Query(sqlStatement)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}
	defer rows.Close()

	// Iterate over the result set and scan each row into an employee struct
	for rows.Next() {
		var employee models.Employee

		// Scan the row data into the employee struct fields
		err = rows.Scan(
			&employee.EmployeeID,
			&employee.CreatedAt,
			&employee.UpdatedAt,
			&employee.Name,
			&employee.JoiningDate,
			&employee.Salary,
			&employee.PanNumber,
			&employee.Year,
			&employee.TaxIncome,
			&employee.Deductions,
			&employee.Designation,
		)

		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
		}

		// Append the employee to the employees slice
		employees = append(employees, employee)
	}

	// Return the employees slice and any error that occurred
	return employees, err
}

// updateEmployee updates an employee record in the database based on the provided employeeID.
// It takes the employeeID and updatedEmployee as parameters and returns the number of rows affected by the update operation.
func updateEmployee(employeeID int64, updatedEmployee models.Employee) int64 {
	// create a new database connection
	db := createConnection()
	defer db.Close()

	// SQL statement to update employee details based on employee_id
	sqlStatement := `
		UPDATE employees
		SET name=$2, joining_date=$3, salary=$4, pan_number=$5, year=$6, tax_income=$7, deductions=$8, designation=$9
		WHERE employee_id=$1
	`

	// Execute the update query with the provided employeeID and updatedEmployee data
	res, err := db.Exec(sqlStatement,
		employeeID,
		updatedEmployee.Name,
		updatedEmployee.JoiningDate,
		updatedEmployee.Salary,
		updatedEmployee.PanNumber,
		updatedEmployee.Year,
		updatedEmployee.TaxIncome,
		updatedEmployee.Deductions,
		updatedEmployee.Designation,
	)

	// Handle any errors that occur during query execution
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	// Get the number of rows affected by the update operation
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}

	// Print the total number of rows/records affected by the update operation
	fmt.Printf("Total rows/records affected: %v", rowsAffected)

	return rowsAffected
}

// deleteEmployee deletes an employee from the database based on the provided employeeID.
// It returns the number of rows affected by the deletion operation.
func deleteEmployee(employeeID int64) int64 {
	// createConnection establishes a connection to the database.
	db := createConnection()
	defer db.Close()

	// SQL statement to delete an employee with the specified employee_id.
	sqlStatement := `DELETE FROM employees WHERE employee_id=$1`

	// Execute the delete query with the provided employeeID.
	res, err := db.Exec(sqlStatement, employeeID)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	// Get the number of rows affected by the delete operation.
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}

	// Print the total number of rows/records affected by the delete operation.
	fmt.Printf("Total rows/records affected: %v", rowsAffected)
	return rowsAffected
}
