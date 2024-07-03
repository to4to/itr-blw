

-- name: CreateEmployee :one
INSERT INTO employee (
    employee_id, created_at, updated_at, name,
    joining_date, salary, pan_number, year,
    tax_income, deductions, designation
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
)
RETURNING *;



-- name: FindEmployee :one
SELECT *
FROM employee
WHERE employee_id = $1;


-- name: UpdateEmployee :exec
UPDATE employee
SET
    name = $1,
    salary = $2,
    joining_date = $3,
    pan_number = $4,
    year = $5,
    tax_income = $6,
    deductions = $7,
    designation = $8,
    updated_at = NOW()
WHERE employee_id = $9;



-- name: DeleteEmployee :exec
DELETE FROM employee
WHERE employee_id = $1;


-- name: FindAllEmployee :many
SELECT *
FROM employee;