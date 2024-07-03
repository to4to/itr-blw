-- +goose Up

CREATE TABLE employee (
    employee_id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name TEXT NOT NULL,
    joining_date TIMESTAMP NOT NULL,
    salary NUMERIC NOT NULL,
    pan_number TEXT NOT NULL,
    year INT NOT NULL,
    tax_income NUMERIC NOT NULL,
    deductions NUMERIC,
    designation TEXT NOT NULL
);

-- +goose Down

DROP TABLE employee;
