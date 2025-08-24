package models

import (
	"database/sql"
	"time"
)

type Employee struct {
	ID         int            `json:"id" db:"id"`
	NIP        string         `json:"nip" db:"nip" binding:"required"`
	Name       string         `json:"name" db:"name" binding:"required"`
	Email      string         `json:"email" db:"email" binding:"required,email"`
	Phone      sql.NullString `json:"phone" db:"phone"`
	Position   string         `json:"position" db:"position" binding:"required"`
	Department string         `json:"department" db:"department" binding:"required"`
	Salary     sql.NullFloat64 `json:"salary" db:"salary"`
	HireDate   time.Time      `json:"hire_date" db:"hire_date" binding:"required"`
	IsActive   bool           `json:"is_active" db:"is_active"`
	CreatedAt  time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at" db:"updated_at"`
}

type CreateEmployeeRequest struct {
	NIP        string  `json:"nip" binding:"required"`
	Name       string  `json:"name" binding:"required"`
	Email      string  `json:"email" binding:"required,email"`
	Phone      string  `json:"phone"`
	Position   string  `json:"position" binding:"required"`
	Department string  `json:"department" binding:"required"`
	Salary     float64 `json:"salary"`
	HireDate   string  `json:"hire_date" binding:"required"` // Format: YYYY-MM-DD
}

type UpdateEmployeeRequest struct {
	Name       string  `json:"name"`
	Email      string  `json:"email"`
	Phone      string  `json:"phone"`
	Position   string  `json:"position"`
	Department string  `json:"department"`
	Salary     float64 `json:"salary"`
	IsActive   *bool   `json:"is_active"` // Pointer to allow null values
}

type EmployeeFilter struct {
	Department string `form:"department"`
	Position   string `form:"position"`
	IsActive   *bool  `form:"is_active"`
	Search     string `form:"search"` // Search by name, email, or NIP
	Limit      int    `form:"limit"`
	Offset     int    `form:"offset"`
}

type EmployeeResponse struct {
	ID         int     `json:"id"`
	NIP        string  `json:"nip"`
	Name       string  `json:"name"`
	Email      string  `json:"email"`
	Phone      string  `json:"phone"`
	Position   string  `json:"position"`
	Department string  `json:"department"`
	Salary     float64 `json:"salary"`
	HireDate   string  `json:"hire_date"`
	IsActive   bool    `json:"is_active"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
}

type EmployeeListResponse struct {
	Employees []EmployeeResponse `json:"employees"`
	Total     int                `json:"total"`
	Limit     int                `json:"limit"`
	Offset    int                `json:"offset"`
}

// ToResponse converts Employee model to EmployeeResponse
func (e *Employee) ToResponse() EmployeeResponse {
	phone := ""
	if e.Phone.Valid {
		phone = e.Phone.String
	}

	salary := 0.0
	if e.Salary.Valid {
		salary = e.Salary.Float64
	}

	return EmployeeResponse{
		ID:         e.ID,
		NIP:        e.NIP,
		Name:       e.Name,
		Email:      e.Email,
		Phone:      phone,
		Position:   e.Position,
		Department: e.Department,
		Salary:     salary,
		HireDate:   e.HireDate.Format("2006-01-02"),
		IsActive:   e.IsActive,
		CreatedAt:  e.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:  e.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}