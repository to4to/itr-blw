package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/to4to/itr-blw/helper"
	"github.com/to4to/itr-blw/internal/db"
	"github.com/to4to/itr-blw/models"
)

// createConnection initializes the database connection and returns an ApiConfig instance.
func createConnection() models.ApiConfig {
	// Load environment variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Retrieve the database URL from environment variables
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not defined in .env")
	}

	// Establish a connection to the database
	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Can't connect to database:", err)
	}

	// Create a new instance of db.Queries
	dbQueries := db.New(conn)

	// Initialize and return an ApiConfig instance
	apiCfg := models.ApiConfig{
		DB: dbQueries,
	}

	return apiCfg
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	apiCfg := createConnection() // Ensure this returns a struct with a DB field of type *db.Queries

	var params struct {
		Name        string `json:"name"`
		Salary      string `json:"salary"`
		PanNumber   string `json:"pan_number"`
		Year        int32  `json:"year"`
		TaxIncome   string `json:"tax_income"`
		Deductions  string `json:"deductions"`
		Designation string `json:"designation"`
	}

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	deductions := sql.NullString{String: params.Deductions, Valid: params.Deductions != ""}
	employee, err := apiCfg.DB.CreateEmployee(r.Context(), db.CreateEmployeeParams{
		EmployeeID:  uuid.New(),
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
		Name:        params.Name,
		JoiningDate: time.Now().UTC(),
		Salary:      params.Salary,
		PanNumber:   params.PanNumber,
		Year:        params.Year,
		TaxIncome:   params.TaxIncome,
		Deductions:  deductions,
		Designation: params.Designation,
	})

	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Couldn't create employee: %v", err))
		return
	}

	helper.RespondWithJSON(w, http.StatusCreated, employee)
}

func FindUser(w http.ResponseWriter, r *http.Request) {
	apiCfg := createConnection() // Ensure this returns a struct with a DB field of type *db.Queries

	// Using Chi to extract the employee ID from the URL path
	employeeID := chi.URLParam(r, "id") // Assuming the URL pattern includes {id} for the employee ID

	if employeeID == "" {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid employee ID")
		return
	}

	employee, err := apiCfg.DB.FindEmployee(r.Context(), employeeID)
	if err != nil {
		if err == sql.ErrNoRows {
			helper.RespondWithError(w, http.StatusNotFound, "Employee not found")
		} else {
			helper.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error finding employee: %v", err))
		}
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, employee)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	apiCfg := createConnection() // Ensure this returns a struct with a DB field of type *db.Queries

	// Extracting the employee ID from the URL path
	employeeID := chi.URLParam(r, "id")
	if employeeID == "" {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid employee ID")
		return
	}

	var params struct {
		Name        string `json:"name"`
		Salary      string `json:"salary"`
		PanNumber   string `json:"pan_number"`
		TaxIncome   string `json:"tax_income"`
		Deductions  string `json:"deductions"`
		Designation string `json:"designation"`
	}

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	deductions := sql.NullString{String: params.Deductions, Valid: params.Deductions != ""}
	err := apiCfg.DB.UpdateEmployee(r.Context(), db.UpdateEmployeeParams{
		EmployeeID:  employeeID,
		UpdatedAt:   time.Now().UTC(),
		Name:        params.Name,
		Salary:      params.Salary,
		PanNumber:   params.PanNumber,
		TaxIncome:   params.TaxIncome,
		Deductions:  deductions,
		Designation: params.Designation,
	})

	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error updating employee: %v", err))
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, map[string]string{"status": "success", "message": "Employee updated successfully"})
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	apiCfg := createConnection() // Ensure this returns a struct with a DB field of type *db.Queries

	// Extracting the employee ID from the URL path
	employeeID := chi.URLParam(r, "id")
	if employeeID == "" {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid employee ID")
		return
	}

	// Call the delete operation on the database
	err := apiCfg.DB.DeleteEmployee(r.Context(), employeeID)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error deleting employee: %v", err))
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, map[string]string{"status": "success", "message": "Employee deleted successfully"})
}

func FindAllUser(w http.ResponseWriter, r *http.Request) {
	apiCfg := createConnection() // Ensure this returns a struct with a DB field of type *db.Queries

	// Query the database for all employees
	employees, err := apiCfg.DB.FindAllEmployees(r.Context())
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error fetching employees: %v", err))
		return
	}

	// Respond with the list of employees in JSON format
	helper.RespondWithJSON(w, http.StatusOK, employees)
}
