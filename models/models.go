package models

import (
	"time"
	//"github.com/google/uuid"
)

type Employee struct {
	EmployeeID  int64     `json:"employee_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Name        string    `json:"name"`
	JoiningDate time.Time `json:"joining_date"`
	Salary      string    `json:"salary"`
	PanNumber   string    `json:"pan_number"`
	Year        int32     `json:"year"`
	TaxIncome   float64   `json:"tax_income"`
	Deductions  float64   `json:"deductions"`
	Designation string    `json:"designations"`
}
