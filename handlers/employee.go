package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go-crud-employee/database"
	"go-crud-employee/models"

	"github.com/gin-gonic/gin"
)

type EmployeeHandler struct {
	db *database.DB
}

func NewEmployeeHandler(db *database.DB) *EmployeeHandler {
	return &EmployeeHandler{
		db: db,
	}
}

// CreateEmployee creates a new employee
func (h *EmployeeHandler) CreateEmployee(c *gin.Context) {
	var req models.CreateEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"Invalid request data",
			err.Error(),
		))
		return
	}

	// Parse hire date
	hireDate, err := time.Parse("2006-01-02", req.HireDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"Invalid hire date format",
			"Use YYYY-MM-DD format",
		))
		return
	}

	// Check if NIP already exists
	var count int
	err = h.db.QueryRow("SELECT COUNT(*) FROM employees WHERE nip = $1", req.NIP).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			"Database error",
			err.Error(),
		))
		return
	}

	if count > 0 {
		c.JSON(http.StatusConflict, models.NewErrorResponse(
			"NIP already exists",
			"Employee with this NIP already exists",
		))
		return
	}

	// Check if email already exists
	err = h.db.QueryRow("SELECT COUNT(*) FROM employees WHERE email = $1", req.Email).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			"Database error",
			err.Error(),
		))
		return
	}

	if count > 0 {
		c.JSON(http.StatusConflict, models.NewErrorResponse(
			"Email already exists",
			"Employee with this email already exists",
		))
		return
	}

	// Insert new employee
	var employee models.Employee
	query := `INSERT INTO employees (nip, name, email, phone, position, department, salary, hire_date, is_active, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) 
			  RETURNING id, nip, name, email, phone, position, department, salary, hire_date, is_active, created_at, updated_at`
	
	now := time.Now()
	var phone sql.NullString
	var salary sql.NullFloat64
	
	if req.Phone != "" {
		phone.String = req.Phone
		phone.Valid = true
	}
	
	if req.Salary > 0 {
		salary.Float64 = req.Salary
		salary.Valid = true
	}

	err = h.db.QueryRow(query, req.NIP, req.Name, req.Email, phone, req.Position, req.Department, salary, hireDate, true, now, now).Scan(
		&employee.ID,
		&employee.NIP,
		&employee.Name,
		&employee.Email,
		&employee.Phone,
		&employee.Position,
		&employee.Department,
		&employee.Salary,
		&employee.HireDate,
		&employee.IsActive,
		&employee.CreatedAt,
		&employee.UpdatedAt,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			"Failed to create employee",
			err.Error(),
		))
		return
	}

	c.JSON(http.StatusCreated, models.NewSuccessResponse(
		"Employee created successfully",
		employee.ToResponse(),
	))
}

// GetEmployees retrieves employees with filtering and pagination
func (h *EmployeeHandler) GetEmployees(c *gin.Context) {
	var filter models.EmployeeFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"Invalid query parameters",
			err.Error(),
		))
		return
	}

	// Set default values
	if filter.Limit <= 0 {
		filter.Limit = 10
	}
	if filter.Limit > 100 {
		filter.Limit = 100
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}

	// Build query
	var conditions []string
	var args []interface{}
	argIndex := 1

	baseQuery := `SELECT id, nip, name, email, phone, position, department, salary, hire_date, is_active, created_at, updated_at FROM employees`
	countQuery := `SELECT COUNT(*) FROM employees`

	if filter.Department != "" {
		conditions = append(conditions, fmt.Sprintf("department = $%d", argIndex))
		args = append(args, filter.Department)
		argIndex++
	}

	if filter.Position != "" {
		conditions = append(conditions, fmt.Sprintf("position = $%d", argIndex))
		args = append(args, filter.Position)
		argIndex++
	}

	if filter.IsActive != nil {
		conditions = append(conditions, fmt.Sprintf("is_active = $%d", argIndex))
		args = append(args, *filter.IsActive)
		argIndex++
	}

	if filter.Search != "" {
		searchPattern := "%" + filter.Search + "%"
		conditions = append(conditions, fmt.Sprintf("(name ILIKE $%d OR email ILIKE $%d OR nip ILIKE $%d)", argIndex, argIndex+1, argIndex+2))
		args = append(args, searchPattern, searchPattern, searchPattern)
		argIndex += 3
	}

	if len(conditions) > 0 {
		whereClause := " WHERE " + strings.Join(conditions, " AND ")
		baseQuery += whereClause
		countQuery += whereClause
	}

	// Get total count
	var total int
	err := h.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			"Database error",
			err.Error(),
		))
		return
	}

	// Add ordering and pagination
	baseQuery += " ORDER BY created_at DESC"
	baseQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, filter.Limit, filter.Offset)

	// Execute query
	rows, err := h.db.Query(baseQuery, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			"Database error",
			err.Error(),
		))
		return
	}
	defer rows.Close()

	var employees []models.EmployeeResponse
	for rows.Next() {
		var emp models.Employee
		err := rows.Scan(
			&emp.ID,
			&emp.NIP,
			&emp.Name,
			&emp.Email,
			&emp.Phone,
			&emp.Position,
			&emp.Department,
			&emp.Salary,
			&emp.HireDate,
			&emp.IsActive,
			&emp.CreatedAt,
			&emp.UpdatedAt,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
				"Database error",
				err.Error(),
			))
			return
		}
		employees = append(employees, emp.ToResponse())
	}

	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			"Database error",
			err.Error(),
		))
		return
	}

	response := models.EmployeeListResponse{
		Employees: employees,
		Total:     total,
		Limit:     filter.Limit,
		Offset:    filter.Offset,
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(
		"Employees retrieved successfully",
		response,
	))
}

// GetEmployee retrieves a single employee by ID
func (h *EmployeeHandler) GetEmployee(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"Invalid employee ID",
			"Employee ID must be a number",
		))
		return
	}

	var employee models.Employee
	query := `SELECT id, nip, name, email, phone, position, department, salary, hire_date, is_active, created_at, updated_at 
			  FROM employees WHERE id = $1`
	
	err = h.db.QueryRow(query, id).Scan(
		&employee.ID,
		&employee.NIP,
		&employee.Name,
		&employee.Email,
		&employee.Phone,
		&employee.Position,
		&employee.Department,
		&employee.Salary,
		&employee.HireDate,
		&employee.IsActive,
		&employee.CreatedAt,
		&employee.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, models.NewErrorResponse(
				"Employee not found",
				"Employee with the specified ID does not exist",
			))
			return
		}
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			"Database error",
			err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(
		"Employee retrieved successfully",
		employee.ToResponse(),
	))
}

// UpdateEmployee updates an existing employee
func (h *EmployeeHandler) UpdateEmployee(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"Invalid employee ID",
			"Employee ID must be a number",
		))
		return
	}

	var req models.UpdateEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"Invalid request data",
			err.Error(),
		))
		return
	}

	// Check if employee exists
	var exists int
	err = h.db.QueryRow("SELECT COUNT(*) FROM employees WHERE id = $1", id).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			"Database error",
			err.Error(),
		))
		return
	}

	if exists == 0 {
		c.JSON(http.StatusNotFound, models.NewErrorResponse(
			"Employee not found",
			"Employee with the specified ID does not exist",
		))
		return
	}

	// Build update query dynamically
	var setParts []string
	var args []interface{}
	argIndex := 1

	if req.Name != "" {
		setParts = append(setParts, fmt.Sprintf("name = $%d", argIndex))
		args = append(args, req.Name)
		argIndex++
	}

	if req.Email != "" {
		// Check if email already exists for other employees
		var count int
		err = h.db.QueryRow("SELECT COUNT(*) FROM employees WHERE email = $1 AND id != $2", req.Email, id).Scan(&count)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
				"Database error",
				err.Error(),
			))
			return
		}
		if count > 0 {
			c.JSON(http.StatusConflict, models.NewErrorResponse(
				"Email already exists",
				"Another employee with this email already exists",
			))
			return
		}

		setParts = append(setParts, fmt.Sprintf("email = $%d", argIndex))
		args = append(args, req.Email)
		argIndex++
	}

	if req.Phone != "" {
		setParts = append(setParts, fmt.Sprintf("phone = $%d", argIndex))
		args = append(args, req.Phone)
		argIndex++
	}

	if req.Position != "" {
		setParts = append(setParts, fmt.Sprintf("position = $%d", argIndex))
		args = append(args, req.Position)
		argIndex++
	}

	if req.Department != "" {
		setParts = append(setParts, fmt.Sprintf("department = $%d", argIndex))
		args = append(args, req.Department)
		argIndex++
	}

	if req.Salary > 0 {
		setParts = append(setParts, fmt.Sprintf("salary = $%d", argIndex))
		args = append(args, req.Salary)
		argIndex++
	}

	if req.IsActive != nil {
		setParts = append(setParts, fmt.Sprintf("is_active = $%d", argIndex))
		args = append(args, *req.IsActive)
		argIndex++
	}

	if len(setParts) == 0 {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"No fields to update",
			"At least one field must be provided for update",
		))
		return
	}

	// Add updated_at
	setParts = append(setParts, fmt.Sprintf("updated_at = $%d", argIndex))
	args = append(args, time.Now())
	argIndex++

	// Add ID for WHERE clause
	args = append(args, id)

	query := fmt.Sprintf("UPDATE employees SET %s WHERE id = $%d RETURNING id, nip, name, email, phone, position, department, salary, hire_date, is_active, created_at, updated_at",
		strings.Join(setParts, ", "), argIndex)

	var employee models.Employee
	err = h.db.QueryRow(query, args...).Scan(
		&employee.ID,
		&employee.NIP,
		&employee.Name,
		&employee.Email,
		&employee.Phone,
		&employee.Position,
		&employee.Department,
		&employee.Salary,
		&employee.HireDate,
		&employee.IsActive,
		&employee.CreatedAt,
		&employee.UpdatedAt,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			"Failed to update employee",
			err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(
		"Employee updated successfully",
		employee.ToResponse(),
	))
}

// DeleteEmployee deletes an employee (soft delete by setting is_active to false)
func (h *EmployeeHandler) DeleteEmployee(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"Invalid employee ID",
			"Employee ID must be a number",
		))
		return
	}

	// Check if employee exists
	var exists int
	err = h.db.QueryRow("SELECT COUNT(*) FROM employees WHERE id = $1", id).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			"Database error",
			err.Error(),
		))
		return
	}

	if exists == 0 {
		c.JSON(http.StatusNotFound, models.NewErrorResponse(
			"Employee not found",
			"Employee with the specified ID does not exist",
		))
		return
	}

	// Soft delete by setting is_active to false
	query := `UPDATE employees SET is_active = false, updated_at = $1 WHERE id = $2`
	_, err = h.db.Exec(query, time.Now(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			"Failed to delete employee",
			err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(
		"Employee deleted successfully",
		nil,
	))
}